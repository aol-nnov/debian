package changelog_test

import (
	"fmt"
	"time"

	"github.com/aol-nnov/debian/changelog"
	"github.com/aol-nnov/debian/fields"
)

func ExampleEntry_SetBody() {
	e := changelog.NewEntry()
	e.PackageName = "coolpkg"
	e.Version = fields.MakeVersion("1.2.3")
	e.Distribution = "next"
	e.Maintainer.Name = "maint"
	e.Maintainer.Email = "qwe@asd.zxc"
	ts, _ := time.Parse(time.RFC1123Z, "Tue, 01 Oct 2024 16:12:39 +0300")
	e.Timestamp = changelog.Timestamp(ts)

	body := `lalala

SrcRef: deadbeef
	`
	e.SetBody(body)

	fmt.Println(e)

	// Output:
	// coolpkg (1.2.3) next; urgency=medium
	//
	//   lalala
	//
	//   SrcRef: deadbeef
	//
	//  -- maint <qwe@asd.zxc>  Tue, 01 Oct 2024 16:12:39 +0300
}

func ExampleEntry_AddTag_maintain_order() {
	e := changelog.NewEntry()
	e.PackageName = "coolpkg"
	e.Version = fields.MakeVersion("1.2.3")
	e.Distribution = "next"
	e.SetBody("lalala")
	e.Maintainer.Name = "maint"
	e.Maintainer.Email = "qwe@asd.zxc"
	ts, _ := time.Parse(time.RFC1123Z, "Tue, 01 Oct 2024 16:12:39 +0300")
	e.Timestamp = changelog.Timestamp(ts)

	e.AddTag(changelog.BuildRefTag, "asd")
	e.AddTag(changelog.SrcRefTag, "qwe")

	fmt.Println(e)

	// Output:
	// coolpkg (1.2.3) next; urgency=medium
	//
	//   lalala
	//
	//   SrcRef: qwe
	//   BuildRef: asd
	//
	//  -- maint <qwe@asd.zxc>  Tue, 01 Oct 2024 16:12:39 +0300
}

func ExampleEntry_GetTag() {
	e := changelog.NewEntry()
	e.PackageName = "coolpkg"
	e.Version = fields.MakeVersion("1.2.3")
	e.Distribution = "next"
	e.SetBody(`lalala

SrcRef: deadbeef`)

	e.Maintainer.Name = "maint"
	e.Maintainer.Email = "qwe@asd.zxc"
	ts, _ := time.Parse(time.RFC1123Z, "Tue, 01 Oct 2024 16:12:39 +0300")
	e.Timestamp = changelog.Timestamp(ts)

	fmt.Println(e.GetTag(changelog.SrcRefTag))

	// Output: deadbeef
}

func ExampleEntry_GetTag_uninitialized() {
	e := &changelog.Entry{}

	fmt.Println(e.GetTag(changelog.SrcRefTag))

	// Output:
}
