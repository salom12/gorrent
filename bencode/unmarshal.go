package bencode

import (
	"errors"
	"fmt"
	"reflect"
)

// Unmarshal takes a Bencode input and populates the struct fields based on the Bencode tags.
func Unmarshal(data []byte, v interface{}) error {
	// Parse the Bencode input first
	decoder := New(data)
	if err := decoder.Parse(); err != nil {
		return fmt.Errorf("failed to parse bencode: %v", err)
	}

	// We expect the top-level result to be a dictionary
	if len(decoder.Result) == 0 {
		return errors.New("no data to unmarshal")
	}

	dict, ok := decoder.Result[0].(map[string]any)
	if !ok {
		return errors.New("top-level element must be a dictionary")
	}

	// Reflect on the struct and populate fields
	return populateStruct(v, dict)
}

// populateStruct fills the struct `v` with the dictionary `data` based on bencode tags.
func populateStruct(v interface{}, data map[string]any) error {
	val := reflect.ValueOf(v).Elem()
	if val.Kind() != reflect.Struct {
		return errors.New("unmarshal target must be a pointer to a struct")
	}

	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := valType.Field(i)

		// Get the bencode tag
		tag := fieldType.Tag.Get("bencode")
		if tag == "" {
			continue
		}

		// Get the corresponding value from the bencode dictionary
		bencodeValue, found := data[tag]
		if !found {
			continue // Ignore missing fields
		}

		// Populate the field based on its type
		// if err := setValuex(field, bencodeValue); err != nil {
		// 	return fmt.Errorf("failed to set field %s: %v", fieldType.Name, err)
		// }

		func(field reflect.Value, bencodeValue any) error {
			fieldType := field.Type()
			switch fieldType.Kind() {
			case reflect.String:
				// Handle strings
				if str, ok := bencodeValue.([]byte); ok {
					field.SetString(string(str))
				} else {
					return fmt.Errorf("expected string, got %T", bencodeValue)
				}
			case reflect.Int, reflect.Int64, reflect.Int32:
				// Handle integers
				if i, ok := bencodeValue.(int64); ok {
					field.SetInt(int64(i))
				} else {
					return fmt.Errorf("expected int, got %T", bencodeValue)
				}
			case reflect.Slice:
				// Handle slices (lists in bencode)
				if fieldType.Elem().Kind() == reflect.String {
					// Handle []string
					if list, ok := bencodeValue.([]any); ok {
						strSlice := make([]string, 0, len(list))
						for _, item := range list {
							if str, ok := item.([]byte); ok {
								strSlice = append(strSlice, string(str))
							} else {
								return fmt.Errorf("expected string in list, got %T", item)
							}
						}
						field.Set(reflect.ValueOf(strSlice))
					} else {
						return fmt.Errorf("expected list for slice, got %T", bencodeValue)
					}
				} else if fieldType.Elem().Kind() == reflect.Uint8 {
					// Handle []byte
					if list, ok := bencodeValue.([]byte); ok {
						field.Set(reflect.ValueOf(list))
					} else {
						return fmt.Errorf("expected list for slice, got %T", bencodeValue)
					}
				} else {
					return fmt.Errorf("unsupported slice type: %s", fieldType.Elem().Kind())
				}
			case reflect.Struct:
				// Handle nested structs
				if dict, ok := bencodeValue.(map[string]any); ok {
					return populateStruct(field.Addr().Interface(), dict)
				} else {
					return fmt.Errorf("expected dictionary for struct, got %T", bencodeValue)
				}
			default:
				return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
			}
			return nil
		}(field, bencodeValue)
	}

	return nil
}

// setValue sets the value of the struct field based on the type of bencode value.
func setValue(field reflect.Value, bencodeValue any) error {
	fieldType := field.Type()

	switch fieldType.Kind() {
	case reflect.String:
		// Handle strings
		if str, ok := bencodeValue.([]byte); ok {
			field.SetString(string(str))
		} else {
			return fmt.Errorf("expected string, got %T", bencodeValue)
		}

	case reflect.Int, reflect.Int64, reflect.Int32:
		// Handle integers
		if i, ok := bencodeValue.(int64); ok {
			field.SetInt(int64(i))
		} else {
			return fmt.Errorf("expected int, got %T", bencodeValue)
		}

	case reflect.Slice:
		// Handle slices (lists in bencode)
		if fieldType.Elem().Kind() == reflect.String {
			// Handle []string
			if list, ok := bencodeValue.([]any); ok {
				strSlice := make([]string, 0, len(list))
				for _, item := range list {
					if str, ok := item.([]byte); ok {
						strSlice = append(strSlice, string(str))
					} else {
						return fmt.Errorf("expected string in list, got %T", item)
					}
				}
				field.Set(reflect.ValueOf(strSlice))
			} else {
				return fmt.Errorf("expected list for slice, got %T", bencodeValue)
			}
		} else {
			return fmt.Errorf("unsupported slice type: %s", fieldType.Elem().Kind())
		}

	case reflect.Struct:
		// Handle nested structs
		if dict, ok := bencodeValue.(map[string]any); ok {
			return populateStruct(field.Addr().Interface(), dict)
		} else {
			return fmt.Errorf("expected dictionary for struct, got %T", bencodeValue)
		}

	default:
		return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}

	return nil
}
