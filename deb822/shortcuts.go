package deb822

import (
	"bytes"
	"strings"
)

func Marshal(v any) (string, error) {
	res := new(bytes.Buffer)

	return res.String(), NewEncoder(res).Encode(v)
}

func Unmarshal(data string, v any) error {
	return NewDecoder(strings.NewReader(data)).Decode(v)
}
