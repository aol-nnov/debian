package changelog_test

import (
	"fmt"
	"os"

	"github.com/aol-nnov/debian/changelog"
)

func ExampleDecoder_Decode_single() {
	changelogReader, _ := os.Open("./debian/changelog")

	var e changelog.Entry

	changelog.NewDecoder(changelogReader).Decode(&e)

	fmt.Println(e)

	// Output:
	// pkg-name (3.2.16) next; urgency=medium
	//
	//   [ Author1 ]
	//   * Добавлен слот шины для инициализации плагинов микшера
	//
	//   SrcRef: deadbeef
	//
	//  -- Package Maintainer <pkg-maint@example.net>  Mon, 19 Dec 2022 11:50:13 +0000
}

func ExampleDecoder_Decode_tags() {
	changelogReader, _ := os.Open("./debian/changelog")

	var e changelog.Entry

	changelog.NewDecoder(changelogReader).Decode(&e)

	fmt.Println(e.GetTag(changelog.SrcRefTag))

	// Output: deadbeef
}

func ExampleDecoder_Decode_slice() {
	changelogReader, _ := os.Open("./debian/changelog")

	var e []changelog.Entry

	changelog.NewDecoder(changelogReader).Decode(&e)

	fmt.Println(len(e))
	fmt.Println(e[1])

	// Output:
	// 5
	// pkg-name (3.2.15) next; urgency=medium
	//
	//   [ Author2 ]
	//   * Исправлен вызов клиента на нестандартный порт
	//
	//  -- Package Maintainer <pkg-maint@example.net>  Thu, 15 Dec 2022 17:27:01 +0000
}
