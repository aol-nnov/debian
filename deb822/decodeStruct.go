package deb822

import (
	"fmt"
	"reflect"
)

func decodeStruct(s stanza, into reflect.Value) error {
	// If we have a pointer, let's follow it
	if into.Type().Kind() == reflect.Ptr {
		return decodeStruct(s, into.Elem())
	}

	// Right, now, we're going to decode a [stanza] into the struct
	// fmt.Printf("%s has %d fields\n", into.Type().Name(), into.NumField())
	for i := 0; i < into.NumField(); i++ {

		field := into.Field(i)
		fieldType := into.Type().Field(i)

		if field.Type().Kind() == reflect.Struct {
			err := decodeStruct(s, field)
			if err != nil {
				return err
			}
		}

		// Get the name of the field as we'd index into the [stanza]
		fieldName := fieldType.Name
		if name, _ := parseTag(fieldType.Tag.Get("deb822")); name != "" {
			fieldName = name
		}

		if fieldName == "-" {
			// If the key is "-", lets go ahead and skip it
			continue
		}

		if value, ok := s[fieldName]; ok {
			if err := decodeStructValue(field, fieldType, value); err != nil {
				return err
			}
			continue
		} else {
			if fieldType.Tag.Get("required") == "true" {
				return fmt.Errorf(
					"%s: required field '%s' is missing",
					into.Type().Name(),
					fieldType.Name,
				)
			}

			// check alternate field value
			if alias := fieldType.Tag.Get("if_missing"); alias != "" {
				if value, ok := s[alias]; ok {
					//fmt.Printf("value from stanza '%s'\n", value)
					if err := decodeStructValue(field, fieldType, value); err != nil {
						return err
					}
				} else {
					// if alias field is also missing and current field is marked as required, bail out with error
					if fieldType.Tag.Get("required") == "true" {
						return fmt.Errorf(
							"%s: required field '%s' is missing, alias %s is missing too",
							into.Type().Name(),
							fieldType.Name,
							alias,
						)
					}
				}
			}
			// TODO(aol): is this diagnostic useful at all?
			// if fieldType.Tag.Get("recommended") == "true" {
			// 	fmt.Printf(
			// 		"%s: recommended field '%s' is missing!\n",
			// 		into.Type().Name(),
			// 		fieldType.Name,
			// 	)
			// }
			continue
		}
	}

	return nil
}
