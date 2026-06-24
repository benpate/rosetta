package convert

import (
	"io"
	"reflect"
)

// Bool forces a conversion from an arbitrary value into a boolean.
// If the value cannot be converted, then the default value for the type is used.
func Bool(value any) bool {

	result, _ := BoolOk(value, false)
	return result
}

// BoolDefault forces a conversion from an arbitrary value into a bool.
// if the value cannot be converted, then the default value is used.
func BoolDefault(value any, defaultValue bool) bool {

	result, _ := BoolOk(value, defaultValue)
	return result
}

// BoolOk converts an arbitrary value (passed in the first parameter) into a boolean, somehow,
// no matter what. The first result is the final converted value, or the default value
// (passed in the second parameter)
// The second result is TRUE if the conversion was lossless (the converted value round-trips
// back to the original input), and FALSE otherwise.
//
// Conversion Rules:
// Nils return default value and Ok=false
// Bools are passed through with Ok=true
// Ints and Floats of exactly 0 or 1 map losslessly to false/true with Ok=true; any other
// numeric value is lossy and returns Ok=false
// String values of "true" and "false" convert losslessly with Ok=true
// All other strings return the default value, with Ok=false
// A slice of length 1 carries the Ok of its single element; an empty or longer slice is lossy (Ok=false)
// Known interfaces (Booler, Inter, Floater, Stringer) are handled like their corresponding types
// All other values return the default value with Ok=false
func BoolOk(value any, defaultValue bool) (bool, bool) {

	if value == nil {
		return defaultValue, false
	}

	switch v := value.(type) {

	case bool:
		return v, true

	case int:
		return boolFromNumber(v)

	case int8:
		return boolFromNumber(v)

	case int16:
		return boolFromNumber(v)

	case int32:
		return boolFromNumber(v)

	case int64:
		return boolFromNumber(v)

	case float32:
		return boolFromNumber(v)

	case float64:
		return boolFromNumber(v)

	case string:

		switch v {
		case "true":
			return true, true
		case "false":
			return false, true
		default:
			return defaultValue, false
		}

	// []string is useful for parsing url.Values data
	case []string:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return BoolOk(v[0], defaultValue)
		}
		return BoolDefault(v[0], defaultValue), false

	case []any:
		if len(v) == 0 {
			return defaultValue, false
		}
		if len(v) == 1 {
			return BoolOk(v[0], defaultValue)
		}
		return BoolDefault(v[0], defaultValue), false

	case reflect.Value:
		return BoolOk(Interface(v), defaultValue)

	// Use standard interfaces, if available
	case Booler:
		return v.Bool(), true

	case Inter:
		return BoolOk(v.Int(), defaultValue)

	case Floater:
		return BoolOk(v.Float(), defaultValue)

	case Stringer:
		return BoolOk(v.String(), defaultValue)

	case io.Reader:
		return BoolOk(String(v), defaultValue)

	}

	return defaultValue, false
}

// boolFromNumber reports a numeric-to-bool conversion. A numeric value of exactly
// 0 or 1 maps losslessly (ok=true); any other value is lossy (ok=false) and yields
// FALSE, since it cannot round-trip back to the original number.
func boolFromNumber[T int | int8 | int16 | int32 | int64 | float32 | float64](value T) (bool, bool) {

	result := (value == 1)
	lossless := (value == 0 || value == 1)

	return result, lossless
}
