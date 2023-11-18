package deb

import (
	"testing"
)

func TestDepencencyNoConstraint(t *testing.T) {

	in := []byte("pkgname")

	var d Dependency

	if err := d.UnmarshalText(in); err != nil {
		t.Fatalf("%v", err)
	}

	t.Log(d)
}

func TestDepencencyVersion(t *testing.T) {
	in := []byte("pkgname (>= 1.2.3)")

	var d Dependency

	if err := d.UnmarshalText(in); err != nil {
		t.Fatalf("%v", err)
	}

	t.Log(d)
}

func TestDepencencyVerArch(t *testing.T) {
	in := []byte("pkgname (>= 1.2.3) [arch] ")

	var d Dependency

	if err := d.UnmarshalText(in); err != nil {
		t.Fatalf("%v", err)
	}

	t.Log(d)
}

func TestDepencencyArchProfiles(t *testing.T) {
	in := []byte("pkgname [arch] <!profile>")

	var d Dependency

	if err := d.UnmarshalText(in); err != nil {
		t.Fatalf("%v", err)
	}

	t.Log(d)
}

func TestDepencencyVerProfiles(t *testing.T) {
	in := []byte("pkgname (>= 1.2.3) <!profile>")

	var d Dependency

	if err := d.UnmarshalText(in); err != nil {
		t.Fatalf("%v", err)
	}

	t.Log(d)

}
func TestDepencencyVerArchProfiles(t *testing.T) {
	in := []byte("pkgname:native (>= 1.2.3) [arch1 arch2 !arch3] <!profile1 profile2> <profile3> | another")

	var d Dependency

	if err := d.UnmarshalText(in); err != nil {
		t.Fatalf("%v", err)
	}

	t.Log(d)

	if d.Name != "pkgname" ||
		d.ArchQualifier != "native" ||
		d.VersionConstraint.Op != VersionConstraintGreaterOrEqual ||
		d.VersionConstraint.Value.Epoch != 0 ||
		d.VersionConstraint.Value.UpstreamVersion != "1.2.3" ||
		d.VersionConstraint.Value.DebianVersion != "" ||
		len(d.ArchitectureConstraints) != 3 ||
		d.ArchitectureConstraints[1].Name.Cpu() != "arch2" ||
		d.ArchitectureConstraints[1].Negate != false ||
		len(d.ProfileConstraints) != 2 ||
		d.ProfileConstraints[1][0].Name != "profile3" ||
		d.ProfileConstraints[1][0].Negate != false ||
		d.Alt == nil {
		t.Fail()
	}
}
