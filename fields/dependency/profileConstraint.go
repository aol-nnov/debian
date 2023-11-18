package deb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aol-nnov/debian/internal/stringspp"
)

/*
https://wiki.debian.org/BuildProfileSpec

*/
// first dimension elements are or-ed (||), second dimension elements are and-ed (&&)
type ProfileConstraints [][]ProfileConstraint

type ProfileConstraint struct {
	Name   string
	Negate bool
}

func unmarshalAndSet(text []byte) ([]ProfileConstraint, error) {
	var res []ProfileConstraint

	profiles := bytes.Split(bytes.Trim(text, "<> "), []byte(" "))
	for _, profile := range profiles {
		if len(profile) == 0 {
			continue
		}

		var pc ProfileConstraint
		if err := pc.UnmarshalText(bytes.TrimSpace(profile)); err != nil {
			return nil, err
		}

		res = append(res, pc)
	}

	return res, nil
}

func (constraints *ProfileConstraints) UnmarshalText(text []byte) (err error) {
	tail := bytes.TrimSpace(text)
	if tail[0] != '<' {
		return fmt.Errorf("ProfileConstraints unmarshal: wrong input string '%s'", text)
	}

	var andSet []byte
	for found := true; found; {
		andSet, found, tail = stringspp.Between(tail, '<', '>', false)

		if andSet = bytes.TrimSpace(andSet); len(andSet) > 0 {

			var andConstraints []ProfileConstraint
			if andConstraints, err = unmarshalAndSet(andSet); err == nil {
				*constraints = append(*constraints, ProfileConstraints{andConstraints}...)
			} else {
				return err
			}
		}
	}

	return err
}

func (pc *ProfileConstraint) UnmarshalText(text []byte) (err error) {

	if text[0] == '!' {
		pc.Negate = true
		pc.Name = string(text[1:])
	} else {
		pc.Negate = false
		pc.Name = string(text)
	}

	return nil
}

func (pc ProfileConstraint) String() string {
	not := ""

	if pc.Negate {
		not = "!"
	}

	return fmt.Sprintf("%s%s", not, pc.Name)
}

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

func (pc ProfileConstraint) satisfies(profiles []string) bool {

	satisfies := pc.Negate

	for _, profile := range profiles {
		constraintSatisfied := (!pc.Negate && pc.Name == profile) || (pc.Negate && !(pc.Name == profile))
		// fmt.Printf("%v %v %v\n", pc, constraintSatisfied, profile)

		if pc.Negate {
			satisfies = satisfies && constraintSatisfied
		} else {
			satisfies = satisfies || constraintSatisfied
		}
	}

	return satisfies
}

func (orSet ProfileConstraints) Satisfies(activeProfiles []string) bool {
	/*
		orSet: [ andSet || andSet || andSet]
		andSet: [ profileConstraint && profileConstraint && profileConstraint ]
	*/
	if len(orSet) == 0 {
		return true
	}

	// groups are OR-ed
	for _, andSet := range orSet {
		andSetSatisfies := true

		// ProfileConstraint-s inside each group are and-ed
		for _, profile := range andSet {
			andSetSatisfies = andSetSatisfies && profile.satisfies(activeProfiles)
		}

		// According to OR-operator truth table, if any of perands are truthy, result is truthy
		// so, we can bail out after first match
		if andSetSatisfies {
			return true
		}
	}

	return false
}
