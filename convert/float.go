package convert

import (
	"io"
	"math"
	"reflect"
	"strconv"
)

// Float forces a conversion from an arbitrary value into a float64.
// If the value cannot be converted, then the zero value for the type (false) is used.
func Float(value any) float64 {
	result, _ := FloatOk(value, 0)
	return result
}

// FloatDefault forces a conversion from an arbitrary value into a float64.
// if the value cannot be converted, then the default value is used.
func FloatDefault(value any, defaultValue float64) float64 {
	result, _ := FloatOk(value, defaultValue)
	return result
}

// FloatOk converts an arbitrary value (passed in the first parameter) into a float64, no matter what.
// The first result is the final converted value, or the default value (passed in the second parameter)
// The second result is TRUE if the value was naturally a floating point number, and FALSE otherwise
//
// Conversion Rules:
// Nils and Bools return default value and Ok=false
// Ints and Floats are converted into float64, with Ok=true
// String values are attempted to parse as a float64.  If unsuccessful, default value is returned.  For all strings, Ok=false
// Known interfaces (Inter, Floater, Stringer) are handled like their corresponding types.
// All other values return the default value with Ok=false
func FloatOk(value any, defaultValue float64) (float64, bool) {

	if value == nil {
		return defaultValue, false
	}

	switch v := value.(type) {

	case bool:
		if v {
			return 1, true
		}
		return 0, true

	case int:
		return float64(v), true

	case int8:
		return float64(v), true

	case int16:
		return float64(v), true

	case int32:
		return float64(v), true

	case int64:
		return float64(v), true

	case float32:
		return float64(v), true

	case float64:
		return v, true

	case string:
		result, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return defaultValue, false
		}

		return result, true

		// []string is useful for parsing url.Values data
	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		return FloatDefault(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		return FloatDefault(v[0], defaultValue), false

	case reflect.Value:
		return FloatOk(Interface(v), defaultValue)

	// Use standard interfaces, if available
	case Inter:
		return float64(v.Int()), true

	case Floater:
		return v.Float(), true

	case Stringer:
		return FloatOk(v.String(), defaultValue)

	case io.Reader:
		return FloatOk(String(v), defaultValue)
	}

	return defaultValue, false
}

func hasDecimal(value float64) bool {

	return (value == math.Floor(value))
}
