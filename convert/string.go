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
// The second result is TRUE if the conversion was lossless (the converted value round-trips back to
// the original input), and FALSE otherwise.
//
// Conversion Rules:
// Nils return default value and Ok=false
// Bools are formatted as "true" or "false" losslessly, with Ok=true
// Ints are formatted as decimal strings losslessly, with Ok=true
// Floats are formatted with two decimal places; Ok=true only when that two-decimal
// string parses back to the original value, otherwise the rounding is lossy (Ok=false)
// Strings are passed through directly, with Ok=true
// A slice of length 1 carries the Ok of its single element; an empty or longer slice is lossy (Ok=false)
// Known interfaces (Inter, Floater, Stringer) are handled like their corresponding types.
// All other values return the default value with Ok=false
func StringOk(value any, defaultValue string) (string, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return defaultValue, false
	}

	switch v := value.(type) {

	case bool:

		if v {
			return "true", true
		}

		return "false", true

	case []byte:
		return string(v), true

	case int:
		return strconv.Itoa(v), true

	case int8:
		return strconv.FormatInt(int64(v), 10), true

	case int16:
		return strconv.FormatInt(int64(v), 10), true

	case int32:
		return strconv.FormatInt(int64(v), 10), true

	case int64:
		return strconv.FormatInt(v, 10), true

	case float32:
		return floatToString(float64(v))

	case float64:
		return floatToString(v)

	case string:
		return v, true

	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return StringOk(v[0], defaultValue)
		}
		return v[0], false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return StringOk(v[0], defaultValue)
		}
		return StringDefault(v[0], defaultValue), false

	case reflect.Value:
		return StringOk(Interface(v), defaultValue)

	case Booler:
		return StringOk(v.Bool(), defaultValue)

	case Inter:
		return strconv.FormatInt(int64(v.Int()), 10), true

	case Floater:
		return floatToString(v.Float())

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

// floatToString formats a float64 with two decimal places. The conversion is
// lossless (ok=true) only when the two-decimal string parses back to the original
// value; a value needing more than two decimals rounds and is reported lossy.
func floatToString(value float64) (string, bool) {
	result := strconv.FormatFloat(value, 'f', 2, 64)
	roundTrip, _ := strconv.ParseFloat(result, 64)
	return result, roundTrip == value
}

// JoinString converts the value into a string.
// If the value is a slice ([]string or []any), then the values are joined with the specified delimiter.
// If the value is not a slice, then it is converted to a string using the String() function.
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
