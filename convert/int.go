package convert

import (
	"io"
	"math"
	"reflect"
	"strconv"
)

// Int forces a conversion from an arbitrary value into an int.
// If the value cannot be converted, then the zero value for the type (0) is used.
func Int(value any) int {

	result, _ := IntOk(value, 0)
	return result
}

// IntDefault forces a conversion from an arbitrary value into a int.
// if the value cannot be converted, then the default value is used.
func IntDefault(value any, defaultValue int) int {

	result, _ := IntOk(value, defaultValue)
	return result
}

// IntOk converts an arbitrary value (passed in the first parameter) into an int, no matter what.
// The first result is the final converted value, or the default value (passed in the second parameter)
// The second result is TRUE if the value was naturally an integer, and FALSE otherwise
//
// Conversion Rules:
// Nils and Bools return default value and Ok=false
// Ints are returned directly with Ok=true
// Floats are truncated into ints.  If there is no decimal value then Ok=true
// String values are attempted to parse as a int.  If unsuccessful, default value is returned.  For all strings, Ok=false
// Known interfaces (Inter, Floater, Stringer) are handled like their corresponding types.
// All other values return the default value with Ok=false
func IntOk(value any, defaultValue int) (int, bool) {

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
		return int(v), true

	case int8:
		return int(v), true

	case int16:
		return int(v), true

	case int32:
		return int(v), true

	case int64:
		if v > math.MaxInt {
			return math.MaxInt, false
		}

		if v < math.MinInt {
			return math.MinInt, false
		}

		return int(v), true

	case float32:
		if v > float32(math.MaxInt) {
			return math.MaxInt, false
		}

		if v < float32(math.MinInt) {
			return math.MinInt, false
		}

		return int(v), hasDecimal(float64(v))

	case float64:
		if v > float64(math.MaxInt) {
			return math.MaxInt, false
		}

		if v < float64(math.MinInt) {
			return math.MinInt, false
		}

		return int(v), hasDecimal(v)

	case string:
		result, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return defaultValue, false
		}

		return IntOk(result, defaultValue)

		// []string is useful for parsing url.Values data
	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		return IntDefault(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		return IntDefault(v[0], defaultValue), false

	case reflect.Value:
		return IntOk(Interface(v), defaultValue)

	// Use standard interfaces, if available
	case Inter:
		return v.Int(), true

	case Floater:
		result := v.Float()
		return Int(result), hasDecimal(result)

	case Hexer:
		result, err := strconv.ParseInt(v.Hex(), 16, 64)
		if err != nil {
			return 0, false
		}

		return IntOk(result, defaultValue)

	case Stringer:
		return IntOk(v.String(), defaultValue)

	case io.Reader:
		return IntOk(String(v), defaultValue)
	}

	return defaultValue, false
}
