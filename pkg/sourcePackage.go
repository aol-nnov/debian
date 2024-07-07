package pkg

import "github.com/aol-nnov/debian/fields"

// https://www.debian.org/doc/debian-policy/ch-controlfields.html#source-package-control-files-debian-control
type SourcePackage struct {
	Name             string             `deb822:"Source" required:"true"`
	Maintainer       string             `required:"true"`
	Section          string             `deb822:",omitempty" recommended:"true"`
	Priority         string             `deb822:",omitempty" recommended:"true"`
	StandardsVersion fields.Version     `deb822:"Standards-Version" required:"true"`
	Description      fields.Description `required:"true"`

	BuildDepends      fields.Dependencies `deb822:"Build-Depends,omitempty" delim:"," strip:" "`
	BuildDependsArch  fields.Dependencies `deb822:"Build-Depends-Arch,omitempty" delim:"," strip:" "`
	BuildDependsIndep fields.Dependencies `deb822:"Build-Depends-Indep,omitempty" delim:"," strip:" "`
}
