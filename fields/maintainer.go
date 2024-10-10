package fields

import (
	"encoding"
	"fmt"
)

type Maintainer struct {
	Name  string
	Email string
}

func (m *Maintainer) UnmarshalText(text []byte) (err error) {
	// strings.Between()
	return
}

func (m *Maintainer) MarshalText() (text []byte, err error) {
	return
}

func (m Maintainer) String() string {
	return fmt.Sprintf("%s <%s>", m.Name, m.Email)
}

var _ encoding.TextMarshaler = (*Maintainer)(nil)
var _ encoding.TextUnmarshaler = (*Maintainer)(nil)
