package changelog

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/aol-nnov/debian/fields"
)

var (
	headerRe  = regexp.MustCompile(`^([\w-]*)\s+\(([^\(\) \t]+)\)\s+(.*);\s+(.*)`)
	trailerRe = regexp.MustCompile(`^ -- (.*) <(.*)>\s+(.*)`)
)

type decoderState int

const (
	idle decoderState = iota
	doneHeader
	doneHeaderSeparator
	doneTrailer
)

type Decoder struct {
	err    error
	reader *bufio.Reader
	atEOF  bool
	state  decoderState
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		reader: bufio.NewReader(reader),
		atEOF:  false,
		state:  idle,
	}
}

func (d *Decoder) readAndDecodeStanza(entry *Entry) bool {
	if d.atEOF {
		return false
	}

	body := ""
	for {
		line, err := d.reader.ReadString('\n')

		if err == io.EOF /*&& line != ""*/ {
			d.atEOF = true
			err = nil
		}

		d.err = err
		if d.err != nil {
			return !d.atEOF
		}

		switch d.state {
		case idle:
			d.err = decodeHeader(line, entry)
			d.state = doneHeader
		case doneHeader:
			if line == "\n" {
				d.state = doneHeaderSeparator
			} else {
				d.err = fmt.Errorf("changelog format error: missing header separator")
			}
		case doneHeaderSeparator:
			if !trailerRe.MatchString(line) {
				body += strings.TrimLeft(line, " ")
				// fmt.Println("adding", strings.TrimLeft(line, " "), ".")
			} else {
				// remove any leading and trailing newlines
				entry.SetBody(body)

				d.err = decodeTrailer(line, entry)
				d.state = doneTrailer
			}
		case doneTrailer:
			if line == "\n" || d.atEOF {
				d.state = idle
			} else {
				d.err = fmt.Errorf("changelog format error: missing trailer separator")
			}
			return !d.atEOF
		}

		if d.err != nil {
			return !d.atEOF
		}
	}
}

func (d *Decoder) Decode(res any) error {
	into := reflect.ValueOf(res)
	if into.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("Decode can only decode into a pointer")
	}

	switch into.Elem().Type().Kind() {
	case reflect.Struct:
		d.readAndDecodeStanza(res.(*Entry))
		return d.err

	case reflect.Slice:
		item := Entry{}
		for d.readAndDecodeStanza(&item) {
			if d.err != nil {
				return d.err
			}
			into.Elem().Set(reflect.Append(into.Elem(), reflect.ValueOf(item)))
			item = Entry{}
		}
	default:
		return fmt.Errorf("unable to decode into a %s", into.Elem().Type().Name())
	}

	return nil
}

func decodeHeader(line string, entry *Entry) error {
	if !headerRe.MatchString(line) {
		entry = &Entry{}
		return fmt.Errorf("changelog entry header format error")
	}

	headerMatches := headerRe.FindAllStringSubmatch(line, -1)[0]

	entry.PackageName = headerMatches[1]
	entry.Version = fields.MakeVersion(headerMatches[2])
	entry.Distribution = headerMatches[3]
	entry.Metadata = headerMatches[4]

	return nil
}

func decodeTrailer(line string, entry *Entry) error {
	trailerMatches := trailerRe.FindAllStringSubmatch(line, -1)[0]
	entry.Maintainer.Name = trailerMatches[1]
	entry.Maintainer.Email = trailerMatches[2]

	if timestamp, err := time.Parse(time.RFC1123Z, trailerMatches[3]); err != nil {
		entry = &Entry{}
		return fmt.Errorf("changelog entry timestamp format error")
	} else {

		entry.Timestamp = Timestamp(timestamp)
	}

	return nil
}
