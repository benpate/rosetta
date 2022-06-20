package path

import (
	"reflect"
	"strconv"
)

func GetWithReflection(object reflect.Value, path string) (interface{}, bool) {

	kind := object.Kind()

	switch kind {
	case reflect.Array:
		return GetFromSlice(object, path)

	case reflect.Map:
		return GetFromMap(object, path)

	case reflect.Slice:
		return GetFromSlice(object, path)

	case reflect.Struct:
		return GetFromStruct(object, path)

	case reflect.Ptr:
		return GetWithReflection(object.Elem(), path)
	}

	// Fall through means we don't support this kind just yet.
	return nil, false
}

// GetFromMap uses reflection to set a value into a map
func GetFromMap(object reflect.Value, path string) (interface{}, bool) {

	head, tail := Split(path)
	index := reflect.ValueOf(head)
	result := object.MapIndex(index).Interface()
	return GetOK(result, tail)
}

// GetFromStructreturns a value from a struct.  The struct MUST have "path" tags
// that identify how each field is to be addressed.
func GetFromStruct(object reflect.Value, path string) (interface{}, bool) {

	// Get reflect meta-data for this value
	head, tail := Split(path)

	if index := findFieldByTag(object, head); index != -1 {
		result := object.Field(index).Interface()
		return GetOK(result, tail)
	}

	// Fall through means that this path does not exist in this struct.
	return nil, false
}

func GetFromSlice(object reflect.Value, path string) (interface{}, bool) {

	head, tail := Split(path)
	index, err := strconv.Atoi(head)

	// If this is not an integer, then fail
	if err != nil {
		return nil, false
	}

	// Bounds check on index
	if (index < 0) || (index >= object.Len()) {
		return nil, false
	}

	// Return the result to the caller
	result := object.Index(index).Interface()
	return GetOK(result, tail)
}

/***********************************
 * Reflection Utilities
 ***********************************/

// findFieldByTag returns the index of the field whose "path"
// tag matches the provided value.  If none is found, then -1 is returned
func findFieldByTag(value reflect.Value, tag string) int {

	if value.Kind() != reflect.Struct {
		return -1
	}

	typeOf := value.Type()
	count := typeOf.NumField()

	for index := 0; index < count; index++ {
		if typeOf.Field(index).Tag.Get("path") == tag {
			return index
		}
	}

	return -1
}
