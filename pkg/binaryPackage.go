package pkg

import "github.com/aol-nnov/debian/fields"

// represents binary package stanaza
// https://www.debian.org/doc/debian-policy/ch-controlfields.html
type BinaryPackage struct {
	Name         string              `control:"Package" required:"true"`
	Architecture fields.Architecture `required:"true"`
	Section      string              `recommended:"true"`
	Priority     string              `recommended:"true"`
	Essential    string

	Depends    fields.Dependencies `delim:","` //
	Recommends fields.Dependencies `delim:","`
	Suggests   fields.Dependencies `delim:","`
	Enhances   fields.Dependencies `delim:","`
	PreDepends fields.Dependencies `control:"Pre-Depends" delim:","`

	Description string `required:"true"`
	Homepage    string

	Provides []string `delim:"," strip:" "`

	MultiArch fields.MultiArch `control:"Multi-Arch"`
}

func (pkg BinaryPackage) Satisfies(dep fields.Dependency, buildArch fields.Architecture,
	hostArch fields.Architecture, profiles []string) bool {
	return false
}
