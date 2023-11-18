package deb

import (
	"bytes"
	"fmt"
)

type MultiArch int

const (
	MultiArchSame MultiArch = iota
	MultiArchForeign
	MultiArchAllowed
)

var string2Type = map[string]MultiArch{
	"":        MultiArchSame,
	"same":    MultiArchSame,
	"foreign": MultiArchForeign,
	"allowed": MultiArchAllowed,
}

func (ma *MultiArch) UnmarshalText(text []byte) (err error) {

	if value, found := string2Type[string(bytes.Trim(text, " "))]; found {
		*ma = value
		return nil
	}

	return fmt.Errorf("wrong multiarch type '%s'", text)
}

func (ma MultiArch) String() string {
	return [...]string{"same", "foreign", "allowed"}[ma]
}
