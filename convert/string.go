package convert

import (
	"bytes"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// String forces a conversion from an arbitrary value into an string.
// If the value cannot be converted, then the default value for the type is used.
func String(value any) string {
	result, _ := StringOk(value, "")
	return result
}

// StringDefault forces a conversion from an arbitrary value into a string.
// if the value cannot be converted, then the default value is used.
func StringDefault(value any, defaultValue string) string {
	result, _ := StringOk(value, defaultValue)
	return result
}

// StringOk converts an arbitrary value (passed in the first parameter) into a string, no matter what.
// The first result is the final converted value, or the default value (passed in the second parameter)
// The second result is TRUE if the value was naturally a string, and FALSE otherwise
//
// Conversion Rules:
// Nils return default value and Ok=false
// Bools are formated as "true" or "false" with Ok=false
// Ints are formated as strings with Ok=false
// Floats are formatted with 2 decimal places, with Ok=false
// String are passed through directly, with Ok=true
// Known interfaces (Inter, Floater, Stringer) are handled like their corresponding types.
// All other values return the default value with Ok=false
func StringOk(value any, defaultValue string) (string, bool) {

	if value == nil {
		return defaultValue, false
	}

	switch v := value.(type) {

	case bool:

		if v {
			return "true", false
		}

		return "false", false

	case []byte:
		return string(v), true

	case int:
		return strconv.Itoa(v), false

	case int8:
		return strconv.FormatInt(int64(v), 10), false

	case int16:
		return strconv.FormatInt(int64(v), 10), false

	case int32:
		return strconv.FormatInt(int64(v), 10), false

	case int64:
		return strconv.FormatInt(v, 10), false

	case float32:
		return strconv.FormatFloat(float64(v), 'f', -2, 64), false

	case float64:
		return strconv.FormatFloat(v, 'f', -2, 64), false

	case string:
		return v, true

	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		return v[0], false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		return StringDefault(v[0], defaultValue), false

	case reflect.Value:
		return StringOk(Interface(v), defaultValue)

	case Booler:
		return StringOk(v.Bool(), defaultValue)

	case Inter:
		return strconv.FormatInt(int64(v.Int()), 10), false

	case Floater:
		return strconv.FormatFloat(v.Float(), 'f', -2, 64), false

	case Hexer:
		return v.Hex(), true

	case Stringer:
		return v.String(), true

	case io.Reader:
		var buffer bytes.Buffer

		if _, err := io.Copy(&buffer, v); err != nil {
			return "", false
		}

		return buffer.String(), true
	}

	return defaultValue, false
}

func JoinString(value any, delimiter string) string {

	if delimiter != "" {

		switch value := value.(type) {

		case []string:
			return strings.Join(value, delimiter)

		case []any:
			result := make([]string, len(value))
			for index, v := range value {
				result[index] = String(v)
			}
			return strings.Join(result, delimiter)

		case reflect.Value:
			return JoinString(Interface(value), delimiter)
		}
	}

	return String(value)
}
