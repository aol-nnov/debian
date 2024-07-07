package fields

import (
	"bytes"
	"fmt"
	"strconv"
)

func (v *Version) UnmarshalText(text []byte) (err error) {
	var epoch, rest []byte
	var found bool

	text = bytes.TrimSpace(text)

	// save templates as is
	if bytes.HasPrefix(text, []byte{'$'}) {
		v.raw = string(text)
		return
	}

	if epoch, rest, found = bytes.Cut(text, []byte{':'}); found {
		if v.Epoch, err = strconv.Atoi(string(epoch)); err != nil {
			return fmt.Errorf("Version: Epoch format error ('%s')", text)
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

func (v Version) MarshalText() ([]byte, error) {
	if v.raw != "" {
		return []byte(v.raw), nil
	}

	return []byte(v.String()), nil
}

// String representation of Version struct
func (v Version) String() string {
	res := v.raw

	if res != "" {
		return res
	}

	if v.Epoch > 0 {
		res = fmt.Sprintf("%d:", v.Epoch)
	}
	res += v.UpstreamVersion

	if v.DebianVersion != "" {
		res += "-" + v.DebianVersion
	}

	return res
}
