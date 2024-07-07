package fields

import (
	"bytes"
	"fmt"
)

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

func (v VersionConstraint) MarshalText() (res []byte, err error) {
	fmt.Println("VersionConstraint MarshalText")
	return fmt.Appendf(res, "(%v %v)", v.Op, v.Value), nil
}
