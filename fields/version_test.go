package fields

import "testing"

func TestEpochIsNumber(t *testing.T) {
	var v Version
	if err := v.UnmarshalText([]byte("1a:1.2.3")); err != nil {
		t.Fail()
	}
}

func TestUpstreamVersionStartsWithNumber(t *testing.T) {
	var v Version
	if err := v.UnmarshalText([]byte("a.1.2.3")); err == nil {
		t.Fail()
	}

	if err := v.UnmarshalText([]byte("a1.2.3")); err == nil {
		t.Fail()
	}
}

func TestUpstreamVersionStartsWithNumberBytMayContainLetters(t *testing.T) {
	var v Version
	if err := v.UnmarshalText([]byte("1a.2.3")); err != nil {
		t.Error(err)
	}
}

func TestVersion_UnmarshalText(t *testing.T) {
	var v Version
	if err := v.UnmarshalText([]byte("1:1.2.3-0.0.1~alpha2")); err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestCompare(t *testing.T) {
	v1 := MakeVersion("1:1.2.3-0.0.1~beta1")
	v2 := MakeVersion("1:1.2.3-0.0.1~alpha2")

	if v1.Compare(v2) != VersionCompareResultGreaterThan {
		t.Logf("%v %v %v", v1, v1.Compare(v2), v2)
		t.Fail()
	}
}

// this test was created as a result of debugging [verrevcmp], which used to return the difference (in numbers) between
// versions, rather than VersionCompareResult, i.e. [-1, 0, 1], as a result calling String() on it resulted in getting out of array bounds, which lead to the comparison algorithm malfunction, for example, when applied in [pkg/slices.Sort]
func TestCompareBigDifference(t *testing.T) {
	v1 := MakeVersion("1.0.38")
	v2 := MakeVersion("1.0.36")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("VersionCompareResult stringer failed")
		}
	}()
	cr := v1.Compare(v2).String()
	t.Logf("%v %s %v", v1, cr, v2)
}

func TestNotLess(t *testing.T) {
	v1 := MakeVersion("3.2")
	v2 := MakeVersion("2.2")

	if v1.Less(v2) {
		t.Fail()
	}

	t.Logf("%v %v %v", v1, v1.Compare(v2), v2)
}

func TestLess(t *testing.T) {
	v1 := MakeVersion("2.2")
	v2 := MakeVersion("2.35")

	if !v1.Less(v2) {
		t.Fail()
	}

	t.Logf("%v %v %v", v1, v1.Compare(v2), v2)
}
