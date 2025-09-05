package convert

import (
	"io"
	"math"
	"reflect"
	"strconv"
)

// Int32 forces a conversion from an arbitrary value into an int.
// If the value cannot be converted, then the zero value for the type (0) is used.
func Int32(value any) int32 {

	result, _ := Int32Ok(value, 0)
	return result
}

// Int32Default forces a conversion from an arbitrary value into a int.
// if the value cannot be converted, then the default value is used.
func Int32Default(value any, defaultValue int32) int32 {

	result, _ := Int32Ok(value, defaultValue)
	return result
}

// Int32Ok converts an arbitrary value (passed in the first parameter) into an int, no matter what.
// The first result is the final converted value, or the default value (passed in the second parameter)
// The second result is TRUE if the value was naturally an integer, and FALSE otherwise
//
// Conversion Rules:
// Nils and Bools return default value and Ok=false
// Int32s are returned directly with Ok=true
// Floats are truncated into ints.  If there is no decimal value then Ok=true
// String values are attempted to parse as a int.  If unsuccessful, default value is returned.  For all strings, Ok=false
// Known interfaces (Inter, Floater, Stringer) are handled like their corresponding types.
// All other values return the default value with Ok=false
func Int32Ok(value any, defaultValue int32) (int32, bool) {

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
		return int32(v), true

	case int8:
		return int32(v), true

	case int16:
		return int32(v), true

	case int32:
		return int32(v), true

	case int64:
		if v > math.MaxInt32 {
			return math.MaxInt32, false
		}

		if v < math.MinInt32 {
			return math.MinInt32, false
		}

		return int32(v), true

	case float32:
		if v > math.MaxInt32 {
			return math.MaxInt32, false
		}

		if v < math.MinInt32 {
			return math.MinInt32, false
		}

		return int32(v), hasDecimal(float64(v))

	case float64:
		if v > math.MaxInt32 {
			return math.MaxInt32, false
		}

		if v < math.MinInt32 {
			return math.MinInt32, false
		}

		return int32(v), hasDecimal(v)

	case string:
		result, err := strconv.Atoi(v)

		if err != nil {
			return defaultValue, false
		}

		return int32(result), true

		// []string is useful for parsing url.Values data
	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		return Int32Default(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		return Int32Default(v[0], defaultValue), false

	case reflect.Value:
		return Int32Ok(Interface(v), defaultValue)

	// Use standard interfaces, if available
	case Inter:
		return int32(v.Int()), true

	case Floater:
		return Int32Ok(v.Float(), defaultValue)

	case Hexer:
		result, err := strconv.ParseInt(v.Hex(), 16, 64)
		if err != nil {
			return 0, false
		}

		return Int32Ok(result, defaultValue)

	case Stringer:
		return Int32Ok(v.String(), defaultValue)

	case io.Reader:
		return Int32Ok(String(v), defaultValue)
	}

	return defaultValue, false
}
