package fields_test

import (
	"fmt"

	"github.com/aol-nnov/debian/fields"
)

func ExampleQuiltBump() {
	v := fields.MakeVersion("1.2.3-1.1~1.gbpasd")
	fmt.Println(v.Modificators)
	fmt.Println(v.DebianRevision)

	v.Bump(fields.ChangeImpactTrivial)
	v.Snapshot("booo")

	fmt.Println(v)

	// Output:
	// [.1 ~1.gbpasd]
	// 1
	// 1.2.3-2~1.gbpbooo
}

func ExampleNativeBump() {
	v := fields.MakeVersion("1.2.3+b5~1.gbpasd")

	v.Bump(fields.ChangeImpactTrivial)
	fmt.Println(v)

	// Output: 1.2.4
}

func ExampleSnapshotBinNmu() {

	v := fields.MakeVersion("3:1.2.3")
	fmt.Println(v)

	v.Snapshot("lala")
	fmt.Println(v)

	v.BinaryNmu()
	fmt.Println(v)

	v.BinaryNmu()
	fmt.Println(v)

	v.Snapshot("moo")
	fmt.Println(v)

	v.Bump(fields.ChangeImpactTrivial)
	fmt.Println(v)

	v.Bump(fields.ChangeImpactMajor)
	fmt.Println(v)

	// Output:
	// 3:1.2.3
	// 3:1.2.3~1.gbplala
	// 3:1.2.3~1.gbplala+b1
	// 3:1.2.3~1.gbplala+b2
	// 3:1.2.3~2.gbpmoo
	// 3:1.2.4
	// 3:2.0.0
}
