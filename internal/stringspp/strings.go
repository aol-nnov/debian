// Strings++: functions missing from `strings` package
package stringspp

import (
	"bytes"
	"encoding/json"
)

// dummy stringer for any marshallable type
func UniversalStringer(s interface{}) string {
	str, _ := json.MarshalIndent(s, "", " ")
	return string(str)
}

/*
Extracts byte slice surrounded by delimiters
returns:
  - found slice (including delimiters)
  - found flag
  - tail slice - the reminder after the rightmost delimiter

returns incoming slice if delimiters are not found
*/
func Between(in []byte, left, right byte, greedy bool) (foundSlice []byte, found bool, tail []byte) {
	var leftIdx, rightIdx int

	if leftIdx = bytes.IndexByte(in, left); leftIdx == -1 {
		return in, false, nil
	}

	if greedy {
		if rightIdx = bytes.LastIndexByte(in, right); rightIdx == -1 {
			return in, false, nil
		}
	} else {
		if rightIdx = bytes.IndexByte(in, right); rightIdx == -1 {
			return in, false, nil
		}
	}

	return in[leftIdx : rightIdx+1], true, in[rightIdx+1:]
}
