package changelog_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/aol-nnov/debian/changelog"
)

func ExampleLoad() {

	c, _ := changelog.Load()

	fmt.Println(len(c.Entries))
	fmt.Println(c.Last())

	// Output:
	// 1
	// pkg-name (3.2.16) next; urgency=medium
	//
	//   [ Author1 ]
	//   * Добавлен слот шины для инициализации плагинов микшера
	//
	//   SrcRef: deadbeef
	//
	//  -- Package Maintainer <pkg-maint@example.net>  Mon, 19 Dec 2022 11:50:13 +0000
}

func ExampleLoadFull() {
	c, _ := changelog.LoadFull()

	fmt.Println(len(c.Entries))

	// Output: 5
}

func TestReplaceEntry(t *testing.T) {
	origFileName := "./debian/changelog.orig"
	fileName := "./debian/changelog"

	changelogFile, err := os.Open(fileName)

	if err != nil {
		t.Fatal(err)
	}

	origFile, err := os.Create(origFileName)

	if err != nil {
		t.Fatal(err)
	}

	io.Copy(origFile, changelogFile)

	defer os.Rename(origFileName, fileName)

	c, err := changelog.Load()

	if err != nil {
		t.Fatal(err)
	}

	entry := c.Last()

	entry.SetBody("mooo")

	if err := c.ReplaceLastEntry(entry); err != nil {
		t.Fatal(err)
	}

	newChangelog, err := changelog.Load()

	if err != nil {
		t.Fatal(err)
	}

	if newChangelog.Last().GetBody() != "mooo" {
		t.Fatal("Replacing entry failed")
	}
}

func TestAddEntry(t *testing.T) {
	origFileName := "./debian/changelog.orig"
	fileName := "./debian/changelog"

	changelogFile, err := os.Open(fileName)

	if err != nil {
		t.Fatal(err)
	}

	origFile, err := os.Create(origFileName)

	if err != nil {
		t.Fatal(err)
	}

	io.Copy(origFile, changelogFile)

	defer os.Rename(origFileName, fileName)

	c, err := changelog.LoadFull()
	entriesCount := len(c.Entries)

	if err != nil {
		t.Fatal(err)
	}

	entry := changelog.NewEntryFromTemplate(c.Last())

	entry.SetBody("mooo")

	if err := c.AddEntry(entry); err != nil {
		t.Fatal(err)
	}

	newChangelog, err := changelog.LoadFull()

	if err != nil {
		t.Fatal(err)
	}

	if newChangelog.Last().GetBody() != "mooo" && len(newChangelog.Entries) != entriesCount+1 {
		t.Fatal("Adding entry failed")
	}
}
