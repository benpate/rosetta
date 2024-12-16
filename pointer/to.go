package pointer

import "reflect"

// To returns the provided value as a pointer to the original.
// If the provided value is already a pointer, it is returned as-is.
func To(value any) any {

	// Some reflect magic to make sure we're working with a pointer
	switch reflect.ValueOf(value).Kind() {

	case reflect.Ptr:
		return value

	case reflect.Interface:
		return value
	}

	// Create a new pointer to the provided value
	ptrValue := reflect.New(reflect.TypeOf(value))
	reflect.Indirect(ptrValue).Set(reflect.ValueOf(value))
	return ptrValue.Interface()
}
