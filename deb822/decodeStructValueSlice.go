package deb822

import (
	"reflect"
	"strings"
)

func decodeStructValueSlice(field reflect.Value, fieldType reflect.StructField, value string) error {
	underlyingType := field.Type().Elem()

	var delim = " "
	if tagDelim := fieldType.Tag.Get("delim"); tagDelim != "" {
		delim = tagDelim
	}

	var strip = ""
	if tagStrip := fieldType.Tag.Get("strip"); tagStrip != "" {
		strip = tagStrip
	}

	for _, el := range strings.Split(strings.Trim(value, strip), delim) {
		el = strings.Trim(el, strip)

		targetValue := reflect.New(underlyingType)

		if err := decodeStructValue(targetValue.Elem(), fieldType, el); err != nil {
			return err
		}
		field.Set(reflect.Append(field, targetValue.Elem()))
	}

	return nil
}
