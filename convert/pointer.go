package convert

import "reflect"

// Pointer returns a pointer to the original value
func Pointer[T any](original T) *T {
	return &original
}

// Element defreferences a pointer, if necessary, and returns the underlying value
func Element(original any) any {

	if valueOf := reflect.ValueOf(original); valueOf.Kind() == reflect.Ptr {
		return valueOf.Elem().Interface()
	}

	return original
}
