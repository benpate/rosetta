package convert

import (
	"net/url"
	"reflect"
)

type Maplike interface {
	MapOfAny() map[string]any
}

// MapOfAny attempts to convert the generic value into a map[string]any
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfAny(value any) map[string]any {
	result, _ := MapOfAnyOk(value)
	return result
}

// MapOfAnyOk attempts to convert the generic value into a map[string]any
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfAnyOk(value any) (map[string]any, bool) {

	switch v := value.(type) {

	case map[string]any:
		return v, true

	case map[string]string:
		result := make(map[string]any, len(v))

		for key, value := range v {
			result[key] = value
		}
		return result, true

	case url.Values:
		result := make(map[string]any, len(v))
		for key, value := range v {
			switch len(value) {
			case 0:
			case 1:
				result[key] = value[0]
			default:
				result[key] = value
			}
		}
		return result, true

	case reflect.Value:
		return MapOfAnyOk(Interface(v))

	case Maplike:
		return v.MapOfAny(), true

	}

	// Fall through means conversion failed
	return make(map[string]any), false
}
