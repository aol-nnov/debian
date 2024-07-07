package pkg

import "github.com/aol-nnov/debian/fields"

// represents binary package stanaza
// https://www.debian.org/doc/debian-policy/ch-controlfields.html
type BinaryPackage struct {
	Name         string              `deb822:"Package" required:"true"`
	Architecture fields.Architecture `required:"true"`
	Section      string              `deb822:",omitempty" recommended:"true"`
	Priority     string              `deb822:",omitempty" recommended:"true"`
	Essential    string

	Depends    fields.Dependencies `deb822:",omitempty" delim:","` //
	Recommends fields.Dependencies `deb822:",omitempty" delim:","`
	Suggests   fields.Dependencies `deb822:",omitempty" delim:","`
	Enhances   fields.Dependencies `deb822:",omitempty" delim:","`
	PreDepends fields.Dependencies `deb822:"Pre-Depends,omitempty" delim:","`

	Description string `required:"true"`
	Homepage    string `deb822:",omitempty"`

	Provides []string `deb822:",omitempty" delim:"," strip:" "`

	MultiArch fields.MultiArch `deb822:"Multi-Arch,omitempty"`
}

func (pkg BinaryPackage) Satisfies(dep fields.Dependency, buildArch fields.Architecture,
	hostArch fields.Architecture, profiles []string) bool {
	return false
}
