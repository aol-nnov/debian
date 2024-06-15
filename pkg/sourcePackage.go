package pkg

import "github.com/aol-nnov/debian/fields"

// https://www.debian.org/doc/debian-policy/ch-controlfields.html#source-package-control-files-debian-control
type SourcePackage struct {
	Name             string         `control:"Source" required:"true"`
	Maintainer       string         `required:"true"`
	Section          string         `recommended:"true"`
	Priority         string         `recommended:"true"`
	StandardsVersion fields.Version `control:"Standards-Version" required:"true"`
	Description      string

	BuildDepends      fields.Dependencies `control:"Build-Depends" delim:"," strip:" "`
	BuildDependsArch  fields.Dependencies `control:"Build-Depends-Arch" delim:"," strip:" "`
	BuildDependsIndep fields.Dependencies `control:"Build-Depends-Indep" delim:"," strip:" "`
}
