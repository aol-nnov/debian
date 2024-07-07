package fields

import (
	"encoding"
	"fmt"
	"strings"
)

type Description string

func (descr *Description) UnmarshalText(text []byte) (err error) {
	*descr = Description(text)
	return nil
}

func (descr *Description) MarshalText() (text []byte, err error) {
	for lineNum, line := range strings.Split(string(*descr), "\n") {
		if lineNum == 0 {
			text = fmt.Append(text, line)
		} else {
			if line == "" {
				line = "."
			}

			text = fmt.Append(text, "\n ", line)
		}

	}

	return
}

var _ encoding.TextMarshaler = (*Description)(nil)
var _ encoding.TextUnmarshaler = (*Description)(nil)
