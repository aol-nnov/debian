package fields

import (
	"bytes"
	"fmt"
)

/*
Architecture constraints are used in different fields of Debian control file, mostly, in Build-Depends

BNF form is as follows:

	<ArchitectureConstraints> ::= "[" <constraints> "]"
	<constraints> ::= <constraint> | <constraints> " " <constraint>
	<constraint> ::= <negate> <name>
	<negate> ::=  E | "!"
	<name> ::= ([0-9] | [a-z])+

Examples:

	[amd64 armhf]
	[armel !kfreebsd]
*/
type ArchitectureConstraints []architectureConstraint

type architectureConstraint struct {
	Negate bool
	Name   Architecture
}

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

func (ac architectureConstraint) satisfiedBy(a Architecture) bool {
	res := (!ac.Negate && ac.Name.Equals(a)) || (ac.Negate && !ac.Name.Equals(a))
	// fmt.Printf("%v %v %v\n", ac, res, a)
	return res
}

// checks if ArchitectureConstraints are satisfied by provided Architecture
func (constraints ArchitectureConstraints) SatisfiedBy(a Architecture) bool {
	// ArchitectureConstraints are OR-ed

	// empty set satisfies any Architecture
	if len(constraints) == 0 {
		return true
	}

	satisfied := false

	for _, architectureConstraint := range constraints {
		satisfied = satisfied || architectureConstraint.satisfiedBy(a)

		// fastlane (baling out on first match)
		// if architectureConstraint.satisfies(a) {
		// 	return true
		// }
	}

	return satisfied
}
