package convert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// IntBitsizeOk returns (result, lossless, inBounds): lossless reports whether the
// underlying numeric conversion round-tripped (no parse failure, no fractional part),
// and inBounds reports whether the value fit the requested bit size without clamping.

// TestIntBitsizeOk_8 exercises narrowing to int8, including in-range values, the exact
// boundaries, clamping of out-of-range values (inBounds=false), and lossy inputs.
func TestIntBitsizeOk_8(t *testing.T) {

	// assert checks the returned value, its concrete type (int8), and both flags.
	assert := func(value any, expected int8, expectedLossless, expectedInBounds bool) {
		t.Helper()
		result, lossless, inBounds := IntBitsizeOk(value, 0, 8)
		require.Equal(t, expected, result, "value: %#v", value)
		require.IsType(t, int8(0), result, "value: %#v", value)
		require.Equal(t, expectedLossless, lossless, "lossless for value: %#v", value)
		require.Equal(t, expectedInBounds, inBounds, "inBounds for value: %#v", value)
	}

	// In-range values: lossless and in bounds.
	assert(0, 0, true, true)
	assert(100, 100, true, true)
	assert(-100, -100, true, true)

	// Exact boundaries are in range.
	assert(127, 127, true, true)
	assert(-128, -128, true, true)
	assert(math.MaxInt8, math.MaxInt8, true, true)
	assert(math.MinInt8, math.MinInt8, true, true)

	// One step past each boundary clamps: the source converts losslessly but is out of bounds.
	assert(128, 127, true, false)
	assert(-129, -128, true, false)
	assert(300, 127, true, false)
	assert(-300, -128, true, false)

	// Source-type variety, all in range.
	assert(int8(5), 5, true, true)
	assert(int16(5), 5, true, true)
	assert(int32(5), 5, true, true)
	assert(int64(5), 5, true, true)
	assert(true, 1, true, true)
	assert("42", 42, true, true)

	// Whole floats are lossless; fractional floats are lossy but still in bounds.
	assert(5.0, 5, true, true)
	assert(1.5, 1, false, true)

	// Out-of-int8-range from other source types: lossless conversion, but out of bounds.
	assert(int64(1000), 127, true, false)
	assert("1000", 127, true, false)
}

// TestIntBitsizeOk_16 exercises narrowing to int16 with the same structure as the int8 case.
func TestIntBitsizeOk_16(t *testing.T) {

	assert := func(value any, expected int16, expectedLossless, expectedInBounds bool) {
		t.Helper()
		result, lossless, inBounds := IntBitsizeOk(value, 0, 16)
		require.Equal(t, expected, result, "value: %#v", value)
		require.IsType(t, int16(0), result, "value: %#v", value)
		require.Equal(t, expectedLossless, lossless, "lossless for value: %#v", value)
		require.Equal(t, expectedInBounds, inBounds, "inBounds for value: %#v", value)
	}

	assert(0, 0, true, true)
	assert(30000, 30000, true, true)
	assert(-30000, -30000, true, true)

	assert(math.MaxInt16, math.MaxInt16, true, true)
	assert(math.MinInt16, math.MinInt16, true, true)

	assert(math.MaxInt16+1, math.MaxInt16, true, false)
	assert(math.MinInt16-1, math.MinInt16, true, false)
	assert(100000, math.MaxInt16, true, false)
	assert(-100000, math.MinInt16, true, false)

	assert(1.5, 1, false, true)
}

// TestIntBitsizeOk_32 exercises narrowing to int32, including values that exceed int32
// but fit within the 64-bit working integer.
func TestIntBitsizeOk_32(t *testing.T) {

	assert := func(value any, expected int32, expectedLossless, expectedInBounds bool) {
		t.Helper()
		result, lossless, inBounds := IntBitsizeOk(value, 0, 32)
		require.Equal(t, expected, result, "value: %#v", value)
		require.IsType(t, int32(0), result, "value: %#v", value)
		require.Equal(t, expectedLossless, lossless, "lossless for value: %#v", value)
		require.Equal(t, expectedInBounds, inBounds, "inBounds for value: %#v", value)
	}

	assert(0, 0, true, true)
	assert(2_000_000_000, 2_000_000_000, true, true)

	assert(math.MaxInt32, math.MaxInt32, true, true)
	assert(math.MinInt32, math.MinInt32, true, true)

	assert(int64(math.MaxInt32)+1, math.MaxInt32, true, false)
	assert(int64(math.MinInt32)-1, math.MinInt32, true, false)
	assert(int64(1)<<40, math.MaxInt32, true, false)
	assert(-(int64(1) << 40), math.MinInt32, true, false)
}

// TestIntBitsizeOk_64 confirms the 64-bit case passes the working integer through as int64
// without any narrowing, so it is always in bounds; only the lossless flag varies.
func TestIntBitsizeOk_64(t *testing.T) {

	assert := func(value any, expected int64, expectedLossless, expectedInBounds bool) {
		t.Helper()
		result, lossless, inBounds := IntBitsizeOk(value, 0, 64)
		require.Equal(t, expected, result, "value: %#v", value)
		require.IsType(t, int64(0), result, "value: %#v", value)
		require.Equal(t, expectedLossless, lossless, "lossless for value: %#v", value)
		require.Equal(t, expectedInBounds, inBounds, "inBounds for value: %#v", value)
	}

	assert(0, 0, true, true)
	assert(5, 5, true, true)
	assert(-5, -5, true, true)
	assert(int64(math.MaxInt64), math.MaxInt64, true, true)
	assert(int64(math.MinInt64), math.MinInt64, true, true)

	// A value larger than int64 is clamped to math.MaxInt by Int64Ok (lossy), then passed
	// through the 64-bit case which never narrows, so inBounds is true.
	assert(math.MaxFloat64, math.MaxInt64, false, true)

	// Fractional float remains lossy even at 64 bits.
	assert(1.5, 1, false, true)
}

// TestIntBitsizeOk_DefaultBitSize confirms that any bit size other than 8/16/32/64 is treated
// as a plain int, clamped to the int range. This documents the catch-all contract.
func TestIntBitsizeOk_DefaultBitSize(t *testing.T) {

	assert := func(bitSize int) {
		t.Helper()
		result, lossless, inBounds := IntBitsizeOk(42, 0, bitSize)
		require.IsType(t, int(0), result, "bitSize: %d", bitSize)
		require.Equal(t, 42, result, "bitSize: %d", bitSize)
		require.True(t, lossless, "bitSize: %d", bitSize)
		require.True(t, inBounds, "bitSize: %d", bitSize)
	}

	// Zero, an unsupported width, a too-large width, and a negative width all fall through
	// to plain int rather than erroring or clamping to a smaller width.
	assert(0)
	assert(7)
	assert(100)
	assert(-1)
}

// TestIntBitsizeOk_NonNumeric confirms that values that cannot be converted fall back to the
// default value and are reported as NOT lossless (the parse failed). Whether they are in bounds
// depends only on whether the default value fits the bit size.
func TestIntBitsizeOk_NonNumeric(t *testing.T) {

	// A non-numeric string uses the default value; parse failed, so lossless is false.
	result, lossless, inBounds := IntBitsizeOk("not-a-number", 99, 8)
	require.Equal(t, int8(99), result)
	require.False(t, lossless)
	require.True(t, inBounds, "default 99 fits int8")

	// nil uses the default value too.
	result, lossless, inBounds = IntBitsizeOk(nil, 5, 8)
	require.Equal(t, int8(5), result)
	require.False(t, lossless)
	require.True(t, inBounds)

	// A default value that itself exceeds the bit size is clamped: not lossless AND out of bounds.
	result, lossless, inBounds = IntBitsizeOk("not-a-number", 1000, 8)
	require.Equal(t, int8(math.MaxInt8), result)
	require.False(t, lossless)
	require.False(t, inBounds, "default 1000 does not fit int8")
}

// TestIntBitsizeOk_DefaultValueInRange confirms that when the input is unconvertible but the
// default value fits the bit size, the default is returned at the right width, in bounds but
// not lossless (no conversion of the input actually succeeded).
func TestIntBitsizeOk_DefaultValueInRange(t *testing.T) {

	result, lossless, inBounds := IntBitsizeOk(struct{}{}, 50, 16)
	require.Equal(t, int16(50), result)
	require.IsType(t, int16(0), result)
	require.False(t, lossless)
	require.True(t, inBounds)
}
