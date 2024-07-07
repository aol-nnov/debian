package fields

import "fmt"

func (d Dependencies) MarshalText() (text []byte, err error) {
	for idx, dep := range d {
		depStr, err := dep.MarshalText()
		if err != nil {
			return nil, err
		}

		if idx == 0 {
			text = depStr
		} else {
			text = fmt.Appendf(text, ",\n %s", depStr)
		}
	}
	return
}
