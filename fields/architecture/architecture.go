package deb

import (
	"bytes"
)

/*
Represents Debian architecture. May be in one of the following forms:
  - `any` (implicitly: any-any-any)
  - kfreebsd-any (implicitly: any-kfreebsd-any)
  - kfreebsd-amd64 (implicitly any-kfreebsd-amd64)
  - bsd-openbsd-i386
*/
/*
abi-libc-os-cpu
*/

const (
	VerAbi int = iota
	VerLibc
	VerOs
	VerCpu
)

func MakeArch(name string) Architecture {
	var a Architecture
	a.UnmarshalText([]byte(name))
	return a
}

type Architecture struct {
	// abi-libc-os-cpu Debian tuple
	raw [4]string
}

func (a Architecture) Abi() string {
	return a.raw[VerAbi]
}

func (a Architecture) Libc() string {
	return a.raw[VerLibc]
}

func (a Architecture) Os() string {
	return a.raw[VerOs]
}

func (a Architecture) Cpu() string {
	return a.raw[VerCpu]
}

var (
	allArch = []byte("all")
)

func (a *Architecture) UnmarshalText(text []byte) (err error) {
	defaultVal := "any"
	if bytes.Equal(text, allArch) {
		defaultVal = "all"
	}
	for rawIdx := VerAbi; rawIdx <= VerCpu; rawIdx++ {
		a.raw[rawIdx] = defaultVal
	}

	specs := bytes.SplitN(text, []byte{'-'}, 4)

	for specIdx, rawIdx := len(specs)-1, VerCpu; specIdx >= 0; specIdx, rawIdx = specIdx-1, rawIdx-1 {
		a.raw[rawIdx] = string(specs[specIdx])
	}

	if len(specs) == 1 { // short form, i.e. amd64
		a.raw[VerAbi] = "gnu"
		a.raw[VerOs] = "linux"
	}

	return nil
}

func (a Architecture) MarshalText() (text []byte, err error) {
	return []byte(a.String()), nil
}

func (a Architecture) String() string {
	res := ""

	for rawIdx := VerAbi; rawIdx <= VerOs; rawIdx++ {
		if a.raw[rawIdx] != "any" {
			res += a.raw[rawIdx] + "-"
		}
	}

	return res + a.raw[VerCpu]
}

/*
 */
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
