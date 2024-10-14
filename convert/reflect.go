package convert

import "reflect"

// ReflectValue returns the reflect.Value of the given argument.
// If the argument is already a reflect.Value, then it is returned as-is.
func ReflectValue(value any) reflect.Value {

	if valueOf, ok := value.(reflect.Value); ok {
		return valueOf
	}

	return reflect.ValueOf(value)
}

// ReflectType returns the reflect.Type of the given argument.
// If the argument is already a reflect.Type, then it is returned as-is.
func ReflectType(value any) reflect.Type {

	if typeOf, ok := value.(reflect.Type); ok {
		return typeOf
	}

	return reflect.TypeOf(value)
}

// Interface returns the value of a reflect.Value.
// If the value is not already a reflect.Value, then it is returned as-is.
func Interface(value any) any {

	// Safe handling of reflection values
	if v, ok := value.(reflect.Value); ok {

		if v.Kind() == reflect.Invalid {
			return nil
		}

		return v.Interface()
	}

	// Otherwise, just return the value
	return value
}
