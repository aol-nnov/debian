package deb

import (
	"testing"
)

func TestUnmarshalArch(t *testing.T) {
	var ac ArchitectureConstraints

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
		if err := ac.UnmarshalText(v); err != nil {
			t.Fatalf("failed to unmarshal '%s'", v)
		}
	}

}

// 'arch' should satisfy 'arch'
func TestAcEqualArch(t *testing.T) {
	var ac ArchitectureConstraints
	err := ac.UnmarshalText([]byte("[arch]"))

	if err != nil || !ac.Satisfies(MakeArch("arch")) {
		t.Fatal(err)
	}
}

// 'wildcard' should satisfy 'arch'
func TestWildcardAcEqualArch(t *testing.T) {
	var ac ArchitectureConstraints
	ac.UnmarshalText([]byte("[linux-any]"))

	if !ac.Satisfies(MakeArch("amd64")) {
		t.Fail()
	}
}

// ['arch', 'arch2'] should NOT satisfy 'another'
func TestAcAnotherArch(t *testing.T) {
	var ac ArchitectureConstraints
	var a Architecture

	ac.UnmarshalText([]byte("[arch arch2]"))
	a.UnmarshalText([]byte("another"))

	if ac.Satisfies(a) {
		t.Fail()
	}
}

// '!arch' should NOT satisfy 'arch'
func TestNegAcArch(t *testing.T) {
	var ac ArchitectureConstraints
	ac.UnmarshalText([]byte("[!arch]"))

	if ac.Satisfies(MakeArch("arch")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', '!amd64'] should satisfy 'armhf'
func TestNegWildcardAcArch(t *testing.T) {
	var ac ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any !amd64]"))

	if !ac.Satisfies(MakeArch("armhf")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', 'amd64'] should satisfy 'armhf'
func TestWildcardAcArchNoMatch(t *testing.T) {
	var ac ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any amd64]"))

	if !ac.Satisfies(MakeArch("armhf")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', 'amd64'] should NOT satisfy 'kfreebsd-i386'
func TestWildcardNegMatchAcArchNoMatch(t *testing.T) {
	var ac ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any amd64]"))

	if ac.Satisfies(MakeArch("kfreebsd-i386")) {
		t.Fail()
	}
}

// ['!kfreebsd-any', 'amd64'] should satisfy 'amd64'
func TestWildcardAcArchMatch(t *testing.T) {
	var ac ArchitectureConstraints
	ac.UnmarshalText([]byte("[!kfreebsd-any amd64]"))

	if !ac.Satisfies(MakeArch("amd64")) {
		t.Fail()
	}
}

// '!arch' should satisfy 'another'
func TestNegAcAnotherArch(t *testing.T) {
	var ac ArchitectureConstraints
	ac.UnmarshalText([]byte("[!arch]"))

	if !ac.Satisfies(MakeArch("another")) {
		t.Fail()
	}
}

func TestConstraintsSatisfy(t *testing.T) {
	// ac := []byte("amd64 i386 kfreebsd-amd64 mips mipsel powerpc ppc64 s390 sparc s390x mipsn32 mipsn32el mipsr6 mipsr6el mipsn32r6 mipsn32r6el mips64 mips64el mips64r6 mips64r6el x32")

	ac := []byte("[i386 amd64]")

	var architectureConstraints ArchitectureConstraints
	architectureConstraints.UnmarshalText(ac)

	if !architectureConstraints.Satisfies(MakeArch("amd64")) {
		t.Fail()
	}
}

func TestConstraintsSatisfyReverse(t *testing.T) {
	// ac := []byte("amd64 i386 kfreebsd-amd64 mips mipsel powerpc ppc64 s390 sparc s390x mipsn32 mipsn32el mipsr6 mipsr6el mipsn32r6 mipsn32r6el mips64 mips64el mips64r6 mips64r6el x32")

	ac := []byte("[amd64 i386]")

	var architectureConstraints ArchitectureConstraints
	architectureConstraints.UnmarshalText(ac)

	if !architectureConstraints.Satisfies(MakeArch("amd64")) {
		t.Fail()
	}
}
