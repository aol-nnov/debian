package fields

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
