package convert

import (
	"net/http"
	"net/url"
	"reflect"
)

// MapOfAny attempts to convert the generic value into a map[string]any
func MapOfAny(value any) map[string]any {
	if result, _ := MapOfAnyOk(value); result != nil {
		return result
	}
	return make(map[string]any)
}

// MapOfAnyOk attempts to convert the generic value into a map[string]any
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfAnyOk(value any) (map[string]any, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return make(map[string]any), false
	}

	switch typed := value.(type) {

	case map[string]any:
		return typed, true

	case map[string]string:
		result := make(map[string]any, len(typed))

		for key, value := range typed {
			result[key] = value
		}
		return result, true

	case map[string][]string:
		result := make(map[string]any, len(typed))
		for key, value := range typed {
			switch len(value) {
			case 0:
			case 1:
				result[key] = value[0]
			default:
				result[key] = value
			}
		}
		return result, true

	case url.Values:
		return MapOfAnyOk(map[string][]string(typed))

	case http.Header:
		return MapOfAnyOk(map[string][]string(typed))

	case reflect.Value:
		return MapOfAnyOk(Interface(typed))

	case MapOfAnyGetter:
		return typed.MapOfAny(), true
	}

	// Last chance, try reflection
	if valueOf := reflect.ValueOf(value); valueOf.Type().Kind() == reflect.Map {

		result := make(map[string]any)
		for _, reflectKey := range valueOf.MapKeys() {
			key := String(reflectKey)
			result[key] = valueOf.MapIndex(reflectKey).Interface()
		}

		return result, true
	}

	// Fall through means conversion failed
	return make(map[string]any), false
}
