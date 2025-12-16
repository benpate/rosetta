package convert

import (
	"net/http"
	"net/url"
	"reflect"
)

// MapOfString attempts to convert the generic value into a map[string]string
func MapOfString(value any) map[string]string {
	if result, _ := MapOfStringOk(value); result != nil {
		return result
	}

	return make(map[string]string)
}

// MapOfStringOk attempts to convert the generic value into a map[string]string
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfStringOk(value any) (map[string]string, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return make(map[string]string), false
	}

	switch typed := value.(type) {

	case map[string]any:
		result := make(map[string]string, len(typed))

		for key, value := range typed {
			result[key] = String(value)
		}
		return result, true

	case map[string]string:
		return typed, true

	case map[string][]string:
		result := make(map[string]string, len(typed))
		for key, value := range typed {
			if len(value) > 0 {
				result[key] = value[0]
			}
		}
		return result, true

	case url.Values:
		return MapOfStringOk(map[string][]string(typed))

	case http.Header:
		return MapOfStringOk(map[string][]string(typed))

	case reflect.Value:
		return MapOfStringOk(Interface(typed))

	case MapOfAnyGetter:
		return MapOfStringOk(typed.MapOfAny())
	}

	// Last chance, try reflection
	if valueOf := reflect.ValueOf(value); valueOf.Kind() == reflect.Map {
		result := make(map[string]string)
		for _, reflectKey := range valueOf.MapKeys() {
			key := String(reflectKey)
			value := String(valueOf.MapIndex(reflectKey).Interface())
			result[key] = value
		}

		return result, true
	}

	// Fall through means conversion failed
	return make(map[string]string), false
}
