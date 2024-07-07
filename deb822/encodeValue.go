package deb822

import (
	"encoding"
	"fmt"
	"io"
	"reflect"
)

func encodeStructValue(w io.Writer, field reflect.Value, fieldType reflect.StructField) (err error) {
	// If we have a pointer, let's follow it
	// if field.Type().Kind() == reflect.Ptr {
	// 	return encodeValue(w, field.Elem(), fieldType)
	// }

	// reflect.PointerTo(field.Type()).Implements(reflect.TypeFor[encoding.TextMarshaler]())
	if marshal, ok := field.Addr().Interface().(encoding.TextMarshaler); ok {
		// fmt.Fprintf(os.Stderr, "encodeStructValue %v is TextMarshaler\n", field.Type().String())
		var err error
		var res []byte
		if res, err = marshal.MarshalText(); err == nil {
			_, err = w.Write(res)
		}
		return err
	}
	// else {
	// 	fmt.Fprintf(os.Stderr, "encodeStructValue %v is NOT TextMarshaler\n", field.Type().String())
	// }

	switch field.Kind() {
	case reflect.String:
		_, err = fmt.Fprint(w, field.String())
	case reflect.Struct:
		err = encodeStruct(w, field)
	case reflect.Slice:
		err = encodeSlice(w, field, fieldType, "")
	default:
		err = fmt.Errorf("unable to encode from a %s", field.Type().String())
	}
	return
}
