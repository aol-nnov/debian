package fields

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ChangeImpact int

const (
	ChangeImpactMajor ChangeImpact = iota
	ChangeImpactMinor
	ChangeImpactTrivial
)

// rightmost number in string like `12local45`
func rightmostNumberInStr(in string) (prefix string, val int, found bool) {
	var err error

	re := regexp.MustCompile(`.*[^\d](\d+)`)
	matches := re.FindStringSubmatchIndex(in)

	if matches == nil {
		if val, err := strconv.Atoi(in); err != nil {
			return "", 0, false
		} else {
			return "", val, true
		}
	}

	startPos := matches[2]
	endPos := matches[3]

	if val, err = strconv.Atoi(in[startPos:endPos]); err != nil {
		return "", 0, false
	}

	return in[:startPos], val, true

}

func (v *Version) bumpSemVer(impact ChangeImpact) {
	// segments in terms `major.minor.patch`
	// segments[0] - Major
	// segments[1] - Minor
	// segments[2] - Patch
	// and
	// ChangeImpactMajor == 0, ChangeImpactMinor == 1...
	// but version might be shorter (fewer segments) or longer (more segments)
	segments := strings.FieldsFunc(v.UpstreamVersion, func(r rune) bool { return r == '.' })

	segmentToChange := int(impact)

	// if it shorter, let's bump next significant segment
	if len(segments)-1 < segmentToChange {
		segmentToChange = len(segments) - 1
	}

	if prefix, val, found := rightmostNumberInStr(segments[segmentToChange]); found {

		// bump segment value
		segments[segmentToChange] = fmt.Sprintf("%s%d", prefix, val+1)

		// ... and zero-out less significant segments
		for p := segmentToChange + 1; p < len(segments); p++ {
			segments[p] = "0"
		}
	}

	v.UpstreamVersion = strings.Join(segments, ".")
}

// smartly bumps Debian version
// removes nmu, snapshots, etc
func (v *Version) Bump(impact ChangeImpact) {
	for v.IsMod() != VersionModNone {
		v.RemoveMod()
	}

	if v.IsQuilt() {
		if prefix, val, found := rightmostNumberInStr(v.DebianRevision); found {
			v.DebianRevision = fmt.Sprintf("%s%d", prefix, val+1)
		}
	} else {
		v.bumpSemVer(impact)
	}

}
