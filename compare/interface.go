package compare

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// Interface tries its best to muscle value1 and value2 into compatible types so that they can be compared.
// If value1 is LESS THAN value2, it returns -1, nil
// If value1 is EQUAL TO value2, it returns 0, nil
// If value1 is GREATER THAN value2, it returns 1, nil
// If the two values are not compatible, then it returns 0, [DERP] with an explanation of the error.
// Numeric values are compared across every signed, unsigned, and floating-point combination; strings
// and Stringers are compared lexically. Other types (and mismatched non-numeric types) are incompatible.
func Interface(value1 any, value2 any) (int, error) {

	const location = "compare.Interface"
	const incompatibleTypes = "Incompatible data types"

	// A bool on either side coerces the other operand to bool, so the comparison
	// succeeds symmetrically (Interface(true, 1) and Interface(1, true) both work).
	// This runs before the numeric block so a bool is never treated as a number.
	if v1, ok := value1.(bool); ok {
		return Bool(v1, convert.Bool(value2)), nil
	}
	if v2, ok := value2.(bool); ok {
		return Bool(convert.Bool(value1), v2), nil
	}

	// Numbers are compared on a common numeric kind, so every signed/unsigned/float
	// combination is handled symmetrically without a per-type-pair switch.
	if n1, ok := toNumber(value1); ok {
		if n2, ok := toNumber(value2); ok {
			return compareNumbers(n1, n2), nil
		}
		return 0, derp.Internal(location, incompatibleTypes, value1, value2)
	}

	switch v1 := value1.(type) {

	case string:
		if v2, ok := value2.(string); ok {
			return String(v1, v2), nil
		}

	case Stringer:
		if v2, ok := value2.(Stringer); ok {
			return String(v1.String(), v2.String()), nil
		}
	}

	return 0, derp.Internal(location, incompatibleTypes, value1, value2)
}

// numberKind identifies which of number's three fields holds the comparable value.
type numberKind int

const (
	kindSigned   numberKind = iota // value is held in the signed field
	kindUnsigned                   // value is held in the unsigned field
	kindFloat                      // value is held in the float field
)

// number is a numeric value decoded into whichever representation keeps it exact:
// signed integers in signed, unsigned integers in unsigned, and floats in float.
type number struct {
	kind     numberKind
	signed   int64
	unsigned uint64
	float    float64
}

// toNumber decodes any of Go's built-in numeric types into a number, reporting FALSE
// for any value that is not numeric. The original width is widened (e.g. int8 to int64)
// but never reinterpreted, so signedness and magnitude are preserved exactly.
func toNumber(value any) (number, bool) {

	switch v := value.(type) {

	case int:
		return number{kind: kindSigned, signed: int64(v)}, true
	case int8:
		return number{kind: kindSigned, signed: int64(v)}, true
	case int16:
		return number{kind: kindSigned, signed: int64(v)}, true
	case int32:
		return number{kind: kindSigned, signed: int64(v)}, true
	case int64:
		return number{kind: kindSigned, signed: v}, true

	case uint:
		return number{kind: kindUnsigned, unsigned: uint64(v)}, true
	case uint8:
		return number{kind: kindUnsigned, unsigned: uint64(v)}, true
	case uint16:
		return number{kind: kindUnsigned, unsigned: uint64(v)}, true
	case uint32:
		return number{kind: kindUnsigned, unsigned: uint64(v)}, true
	case uint64:
		return number{kind: kindUnsigned, unsigned: v}, true

	case float32:
		return number{kind: kindFloat, float: float64(v)}, true
	case float64:
		return number{kind: kindFloat, float: v}, true
	}

	return number{}, false
}

// compareNumbers returns -1, 0, or 1 comparing two decoded numbers across any kind pairing.
func compareNumbers(a, b number) int {

	switch {

	// Same kind: compare directly with no conversion.
	case (a.kind == kindSigned) && (b.kind == kindSigned):
		return Int64(a.signed, b.signed)

	case (a.kind == kindUnsigned) && (b.kind == kindUnsigned):
		return UInt64(a.unsigned, b.unsigned)

	// Any float operand widens both sides to float64. This can lose precision for very
	// large int64/uint64 values, but it matches the long-standing behavior of this function.
	case (a.kind == kindFloat) || (b.kind == kindFloat):
		return Float64(asFloat(a), asFloat(b))

	// One signed, one unsigned: a negative signed value is always the smaller of the two;
	// otherwise the signed value is non-negative and fits in uint64, so compare as unsigned.
	// This avoids the wrap that a naive uint64(negativeSigned) conversion would cause.
	case a.kind == kindSigned:
		if a.signed < 0 {
			return -1
		}
		return UInt64(uint64(a.signed), b.unsigned)

	default:
		if b.signed < 0 {
			return 1
		}
		return UInt64(a.unsigned, uint64(b.signed))
	}
}

// asFloat returns a number's value as a float64, regardless of which kind it holds.
func asFloat(n number) float64 {

	switch n.kind {
	case kindSigned:
		return float64(n.signed)
	case kindUnsigned:
		return float64(n.unsigned)
	default:
		return n.float
	}
}
