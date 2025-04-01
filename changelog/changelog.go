package changelog

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/aol-nnov/debian/fields"
)

type Changelog struct {
	changelogFileName string
	Entries           []Entry
	lastParsed        Entry
	changelogReader   *os.File
	decoder           *Decoder
}

func New() *Changelog {
	return &Changelog{
		changelogFileName: "./debian/changelog",
		Entries:           []Entry{},
		changelogReader:   nil,
		decoder:           nil,
	}
}

func Load() (*Changelog, error) {
	var err error

	c := New()
	c.changelogReader, err = os.Open(c.changelogFileName)

	if err != nil {
		return nil, err
	}

	c.decoder = NewDecoder(c.changelogReader)

	if err := c.decoder.Decode(&c.lastParsed); err != nil {
		return nil, err
	}

	return c, nil
}

func LoadFull() (*Changelog, error) {
	tmp := New()

	changelogReader, err := os.Open(tmp.changelogFileName)

	if err != nil {
		return nil, err
	}

	// read and parse the whole file
	if err := NewDecoder(changelogReader).Decode(&tmp.Entries); err != nil {
		return nil, err
	}

	// set the whole machinery to the initial position:
	// one record parsed, c.decoder.reader is pointing to the second record
	c, err := Load()
	if err != nil {
		return nil, err
	}

	c.Entries = tmp.Entries

	return c, nil
}

func (c *Changelog) finalize() error {
	if err := c.changelogReader.Close(); err != nil {
		return err
	}

	c.decoder = nil
	c.lastParsed = Entry{}
	c.Entries = nil

	return nil
}

func (c *Changelog) Last() Entry {

	if len(c.Entries) > 0 {
		return c.Entries[0]
	}

	return c.lastParsed
}

func (c *Changelog) SkipSnapshotOrDistribution(extraDistributionsToSkip []string) error {
	distributionsToSkip := []string{"UNRELEASED"}
	distributionsToSkip = append(distributionsToSkip, extraDistributionsToSkip...)

	if err := c.decoder.Decode(&c.lastParsed); err != nil {
		return err
	}

	// start peeking records and skip not needed ones
	for c.lastParsed.Version.IsMod() == fields.VersionModSnapshot ||
		slices.Contains(distributionsToSkip, c.lastParsed.Distribution) {
		if err := c.decoder.Decode(&c.lastParsed); err != nil {
			return err
		}
	}

	return nil
}

func (c *Changelog) AddEntry(e Entry) error {
	tmpFileName := fmt.Sprintf("%s.tmp", c.changelogFileName)

	tmpFile, err := os.Create(tmpFileName)

	if err != nil {
		return err
	}

	// write new record
	if _, err := tmpFile.WriteString(e.String()); err != nil {
		tmpFile.Close()
		os.Remove(tmpFileName)

		return err
	}

	// write last already parsed recodr
	if _, err := tmpFile.WriteString(c.lastParsed.String()); err != nil {
		tmpFile.Close()
		os.Remove(tmpFileName)

		return err
	}

	// ... then write original changelog tail
	if _, err := io.Copy(tmpFile, c.decoder.reader); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := c.finalize(); err != nil {
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

	// write new entry
	if _, err := tmpFile.WriteString(e.String()); err != nil {
		tmpFile.Close()
		os.Remove(tmpFileName)

		return err
	}

	// at this point we have one record read out already (by changelog.Load())

	// then copy the rest
	// !!! use c.decoder.reader, as there is a bufio.Reader in decoder and it already advances underlying descriptor by
	// the buffer size (defaults to 4k)
	if _, err := io.Copy(tmpFile, c.decoder.reader); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := c.finalize(); err != nil {
		return err
	}

	return os.Rename(tmpFileName, c.changelogFileName)
}
