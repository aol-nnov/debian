package changelog

import (
	"fmt"
	"io"
	"os"
)

type Changelog struct {
	changelogFileName string
	Entries           []Entry
}

func New() *Changelog {
	return &Changelog{
		changelogFileName: "./debian/changelog",
		Entries:           []Entry{},
	}
}

func Load() (*Changelog, error) {
	c := New()

	changelogReader, err := os.Open(c.changelogFileName)

	if err != nil {
		return nil, err
	}
	defer changelogReader.Close()

	c.Entries = make([]Entry, 1)
	if err := NewDecoder(changelogReader).Decode(&c.Entries[0]); err != nil {
		return nil, err
	}

	return c, nil
}

func LoadFull() (*Changelog, error) {
	c := New()

	changelogReader, err := os.Open(c.changelogFileName)

	if err != nil {
		return nil, err
	}
	defer changelogReader.Close()

	if err := NewDecoder(changelogReader).Decode(&c.Entries); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Changelog) Last() Entry {

	return c.Entries[0]
}

func (c *Changelog) AddEntry(e Entry) error {
	tmpFileName := fmt.Sprintf("%s.tmp", c.changelogFileName)

	tmpFile, err := os.Create(tmpFileName)

	if err != nil {
		return err
	}

	origFile, err := os.Open(c.changelogFileName)

	if err != nil {
		return err
	}

	defer origFile.Close()

	if _, err := tmpFile.WriteString(e.String()); err != nil {
		tmpFile.Close()
		os.Remove(tmpFileName)

		return err
	}

	if _, err := io.Copy(tmpFile, origFile); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFileName, c.changelogFileName)
}

func (c *Changelog) ReplaceLastEntry(e Entry) error {
	tmpFileName := fmt.Sprintf("%s.tmp", c.changelogFileName)

	tmpFile, err := os.Create(tmpFileName)

	if err != nil {
		return err
	}

	origFile, err := os.Open(c.changelogFileName)

	if err != nil {
		return err
	}

	defer origFile.Close()

	if _, err := tmpFile.WriteString(e.String()); err != nil {
		tmpFile.Close()
		os.Remove(tmpFileName)

		return err
	}

	// now move file descriptor of original file to the second record.
	// just read the first record out for simplicity

	var unused Entry
	d := NewDecoder(origFile)
	d.Decode(&unused)

	// then copy the rest
	// !!! use d.reader, as there is a bufio.Reader in decoder and it already advances underlying descriptor by the
	// buffer size (defaults to 4k)
	if _, err := io.Copy(tmpFile, d.reader); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFileName, c.changelogFileName)
}
