package deb822

import (
	"fmt"
	"io"
	"reflect"
)

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: w,
	}
}

func (enc *Encoder) Encode(v any) (err error) {
	from := reflect.ValueOf(v)

	switch from.Kind() {
	case reflect.Struct:
		err = encodeStruct(enc.writer, from)
	case reflect.Slice:
		err = encodeSlice(enc.writer, from, reflect.StructField{}, "\n")
	default:
		err = fmt.Errorf("unable to encode from a %s", from.Type().String())
	}

	return
}
