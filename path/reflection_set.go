package path

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// SetWithReflection uses a "path" to apply a value to a generic variable.
func SetWithReflection(object reflect.Value, path string, value interface{}) error {

	kind := object.Kind()

	// Dereference pointers
	if kind == reflect.Ptr {
		return SetWithReflection(object.Elem(), path, value)
	}

	// If the path is empty, then this is the final property to actually set.
	if path == "" {

		switch kind {

		case reflect.Bool:
			object.SetBool(convert.Bool(value))
			return nil

		case
			reflect.Float32,
			reflect.Float64:

			object.SetFloat(convert.Float(value))
			return nil

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
			return nil

		case reflect.String:
			object.SetString(convert.String(value))
			return nil

		}
	}

	// Fall through means that we need to dig further before we can set a value.

	switch kind {

	case reflect.Array, reflect.Slice:
		return SetToSlice(object, path, value)

	case reflect.Map:
		return SetToMap(object, path, value)

	case reflect.Struct:
		return SetToStruct(object, path, value)

	case reflect.Ptr:
		return SetWithReflection(object.Elem(), path, value)
	}

	// Fall through means we don't support this kind just yet.
	return derp.NewInternalError("path.SetWithReflection", "Cannot set path on this variable", object, path, value)
}

// SetToMap uses reflection to set a value into a map
func SetToMap(object reflect.Value, path string, value interface{}) error {
	head, tail := Split(path)

	if tail != "" {
		return derp.NewInternalError("path.SetToMap", "Cannot use Reflection to create sub-element in a map.  Try making a setter for this value", path, value)
	}

	key := reflect.ValueOf(head)
	object.SetMapIndex(key, reflect.ValueOf(value))
	return nil
}

// SetToSlice uses reflection to set a value into a slice/array variable.
func SetToSlice(object reflect.Value, path string, value interface{}) error {
	head, tail := Split(path)
	max := object.Len()

	index, err := Index(head, max)

	if err != nil {
		return derp.Wrap(err, "path.SetToSlice", "Cannot get index", object, path, value)
	}

	return SetWithReflection(object.Index(index), tail, value)
}

// SetToStruct sets a map reflection value
func SetToStruct(object reflect.Value, path string, value interface{}) error {

	head, tail := Split(path)

	if index := findFieldByTag(object, head); index != -1 {
		field := object.Field(index)
		return SetWithReflection(field, tail, value)
	}

	return derp.NewInternalError("path.SetToStruct", "Field not found", head, value)
}
