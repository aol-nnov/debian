package fields

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aol-nnov/debian/internal/stringspp"
)

/*
Build profile constraints as described in https://wiki.debian.org/BuildProfileSpec

BNF descriptor:

	<ProfileConstraints> ::= <AndSet> | <ProfileConstraints> " " <AndSet>
	<AndSet> ::= "<" <andConstraints> ">"
	<andConstraints> ::= <profileConstraint> | <andConstraints> " " <profileConstraint>
	<profileConstraint> ::= <Negate> <Name>
	<Negate> ::= E | "!"
	<Name> ::= ([0-9] | [a-z] | [A-Z])+

Example:

	<p1 !p2> <p3 p4> === (p1 & !p2) | (p3 & p4)

After unmarshalling

		ProfileConstraints[0]: [p1, !p2]
		ProfileConstraints[1]: [p3, p4]

	    Total value of ProfileConstraints is calculated as follows: first dimension elements are OR-ed (||), second dimension elements are AND-ed (&&)

[Online editor] for BNF checking and fiddling.

[Online editor]: http://bnfplayground.pauliankline.com/?bnf=%3CProfileConstraints%3E%20%3A%3A%3D%20%3CAndSet%3E%20%7C%20%3CProfileConstraints%3E%20%22%20%22%20%3CAndSet%3E%0A%3CAndSet%3E%20%3A%3A%3D%20%22%3C%22%20%3CandConstraints%3E%20%22%3E%22%0A%3CandConstraints%3E%20%3A%3A%3D%20%3CprofileConstraint%3E%20%7C%20%3CandConstraints%3E%20%22%20%22%20%3CprofileConstraint%3E%0A%3CprofileConstraint%3E%20%3A%3A%3D%20%3CNegate%3E%20%3CName%3E%0A%3CNegate%3E%20%3A%3A%3D%20E%20%7C%20%22!%22%0A%3CName%3E%20%3A%3A%3D%20(%5B0-9%5D%20%7C%20%5Ba-z%5D%20%7C%20%5BA-Z%5D)%2B&name=Profile%20Constraints
*/
type ProfileConstraints []andConstraints
type andConstraints []profileConstraint

type profileConstraint struct {
	Negate bool
	Name   string
}

// unmarshals byte slice to sets, that are AND-ed
func unmarshalAndSet(text []byte) (andConstraints, error) {
	var res andConstraints

	profiles := bytes.Split(bytes.Trim(text, "<> "), []byte(" "))
	for _, profile := range profiles {
		if len(profile) == 0 {
			continue
		}

		var pc profileConstraint
		if err := pc.UnmarshalText(bytes.TrimSpace(profile)); err != nil {
			return nil, err
		}

		res = append(res, pc)
	}

	return res, nil
}

// [pkg/encoding.TextUnmarshaler] implementation
func (constraints *ProfileConstraints) UnmarshalText(text []byte) (err error) {
	tail := bytes.TrimSpace(text)
	if !bytes.HasPrefix(tail, []byte{'<'}) ||
		!bytes.HasSuffix(tail, []byte{'>'}) {
		return fmt.Errorf("ProfileConstraints unmarshal: wrong input string '%s'", text)
	}

	var andSet []byte
	for found := true; found; {
		andSet, found, tail = stringspp.Between(tail, '<', '>', false)

		if andSet = bytes.TrimSpace(andSet); len(andSet) > 0 {

			var andConstraints []profileConstraint
			if andConstraints, err = unmarshalAndSet(andSet); err == nil {
				*constraints = append(*constraints, ProfileConstraints{andConstraints}...)
			} else {
				return err
			}
		}
	}

	return err
}

func (pc *profileConstraint) UnmarshalText(text []byte) (err error) {
	//	fmt.Printf("unm prof '%s'\n", text)

	if text[0] == '!' {
		pc.Negate = true
		pc.Name = string(text[1:])
	} else {
		pc.Negate = false
		pc.Name = string(text)
	}

	return nil
}

// [pkg/fmt.Stringer] implementation
func (pc profileConstraint) String() string {
	not := ""

	if pc.Negate {
		not = "!"
	}

	return fmt.Sprintf("%s%s", not, pc.Name)
}

// [pkg/fmt.Stringer] implementation
func (pc ProfileConstraints) String() string {
	res := ""

	for _, andSet := range pc {
		res += "<"
		for _, profile := range andSet {
			res += profile.String() + " "
		}
		res = strings.TrimSpace(res) + "> "
	}

	return strings.TrimSpace(res)
}

/*
Check

truth table for comparing single profile with single ProfileConstraint:

	neg -> pc.Negate
	cmp -> pc.Name == profile

	neg | cmp | res | example
	====+=====+=====+========
	 f |   f  |  f  |  profile1 vs profile2
	 f |   t  |  t  |  profile1 vs profile1
	 t |   f  |  t  | !profile1 vs profile2
	 t |   t  |  f  | !profile1 vs profile1

	 Synthesis:
	 for each row with truthy RESult, falsy components (neg or cmp) are negated, then components are and-ed.
	 Tuples then or-ed.

	 (!neg && cmp) || (neg && !cmp) => unable to simplify
*/
func (pc profileConstraint) satisfiedBy(activeProfiles []string) bool {

	fullySatisfied := pc.Negate

	for _, profile := range activeProfiles {
		constraintSatisfied := (!pc.Negate && pc.Name == profile) || (pc.Negate && !(pc.Name == profile))
		// fmt.Printf("%v %v %v\n", pc, constraintSatisfied, profile)

		if pc.Negate {
			fullySatisfied = fullySatisfied && constraintSatisfied
		} else {
			fullySatisfied = fullySatisfied || constraintSatisfied
		}
	}

	return fullySatisfied
}

func (orSet ProfileConstraints) SatisfiedBy(activeProfiles []string) bool {
	/*
		orSet: [ andSet || andSet || andSet]
		andSet: [ profileConstraint && profileConstraint && profileConstraint ]
	*/
	if len(orSet) == 0 {
		return true
	}

	// `profileConstraint`-s stored in each item of `orSet` are OR-ed
	for _, andSet := range orSet {
		andSetSatisfies := true

		// ProfileConstraint-s inside each group are and-ed
		for _, profile := range andSet {
			andSetSatisfies = andSetSatisfies && profile.satisfiedBy(activeProfiles)
		}

		// According to OR-operator truth table, if any of perands are truthy, result is truthy
		// so, we can bail out after first match
		if andSetSatisfies {
			return true
		}
	}

	return false
}
