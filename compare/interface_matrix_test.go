package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestInterface_NumericMatrix exhaustively compares every numeric type against every
// other numeric type, asserting less-than, equal, and greater-than in all directions.
// The "ordered" closure asserts a < b (and, by symmetry, b > a); "equal" asserts a == b.
func TestInterface_NumericMatrix(t *testing.T) {

	// ordered asserts that low < high and high > low.
	ordered := func(low any, high any) {
		t.Helper()

		result, err := Interface(low, high)
		require.Nil(t, err)
		require.Equal(t, -1, result, "expected %#v < %#v", low, high)

		result, err = Interface(high, low)
		require.Nil(t, err)
		require.Equal(t, 1, result, "expected %#v > %#v", high, low)
	}

	// equal asserts that two (possibly different-typed) values compare as equal.
	equal := func(a any, b any) {
		t.Helper()

		result, err := Interface(a, b)
		require.Nil(t, err)
		require.Equal(t, 0, result, "expected %#v == %#v", a, b)
	}

	// --- int(1) vs every type ---
	ordered(int(1), int(2))
	ordered(int(1), int8(2))
	ordered(int(1), int16(2))
	ordered(int(1), int32(2))
	ordered(int(1), int64(2))
	ordered(int(1), uint(2))
	ordered(int(1), uint8(2))
	ordered(int(1), uint16(2))
	ordered(int(1), uint32(2))
	ordered(int(1), uint64(2))
	ordered(int(1), float32(2))
	ordered(int(1), float64(2))

	// --- int8(1) vs every type ---
	ordered(int8(1), int(2))
	ordered(int8(1), int8(2))
	ordered(int8(1), int16(2))
	ordered(int8(1), int32(2))
	ordered(int8(1), int64(2))
	ordered(int8(1), uint(2))
	ordered(int8(1), uint8(2))
	ordered(int8(1), uint16(2))
	ordered(int8(1), uint32(2))
	ordered(int8(1), uint64(2))
	ordered(int8(1), float32(2))
	ordered(int8(1), float64(2))

	// --- int16(1) vs every type ---
	ordered(int16(1), int(2))
	ordered(int16(1), int8(2))
	ordered(int16(1), int16(2))
	ordered(int16(1), int32(2))
	ordered(int16(1), int64(2))
	ordered(int16(1), uint(2))
	ordered(int16(1), uint8(2))
	ordered(int16(1), uint16(2))
	ordered(int16(1), uint32(2))
	ordered(int16(1), uint64(2))
	ordered(int16(1), float32(2))
	ordered(int16(1), float64(2))

	// --- int32(1) vs every type ---
	ordered(int32(1), int(2))
	ordered(int32(1), int8(2))
	ordered(int32(1), int16(2))
	ordered(int32(1), int32(2))
	ordered(int32(1), int64(2))
	ordered(int32(1), uint(2))
	ordered(int32(1), uint8(2))
	ordered(int32(1), uint16(2))
	ordered(int32(1), uint32(2))
	ordered(int32(1), uint64(2))
	ordered(int32(1), float32(2))
	ordered(int32(1), float64(2))

	// --- int64(1) vs every type ---
	ordered(int64(1), int(2))
	ordered(int64(1), int8(2))
	ordered(int64(1), int16(2))
	ordered(int64(1), int32(2))
	ordered(int64(1), int64(2))
	ordered(int64(1), uint(2))
	ordered(int64(1), uint8(2))
	ordered(int64(1), uint16(2))
	ordered(int64(1), uint32(2))
	ordered(int64(1), uint64(2))
	ordered(int64(1), float32(2))
	ordered(int64(1), float64(2))

	// --- uint(1) vs every type ---
	ordered(uint(1), int(2))
	ordered(uint(1), int8(2))
	ordered(uint(1), int16(2))
	ordered(uint(1), int32(2))
	ordered(uint(1), int64(2))
	ordered(uint(1), uint(2))
	ordered(uint(1), uint8(2))
	ordered(uint(1), uint16(2))
	ordered(uint(1), uint32(2))
	ordered(uint(1), uint64(2))
	ordered(uint(1), float32(2))
	ordered(uint(1), float64(2))

	// --- uint8(1) vs every type ---
	ordered(uint8(1), int(2))
	ordered(uint8(1), int8(2))
	ordered(uint8(1), int16(2))
	ordered(uint8(1), int32(2))
	ordered(uint8(1), int64(2))
	ordered(uint8(1), uint(2))
	ordered(uint8(1), uint8(2))
	ordered(uint8(1), uint16(2))
	ordered(uint8(1), uint32(2))
	ordered(uint8(1), uint64(2))
	ordered(uint8(1), float32(2))
	ordered(uint8(1), float64(2))

	// --- uint16(1) vs every type ---
	ordered(uint16(1), int(2))
	ordered(uint16(1), int8(2))
	ordered(uint16(1), int16(2))
	ordered(uint16(1), int32(2))
	ordered(uint16(1), int64(2))
	ordered(uint16(1), uint(2))
	ordered(uint16(1), uint8(2))
	ordered(uint16(1), uint16(2))
	ordered(uint16(1), uint32(2))
	ordered(uint16(1), uint64(2))
	ordered(uint16(1), float32(2))
	ordered(uint16(1), float64(2))

	// --- uint32(1) vs every type ---
	ordered(uint32(1), int(2))
	ordered(uint32(1), int8(2))
	ordered(uint32(1), int16(2))
	ordered(uint32(1), int32(2))
	ordered(uint32(1), int64(2))
	ordered(uint32(1), uint(2))
	ordered(uint32(1), uint8(2))
	ordered(uint32(1), uint16(2))
	ordered(uint32(1), uint32(2))
	ordered(uint32(1), uint64(2))
	ordered(uint32(1), float32(2))
	ordered(uint32(1), float64(2))

	// --- uint64(1) vs every type ---
	ordered(uint64(1), int(2))
	ordered(uint64(1), int8(2))
	ordered(uint64(1), int16(2))
	ordered(uint64(1), int32(2))
	ordered(uint64(1), int64(2))
	ordered(uint64(1), uint(2))
	ordered(uint64(1), uint8(2))
	ordered(uint64(1), uint16(2))
	ordered(uint64(1), uint32(2))
	ordered(uint64(1), uint64(2))
	ordered(uint64(1), float32(2))
	ordered(uint64(1), float64(2))

	// --- float32(1) vs every type ---
	ordered(float32(1), int(2))
	ordered(float32(1), int8(2))
	ordered(float32(1), int16(2))
	ordered(float32(1), int32(2))
	ordered(float32(1), int64(2))
	ordered(float32(1), uint(2))
	ordered(float32(1), uint8(2))
	ordered(float32(1), uint16(2))
	ordered(float32(1), uint32(2))
	ordered(float32(1), uint64(2))
	ordered(float32(1), float32(2))
	ordered(float32(1), float64(2))

	// --- float64(1) vs every type ---
	ordered(float64(1), int(2))
	ordered(float64(1), int8(2))
	ordered(float64(1), int16(2))
	ordered(float64(1), int32(2))
	ordered(float64(1), int64(2))
	ordered(float64(1), uint(2))
	ordered(float64(1), uint8(2))
	ordered(float64(1), uint16(2))
	ordered(float64(1), uint32(2))
	ordered(float64(1), uint64(2))
	ordered(float64(1), float32(2))
	ordered(float64(1), float64(2))

	// --- equal magnitudes across mixed types ---
	equal(int(1), uint64(1))
	equal(uint(1), int64(1))
	equal(int8(1), float64(1))
	equal(uint8(1), float32(1))
	equal(int64(0), uint64(0))
	equal(float32(1), float64(1))
}

// TestInterface_NegativeAcrossSignedTypes confirms that a negative value of every signed
// type compares as less than every unsigned type, in both directions, with no wrap to a
// large uint64. This is the regression guard for the signed/unsigned comparison bug.
func TestInterface_NegativeAcrossSignedTypes(t *testing.T) {

	// lessThanUnsigned asserts that the (negative) signed value is less than the
	// unsigned value, regardless of operand order.
	lessThanUnsigned := func(negative any, unsigned any) {
		t.Helper()

		result, err := Interface(negative, unsigned)
		require.Nil(t, err)
		require.Equal(t, -1, result, "negative %#v must be < unsigned %#v", negative, unsigned)

		result, err = Interface(unsigned, negative)
		require.Nil(t, err)
		require.Equal(t, 1, result, "unsigned %#v must be > negative %#v", unsigned, negative)
	}

	// Every signed type at -1 against every unsigned type at 0.
	lessThanUnsigned(int(-1), uint(0))
	lessThanUnsigned(int(-1), uint8(0))
	lessThanUnsigned(int(-1), uint16(0))
	lessThanUnsigned(int(-1), uint32(0))
	lessThanUnsigned(int(-1), uint64(0))

	lessThanUnsigned(int8(-1), uint(0))
	lessThanUnsigned(int8(-1), uint8(0))
	lessThanUnsigned(int8(-1), uint16(0))
	lessThanUnsigned(int8(-1), uint32(0))
	lessThanUnsigned(int8(-1), uint64(0))

	lessThanUnsigned(int16(-1), uint(0))
	lessThanUnsigned(int16(-1), uint8(0))
	lessThanUnsigned(int16(-1), uint16(0))
	lessThanUnsigned(int16(-1), uint32(0))
	lessThanUnsigned(int16(-1), uint64(0))

	lessThanUnsigned(int32(-1), uint(0))
	lessThanUnsigned(int32(-1), uint8(0))
	lessThanUnsigned(int32(-1), uint16(0))
	lessThanUnsigned(int32(-1), uint32(0))
	lessThanUnsigned(int32(-1), uint64(0))

	lessThanUnsigned(int64(-1), uint(0))
	lessThanUnsigned(int64(-1), uint8(0))
	lessThanUnsigned(int64(-1), uint16(0))
	lessThanUnsigned(int64(-1), uint32(0))
	lessThanUnsigned(int64(-1), uint64(0))
}

// TestInterface_Boundaries checks the extreme values that the signed/unsigned widening
// logic is most likely to get wrong: max uint64 (which overflows int64) and the signed
// 64-bit limits. The unsigned-vs-signed path must stay exact, not detour through float.
func TestInterface_Boundaries(t *testing.T) {

	const maxUint64 = uint64(1<<64 - 1)
	const maxInt64 = int64(1<<63 - 1)
	const minInt64 = int64(-1 << 63)

	ordered := func(low any, high any) {
		t.Helper()

		result, err := Interface(low, high)
		require.Nil(t, err)
		require.Equal(t, -1, result, "expected %#v < %#v", low, high)

		result, err = Interface(high, low)
		require.Nil(t, err)
		require.Equal(t, 1, result, "expected %#v > %#v", high, low)
	}

	// max uint64 is larger than the largest signed value, preserved exactly.
	ordered(maxInt64, maxUint64)

	// min int64 is less than zero unsigned.
	ordered(minInt64, uint64(0))

	// Equal magnitudes spanning the signed/unsigned boundary still compare equal.
	result, err := Interface(maxInt64, uint64(maxInt64))
	require.Nil(t, err)
	require.Equal(t, 0, result, "maxInt64 == uint64(maxInt64)")
}

// TestInterface_FloatFraction confirms fractional floats compare correctly against
// integers, which are widened to float64 (matching the documented behavior).
func TestInterface_FloatFraction(t *testing.T) {

	ordered := func(low any, high any) {
		t.Helper()

		result, err := Interface(low, high)
		require.Nil(t, err)
		require.Equal(t, -1, result, "expected %#v < %#v", low, high)
	}

	ordered(int(1), 1.5)
	ordered(1.5, int(2))
	ordered(float32(1.5), float64(1.6))

	result, err := Interface(float32(1.5), float64(1.5))
	require.Nil(t, err)
	require.Equal(t, 0, result, "float32(1.5) == float64(1.5)")
}

// TestInterface_NumericVsNonNumeric confirms that pairing a number with a string, bool,
// or unknown type is reported as incompatible rather than silently coerced.
func TestInterface_NumericVsNonNumeric(t *testing.T) {

	incompatible := func(value1 any, value2 any) {
		t.Helper()

		_, err := Interface(value1, value2)
		require.NotNil(t, err, "expected %#v vs %#v to be incompatible", value1, value2)
	}

	incompatible(1, "1")
	incompatible("1", 1)
	incompatible(1, true) // a numeric value1 only compares against another number
	incompatible(float64(1), struct{}{})

	// Note: Interface(true, 1) is NOT incompatible. When value1 is a bool, value2 is
	// coerced to bool, so a bool-vs-number comparison succeeds in that direction only.
	result, err := Interface(true, 1)
	require.Nil(t, err)
	require.Equal(t, 0, result, "bool(true) coerces 1 to true, so they compare equal")
}
