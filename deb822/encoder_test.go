package deb822_test

import (
	"fmt"
	"os"

	"github.com/aol-nnov/debian/deb822"
	"github.com/aol-nnov/debian/fields"
	"github.com/aol-nnov/debian/pkg"
)

func ExampleEncoder_string_noFieldName() {
	fmt.Println(deb822.Marshal("Boo"))

	// Output:
	// unable to encode from a string
}

func ExampleEncoder_Description() {
	dStr := `First line
second line

Detailed description`

	d := fields.Description(dStr)

	fmt.Println(deb822.Marshal(&d))

	// Output:
	// unable to encode from a *fields.Description
}

func ExampleEncoder_slice_pkg_SourcePackage() {
	in := []pkg.SourcePackage{
		{
			Name:       "name",
			Maintainer: "maint <qwe@asd.zxc>",
			Section:    "libs",
			Priority:   "optional",
			StandardsVersion: fields.Version{
				UpstreamVersion: "1.2.3",
			},
			Description: `descr
Second line

After empty line`,
			BuildDepends: []fields.Dependency{
				{
					Name: "one",
				},
				{
					Name: "two",
				},
			},
			BuildDependsArch:  []fields.Dependency{},
			BuildDependsIndep: []fields.Dependency{},
		},
		{
			Name:       "another",
			Maintainer: "maint <qwe@asd.zxc>",
			Section:    "libs",
			Priority:   "optional",
			StandardsVersion: fields.Version{
				UpstreamVersion: "1.2.3",
			},
			// Description:       "descr",
			BuildDepends:      []fields.Dependency{},
			BuildDependsArch:  []fields.Dependency{},
			BuildDependsIndep: []fields.Dependency{},
		},
	}
	fmt.Println(deb822.NewEncoder(os.Stdout).Encode(in))

	// Output:
	// Source: name
	// Maintainer: maint <qwe@asd.zxc>
	// Section: libs
	// Priority: optional
	// Standards-Version: 1.2.3
	// Description: descr
	//  Second line
	//  .
	//  After empty line
	// Build-Depends: one,
	//  two
	//
	// Source: another
	// Maintainer: maint <qwe@asd.zxc>
	// Section: libs
	// Priority: optional
	// Standards-Version: 1.2.3
	// encodeStruct: missing value for required field 'Description'
}
