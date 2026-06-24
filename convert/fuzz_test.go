package convert

import (
	"math"
	"testing"
)

// FuzzStringParsers feeds arbitrary strings to the string-based parsers to confirm
// that none of them panic, regardless of input.
func FuzzStringParsers(f *testing.F) {

	f.Add("")
	f.Add("123")
	f.Add("-456")
	f.Add("3.14159")
	f.Add("true")
	f.Add("false")
	f.Add("ff")
	f.Add("2026-03-04T13:02:00Z")
	f.Add("not a number")
	f.Add("99999999999999999999999999999")

	f.Fuzz(func(t *testing.T, value string) {
		// None of these conversions should ever panic.
		_ = Int(value)
		_ = Int32(value)
		_ = Int64(value)
		_ = Float(value)
		_ = Bool(value)
		_ = Bytes(value)
		_ = String(value)
		_ = Time(value)
		_ = SliceOfString(value)
		_ = SliceOfInt(value)
		_ = SliceOfFloat(value)
		_ = SliceOfAny(value)
	})
}

// FuzzIntNarrowing feeds arbitrary int64 values to the narrowing integer converters and asserts
// the no-silent-truncation invariant: Ok=true only when the value fits the target type exactly.
func FuzzIntNarrowing(f *testing.F) {

	f.Add(int64(0))
	f.Add(int64(math.MaxInt32))
	f.Add(int64(math.MinInt32))
	f.Add(int64(math.MaxInt32) + 1)
	f.Add(int64(math.MinInt32) - 1)
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))

	f.Fuzz(func(t *testing.T, value int64) {

		// int32: Ok exactly when the value is within int32 range, and the result round-trips.
		if result32, ok32 := Int32Ok(value, -1); ok32 != (value >= math.MinInt32 && value <= math.MaxInt32) {
			t.Fatalf("Int32Ok(%d): ok=%v but in-range=%v", value, ok32, value >= math.MinInt32 && value <= math.MaxInt32)
		} else if ok32 && int64(result32) != value {
			t.Fatalf("Int32Ok(%d) reported ok but returned %d", value, result32)
		}

		// int64: every int64 fits, so it is always natural and exact.
		if result64, ok64 := Int64Ok(value, -1); !ok64 || result64 != value {
			t.Fatalf("Int64Ok(%d) = (%d, %v); want (%d, true)", value, result64, ok64, value)
		}

		// int: identical to int64 on 64-bit; on 32-bit, Ok only within int range.
		if resultInt, okInt := IntOk(value, -1); okInt != (value >= math.MinInt && value <= math.MaxInt) {
			t.Fatalf("IntOk(%d): ok=%v but in-range=%v", value, okInt, value >= math.MinInt && value <= math.MaxInt)
		} else if okInt && int64(resultInt) != value {
			t.Fatalf("IntOk(%d) reported ok but returned %d", value, resultInt)
		}
	})
}

// FuzzFloatNarrowing feeds arbitrary float64 values to the integer converters and asserts that an
// out-of-range float never produces an in-range integer with Ok=true (no silent overflow/wrap).
func FuzzFloatNarrowing(f *testing.F) {

	f.Add(0.0)
	f.Add(3.14)
	f.Add(float64(1 << 31))
	f.Add(float64(1 << 63))
	f.Add(-float64(1 << 63))
	f.Add(math.MaxFloat64)
	f.Add(-math.MaxFloat64)

	f.Fuzz(func(t *testing.T, value float64) {

		// NaN/Inf are not meaningful integers; just confirm no panic and skip the invariant.
		if math.IsNaN(value) || math.IsInf(value, 0) {
			_, _ = Int32Ok(value, -1)
			_, _ = Int64Ok(value, -1)
			_, _ = IntOk(value, -1)
			return
		}

		// int32: Ok only when the (truncated) value is representable in int32.
		if result32, ok32 := Int32Ok(value, -1); ok32 && (float64(result32) > value || value >= float64(1<<31) || value < -float64(1<<31)) {
			t.Fatalf("Int32Ok(%v) reported ok with out-of-range result %d", value, result32)
		}

		// int64: an out-of-range float must report Ok=false rather than a wrapped value.
		if result64, ok64 := Int64Ok(value, -1); ok64 && (value >= float64(1<<63) || value < -float64(1<<63)) {
			t.Fatalf("Int64Ok(%v) reported ok=true for an out-of-range float (result %d)", value, result64)
		}
	})
}

// FuzzBoolRoundTrip asserts the BoolOk lossless rule for arbitrary integers:
// Ok=true exactly when the value is 0 or 1 (the only integers that round-trip
// through a bool), and the result then matches that value.
func FuzzBoolRoundTrip(f *testing.F) {

	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(2))
	f.Add(int64(-1))
	f.Add(int64(math.MaxInt64))

	f.Fuzz(func(t *testing.T, value int64) {

		result, ok := BoolOk(value, false)

		wantOk := (value == 0 || value == 1)
		if ok != wantOk {
			t.Fatalf("BoolOk(%d): ok=%v, want %v", value, ok, wantOk)
		}

		// When lossless, the bool must encode the original value exactly.
		if ok && result != (value == 1) {
			t.Fatalf("BoolOk(%d) reported lossless but result=%v", value, result)
		}
	})
}

// FuzzStringIntRoundTrip asserts that converting an int to a string is lossless
// and that the resulting string parses back to the original integer.
func FuzzStringIntRoundTrip(f *testing.F) {

	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(123456789))
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))

	f.Fuzz(func(t *testing.T, value int64) {

		str, ok := StringOk(value, "DEFAULT")
		if !ok {
			t.Fatalf("StringOk(%d) reported lossy; int-to-string must be lossless", value)
		}

		// The string must round-trip back to the same integer.
		back, backOk := Int64Ok(str, -1)
		if !backOk || back != value {
			t.Fatalf("StringOk(%d)=%q did not round-trip back (got %d, ok=%v)", value, str, back, backOk)
		}
	})
}

// FuzzFloatIntRoundTrip asserts that FloatOk reports Ok=true for an int64 only when
// the value is representable exactly in a float64 (it round-trips back unchanged).
func FuzzFloatIntRoundTrip(f *testing.F) {

	f.Add(int64(0))
	f.Add(int64(1) << 53)
	f.Add((int64(1) << 53) + 1)
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))

	f.Fuzz(func(t *testing.T, value int64) {

		_, ok := FloatOk(value, -1)

		// An int64 is exactly representable in float64 exactly when its magnitude
		// is at most 2^53. Ok must agree with that, in both directions.
		const maxExact = int64(1) << 53
		exact := value >= -maxExact && value <= maxExact

		if ok != exact {
			t.Fatalf("FloatOk(%d): ok=%v, want %v (exactly representable)", value, ok, exact)
		}
	})
}
