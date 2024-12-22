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

func BaseTypeOK(value any) (any, bool) {

	reflectValue := ReflectValue(value)

	switch reflectValue.Kind() {

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
		result := make(map[string]any, reflectValue.Len())
		for _, key := range reflectValue.MapKeys() {
			if value, ok := BaseTypeOK(reflectValue.MapIndex(key)); ok {
				result[key.String()] = value
			} else {
				return nil, false
			}
		}

		return result, true

	default:
		return nil, false
	}
}
