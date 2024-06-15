package fields

import "testing"

func TestArch(t *testing.T) {

	var a Architecture

	a.UnmarshalText([]byte("all"))
	t.Log(a)

	a.UnmarshalText([]byte("any"))
	t.Log(a)

	a.UnmarshalText([]byte("linux-any"))
	t.Log(a)

	a.UnmarshalText([]byte("amd64"))
	t.Log(a)

	a.UnmarshalText([]byte("musl-linux-arm64"))
	t.Log(a)

	a.UnmarshalText([]byte("arch"))
	t.Log(a)

	s, _ := a.MarshalText()
	t.Log(string(s))

}

func TestEqualsDifferent(t *testing.T) {
	var a, b Architecture
	a.UnmarshalText([]byte("gnueabi-musl-linux-any"))
	b.UnmarshalText([]byte("amd64"))

	if a.Equals(b) {
		t.Fail()
	}
}

func TestEqualsWildcard(t *testing.T) {
	var a, b Architecture
	a.UnmarshalText([]byte("linux-any"))
	b.UnmarshalText([]byte("amd64"))

	if !a.Equals(b) {
		t.Fail()
	}
}

func TestEqualsSame(t *testing.T) {
	var a, b Architecture
	a.UnmarshalText([]byte("amd64"))
	b.UnmarshalText([]byte("amd64"))

	if !a.Equals(b) {
		t.Fail()
	}
}
