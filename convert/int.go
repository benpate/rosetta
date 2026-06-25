package convert

import (
	"io"
	"math"
	"reflect"
	"strconv"
)

// Float boundaries for safe int conversion. float(math.MaxInt) rounds up to the next power of two,
// so we compare against the exact 2^(N-1) boundary instead: any float >= that or < its negation is
// out of int range. The boundary depends on the platform int width (32- or 64-bit).
var (
	maxIntAsFloat = boundaryIf(math.MaxInt == math.MaxInt32, maxInt32AsFloat, maxInt64AsFloat)
	minIntAsFloat = boundaryIf(math.MaxInt == math.MaxInt32, minInt32AsFloat, minInt64AsFloat)
)

// boundaryIf selects between two float boundaries based on the platform int width.
func boundaryIf(is32 bool, on32, on64 float64) float64 {
	if is32 {
		return on32
	}
	return on64
}

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

// IntOk converts an arbitrary value (passed in the first parameter) into an int,
// no matter what. The first result is the final converted value, or the default
// value (passed in the second parameter). The second result is TRUE if the
// conversion was lossless (the converted value round-trips back to the original
// input), and FALSE otherwise.
//
// Conversion Rules:
// Nils return default value and Ok=false
// Bools map losslessly to 0/1 with Ok=true
// Ints are returned directly with Ok=true
// Floats with no fractional part convert with Ok=true; a fractional part is lossy (Ok=false)
// Out-of-range values are clamped to the int bounds and reported as lossy (Ok=false)
// String values are parsed as an int; a clean parse is lossless (Ok=true), otherwise the
// default value is returned with Ok=false
// A slice of length 1 carries the Ok of its single element; an empty or longer slice is lossy (Ok=false)
// Known interfaces (Inter, Floater, Stringer) are handled like their types.
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
		// float(math.MaxInt) rounds UP to the next power of two, so a plain `> MaxInt` lets that
		// boundary value slip through and overflow. Compare against the exact 2^(N-1) boundary with >=.
		if float64(v) >= maxIntAsFloat {
			return math.MaxInt, false
		}

		if float64(v) < minIntAsFloat {
			return math.MinInt, false
		}

		return int(v), hasDecimal(float64(v))

	case float64:
		if v >= maxIntAsFloat {
			return math.MaxInt, false
		}

		if v < minIntAsFloat {
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
		if len(v) == 1 {
			return IntOk(v[0], defaultValue)
		}
		return IntDefault(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return IntOk(v[0], defaultValue)
		}
		return IntDefault(v[0], defaultValue), false

	case reflect.Value:
		return IntOk(Interface(v), defaultValue)

	// Use standard interfaces, if available
	case Inter:
		return v.Int(), true

	case Floater:
		// Delegate to the float64 case so out-of-range values are clamped and reported as Ok=false.
		return IntOk(v.Float(), defaultValue)

	case Hexer:
		result, err := strconv.ParseInt(v.Hex(), 16, 64)
		if err != nil {
			return defaultValue, false
		}

		return IntOk(result, defaultValue)

	case Stringer:
		return IntOk(v.String(), defaultValue)

	case io.Reader:
		return IntOk(String(v), defaultValue)
	}

	return defaultValue, false
}
