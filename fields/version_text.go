package fields

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func (v *Version) UnmarshalText(text []byte) (err error) {
	var mayBeEpoch, rest []byte
	var found bool

	text = bytes.TrimSpace(text)

	// save templates as is
	if bytes.HasPrefix(text, []byte{'$'}) {
		v.raw = string(text)
		return
	}

	/*
		man deb-version: If there is no debian-revision then hyphens are not allowed; if there  is no epoch then colons
		are not allowed.

		So, `:` is epoch delimiter, and `-` is debian version delimiter
	*/
	if mayBeEpoch, rest, found = bytes.Cut(text, []byte{':'}); found {
		if v.Epoch, err = strconv.Atoi(string(mayBeEpoch)); err != nil {
			return fmt.Errorf("Version: Epoch format error ('%s')", text)
		}
	} else {
		rest = mayBeEpoch
	}

	var leftPart, rightPart []byte
	leftPart, rightPart, found = bytes.Cut(rest, []byte{'-'})

	if !cisdigit(rune(leftPart[0])) {
		return fmt.Errorf("Version: UpstreamVersion must start with a number")
	}

	if found {
		// upstreamVer-debVer (quilt package)
		if bytes.Contains(rightPart, []byte("-")) {
			return fmt.Errorf(`DebianRevision format error. Dash must delimit upstream version and debian revision. Got debian revision '%s'`, rightPart)
		}
		v.UpstreamVersion = string(leftPart)
		v.DebianRevision, v.Modificators = extractVersionModificators(string(rightPart), "+~")

		if rev, nmu, found := strings.Cut(v.DebianRevision, "."); found {
			v.DebianRevision = rev
			v.Modificators = append(VersionModificators{"." + nmu}, v.Modificators...)
		}
	} else {
		v.UpstreamVersion, v.Modificators = extractVersionModificators(string(leftPart), "+~")
	}

	return
}

func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// String representation of Version struct
func (v Version) String() string {
	if v.raw != "" {
		return v.raw
	}

	res := ""
	if v.Epoch > 0 {
		res = fmt.Sprintf("%d:", v.Epoch)
	}

	res += v.UpstreamVersion

	if v.DebianRevision != "" {
		res += "-"
		res += string(v.DebianRevision)
	}

	res += strings.Join(v.Modificators, "")

	return res
}
