package convert

import "reflect"

// Pointer returns a pointer to the original value
func Pointer[T any](original T) *T {
	return &original
}

// Element defreferences a pointer, if necessary, and returns the underlying value
func Element(original any) any {

	if valueOf := reflect.ValueOf(original); valueOf.Kind() == reflect.Pointer {

		// RULE: A nil pointer dereferences to nothing, not to a zero Value
		if valueOf.IsNil() {
			return nil
		}

		return valueOf.Elem().Interface()
	}

	return original
}
