package compare

import (
	"reflect"
)

// IsNil returns TRUE if the provided value is nil.  This uses
func IsNil(value any) bool {

	if value == nil {
		return true
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
		return reflect.ValueOf(value).IsNil()
	}

	return false
}

// NotNil returns TRUE if the provided value is NOT nil
func NotNil(value any) bool {
	return !IsNil(value)
}
