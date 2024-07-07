package changes

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
	"github.com/aol-nnov/debian/deb822"
)

type Changes struct {
	Format       string
	Date         string
	Source       string
	Binary       string
	Architecture string
	Version      string
	Distribution string
	Urgency      string
	Maintainer   string
	ChangedBy    string `deb822:"Changed-By"`
	Description  string
	Changes      string
	Sha1Files    []ChecksummedFile `deb822:"Checksums-Sha1" delim:"\n" s_trip:" \n"`
	Sha256Files  []ChecksummedFile `deb822:"Checksums-Sha256" delim:"\n" s_trip:" \n"`
	Files        []ChangesFile     `delim:"\n" s_trip:" \n"`
}

func FromStream(r io.Reader) (*Changes, error) {
	var c Changes
	var buf bytes.Buffer
	buf.ReadFrom(r)
	block, plain := clearsign.Decode(buf.Bytes())

	var changelog []byte

	if block != nil {
		changelog = block.Plaintext
	} else {
		changelog = plain
	}

	if err := deb822.NewDecoder(bytes.NewReader(changelog)).Decode(&c); err == nil {
		return &c, nil
	} else {

		return nil, err
	}
}

type ChecksummedFile struct {
	Checksum string
	Filesize string
	Path     string
}

func (rf *ChecksummedFile) UnmarshalText(text []byte) (err error) {
	tmp := bytes.Fields(text)
	if len(tmp) != 3 {
		return fmt.Errorf("unable to unmarshal ChecksummedFile record '%s'", text)
	}
	rf.Checksum = string(tmp[0])
	rf.Filesize = string(tmp[1])
	rf.Path = string(tmp[2])

	return nil
}

type ChangesFile struct {
	ChecksummedFile
	Section  string
	Priority string
}

func (rf *ChangesFile) UnmarshalText(text []byte) (err error) {
	tmp := bytes.Fields(text)
	if len(tmp) != 5 {
		return fmt.Errorf("unable to unmarshal ChangesFile record '%s'", text)
	}

	rf.Checksum = string(tmp[0])
	rf.Filesize = string(tmp[1])
	rf.Section = string(tmp[2])
	rf.Priority = string(tmp[3])
	rf.Path = string(tmp[4])

	return nil
}
