package fields

import (
	"bytes"
	"fmt"
)

type VersionConstraintOperator int

const (
	VersionConstraintNotSet VersionConstraintOperator = iota
	VersionConstraintLessThan
	VersionConstraintLessOrEqual
	VersionConstraintEqual
	VersionConstraintGreaterOrEqual
	VersionConstraintGreaterThan
)

/*
BNF descriptor:

	<VersionConstraint> ::= "(" <Op> " " <Value> ")" | E
	<Op> ::= "<<" | "<=" | "=" | ">=" | ">>"
	<Value> ::= <Version>
*/
type VersionConstraint struct {
	Op    VersionConstraintOperator
	Value Version
}

func (v VersionConstraint) String() string {
	return fmt.Sprintf("(%v %v)", v.Op, v.Value)
}

func (v *VersionConstraint) UnmarshalText(text []byte) error {
	if bytes.TrimSpace(text)[0] != '(' {
		return fmt.Errorf("VersionConstraint unmarshal: wrong input string '%s'", text)
	}

	if opBytes, verBytes, found := bytes.Cut(bytes.Trim(text, "() "), []byte{' '}); found {

		if err := v.Op.UnmarshalText(opBytes); err != nil {
			return err
		}

		if err := v.Value.UnmarshalText(verBytes); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("VersionConstraint unmarshal failed. Got '%s'", string(text))
}

func (v *VersionConstraint) SatisfiedBy(another Version) bool {
	if v.Op == VersionConstraintNotSet {
		return true
	}

	cmpRes := v.Value.Compare(another)

	switch v.Op {
	case VersionConstraintGreaterThan, VersionConstraintGreaterOrEqual:
		return cmpRes == VersionCompareResultGreaterThan || cmpRes == VersionCompareResultEquals
	case VersionConstraintEqual:
		return cmpRes == VersionCompareResultEquals
	case VersionConstraintLessThan, VersionConstraintLessOrEqual:
		return cmpRes == VersionCompareResultLessThan || cmpRes == VersionCompareResultEquals
	}

	return false
}

var verConstraintDecoder = map[byte]map[byte]VersionConstraintOperator{
	'<': {
		'<': VersionConstraintLessThan,
		'=': VersionConstraintLessOrEqual,
	},
	'>': {
		'>': VersionConstraintGreaterThan,
		'=': VersionConstraintGreaterOrEqual,
	},
}

func (v *VersionConstraintOperator) UnmarshalText(text []byte) error {
	text = bytes.TrimSpace(text)

	if text[0] == '=' && len(text) == 1 {
		*v = VersionConstraintEqual
		return nil
	}

	if len(text) == 2 {
		if vc, found := verConstraintDecoder[text[0]][text[1]]; found {
			*v = vc
			return nil
		}
	}

	return fmt.Errorf("VersionConstraintOperator unmarshal failed. Got '%s'", string(text))
}

func (v VersionConstraintOperator) MarshalText() (text []byte, err error) {
	return fmt.Appendf(text, "%v", v), nil
}

func (v VersionConstraintOperator) String() string {
	return [...]string{"-", "<<", "<=", "=", ">=", ">>"}[v]
	// return [...]string{"lt", "le", "eq", "ge", "gt"}[v]
}
