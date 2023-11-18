package deb

import (
	"bytes"
	"fmt"
)

type ArchitectureConstraints []ArchitectureConstraint

type ArchitectureConstraint struct {
	Name   Architecture
	Negate bool
}

func (constraints *ArchitectureConstraints) UnmarshalText(text []byte) (err error) {
	if bytes.TrimSpace(text)[0] != '[' {
		return fmt.Errorf("ArchitectureConstraints unmarshal: wrong input string '%s'", text)
	}

	arches := bytes.Split(bytes.Trim(text, "[] "), space)

	for _, arch := range arches {
		if len(arch) > 0 {
			var ac ArchitectureConstraint
			if err = ac.UnmarshalText(bytes.TrimSpace(arch)); err != nil {
				return
			}

			*constraints = append(*constraints, ac)
		}
	}

	return
}

func (ac *ArchitectureConstraint) UnmarshalText(text []byte) (err error) {
	if text[0] == '!' {
		ac.Negate = true
		return ac.Name.UnmarshalText(text[1:])
	} else {
		ac.Negate = false
		return ac.Name.UnmarshalText(text)
	}
}

func (ac ArchitectureConstraint) String() string {
	not := ""

	if ac.Negate {
		not = "!"
	}

	return fmt.Sprintf("%s%v", not, ac.Name)
}

func (ac ArchitectureConstraint) satisfies(a Architecture) bool {
	res := (!ac.Negate && ac.Name.Equals(a)) || (ac.Negate && !ac.Name.Equals(a))
	// fmt.Printf("%v %v %v\n", ac, res, a)
	return res
}

// Architecture constraints are OR-ed
func (constraints ArchitectureConstraints) Satisfies(a Architecture) bool {

	// empty set satisfied any Architecture
	if len(constraints) == 0 {
		return true
	}

	satisfied := false

	for _, architectureConstraint := range constraints {
		satisfied = satisfied || architectureConstraint.satisfies(a)

		// quick path (first match bail out)
		// if architectureConstraint.satisfies(a) {
		// 	return true
		// }
	}

	return satisfied
}
