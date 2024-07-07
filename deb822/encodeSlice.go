package deb822

import (
	"io"
	"reflect"
)

func encodeSlice(w io.Writer, items reflect.Value, fieldType reflect.StructField, delimiter string) (err error) {
	if it := fieldType.Tag.Get("delim"); it != "" {
		delimiter = it
		if delimiter != "\n" {
			delimiter = it + " "
		}
	}

	for idx := 0; idx < items.Len(); idx++ {
		if idx > 0 {
			w.Write([]byte(delimiter))
		}

		item := items.Index(idx)
		switch item.Kind() {
		case reflect.Struct:
			err = encodeStruct(w, reflect.Indirect(item))
		}

		if err != nil {
			return
		}

	}
	return
}
