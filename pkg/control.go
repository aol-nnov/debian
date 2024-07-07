package pkg

import (
	"io"

	"github.com/aol-nnov/debian/deb822"
	"github.com/aol-nnov/debian/fields"
	"github.com/aol-nnov/debian/internal/universalreader"
)

type binaryPackageInSrc struct {
	Name         string              `deb822:"Package" required:"true"`
	Architecture fields.Architecture `required:"true"`
	Section      string              `recommended:"true"`
	Priority     string              `recommended:"true"`
	Essential    string

	Depends    []string `delim:"," strip:"\n "`
	Recommends []string `delim:"," strip:"\n "`
	Suggests   []string `delim:"," strip:"\n "`
	Enhances   []string `delim:"," strip:"\n "`
	PreDepends []string `deb822:"Pre-Depends" delim:"," strip:"\n "`

	Description string `required:"true"`
	Homepage    string

	Provides []string `delim:"," strip:"\n "`

	MultiArch fields.MultiArch `deb822:"Multi-Arch"`
}

type Control struct {
	DebSrc SourcePackage
	Deb    []binaryPackageInSrc
}

func (c *Control) Decode(in io.Reader) error {
	defer universalreader.MaybeClose(in)

	decoder := deb822.NewDecoder(in)

	if err := decoder.Decode(&c.DebSrc); err != nil {
		return err
	}

	if err := decoder.Decode(&c.Deb); err != nil {
		return err
	}

	return nil
}
