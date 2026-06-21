package convert

import (
	"io"
	"math"
	"reflect"
	"strconv"
)

// Float boundaries for safe int64 conversion. math.MaxInt64 is not exactly representable as a
// float64 (it rounds up to 2^63), so we compare against the exact powers of two instead: any float
// >= 2^63 or < -2^63 is out of int64 range.
const (
	maxInt64AsFloat = float64(1 << 63)  // 2^63, the first float64 above math.MaxInt64
	minInt64AsFloat = -float64(1 << 63) // -2^63, exactly math.MinInt64
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
		// math.MaxInt64 rounds UP to 2^63 as a float, so a plain `> MaxInt64` lets 2^63 slip
		// through and overflow. Compare against the 2^63 power-of-two boundary with >= instead.
		if float64(v) >= maxInt64AsFloat {
			return math.MaxInt64, false
		}

		if float64(v) < minInt64AsFloat {
			return math.MinInt64, false
		}

		return int64(v), hasDecimal(float64(v))

	case float64:
		if v >= maxInt64AsFloat {
			return math.MaxInt64, false
		}

		if v < minInt64AsFloat {
			return math.MinInt64, false
		}

		return int64(v), hasDecimal(v)

	case string:
		if result, err := strconv.ParseInt(v, 10, 64); err == nil {
			return result, true
		}

		return defaultValue, false

	// []string is useful for parsing url.Values data
	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return Int64Ok(v[0], defaultValue)
		}
		return Int64Default(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return Int64Ok(v[0], defaultValue)
		}
		return Int64Default(v[0], defaultValue), false

	case reflect.Value:
		return Int64Ok(Interface(v), defaultValue)

	// Use standard interfaces, if available
	case Inter:
		return Int64Ok(v.Int(), defaultValue)

	case Floater:
		return Int64Ok(v.Float(), defaultValue)

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
