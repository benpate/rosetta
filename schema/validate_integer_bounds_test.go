package schema

import (
	"strings"
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// TestValidateInteger_Gates pins the two distinct rejection gates in validate_Integer: a value
// that is not a clean integer (non-numeric or fractional) must report "must be an integer", while
// a clean integer that exceeds the declared bit size must report "fit within the specified bit
// size". The two gates must not be conflated, and neither may silently accept a bad value.
func TestValidateInteger_Gates(t *testing.T) {

	// notInteger asserts the value is rejected with the "must be an integer" message.
	notInteger := func(value any) {
		t.Helper()
		_, _, err := validate_Integer(Integer{}, value)
		require.Error(t, err, "value: %#v", value)
		require.True(t, strings.Contains(err.Error(), "must be an integer"),
			"value %#v should report 'must be an integer', got: %v", value, err)
	}

	// outOfBounds asserts the value is rejected with the "fit within the specified bit size" message.
	outOfBounds := func(value any, bitSize int) {
		t.Helper()
		_, _, err := validate_Integer(Integer{BitSize: bitSize}, value)
		require.Error(t, err, "value: %#v bitSize: %d", value, bitSize)
		require.True(t, strings.Contains(err.Error(), "fit within the specified bit size"),
			"value %#v should report bit-size overflow, got: %v", value, err)
	}

	// The lossless gate: non-numeric and fractional inputs are not integers.
	notInteger("abc")
	notInteger("not-a-number")
	notInteger(struct{}{})
	notInteger(3.5)
	notInteger(nil)

	// The inBounds gate: clean integers that overflow the declared width.
	outOfBounds(300, 8)
	outOfBounds(-200, 8)
	outOfBounds(70000, 16)
	outOfBounds(5_000_000_000, 32)
}

// TestValidateInteger_OutOfRangeBound demonstrates the bounds-narrowing guard: a schema
// bound that does not fit the declared integer type must produce an error rather than
// silently wrapping. Before the guard, an int8 Minimum of 300 wrapped to 44 (int8(300)
// == 44), so a value like 100 would clamp UP to the wrapped floor instead of being
// recognized as below the real minimum.
func TestValidateInteger_OutOfRangeBound(t *testing.T) {

	// An int8 field with a minimum of 300, which cannot fit in an int8.
	element := Integer{
		BitSize: 8,
		Minimum: null.NewInt64(300),
	}

	// An int8 value of 44 is exactly what the old wrapped minimum (int8(300) == 44)
	// would have compared against. With the guard, narrowing the bound fails instead.
	_, _, err := validate_Integer_Generic(element, int8(44), false)
	require.Error(t, err, "an out-of-range minimum must produce an error, not a wrapped comparison")

	// A value above the wrapped floor (which the old code would have accepted) also errors,
	// because the bound itself is rejected before any comparison can succeed.
	_, _, err = validate_Integer_Generic(element, int8(100), false)
	require.Error(t, err, "the bound is invalid regardless of the value")
}

// TestValidateInteger_OutOfRangeMaximum demonstrates the more dangerous direction: int8(1000)
// wraps to -24, so the OLD code compared every value against -24 and clamped legitimate values
// DOWN to that wrapped ceiling. The guard now rejects the impossible bound eagerly: a maximum
// that cannot fit the declared bit size is a schema error, not a wrapped comparison.
func TestValidateInteger_OutOfRangeMaximum(t *testing.T) {

	// An int8 field with a maximum of 1000, which cannot fit in an int8.
	element := Integer{
		BitSize: 8,
		Maximum: null.NewInt64(1000),
	}

	// Even a value well within the real maximum surfaces the invalid-bound error, because the
	// bound is validated whenever it is present (not only when a value happens to exceed it).
	_, _, err := validate_Integer_Generic(element, int8(50), false)
	require.Error(t, err, "a maximum that does not fit the declared int8 type must be rejected")
}

// TestValidateInteger_InRangeBoundsStillWork confirms the guard does not break the normal
// case: bounds that fit the declared type still clamp values as before.
func TestValidateInteger_InRangeBoundsStillWork(t *testing.T) {

	element := Integer{
		BitSize: 8,
		Minimum: null.NewInt64(10),
		Maximum: null.NewInt64(100),
	}

	// Below the minimum: clamps up to 10, reported as changed.
	result, changed, err := validate_Integer_Generic(element, int8(5), false)
	require.NoError(t, err)
	require.True(t, changed)
	require.Equal(t, int8(10), result)

	// Above the maximum: clamps down to 100, reported as changed.
	result, changed, err = validate_Integer_Generic(element, int8(120), false)
	require.NoError(t, err)
	require.True(t, changed)
	require.Equal(t, int8(100), result)

	// Within range: passes through unchanged.
	result, changed, err = validate_Integer_Generic(element, int8(50), false)
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, int8(50), result)
}

// TestValidateInteger_OutOfRangeBoundEndToEnd exercises the full Validate flow to confirm
// the guard fires through the public entry point, not just the internal generic helper.
func TestValidateInteger_OutOfRangeBoundEndToEnd(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Integer{BitSize: 8, Minimum: null.NewInt64(300)},
		},
	})

	var getter testIntPointer

	// Setting a value runs Set+Validate; the impossible int8 minimum must surface as an error.
	err := schema.Set(&getter, "value", int(50))
	require.Error(t, err, "an int8 field with a minimum of 300 must reject, not silently clamp")
}
