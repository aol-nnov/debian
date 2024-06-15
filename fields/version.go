package fields

import (
	"bytes"
	"fmt"
	"strconv"
)

/*
Represents Debian package version number format

  - [Man page](https://linux.die.net/man/5/deb-version)
  - [Reference implementation](https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/lib/dpkg/version.c)
  - Other good stuff: https://github.com/sdumetz/node-deb-version-compare/blob/master/lib/Version.js

BNF descriptor:

	<Version> ::= (<Epoch> ":")? <UpstreamVersion> ("-" <DebianVersion>)?
	<Epoch> ::= <Num>+
	<UpstreamVersion> ::= <VersionPart> | <UpstreamVersion>  <Separator> (<Num> | <Alpha>)+
	<DebianVersion> ::= <VersionPart> | <UpstreamVersion>  <Separator> (<Num> | <Alpha>)+
	<VersionPart> ::= <Num> (<Num> | <Alpha>)*
	<Num> ::= [0-9]
	<Alpha> ::= ([A-Z] | [a-z])
	<Separator> ::= "." | "+" | "~"
*/
type Version struct {

	/*
	   This is a single (generally small) unsigned integer. It may be omitted in which case zero is assumed. If it is
	   omitted then the upstream_version may not contain any colons. It is provided to allow mistakes in the version
	   numbers of older versions of a package, and also a package's previous version numbering schemes, to be left
	   behind.
	*/
	Epoch int

	/*
		This is the main part of the version number. It is usually the version number of the original ("upstream")
		package from which the .deb file has been made, if this is applicable. Usually this will be in the same
		format as that specified by the upstream author(s); however, it may need to be reformatted to fit into the
		package management system's format and comparison scheme. The comparison behavior of the package management
		system with respect to the upstream_version is described below. The upstream_version portion of the version
		number is mandatory.

		The upstream_version may contain only alphanumerics ("A-Za-z0-9") and the characters . + - : ~ (full stop,
		plus, hyphen, colon, tilde) and should start with a digit. If there is no debian_revision then hyphens are
		not allowed; if there is no epoch then colons are not allowed.
	*/
	UpstreamVersion string

	/*
		This part of the version number specifies the version of the Debian package based on the upstream version. It
		may contain only alphanumerics and the characters + . ~ (plus, full stop, tilde) and is compared in the same way
		as the upstream_version is. It is optional; if it isn't present then the upstream_version may not contain a
		hyphen. This format represents the case where a piece of software was written specifically to be turned into a
		Debian package, and so there is only one "debianisation" of it and therefore no revision indication is required.

		It is conventional to restart the debian_revision at '1' each time time the upstream_version is increased.

		Dpkg will break the version number apart at the last hyphen in the string (if there is one) to determine the
		upstream_version and debian_revision. The absence of a debian_revision compares earlier than the presence of one
		(but note that the debian_revision is the least significant part of the version number).
	*/
	DebianVersion string
}

func MakeVersion(v string) Version {
	var res Version
	res.UnmarshalText([]byte(v))

	return res
}

// String representation of Version struct
func (v Version) String() string {
	res := ""
	if v.Epoch > 0 {
		res = fmt.Sprintf("%d:", v.Epoch)
	}
	res += v.UpstreamVersion

	if v.DebianVersion != "" {
		res += "-" + v.DebianVersion
	}

	return res
}

func (v *Version) UnmarshalText(text []byte) (err error) {
	var epoch, rest []byte
	var found bool

	text = bytes.TrimSpace(text)

	if epoch, rest, found = bytes.Cut(text, []byte{':'}); found {
		if v.Epoch, err = strconv.Atoi(string(epoch)); err != nil {
			return fmt.Errorf("Version: Epoch format error")
		}
	} else {
		rest = epoch
	}

	upstreamVer, debVer, _ := bytes.Cut(rest, []byte{'-'})

	if !cisdigit(rune(upstreamVer[0])) {
		return fmt.Errorf("Version: UpstreamVersion must start with a number")
	}
	v.UpstreamVersion = string(upstreamVer)
	v.DebianVersion = string(debVer)
	return
}

type VersionCompareResult int

const (
	VersionCompareResultLessThan VersionCompareResult = iota - 1
	VersionCompareResultEquals
	VersionCompareResultGreaterThan
)

func (c VersionCompareResult) String() string {
	return [...]string{"<<", "=", ">>"}[c+1]
}

/*
Compares two debian `Version`s

The `upstream_version` and `debian_revision` parts are compared by the package management system using the same algorithm:

1. The strings are compared from left to right.

2. First the initial part of each string consisting entirely of non-digit characters is determined. These two parts (one of
which may be empty) are compared lexically. If a difference is found it is returned. The lexical comparison is a
comparison of ASCII values modified so that all the letters sort earlier than all the non-letters and so that a tilde
sorts before anything, even the end of a part. For example, the following parts are in sorted order: '~~', '~~a', '~',
the empty part, 'a'.

3. Then the initial part of the remainder of each string which consists entirely of digit characters is determined. The
numerical values of these two parts are compared, and any difference found is returned as the result of the comparison.
For these purposes an empty string (which can only occur at the end of one or both version strings being compared)
counts as zero.

4. These two steps (comparing and removing initial non-digit strings and initial digit strings) are repeated until a
difference is found or both strings are exhausted.

*Note* that the purpose of epochs is to allow us to leave behind mistakes in version numbering, and to cope with
situations where the version numbering scheme changes. It is not intended to cope with version numbers containing
strings of letters which the package management system cannot interpret (such as 'ALPHA' or 'pre-'), or with silly
orderings.
*/
func (v Version) Compare(another Version) VersionCompareResult {
	if v.Epoch > another.Epoch {
		return VersionCompareResultGreaterThan
	}
	if v.Epoch < another.Epoch {
		return VersionCompareResultLessThan
	}

	rc := verrevcmp(v.UpstreamVersion, another.UpstreamVersion)
	if rc != VersionCompareResultEquals {
		return rc
	}

	return verrevcmp(v.DebianVersion, another.DebianVersion)
}

func (v Version) Less(another Version) bool {
	return v.Compare(another) == VersionCompareResultLessThan
}

func cisdigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func cisalpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func order(r rune) int {
	if cisdigit(r) {
		return 0
	}
	if cisalpha(r) {
		return int(r)
	}
	if r == '~' {
		return -1
	}
	if int(r) != 0 {
		return int(r) + 256
	}
	return 0
}

func intToVersionCompareResult(in int) VersionCompareResult {
	if in < 0 {
		return VersionCompareResultLessThan
	} else {
		return VersionCompareResultGreaterThan
	}
}

func verrevcmp(a string, b string) VersionCompareResult {
	i := 0
	j := 0
	for i < len(a) || j < len(b) {
		var first_diff int
		for (i < len(a) && !cisdigit(rune(a[i]))) ||
			(j < len(b) && !cisdigit(rune(b[j]))) {
			ac := 0
			if i < len(a) {
				ac = order(rune(a[i]))
			}
			bc := 0
			if j < len(b) {
				bc = order(rune(b[j]))
			}
			if ac != bc {
				return intToVersionCompareResult(ac - bc)
			}
			i++
			j++
		}

		for i < len(a) && a[i] == '0' {
			i++
		}
		for j < len(b) && b[j] == '0' {
			j++
		}

		for i < len(a) && cisdigit(rune(a[i])) && j < len(b) && cisdigit(rune(b[j])) {
			if first_diff == 0 {
				first_diff = int(rune(a[i]) - rune(b[j]))
			}
			i++
			j++
		}

		if i < len(a) && cisdigit(rune(a[i])) {
			return VersionCompareResultGreaterThan
		}
		if j < len(b) && cisdigit(rune(b[j])) {
			return VersionCompareResultLessThan
		}
		if first_diff != 0 {
			return intToVersionCompareResult(first_diff)
		}
	}
	return VersionCompareResultEquals
}
