package pointer

import "reflect"

// To returns a pointer to a copy of the provided value.
// A nil, or a value that is already a pointer, is returned as-is.
func To(value any) any {

	// RULE: A nil has no type to build a pointer to. reflect.TypeOf returns nil for it, and
	// reflect.New(nil) panics, so nil is handed straight back.
	if value == nil {
		return nil
	}

	// Pointers and interfaces already carry a reference, so return them as-is.
	switch reflect.ValueOf(value).Kind() {

	case reflect.Pointer, reflect.Interface:
		return value
	}

	// Create a new pointer to the provided value
	ptrValue := reflect.New(reflect.TypeOf(value))
	reflect.Indirect(ptrValue).Set(reflect.ValueOf(value))
	return ptrValue.Interface()
}
