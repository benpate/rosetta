package convert

import (
	"net/http"
	"net/url"
	"reflect"
)

// MapOfSliceOfString converts the given value to a map[string][]string.
// If conversion is not possible, then an empty map is returned.
func MapOfSliceOfString(value any) map[string][]string {
	result, _ := MapOfSliceOfStringOk(value)
	return result
}

// MapOfSliceOfStringOk converts the given value to a map[string][]string.
// It returns TRUE if the conversion was successful.  If conversion is not
// possible, then it returns an empty map and FALSE.
func MapOfSliceOfStringOk(value any) (map[string][]string, bool) {

	switch typed := value.(type) {

	case map[string]any:
		result := make(map[string][]string, len(typed))

		for key, value := range typed {
			result[key] = SliceOfString(value)
		}

		return result, true

	case map[string]string:
		result := make(map[string][]string, len(typed))

		for key, value := range typed {
			result[key] = SliceOfString(value)
		}

		return result, true

	case map[string][]string:
		return typed, true

	case url.Values:
		return map[string][]string(typed), true

	case http.Header:
		return map[string][]string(typed), true

	case reflect.Value:
		return MapOfSliceOfStringOk(Interface(typed))

	case MapOfAnyGetter:
		return MapOfSliceOfStringOk(typed.MapOfAny())

	}

	// Fall through means conversion failed
	return map[string][]string{}, false
}

// URLValues converts a data structure into a url.Values object,
// which is a specialized instance of a map[string][]string.
func URLValues(value any) url.Values {
	return url.Values(MapOfSliceOfString(value))
}

// URLValues converts a data structure into a url.Values object,
// which is a specialized instance of a map[string][]string.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func URLValuesOk(value any) (url.Values, bool) {
	result, ok := MapOfSliceOfStringOk(value)
	return url.Values(result), ok
}

// HTTPHeader converts a data structure into a http.Header object,
// which is a specialized instance of a map[string][]string.
func HTTPHeader(value any) http.Header {
	return http.Header(MapOfSliceOfString(value))
}

// HTTPHeader converts a data structure into a http.Header object,
// which is a specialized instance of a map[string][]string.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func HTTPHeaderOk(value any) (http.Header, bool) {
	result, ok := MapOfSliceOfStringOk(value)
	return http.Header(result), ok
}
