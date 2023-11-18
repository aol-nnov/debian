package deb822

import (
	"fmt"
	"reflect"
	"strconv"
)

func decodeStructValue(field reflect.Value, fieldType reflect.StructField, value string) error {
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Int:
		if value == "" {
			field.SetInt(0)
			return nil
		}
		value, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(value))
		return nil
	case reflect.Slice:
		return decodeStructValueSlice(field, fieldType, value)
	case reflect.Struct:
		return decodeStructValueStruct(field, fieldType, value)
	case reflect.Bool:
		field.SetBool(value == "yes")
		return nil
	}

	return fmt.Errorf("decodeStructValue: decoding field of type %s is not supported", field.Type())
}
