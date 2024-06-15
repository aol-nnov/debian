package fields

import (
	"bytes"
)

const (
	VerAbi int = iota
	VerLibc
	VerOs
	VerCpu
)

// Creates new [Architecture] instance from [pkg/string]
func MakeArch(name string) Architecture {
	var a Architecture
	a.UnmarshalText([]byte(name))
	return a
}

/*
Represents a Debian architecture tuple in the fully qualified architecture with all its components spelled out. The current tuple has the form abi-libc-os-cpu.

May be in one of the following forms:

  - `any` (implicitly: any-any-any-any)

  - kfreebsd-any (implicitly: any-any-kfreebsd-any)

  - kfreebsd-amd64 (implicitly any-any-kfreebsd-amd64)

  - bsd-openbsd-i386

More Examples:

  - base-gnu-linux-amd64

  - eabihf-musl-linux-arm.
*/
type Architecture struct {
	// abi-libc-os-cpu Debian tuple
	raw [4]string
}

// Abi part getter
func (a Architecture) Abi() string {
	return a.raw[VerAbi]
}

// Libc part getter
func (a Architecture) Libc() string {
	return a.raw[VerLibc]
}

// Os part getter
func (a Architecture) Os() string {
	return a.raw[VerOs]
}

// Cpu part getter
func (a Architecture) Cpu() string {
	return a.raw[VerCpu]
}

var (
	allArch = []byte("all")
)

// [pkg/encoding.TextUnmarshaler] interface implementation
func (a *Architecture) UnmarshalText(text []byte) (err error) {
	// initialize internal structure with corresponding wildcard first
	// `all`` Architecture is arch-indep wildcard (full form all-all-all-all)
	// `any` is binary arch wildcard
	defaultVal := "any"
	if bytes.Equal(text, allArch) {
		defaultVal = "all"
	}

	for rawIdx := VerAbi; rawIdx <= VerCpu; rawIdx++ {
		a.raw[rawIdx] = defaultVal
	}

	// tokenize string to 4 tokens at most (may be less, if short form is provided)
	specs := bytes.SplitN(text, []byte{'-'}, 4)

	specIdx := len(specs) - 1
	rawIdx := VerCpu

	// fill the internal structure starting from VerCpu, overrriding default values with actual ones
	for specIdx >= 0 {
		a.raw[rawIdx] = string(specs[specIdx])

		specIdx--
		rawIdx--
	}

	// set debian implications for the shortest form, i.e. amd64
	if len(specs) == 1 {
		a.raw[VerAbi] = "gnu"
		a.raw[VerOs] = "linux"
	}

	return nil
}

// [pkg/encoding.TextMarshaler] interface implementation
func (a Architecture) MarshalText() (text []byte, err error) {
	return []byte(a.String()), nil
}

// [pkg/fmt.Stringer] interface implementations
func (a Architecture) String() string {
	res := ""

	for rawIdx := VerAbi; rawIdx <= VerOs; rawIdx++ {
		if a.raw[rawIdx] != "any" {
			res += a.raw[rawIdx] + "-"
		}
	}

	return res + a.raw[VerCpu]
}

// Compares two `Architecture`s. (Wildcard comparison included)
func (a Architecture) Equals(another Architecture) bool {
	matches := 0
	for rawIdx := VerAbi; rawIdx <= VerCpu; rawIdx++ {
		if a.raw[rawIdx] == another.raw[rawIdx] ||
			(a.raw[rawIdx] == "any" || another.raw[rawIdx] == "any") {
			matches++
		}
	}

	return matches == 4 // all 4 parts satisfy
}
