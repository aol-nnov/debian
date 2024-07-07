package fields_test

import (
	"fmt"
	"testing"

	"github.com/aol-nnov/debian/fields"
)

func TestEpochIsNumber(t *testing.T) {
	var v fields.Version
	if err := v.UnmarshalText([]byte("1a:1.2.3")); err == nil {
		t.Fail()
	}
}

func TestUpstreamVersionStartsWithNumber(t *testing.T) {
	var v fields.Version
	if err := v.UnmarshalText([]byte("a.1.2.3")); err == nil {
		t.Fail()
	}

	if err := v.UnmarshalText([]byte("a1.2.3")); err == nil {
		t.Fail()
	}
}

func TestUpstreamVersionStartsWithNumberButMayContainLetters(t *testing.T) {
	var v fields.Version
	if err := v.UnmarshalText([]byte("1a.2.3")); err != nil {
		t.Error(err)
	}
}

func TestVersion_UnmarshalText(t *testing.T) {
	var v fields.Version
	if err := v.UnmarshalText([]byte("1:1.2.3-0.0.1~alpha2")); err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestCompare(t *testing.T) {
	v1 := fields.MakeVersion("1:1.2.3-0.0.1~beta1")
	v2 := fields.MakeVersion("1:1.2.3-0.0.1~alpha2")

	if v1.Compare(v2) != fields.VersionCompareResultGreaterThan {
		t.Logf("%v %v %v", v1, v1.Compare(v2), v2)
		t.Fail()
	}
}

func TestMakeSnapshot(t *testing.T) {
	v := fields.MakeVersion("1.2.3")
	v.Snapshot("deadbeef")
	if v.String() != "1.2.3~1.gbpdeadbeef" {
		t.Fail()
	}

	v.Snapshot("deafc0de")
	if v.String() != "1.2.3~2.gbpdeafc0de" {
		t.Fail()
	}
}

func TestCompareGbpSnapshot(t *testing.T) {
	// as per git-buildpackage sources
	// Format is <some-version>~<buildNum>.gbp<short-commit-id>
	v1 := fields.MakeVersion("1:1.2.3-0.0.1~1.gbpdeadbeef")
	v2 := fields.MakeVersion("1:1.2.3-0.0.1~2.gbp12343455")

	t.Log(v2.DebianRevision)

	if v2.Compare(v1) != fields.VersionCompareResultGreaterThan {
		t.Logf("%v %v %v", v1, v1.Compare(v2), v2)
		t.Fail()
	}
}

// тест появился в результате отладки [verrevcmp], которая раньше возвращала разницу (в числах) между версиями, а не
// fields.VersionCompareResult, т.е. [-1, 0, 1], в результате вызов String() для него приводил к выходу за границу
// массива, а сам алгоритм сравнения работал не верно, например, будучи примененным в [pkg/slices.Sort]
func TestCompareBigDifference(t *testing.T) {
	v1 := fields.MakeVersion("1.0.38")
	v2 := fields.MakeVersion("1.0.36")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("VersionCompareResult stringer failed")
		}
	}()
	cr := v1.Compare(v2).String()
	t.Logf("%v %s %v", v1, cr, v2)
}

func TestNotLess(t *testing.T) {
	v1 := fields.MakeVersion("3.2")
	v2 := fields.MakeVersion("2.2")

	if v1.Less(v2) {
		t.Fail()
	}

	t.Logf("%v %v %v", v1, v1.Compare(v2), v2)
}

func TestLess(t *testing.T) {
	v1 := fields.MakeVersion("2.2")
	v2 := fields.MakeVersion("2.35")

	if !v1.Less(v2) {
		t.Fail()
	}

	t.Logf("%v %v %v", v1, v1.Compare(v2), v2)
}

func TestRollBack(t *testing.T) {
	v1 := fields.MakeVersion("2.3.9-5")
	v2 := v1

	v2.RollBackTo("2.2.0")

	fmt.Printf("%v %v %v\n", v1, v1.Compare(v2), v2)
}

func TestRebuild(t *testing.T) {
	v1 := fields.MakeVersion("2.9-3")
	v2 := fields.MakeVersion("2.9-3+b1")

	fmt.Printf("%v %v %v\n", v1, v1.Compare(v2), v2)
}

func TestSnapshotCompare(t *testing.T) {
	v1 := fields.MakeVersion("2.9-3")
	v2 := fields.MakeVersion("2.9-3~")

	fmt.Printf("%v %v %v\n", v1, v1.Compare(v2), v2)
}

func TestDebianSnapshot(t *testing.T) {
	v1 := fields.MakeVersion("2.9-3")

	v1.Snapshot("deadbeef")
	fmt.Println(v1)
	if v1.String() != "2.9-3~1.gbpdeadbeef" {
		t.Fail()
	}

	v1.Snapshot("badf00d")
	fmt.Println(v1)
	if v1.String() != "2.9-3~2.gbpbadf00d" {
		t.Fail()
	}

}

func TestNativeSnapshot(t *testing.T) {
	v1 := fields.MakeVersion("2.9")

	v1.Snapshot("deadbeef")
	fmt.Println(v1)
	if v1.String() != "2.9~1.gbpdeadbeef" {
		t.Fail()
	}

	v1.Snapshot("badf00d")
	fmt.Println(v1)
	if v1.String() != "2.9~2.gbpbadf00d" {
		t.Fail()
	}

}

func TestNativeVsDebian(t *testing.T) {
	v1 := fields.MakeVersion("2.9")
	v2 := fields.MakeVersion("2.9-1")

	fmt.Printf("%v %v %v\n", v1, v1.Compare(v2), v2)
}
