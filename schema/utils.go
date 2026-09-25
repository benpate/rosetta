package schema

import "fmt"

// getLength returns the length of an object, if it is an ArrayGetter
func getLength(object any) (int, bool) {

	if getter, ok := object.(LengthGetter); ok {
		return getter.Length(), true
	}

	return 0, false
}

// getIndex returns the value at a specific index, if the object is an ArrayGetter
func getIndex(object any, index int) (any, bool) {

	if getter, ok := object.(ArrayGetter); ok {
		return getter.GetIndex(index)
	}

	return nil, false
}

// typeName names a value's Go type for an error detail
func typeName(value any) string {

	// Errors name the type and never the value, because an object being validated or set can
	// hold a password, key, or token.
	return fmt.Sprintf("%T", value)
}
