package convert

import (
	"net/url"
	"reflect"
)

type Maplike interface {
	AsMapOfInterface() map[string]interface{}
}

// MapOfInterface attempts to convert the generic value into a map[string]interface{}
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfInterface(value interface{}) map[string]interface{} {
	result, _ := MapOfInterfaceOk(value)
	return result
}

// MapOfInterfaceOk attempts to convert the generic value into a map[string]interface{}
// The boolean result value returns TRUE if successful.  FALSE otherwise
func MapOfInterfaceOk(value interface{}) (map[string]interface{}, bool) {

	switch v := value.(type) {

	case map[string]interface{}:
		return v, true

	case map[string]string:
		result := make(map[string]interface{})

		for key, value := range v {
			result[key] = value
		}
		return result, true

	case url.Values:
		result := make(map[string]interface{})
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
		return MapOfInterfaceOk(Interface(v))

	case Maplike:
		return v.AsMapOfInterface(), true

	}

	// Fall through means conversion failed
	return make(map[string]interface{}), false
}
