package fields_test

import (
	"testing"

	"github.com/aol-nnov/debian/fields"
)

func TestUnmarshalArch(t *testing.T) {
	variants := [][]byte{
		[]byte("[arch]"),
		[]byte("   [  arch  ]  "),
		[]byte("  [ arch     arch ] "),
		[]byte("[   !arch]"),
		[]byte("[!arch]"),
		[]byte("[!arch arch]"),
		[]byte("[arch !arch]"),
	}

	for _, v := range variants {
		var ac fields.ArchitectureConstraints
		if err := ac.UnmarshalText(v); err != nil {
			t.Fatalf("failed to unmarshal '%s': %v", v, err)
		}
	}

}

func TestUnmarshalArchNegative(t *testing.T) {
	variants := [][]byte{
		[]byte("  [ arch     arch ] garbage"),
		[]byte("qwe [arch]"),
		[]byte("  a [  arch  ]  "),
		[]byte("!arch"),
		[]byte("[]"),
	}

	for _, v := range variants {
		var ac fields.ArchitectureConstraints
		if err := ac.UnmarshalText(v); err == nil {
			t.Fatalf("must fail to unmarshal '%s'", v)
		}
	}

}

// 'arch' should satisfy 'arch'
func TestAcEqualArch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	err := ac.UnmarshalText([]byte("[arch]"))

	if err != nil || !ac.SatisfiedBy(fields.MakeArch("arch")) {
		t.Fatal(err)
	}
}

// 'wildcard' should satisfy 'arch'
func TestWildcardAcEqualArch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	ac.UnmarshalText([]byte("[linux-any]"))

	if !ac.SatisfiedBy(fields.MakeArch("amd64")) {
		t.Fail()
	}
}

// ['arch', 'arch2'] should NOT satisfy 'another'
func TestAcAnotherArch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	var a fields.Architecture

	ac.UnmarshalText([]byte("[arch arch2]"))
	a.UnmarshalText([]byte("another"))

	if ac.SatisfiedBy(a) {
		t.Fail()
	}
}

// '!arch' should NOT satisfy 'arch'
func TestNegAcArch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	ac.UnmarshalText([]byte("[!arch]"))

	if ac.SatisfiedBy(fields.MakeArch("arch")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', '!amd64'] should satisfy 'armhf'
func TestNegWildcardAcArch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any !amd64]"))

	if !ac.SatisfiedBy(fields.MakeArch("armhf")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', 'amd64'] should satisfy 'armhf'
func TestWildcardAcArchNoMatch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any amd64]"))

	if !ac.SatisfiedBy(fields.MakeArch("armhf")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', 'amd64'] should NOT satisfy 'kfreebsd-i386'
func TestWildcardNegMatchAcArchNoMatch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any amd64]"))

	if ac.SatisfiedBy(fields.MakeArch("kfreebsd-i386")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', 'amd64'] should satisfy 'amd64'
func TestWildcardAcArchMatch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any amd64]"))

	if !ac.SatisfiedBy(fields.MakeArch("amd64")) {
		t.Fail()
	}
}

// '!arch' should satisfy 'another'
func TestNegAcAnotherArch(t *testing.T) {
	var ac fields.ArchitectureConstraints
	ac.UnmarshalText([]byte("[!arch]"))

	if !ac.SatisfiedBy(fields.MakeArch("another")) {
		t.Fail()
	}
}

func TestConstraintsSatisfy(t *testing.T) {
	// ac := []byte("amd64 i386 kfreebsd-amd64 mips mipsel powerpc ppc64 s390 sparc s390x mipsn32 mipsn32el mipsr6 mipsr6el mipsn32r6 mipsn32r6el mips64 mips64el mips64r6 mips64r6el x32")

	ac := []byte("[i386 amd64]")

	var architectureConstraints fields.ArchitectureConstraints
	architectureConstraints.UnmarshalText(ac)

	if !architectureConstraints.SatisfiedBy(fields.MakeArch("amd64")) {
		t.Fail()
	}
}

func TestConstraintsSatisfyReverse(t *testing.T) {
	// ac := []byte("amd64 i386 kfreebsd-amd64 mips mipsel powerpc ppc64 s390 sparc s390x mipsn32 mipsn32el mipsr6 mipsr6el mipsn32r6 mipsn32r6el mips64 mips64el mips64r6 mips64r6el x32")

	ac := []byte("[amd64 i386]")

	var architectureConstraints fields.ArchitectureConstraints
	architectureConstraints.UnmarshalText(ac)

	if !architectureConstraints.SatisfiedBy(fields.MakeArch("amd64")) {
		t.Fail()
	}
}
