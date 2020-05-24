package obi

import (
	"errors"
	"fmt"
	"reflect"
)

func getSchemaImpl(t reflect.Type) ([]byte, error) {
	switch t.Kind() {
	case reflect.Uint8:
		return []byte("u8"), nil
	case reflect.Uint16:
		return []byte("u16"), nil
	case reflect.Uint32:
		return []byte("u32"), nil
	case reflect.Uint64:
		return []byte("u64"), nil
	case reflect.Int8:
		return []byte("i8"), nil
	case reflect.Int16:
		return []byte("i16"), nil
	case reflect.Int32:
		return []byte("i32"), nil
	case reflect.Int64:
		return []byte("i64"), nil
	case reflect.String:
		return []byte("string"), nil
	case reflect.Slice:
		inner, err := getSchemaImpl(t.Elem())
		if err != nil {
			return nil, err
		}
		return append(append([]byte("["), inner...), byte(']')), nil
	case reflect.Struct:
		if t.NumField() == 0 {
			return nil, errors.New("obi: empty struct is not supported")
		}
		res := []byte("{")
		for idx := 0; idx < t.NumField(); idx++ {
			field := t.Field(idx)
			inner, err := getSchemaImpl(field.Type)
			if err != nil {
				return nil, err
			}
			name, ok := field.Tag.Lookup("obi")
			if !ok {
				return nil, fmt.Errorf("obi: no obi tag found for field %s of %s", field.Name, t.Name())
			}
			res = append(append(append(append(res, []byte(name)...), byte(':')), inner...), byte(','))
		}
		res[len(res)-1] = byte('}')
		return res, nil
	default:
		return nil, fmt.Errorf("obi: unsupported value type: %s", t.Kind())
	}
}

// GetSchema returns the compact OBI individual schema of the given value.
func GetSchema(v interface{}) (string, error) {
	schemaBytes, err := getSchemaImpl(reflect.TypeOf(v))
	if err != nil {
		return "", err
	}
	return string(schemaBytes), nil
}

// MustGetSchema returns the compact OBI individual schema of the given value. Panics on error.
func MustGetSchema(v interface{}) string {
	schema, err := GetSchema(v)
	if err != nil {
		panic(err)
	}
	return schema
}
