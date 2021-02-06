package forge

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

var (
	// ErrInvalidValue returned when the value passed to Unmarshal is nil or not a
	// pointer to a struct.
	ErrInvalidValue = errors.New("value must be a non-nil pointer to a struct")

	// ErrUnsupportedType returned when a field with tag "env" is unsupported.
	ErrUnsupportedType = errors.New("field is an unsupported type")

	// ErrUnexportedField returned when a field with tag "env" is not exported.
	ErrUnexportedField = errors.New("field must be exported")
)

// ParseEnvironment variables into a existing struct
func ParseEnvironment(target interface{}) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidValue
	}

	rv = rv.Elem()

	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		valueField := rv.Field(i)
		switch valueField.Kind() {
		case reflect.Struct:
			valueInterface := valueField.Addr().Interface()
			err := ParseEnvironment(valueInterface)
			if err != nil {
				return err
			}
		}

		typeField := t.Field(i)
		tag := typeField.Tag.Get("env")
		if tag == "" {
			continue
		}

		if !valueField.CanSet() {
			return ErrUnexportedField
		}

		envVar, ok := os.LookupEnv(tag)
		if !ok {
			continue
		}

		err := reflectSet(typeField.Type, valueField, envVar)
		if err != nil {
			return err
		}
	}

	return nil
}

func reflectSet(t reflect.Type, f reflect.Value, value string) error {
	switch t.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(t.Elem())
		err := reflectSet(t.Elem(), ptr.Elem(), value)
		if err != nil {
			return err
		}
		f.Set(ptr)
	case reflect.String:
		f.SetString(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		f.SetBool(v)
	case reflect.Int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		f.SetInt(int64(v))
	default:
		return ErrUnsupportedType
	}

	return nil
}
