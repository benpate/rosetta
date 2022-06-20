package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// setWithReflection applies apply a value to a generic variable.
func setWithReflection(object reflect.Value, value interface{}) error {

	var err error

	kind := object.Kind()

	defer func() {
		if r := recover(); r != nil {
			err = derp.NewInternalError("schema.setWithReflection", "Could not set value", r, kind.String(), value)
		}
	}()

	switch kind {

	// dereference pointers
	case reflect.Pointer:
		return setWithReflection(object.Elem(), value)

	// Simple Variable Types
	case reflect.Bool:
		object.SetBool(convert.Bool(value))

	case
		reflect.Float32,
		reflect.Float64:
		object.SetFloat(convert.Float(value))

	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		object.SetInt(convert.Int64(value))

	case reflect.String:
		object.SetString(convert.String(value))

	// Complex Variables Here

	default:
		return derp.NewInternalError("schema.setWithReflection", "Cannot set complex types with this function", value)
		// object.Set(reflect.ValueOf(value))
	}

	return err
}
