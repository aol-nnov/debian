package fields

import "testing"

func TestArch(t *testing.T) {

	cases := [][]byte{
		[]byte("all"),
		[]byte("any"),
		[]byte("linux-any"),
		[]byte("amd64"),
		[]byte("musl-linux-arm64"),
		[]byte("arch"),
	}

	for _, tc := range cases {
		t.Run(string(tc), func(t *testing.T) {
			var a Architecture
			a.UnmarshalText(tc)
			if a.String() != string(tc) {
				t.Fatalf("%s != %s", a.String(), tc)
			}
		})
	}
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
