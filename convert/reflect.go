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
	if valueOf, ok := value.(reflect.Value); ok {

		if valueOf.Kind() == reflect.Invalid {
			return nil
		}

		return valueOf.Interface()
	}

	// Otherwise, just return the value
	return value
}

// BaseTypeOK attempts to convert a value into a base type (bool, int, float, string, slice, map).
// The boolean result value returns TRUE if successful.  FALSE otherwise
func BaseTypeOK(value any) (any, bool) {

	switch reflectValue := ReflectValue(value); reflectValue.Kind() {

	case reflect.Bool:
		return reflectValue.Bool(), true

	case reflect.Int:
		return int(reflectValue.Int()), true

	case reflect.Int64:
		return reflectValue.Int(), true

	case reflect.Float32:
		return float32(reflectValue.Float()), true

	case reflect.Float64:
		return reflectValue.Float(), true

	case reflect.String:
		return reflectValue.String(), true

	case reflect.Slice, reflect.Array:
		result := make([]any, reflectValue.Len())
		for i := 0; i < reflectValue.Len(); i++ {
			if value, ok := BaseTypeOK(reflectValue.Index(i)); ok {
				result[i] = value
			} else {
				return nil, false
			}
		}
		return result, true

	case reflect.Map:
		result := make(map[string]any)
		for _, key := range reflectValue.MapKeys() {
			if value, ok := BaseTypeOK(reflectValue.MapIndex(key)); ok {
				result[key.String()] = value
			} else {
				return nil, false
			}
		}

		return result, true

	}

	// Fall through is failure
	return nil, false
}
