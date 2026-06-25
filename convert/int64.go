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

// Int64Ok converts an arbitrary value (passed in the first parameter) into an int64, no matter what.
// The first result is the final converted value, or the default value (passed in the second parameter)
// The second result is TRUE if the conversion was lossless (the converted value round-trips back to
// the original input), and FALSE otherwise.
//
// Conversion Rules:
// Nils return default value and Ok=false
// Bools map losslessly to 0/1 with Ok=true
// Int64s are returned directly with Ok=true
// Floats with no fractional part convert with Ok=true; a fractional part is lossy (Ok=false)
// Out-of-range values are clamped to the int64 bounds and reported as lossy (Ok=false)
// String values are parsed as an int; a clean parse is lossless (Ok=true), otherwise the
// default value is returned with Ok=false
// A slice of length 1 carries the Ok of its single element; an empty or longer slice is lossy (Ok=false)
// Known interfaces (Inter, Floater, Stringer) are handled like their corresponding types.
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
			return defaultValue, false
		}

		return Int64Ok(result, defaultValue)

	case Stringer:
		return Int64Ok(v.String(), defaultValue)

	case io.Reader:
		return Int64Ok(String(v), defaultValue)
	}

	return defaultValue, false
}
