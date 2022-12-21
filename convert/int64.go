package convert

import (
	"io"
	"reflect"
	"strconv"
)

// Int64 forces a conversion from an arbitrary value into an int.
// If the value cannot be converted, then the zero value for the type (0) is used.
func Int64(value any) int64 {

	result, _ := Int64Ok(value, 0)
	return result
}

// Int64Default forces a conversion from an arbitrary value into a int.
// if the value cannot be converted, then the default value is used.
func Int64Default(value any, defaultValue int64) int64 {

	result, _ := Int64Ok(value, defaultValue)
	return result
}

// Int64Ok converts an arbitrary value (passed in the first parameter) into an int, no matter what.
// The first result is the final converted value, or the default value (passed in the second parameter)
// The second result is TRUE if the value was naturally an integer, and FALSE otherwise
//
// Conversion Rules:
// Nils and Bools return default value and Ok=false
// Int64s are returned directly with Ok=true
// Floats are truncated into ints.  If there is no decimal value then Ok=true
// String values are attempted to parse as a int.  If unsuccessful, default value is returned.  For all strings, Ok=false
// Known interfaces (Int64er, Floater, Stringer) are handled like their corresponding types.
// All other values return the default value with Ok=false
func Int64Ok(value any, defaultValue int64) (int64, bool) {

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
		return int64(v), true

	case int8:
		return int64(v), true

	case int16:
		return int64(v), true

	case int32:
		return int64(v), true

	case int64:
		return int64(v), true

	case float32:
		return int64(v), hasDecimal(float64(v))

	case float64:
		return int64(v), hasDecimal(v)

	case string:
		result, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return defaultValue, false
		}

		return result, true

		// []string is useful for parsing url.Values data
	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		return Int64Default(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		return Int64Default(v[0], defaultValue), false

	case reflect.Value:
		return Int64Ok(Interface(v), defaultValue)

		// Use standard interfaces, if available
	case Inter:
		return int64(v.Int()), true

	case Floater:
		result := v.Float()
		return int64(result), hasDecimal(result)

	case Hexer:
		result, err := strconv.ParseInt(v.Hex(), 16, 64)
		if err != nil {
			return 0, false
		}

		return Int64Ok(result, defaultValue)

	case Stringer:
		return Int64Ok(v.String(), defaultValue)

	case io.Reader:
		return Int64Ok(String(v), defaultValue)
	}

	return defaultValue, false
}
