package fields_test

import (
	"testing"

	"github.com/aol-nnov/debian/fields"
)

func TestDependency_UnmarshalText(t *testing.T) {
	input := []byte(`pkgname8 (>> 2.2.wrong) [arch1 arch2 arch3] <profile1> <!profile2 profile3>`)

	var d fields.Dependency
	if err := d.UnmarshalText(input); err != nil {
		t.Fatal(err)
	}

	t.Log(d)
}

func TestUnmarshalVersionConstraint(t *testing.T) {
	input := []byte("   (   >=    1:2.2.1-3deb1   )    ")

	var vc fields.VersionConstraint

	if err := vc.UnmarshalText(input); err != nil {
		t.Fatal(err)
	}

	t.Log(vc)
}
