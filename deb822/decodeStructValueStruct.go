package deb822

import (
	"encoding"
	"fmt"
	"reflect"
)

func decodeStructValueStruct(incoming reflect.Value, incomingField reflect.StructField, data string) error {
	// We've got a complex type to decode into. If it supports encoding.TextUnmarshaler, it's up to developer to
	// implement the parsing algorythm. Otherwise, you know, bail out with an error.
	elem := incoming.Addr()

	if unmarshal, ok := elem.Interface().(encoding.TextUnmarshaler); ok {
		return unmarshal.UnmarshalText([]byte(data))
	}

	return fmt.Errorf(
		"decodeStructValueStruct: type '%s' does not implement encoding.TextUnmarshaler",
		incomingField.Name,
	)
}
