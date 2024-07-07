/*
# Decoder for deb822 format

Deb822 (https://manpages.debian.org/unstable/dpkg-dev/deb822.5.en.html) is based on Internet Message Format (RFC 5322).

This decoder follows `encoding` package principles. It also supports `TextUnmarshaler` interface for custom types.

Following struct tags are supported:

  - `deb822:"field-name"` - denotes debian control field name to be decoded into the struct field. Use minus sign to ignore the field;

  - `required:"true|false"` - denotes a required field, default false. If field is missing, decoding fails;

  - `if_missing:"otherFieldName"` - if tagged field is missing in the source stream, assign the value of `otherFieldName` to it;

  - `strip:"cutset"` - while decoding to a slice, removes all leading and trailing Unicode code points contained in `cutset` first hand. See [pkg/strings.Trim] for more info. We do not strip anything by default.;

  - `delim:"sep"` - while decoding to a slice, substrings (future slice elements) are separated by `sep`. See [pkg/strings.Split] for more info. Default delimiter is " " (single space).
*/
package deb822

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"
)

type stanza map[string]string

// A Decoder reads and decodes deb822 values from an input stream.
type Decoder struct {
	stanza stanza
	err    error
	reader *bufio.Reader
	atEOF  bool
}

/*
NewDecoder returns a new decoder that reads from reader.

The decoder introduces its own buffering and may read data from r beyond the JSON values requested.

Note: it drains supplied io.Reader, so do not use it after decoding!
*/
func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		reader: bufio.NewReader(reader),
		atEOF:  false,
	}
}

// reads single stanza from reader
// returns true if there are more stanzas left
func (d *Decoder) readStanza() bool {
	if d.atEOF {
		return false
	}

	d.stanza = make(map[string]string)
	d.err = nil

	var lastKey string

	for {
		line, err := d.reader.ReadString('\n')
		if err == io.EOF && line != "" {
			err = nil
			// each paragraph is terminated by a newline, if it is not the case for the last one, let's add it by
			// ourselves
			line = line + "\n"
		}
		if err == io.EOF {
			// if we have a stanza after reaching EOF, it's not an error
			// subsequent [readStanza] call will return false. We're fine.
			if len(d.stanza) > 0 {
				d.atEOF = true
				return true
			}
			// It's an error otherwise. Propagate it.
			d.err = err
			return false
		} else if err != nil {
			d.err = err
			return false
		}

		if line == "\n" || line == "\r\n" {
			if len(d.stanza) == 0 {
				// Skip over any number of blank lines between paragraphs.
				continue
			}
			// Stanza is parsed otherwise. Ready for field assignment.
			return true
		}

		if strings.HasPrefix(line, "#") {
			continue // skip comments
		}

		/*
			So we have a line in one of the following formats:

			  "Key: Value"
			  " Foobar"

			 Foobar is seen as a continuation of the last line, and the Key line is a Key/Value mapping.
		*/

		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			/*
				This is a continuation line; so we're going to go ahead and clean it up, and add it into the list. We're
				going to remove the first character (which we now know is whitespace), and if it's a line that only has
				a dot on it, we'll remove that too (since " .\n" is actually "\n"). We only trim off space on the right
				hand, because indentation under the whitespace is up to the data format. Not us.
			*/

			// TrimFunc(line[1:], unicode.IsSpace) is identical to calling TrimSpace.
			line = strings.TrimRightFunc(line[1:], unicode.IsSpace)

			if d.stanza[lastKey] == "" {
				d.stanza[lastKey] = line + "\n"
			} else {
				if !strings.HasSuffix(d.stanza[lastKey], "\n") {
					d.stanza[lastKey] = d.stanza[lastKey] + "\n"
				}

				// no need to add a newline to the multiline temporary storage,
				// otherwise we'll end up witn an extra paragraph during parsing
				d.stanza[lastKey] = d.stanza[lastKey] + line
			}
			continue
		}

		// Key: Value line parsing
		els := strings.SplitN(line, ":", 2)
		if len(els) != 2 {
			d.err = fmt.Errorf("bad line: '%s' has no ':'", line)
			return false
		}

		/* We'll go ahead and take off any leading spaces */
		lastKey = strings.TrimSpace(els[0])
		value := strings.TrimSpace(els[1])

		d.stanza[lastKey] = value
	}
}

// decodes single stanza into res
// decodes all stanzas till EOF if res is a slice
func (d *Decoder) Decode(res any) error {
	into := reflect.ValueOf(res)

	if into.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("Decode can only decode into a pointer")
	}

	switch into.Elem().Type().Kind() {
	case reflect.Struct:
		d.readStanza()
		if d.err != nil {
			return d.err
		}
		return decodeStruct(d.stanza, into)
	case reflect.Slice:
		for d.readStanza() {
			if d.err != nil {
				return d.err
			}
			itemType := into.Elem().Type().Elem()

			item := reflect.New(itemType)

			if err := decodeStruct(d.stanza, item); err == nil {
				into.Elem().Set(reflect.Append(into.Elem(), item.Elem()))
			} else {
				return err
			}
		}
	default:
		return fmt.Errorf("unable to decode into a %s", into.Elem().Type().Name())
	}

	return nil
}
