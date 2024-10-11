package convert

import (
	"net/url"
)

func URLValues(value any) url.Values {
	result, _ := URLValuesOk(value)
	return result
}

func URLValuesOk(value any) (url.Values, bool) {

	switch typed := value.(type) {

	case url.Values:
		return typed, true

	case map[string]any:
		result := url.Values{}
		for key, value := range typed {
			result[key] = SliceOfString(value)
		}
		return result, true

	case Maplike:
		return URLValuesOk(typed.MapOfAny())
	}

	return url.Values{}, false
}
