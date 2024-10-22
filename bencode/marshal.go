package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Marshal takes an interface and returns its Bencoded equivalent.
func Marshal(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer

	err := marshalValue(&buffer, reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// marshalValue handles marshaling of individual values by recursively checking the type
func marshalValue(buffer *bytes.Buffer, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buffer.WriteString("i")
		buffer.WriteString(strconv.FormatInt(value.Int(), 10))
		buffer.WriteString("e")

	case reflect.String:
		strValue := value.String()
		buffer.WriteString(strconv.Itoa(len(strValue)))
		buffer.WriteByte(':')
		buffer.WriteString(strValue)

	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			bytesValue := value.Bytes() // Extract []byte directly
			buffer.WriteString(strconv.Itoa(len(bytesValue)))
			buffer.WriteByte(':')
			buffer.Write(bytesValue) // Write byte slice as a string
		} else {
			buffer.WriteByte('l')
			for i := 0; i < value.Len(); i++ {
				err := marshalValue(buffer, value.Index(i))
				if err != nil {
					return err
				}
			}
			buffer.WriteByte('e')
		}

	case reflect.Map:
		buffer.WriteByte('d')
		keys := value.MapKeys()
		for _, key := range keys {
			if key.Kind() != reflect.String {
				return errors.New("Bencode dictionary keys must be strings")
			}

			// Write key (must be a string)
			keyStr := key.String()
			buffer.WriteString(strconv.Itoa(len(keyStr)))
			buffer.WriteByte(':')
			buffer.WriteString(keyStr)

			// Write value
			err := marshalValue(buffer, value.MapIndex(key))
			if err != nil {
				return err
			}
		}
		buffer.WriteByte('e')

	case reflect.Struct:
		buffer.WriteByte('d')
		structType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := structType.Field(i)
			tag := field.Tag.Get("bencode")
			if tag == "" {
				continue
			}

			// Write key (bencode tag value)
			buffer.WriteString(strconv.Itoa(len(tag)))
			buffer.WriteByte(':')
			buffer.WriteString(tag)

			// Write value (recursively marshal field value)
			err := marshalValue(buffer, value.Field(i))
			if err != nil {
				return err
			}
		}
		buffer.WriteByte('e')

	default:
		fmt.Println(value.Kind())
		return errors.New("unsupported type")
	}

	return nil
}
