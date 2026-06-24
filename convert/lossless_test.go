package convert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// The tests in this file pin the "lossless round-trip" contract for the *Ok
// converters: the boolean result is TRUE if and only if the converted value
// maps back to the original input with no loss of information.

// TestLossless_BoolOk covers numeric, string, and slice inputs to BoolOk.
func TestLossless_BoolOk(t *testing.T) {

	// value, defaultValue -> expected result, expected lossless
	test := func(value any, def bool, wantResult bool, wantLossless bool) {
		result, lossless := BoolOk(value, def)
		require.Equal(t, wantResult, result, "result for %#v", value)
		require.Equal(t, wantLossless, lossless, "lossless for %#v", value)
	}

	// Native bool is always lossless.
	test(true, false, true, true)
	test(false, true, false, true)

	// Exactly 0 or 1 round-trips through bool; any other number is lossy.
	test(0, true, false, true)
	test(1, false, true, true)
	test(int8(1), false, true, true)
	test(int64(0), true, false, true)
	test(float32(1), false, true, true)
	test(float64(0), true, false, true)
	test(5, true, false, false)    // lossy: 5 -> true -> 1 != 5
	test(-1, true, false, false)   // lossy
	test(2.0, false, false, false) // lossy
	test(0.5, true, false, false)  // lossy: fractional
	test(math.MaxInt64, true, false, false)

	// "true"/"false" round-trip; other strings are lossy and fall back to default.
	test("true", false, true, true)
	test("false", true, false, true)
	test("banana", true, true, false) // lossy: default used
	test("banana", false, false, false)
	test("1", false, false, false) // "1" is not "true"/"false" -> lossy

	// Slices: length 1 carries the element's result; empty/longer are lossy.
	test([]string{"true"}, false, true, true)
	test([]string{}, true, true, false)
	test([]string{"true", "false"}, false, true, false)
	test([]any{true}, false, true, true)

	// nil is always lossy and returns the default.
	test(nil, true, true, false)
	test(nil, false, false, false)
}

// TestLossless_IntOk covers IntOk's round-trip rule across numeric, string, and float inputs.
func TestLossless_IntOk(t *testing.T) {

	test := func(value any, wantResult int, wantLossless bool) {
		result, lossless := IntOk(value, -999)
		require.Equal(t, wantResult, result, "result for %#v", value)
		require.Equal(t, wantLossless, lossless, "lossless for %#v", value)
	}

	test(42, 42, true)
	test(int8(-7), -7, true)
	test(int64(123), 123, true)
	test(true, 1, true) // bool 0/1 round-trips to int
	test(false, 0, true)
	test(2.0, 2, true)    // whole float
	test(2.5, 2, false)   // lossy: fractional
	test(-3.9, -3, false) // lossy: fractional
	test("100", 100, true)
	test("-100", -100, true)
	test("not a number", -999, false)
	test([]string{"7"}, 7, true)
	test([]string{"7", "8"}, 7, false) // lossy: longer slice
	test(nil, -999, false)
}

// TestLossless_Int32Ok covers the narrowing rule: out-of-range values clamp and report lossy.
func TestLossless_Int32Ok(t *testing.T) {

	test := func(value any, wantResult int32, wantLossless bool) {
		result, lossless := Int32Ok(value, -999)
		require.Equal(t, wantResult, result, "result for %#v", value)
		require.Equal(t, wantLossless, lossless, "lossless for %#v", value)
	}

	test(int32(42), 42, true)
	test(int64(math.MaxInt32), math.MaxInt32, true)
	test(int64(math.MaxInt32)+1, math.MaxInt32, false) // clamp -> lossy
	test(int64(math.MinInt32)-1, math.MinInt32, false) // clamp -> lossy
	test(int(math.MaxInt32)+1, math.MaxInt32, false)
	test(2.0, 2, true)
	test(2.5, 2, false)
	test("123", 123, true)
	test("99999999999", -999, false) // overflows int32 parse -> default, lossy
}

// TestLossless_Int64Ok covers Int64Ok, where every int fits but floats and strings can be lossy.
func TestLossless_Int64Ok(t *testing.T) {

	test := func(value any, wantResult int64, wantLossless bool) {
		result, lossless := Int64Ok(value, -999)
		require.Equal(t, wantResult, result, "result for %#v", value)
		require.Equal(t, wantLossless, lossless, "lossless for %#v", value)
	}

	test(int64(math.MaxInt64), math.MaxInt64, true)
	test(int64(math.MinInt64), math.MinInt64, true)
	test(42, 42, true)
	test(2.0, 2, true)
	test(2.5, 2, false)
	test(float64(1<<63), math.MaxInt64, false)  // at/above 2^63 -> clamp, lossy
	test(-float64(1<<64), math.MinInt64, false) // well below range -> clamp, lossy
	test("9223372036854775807", math.MaxInt64, true)
	test("bad", -999, false)
}

// TestLossless_FloatOk covers the 2^53 integer-precision boundary for float64.
func TestLossless_FloatOk(t *testing.T) {

	test := func(value any, wantLossless bool) {
		_, lossless := FloatOk(value, -1)
		require.Equal(t, wantLossless, lossless, "lossless for %#v", value)
	}

	test(3.14, true)
	test(float32(1.5), true)
	test(42, true)
	test(int8(100), true)
	test(int32(math.MaxInt32), true) // fits exactly in float64
	test(int64(1)<<53, true)         // 2^53 is exactly representable
	test(-(int64(1) << 53), true)
	test((int64(1)<<53)+1, false)     // 2^53+1 cannot be represented exactly -> lossy
	test(int64(math.MaxInt64), false) // far beyond 2^53 -> lossy
	test(int(math.MaxInt64), false)
	test("3.14159", true)
	test("not a float", false)
	test(nil, false)
}

// TestLossless_StringOk covers int/bool/float formatting, all of which are lossless,
// versus slice and nil inputs.
func TestLossless_StringOk(t *testing.T) {

	test := func(value any, wantResult string, wantLossless bool) {
		result, lossless := StringOk(value, "DEFAULT")
		require.Equal(t, wantResult, result, "result for %#v", value)
		require.Equal(t, wantLossless, lossless, "lossless for %#v", value)
	}

	test("hello", "hello", true)
	test([]byte("hi"), "hi", true)
	test(42, "42", true) // int -> string is lossless
	test(int8(-5), "-5", true)
	test(int64(1000), "1000", true)
	test(true, "true", true) // bool -> string is lossless
	test(false, "false", true)
	test(3.5, "3.50", true)      // <=2 decimals -> lossless
	test(12.34, "12.34", true)   // exactly 2 decimals -> lossless
	test(123.45, "123.45", true) // 2 decimals -> lossless
	test(3.14159, "3.14", false) // needs more than 2 decimals -> lossy (rounded)
	test(123.4567890, "123.46", false)
	test([]string{"x"}, "x", true)
	test([]string{"x", "y"}, "x", false) // longer slice -> lossy
	test(nil, "DEFAULT", false)
}

// TestLossless_FloatToString pins the two-decimal formatting rule: a float that
// fits in two decimals is lossless, while one that needs more precision rounds and
// is reported lossy.
func TestLossless_FloatToString(t *testing.T) {

	test := func(value float64, wantResult string, wantLossless bool) {
		result, lossless := StringOk(value, "DEFAULT")
		require.Equal(t, wantResult, result, "result for %v", value)
		require.Equal(t, wantLossless, lossless, "lossless for %v", value)
	}

	// Lossless: representable exactly in two decimal places.
	test(123.45, "123.45", true)
	test(12.34, "12.34", true)
	test(3.5, "3.50", true)
	test(100.0, "100.00", true)
	test(0.1, "0.10", true)
	test(2.5, "2.50", true)
	test(-98.76, "-98.76", true)

	// Lossy: needs more than two decimals, so the two-decimal form rounds.
	test(123.4567890, "123.46", false)
	test(3.14159, "3.14", false)
	test(1.0/3.0, "0.33", false)
	test(-98.765, "-98.77", false)
}

// TestLossless_Interfaces covers the Inter/Floater/Stringer/Hexer arms of the
// scalar converters, confirming the lossless rule applies through interfaces too.
func TestLossless_Interfaces(t *testing.T) {

	// Inter -> int/int64/float is lossless for in-range values.
	{
		result, lossless := IntOk(customInter(42), -1)
		require.Equal(t, 42, result)
		require.True(t, lossless)
	}

	// Floater with a whole value -> int is lossless; a fraction is lossy.
	{
		result, lossless := IntOk(customFloater(7), -1)
		require.Equal(t, 7, result)
		require.True(t, lossless)
	}
	{
		result, lossless := IntOk(customFloater(7.5), -1)
		require.Equal(t, 7, result)
		require.False(t, lossless) // fractional -> lossy
	}

	// Stringer holding a clean number parses losslessly.
	{
		result, lossless := Int64Ok(customStringer("123"), -1)
		require.Equal(t, int64(123), result)
		require.True(t, lossless)
	}

	// Inter -> string and Floater -> string are both lossless (shortest exact form).
	{
		result, lossless := StringOk(customInter(99), "DEF")
		require.Equal(t, "99", result)
		require.True(t, lossless)
	}
	{
		result, lossless := StringOk(customFloater(1.25), "DEF")
		require.Equal(t, "1.25", result)
		require.True(t, lossless)
	}
}

// TestLossless_FloatExactBoundary brute-forces the exact integer-representability
// boundary for FloatOk around 2^53, where precision is first lost.
func TestLossless_FloatExactBoundary(t *testing.T) {

	const twoTo53 = int64(1) << 53

	// At and below 2^53, every integer is exactly representable.
	for _, v := range []int64{0, 1, -1, twoTo53 - 1, twoTo53, -twoTo53} {
		_, lossless := FloatOk(v, -1)
		require.True(t, lossless, "FloatOk(%d) should be lossless", v)
	}

	// Just above 2^53, the odd integer 2^53+1 cannot be represented.
	for _, v := range []int64{twoTo53 + 1, -(twoTo53 + 1)} {
		_, lossless := FloatOk(v, -1)
		require.False(t, lossless, "FloatOk(%d) should be lossy", v)
	}
}

// TestLossless_SliceScalars confirms that the scalar/interface arms of the
// SliceOf*Ok converters propagate the element's lossiness rather than fabricating
// Ok=true. A single lossy element must make the whole slice conversion lossy.
func TestLossless_SliceScalars(t *testing.T) {

	// Lossless single scalars keep Ok=true.
	{
		result, ok := SliceOfIntOk(int(7))
		require.Equal(t, []int{7}, result)
		require.True(t, ok)
	}
	{
		result, ok := SliceOfInt64Ok(int64(42))
		require.Equal(t, []int64{42}, result)
		require.True(t, ok)
	}

	// A fractional Floater truncates -> the slice conversion is lossy.
	{
		result, ok := SliceOfIntOk(customFloater(1.5))
		require.Equal(t, []int{1}, result)
		require.False(t, ok)
	}
	{
		result, ok := SliceOfInt64Ok(customFloater(2.5))
		require.Equal(t, []int64{2}, result)
		require.False(t, ok)
	}

	// An Inter beyond float64 integer precision -> lossy float slice.
	{
		_, ok := SliceOfFloatOk(customInter(1 << 60))
		require.False(t, ok)
	}

	// A small Inter is lossless.
	{
		result, ok := SliceOfFloatOk(customInter(8))
		require.Equal(t, []float64{8}, result)
		require.True(t, ok)
	}
}
