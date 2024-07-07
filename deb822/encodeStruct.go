package deb822

import (
	"encoding"
	"fmt"
	"io"
	"reflect"
)

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	case reflect.Struct:
		return reflect.Zero(v.Type()) == v
	}
	return false
}

func encodeStruct(w io.Writer, from reflect.Value) error {
	// If we have a pointer, let's follow it
	if from.Type().Kind() == reflect.Ptr {
		return encodeStruct(w, from.Elem())
	}

	if marshal, ok := from.Addr().Interface().(encoding.TextMarshaler); ok {
		// fmt.Fprintf(os.Stderr, "encodeStruct %T is TextMarshaler\n", from.Addr().Interface())
		var err error
		var res []byte
		if res, err = marshal.MarshalText(); err == nil {
			_, err = w.Write(res)
		}
		return err
	}
	// else {
	// 	fmt.Fprintf(os.Stderr, "encodeStruct %T is NOT TextMarshaler\n", from.Addr().Interface())
	// }

	for idx, field := range reflect.VisibleFields(from.Type()) {
		fieldName := field.Name

		nameFromTag, opts := parseTag(field.Tag.Get("deb822"))

		if nameFromTag != "" {
			fieldName = nameFromTag
		}

		if fieldName == "-" {
			continue
		}

		if isEmptyValue(from.Field(idx)) {
			if opts.Contain("omitempty") {
				continue
			}

			if field.Tag.Get("required") == "true" {
				return fmt.Errorf("encodeStruct: missing value for required field '%s'", fieldName)
			}
		}

		if fieldName != "" {
			fmt.Fprintf(w, "%s: ", fieldName)
		}

		encodeStructValue(w, from.Field(idx), field)

		fmt.Fprint(w, "\n")
	}

	return nil
}
