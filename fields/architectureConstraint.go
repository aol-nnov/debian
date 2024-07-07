package fields

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
