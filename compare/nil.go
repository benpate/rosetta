package compare

import (
	"reflect"
)

// IsNil returns TRUE if the value is nil, or is a nil pointer, slice, map, channel, or func.
func IsNil(value any) bool {

	if value == nil {
		return true
	}

	// Only these kinds can hold nil; asking an array panics
	switch reflect.TypeOf(value).Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Chan, reflect.Map,
		reflect.Func, reflect.UnsafePointer:
		return reflect.ValueOf(value).IsNil()
	}

	return false
}

// NotNil returns TRUE if the provided value is NOT nil
func NotNil(value any) bool {
	return !IsNil(value)
}
