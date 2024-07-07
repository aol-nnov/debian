package fields

import "bytes"

// https://manpages.debian.org/unstable/dpkg-dev/dpkg-architecture.1.en.html#Debian~2

var (
	allArchShort = []byte("all")
	anyArchShort = []byte("any")
	defaults     = []string{
		"base",  // abi
		"gnu",   // libc
		"linux", // os
		// "",      //cpu
	}
)

// [pkg/encoding.TextUnmarshaler] interface implementation
func (a *Architecture) UnmarshalText(text []byte) (err error) {

	if bytes.Equal(text, allArchShort) {
		a.raw = [4]string{"all", "all", "all", "all"}
	}

	if bytes.Contains(text, anyArchShort) {
		a.raw = [4]string{"any", "any", "any", "any"}
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

	for idx, part := range a.raw {
		if part == "" {
			a.raw[idx] = defaults[idx]
		}
	}

	return nil
}

// [pkg/encoding.TextMarshaler] interface implementation
func (a Architecture) MarshalText() (text []byte, err error) {
	return []byte(a.String()), nil
}

// [pkg/fmt.Stringer] interface implementations
func (a Architecture) String() string {
	if a.raw[VerCpu] == "all" {
		return "all"
	}

	res := ""
	partIsDefault := true
	isWildcard := false

	/*
		wildcard: skip all leading `any`, print all the rest
			any-any-linux-any => <skip> - <skip> - linux (despite it's default!) - any

		normal: skip all leading `defaults`, print all the rest
			base-gnu-linux-amd64 => <skip> - <skip> - <skip> - amd64
			base-musl-linux-arm64 => <skip> - musl - linux (despite it's default!) - amd64
	*/
	for rawIdx := VerAbi; rawIdx <= VerOs; rawIdx++ {
		if a.raw[rawIdx] == "any" {
			isWildcard = true
			continue
		}

		if isWildcard ||
			!partIsDefault || a.raw[rawIdx] != defaults[rawIdx] {

			partIsDefault = false
			res += a.raw[rawIdx] + "-"
		}
	}

	return res + a.raw[VerCpu]
}
