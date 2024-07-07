package fields

import (
	"bytes"
	"fmt"
)

// [pkg/encoding.TextUnmarshaler] interface implementation
func (constraints *ArchitectureConstraints) UnmarshalText(text []byte) (err error) {
	trimmedInput := bytes.TrimSpace(text)
	if !bytes.HasPrefix(trimmedInput, []byte{'['}) ||
		!bytes.HasSuffix(trimmedInput, []byte{']'}) {
		return fmt.Errorf("ArchitectureConstraints unmarshal: wrong input string '%s'", text)
	}

	arches := bytes.Split(bytes.Trim(trimmedInput, "[] "), []byte(" "))

	for _, arch := range arches {
		if len(arch) > 0 {
			var ac architectureConstraint
			if err = ac.UnmarshalText(bytes.TrimSpace(arch)); err != nil {
				return
			}

			*constraints = append(*constraints, ac)
		}
	}

	if len(*constraints) == 0 {
		return fmt.Errorf("empty ArchitectureConstraints")
	}

	return
}

func (constraints *ArchitectureConstraints) MarshalText() (text []byte, err error) {
	panic("unimplemented")
}

func (constraints ArchitectureConstraints) String() (res string) {
	fmt.Println("ArchitectureConstraints String")

	if len(constraints) == 0 {
		return ""
	}

	res = "["
	for idx, c := range constraints {
		if idx == 0 {
			res += c.String()
		} else {
			res = fmt.Sprint(res, " %s", c)
		}
	}
	res += "]"

	return

}

// [pkg/encoding.TextUnmarshaler] interface implementation
func (ac *architectureConstraint) UnmarshalText(text []byte) (err error) {
	if text[0] == '!' {
		ac.Negate = true
		return ac.Name.UnmarshalText(text[1:])
	} else {
		ac.Negate = false
		return ac.Name.UnmarshalText(text)
	}
}

// [pkg/fmt.Stringer] interface implementations
func (ac architectureConstraint) String() string {
	not := ""

	if ac.Negate {
		not = "!"
	}

	return fmt.Sprintf("%s%v", not, ac.Name)
}
