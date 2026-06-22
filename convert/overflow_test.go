package convert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests verify that narrowing integer conversions never silently truncate: a value outside
// the target range must clamp to the boundary AND report Ok=false. In-range values report Ok=true.

// floatInter / floatFloater exercise the Inter and Floater interface paths with custom types.
type testInter int64

func (i testInter) Int() int { return int(i) }

type testFloater float64

func (f testFloater) Float() float64 { return float64(f) }

func TestInt32Ok_Overflow(t *testing.T) {

	cases := []struct {
		name   string
		input  any
		expect int32
		ok     bool
	}{
		// In-range stays natural
		{"int in range", int(123), 123, true},
		{"int64 in range", int64(-123), -123, true},
		{"int8 max", int8(math.MaxInt8), math.MaxInt8, true},
		{"int16 min", int16(math.MinInt16), math.MinInt16, true},
		{"int32 max", int32(math.MaxInt32), math.MaxInt32, true},
		{"int32 min", int32(math.MinInt32), math.MinInt32, true},

		// int (64-bit) past int32 range — the originally-unguarded path
		{"int above max", int(math.MaxInt32) + 1, math.MaxInt32, false},
		{"int below min", int(math.MinInt32) - 1, math.MinInt32, false},

		// int64 past int32 range
		{"int64 above max", int64(math.MaxInt32) + 1, math.MaxInt32, false},
		{"int64 below min", int64(math.MinInt32) - 1, math.MinInt32, false},
		{"int64 far above", int64(math.MaxInt64), math.MaxInt32, false},
		{"int64 far below", int64(math.MinInt64), math.MinInt32, false},

		// float64 / float32 at the exact 2^31 boundary that the naive guard let slip through
		{"float64 == 2^31", float64(1 << 31), math.MaxInt32, false},
		{"float64 below -2^31", -float64(1<<31) - 1, math.MinInt32, false},
		{"float32 == 2^31", float32(1 << 31), math.MaxInt32, false},
		{"float64 in range", float64(1000), 1000, true},

		// string overflow returns default + false
		{"string overflow", "99999999999999999", int32(-1), false},
		{"string in range", "1000", 1000, true},

		// interface paths
		{"Inter above max", testInter(int64(math.MaxInt32) + 1), math.MaxInt32, false},
		{"Floater 2^31", testFloater(float64(1 << 31)), math.MaxInt32, false},
	}

	for _, c := range cases {
		result, ok := Int32Ok(c.input, -1)
		require.Equal(t, c.ok, ok, c.name)
		require.Equal(t, c.expect, result, c.name)
	}
}

func TestInt64Ok_Overflow(t *testing.T) {

	cases := []struct {
		name   string
		input  any
		expect int64
		ok     bool
	}{
		// Every signed/widening source is in range for int64
		{"int", int(math.MaxInt32), math.MaxInt32, true},
		{"int64 max", int64(math.MaxInt64), math.MaxInt64, true},
		{"int64 min", int64(math.MinInt64), math.MinInt64, true},

		// float64 at exactly 2^63 (which math.MaxInt64 rounds up to) must NOT overflow through
		{"float64 == 2^63", float64(1 << 63), math.MaxInt64, false},
		{"float64 below -2^63", -float64(1<<63) * 2, math.MinInt64, false},
		{"float32 == 2^63", float32(1 << 63), math.MaxInt64, false},
		{"float64 in range", float64(1000), 1000, true},

		// string overflow
		{"string overflow", "99999999999999999999999999", int64(-1), false},
		{"string in range", "1234567890123", 1234567890123, true},

		// interface paths
		{"Floater 2^63", testFloater(float64(1 << 63)), math.MaxInt64, false},
	}

	for _, c := range cases {
		result, ok := Int64Ok(c.input, -1)
		require.Equal(t, c.ok, ok, c.name)
		require.Equal(t, c.expect, result, c.name)
	}
}

func TestIntOk_Overflow(t *testing.T) {

	cases := []struct {
		name   string
		input  any
		expect int
		ok     bool
	}{
		{"int", int(123), 123, true},
		{"int64 in range", int64(123), 123, true},

		// On 64-bit, float64 at 2^63 must not overflow through the guard.
		{"float64 == 2^63", float64(1 << 63), math.MaxInt, false},
		{"float64 below -2^63", -float64(1<<63) * 2, math.MinInt, false},
		{"float64 in range", float64(1000), 1000, true},

		// Floater interface previously reported ok=true for out-of-range values
		{"Floater 2^63", testFloater(float64(1 << 63)), math.MaxInt, false},
		{"Floater in range", testFloater(42.0), 42, true},

		// string overflow
		{"string overflow", "99999999999999999999999999", int(-1), false},
	}

	for _, c := range cases {
		result, ok := IntOk(c.input, -1)
		require.Equal(t, c.ok, ok, c.name)
		require.Equal(t, c.expect, result, c.name)
	}
}

// TestInt_NoSilentTruncation is a focused invariant check: for any int64 input, converting to a
// narrower type must either preserve the exact value (ok=true) or report ok=false. It must never
// return a different in-range value with ok=true.
func TestInt_NoSilentTruncation(t *testing.T) {

	inputs := []int64{
		0, 1, -1,
		math.MaxInt8, math.MinInt8,
		math.MaxInt16, math.MinInt16,
		math.MaxInt32, math.MinInt32,
		int64(math.MaxInt32) + 1, int64(math.MinInt32) - 1,
		math.MaxInt64, math.MinInt64,
	}

	for _, in := range inputs {
		if result, ok := Int32Ok(in, -999); ok {
			// Natural conversion must round-trip exactly.
			require.Equal(t, in, int64(result), "Int32Ok(%d) reported ok but changed the value", in)
		} else {
			// Rejected values must be exactly the values that don't fit in int32.
			require.True(t, in > math.MaxInt32 || in < math.MinInt32,
				"Int32Ok(%d) reported NOT-ok for an in-range value", in)
		}
	}
}
