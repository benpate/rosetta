package convert

import (
	"reflect"
)

// MapOfInt attempts to convert the generic value into a map[string]string
func MapOfInt(value any) map[string]int {
	if result, _ := MapOfIntOk(value); result != nil {
		return result
	}

	return make(map[string]int)
}

// MapOfIntOk attempts to convert the generic value into a map[string]string
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfIntOk(value any) (map[string]int, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return make(map[string]int), false
	}

	switch typed := value.(type) {

	case map[string]int:
		return typed, true

	case map[string]int8:
		result := make(map[string]int, len(typed))
		for key, value := range typed {
			result[key] = int(value)
		}
		return result, true

	case map[string]int16:
		result := make(map[string]int, len(typed))
		for key, value := range typed {
			result[key] = int(value)
		}
		return result, true

	case map[string]int32:
		result := make(map[string]int, len(typed))
		for key, value := range typed {
			result[key] = int(value)
		}
		return result, true

	case map[string]int64:
		result := make(map[string]int, len(typed))
		for key, value := range typed {
			result[key] = int(value)
		}
		return result, true

	case map[string]any:
		result := make(map[string]int, len(typed))
		for key, value := range typed {
			result[key] = Int(value)
		}
		return result, false

	case map[string]string:
		result := make(map[string]int, len(typed))
		for key, value := range typed {
			result[key] = Int(value)
		}
		return result, false

	case reflect.Value:
		return MapOfIntOk(Interface(typed))

	case MapOfAnyGetter:
		return MapOfIntOk(typed.MapOfAny())
	}

	// Last chance, try reflection
	v := reflect.ValueOf(value)
	if v.Type().Kind() == reflect.Map {

		result := make(map[string]int)

		for _, reflectKey := range v.MapKeys() {
			key := String(reflectKey)
			value := Int(v.MapIndex(reflectKey).Interface())
			result[key] = value
		}

		return result, true
	}

	// Fall through means conversion failed
	return make(map[string]int), false
}
