package convert

import (
	"io"
	"math"
	"reflect"
	"strconv"
)

// Float boundaries for safe int32 conversion. float32(math.MaxInt32) rounds up to 2^31, so we
// compare against the exact 2^31 power-of-two boundary instead: any float >= 2^31 or < -2^31 is
// out of int32 range. Both boundaries are exactly representable in float32 and float64.
const (
	maxInt32AsFloat = float64(1 << 31)  // 2^31, the first float above math.MaxInt32
	minInt32AsFloat = -float64(1 << 31) // -2^31, exactly math.MinInt32
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

// Int32Ok converts an arbitrary value (passed in the first parameter) into an int32, no matter what.
// The first result is the final converted value, or the default value (passed in the second parameter)
// The second result is TRUE if the conversion was lossless (the converted value round-trips back to
// the original input), and FALSE otherwise.
//
// Conversion Rules:
// Nils return default value and Ok=false
// Bools map losslessly to 0/1 with Ok=true
// Int32s are returned directly with Ok=true
// Floats with no fractional part convert with Ok=true; a fractional part is lossy (Ok=false)
// Out-of-range values are clamped to the int32 bounds and reported as lossy (Ok=false)
// String values are parsed as an int; a clean parse is lossless (Ok=true), otherwise the
// default value is returned with Ok=false
// A slice of length 1 carries the Ok of its single element; an empty or longer slice is lossy (Ok=false)
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
		// int is 64-bit on most platforms, so it must be range-checked before narrowing to int32.
		if int64(v) > math.MaxInt32 {
			return math.MaxInt32, false
		}
		if int64(v) < math.MinInt32 {
			return math.MinInt32, false
		}
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
		// float32(math.MaxInt32) rounds UP to 2^31, so a plain `> MaxInt32` lets 2^31 slip through
		// and overflow. Compare against the exact 2^31 boundary with >= instead.
		if float64(v) >= maxInt32AsFloat {
			return math.MaxInt32, false
		}

		if float64(v) < minInt32AsFloat {
			return math.MinInt32, false
		}

		return int32(v), hasDecimal(float64(v))

	case float64:
		if v >= maxInt32AsFloat {
			return math.MaxInt32, false
		}

		if v < minInt32AsFloat {
			return math.MinInt32, false
		}

		return int32(v), hasDecimal(v)

	case string:
		result, err := strconv.ParseInt(v, 10, 32)

		if err != nil {
			return defaultValue, false
		}

		return Int32Ok(result, defaultValue)

		// []string is useful for parsing url.Values data
	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return Int32Ok(v[0], defaultValue)
		}
		return Int32Default(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return Int32Ok(v[0], defaultValue)
		}
		return Int32Default(v[0], defaultValue), false

	case reflect.Value:
		return Int32Ok(Interface(v), defaultValue)

	// Use standard interfaces, if available
	case Inter:
		return Int32Ok(v.Int(), defaultValue)

	case Floater:
		return Int32Ok(v.Float(), defaultValue)

	case Hexer:
		result, err := strconv.ParseInt(v.Hex(), 16, 32)
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
