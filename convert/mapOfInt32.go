package convert

import (
	"reflect"
)

// MapOfInt32 attempts to convert the generic value into a map[string]string
func MapOfInt32(value any) map[string]int32 {
	if result, _ := MapOfInt32Ok(value); result != nil {
		return result
	}

	return make(map[string]int32)
}

// MapOfInt32Ok attempts to convert the generic value into a map[string]string
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfInt32Ok(value any) (map[string]int32, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return make(map[string]int32), false
	}

	switch typed := value.(type) {

	case map[string]int:
		result := make(map[string]int32, len(typed))
		for key, value := range typed {
			result[key] = int32(value)
		}
		return result, true

	case map[string]int8:
		result := make(map[string]int32, len(typed))
		for key, value := range typed {
			result[key] = int32(value)
		}
		return result, true

	case map[string]int16:
		result := make(map[string]int32, len(typed))
		for key, value := range typed {
			result[key] = int32(value)
		}
		return result, true

	case map[string]int32:
		return typed, true

	case map[string]int64:
		result := make(map[string]int32, len(typed))
		for key, value := range typed {
			result[key] = Int32(value)
		}
		return result, true

	case map[string]any:
		result := make(map[string]int32, len(typed))
		for key, value := range typed {
			result[key] = Int32(value)
		}
		return result, false

	case map[string]string:
		result := make(map[string]int32, len(typed))
		for key, value := range typed {
			result[key] = Int32(value)
		}
		return result, false

	case reflect.Value:
		return MapOfInt32Ok(Interface(typed))

	case MapOfAnyGetter:
		return MapOfInt32Ok(typed.MapOfAny())
	}

	// Fall through means conversion failed
	return make(map[string]int32), false
}
