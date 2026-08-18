package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfInt64_Nil(t *testing.T) {
	expected := []int64{}

	actual, ok := SliceOfInt64Ok(nil)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_Float64(t *testing.T) {
	input := float64(3)
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_Int(t *testing.T) {
	input := 42
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_Int64(t *testing.T) {
	input := int64(42)
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_String(t *testing.T) {
	input := "3"
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_StringDelimited(t *testing.T) {
	input := "3,4,5,6"
	expected := []int64{3, 4, 5, 6}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_ReflectValue(t *testing.T) {
	input := ReflectValue(42)
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfAny(t *testing.T) {
	input := []any{42}
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfInt(t *testing.T) {
	input := []int{42}
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfInt64(t *testing.T) {
	input := []int64{42}
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfFloat64(t *testing.T) {
	input := []float64{3}
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfString(t *testing.T) {
	input := []string{"3"}
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

// SliceOfInt64 is the error-swallowing wrapper around SliceOfInt64Ok: it returns whatever
// Ok produced when that is non-nil, and an empty slice otherwise. It never returns nil.
func TestSliceOfInt64(t *testing.T) {

	check := func(name string, input any, expected []int64) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			result := SliceOfInt64(input)
			require.NotNil(t, result, "SliceOfInt64 must never return nil")
			require.Equal(t, expected, result)
		})
	}

	// Scalars widen into a one-item slice
	check("int", 42, []int64{42})
	check("int64", int64(42), []int64{42})
	check("float64", float64(3), []int64{3})
	check("single string", "3", []int64{3})
	check("comma-delimited string", "1,2,3", []int64{1, 2, 3})

	// Slices convert element by element
	check("slice of int", []int{1, 2}, []int64{1, 2})
	check("slice of int64", []int64{1, 2}, []int64{1, 2})
	check("slice of float64", []float64{1, 2}, []int64{1, 2})
	check("slice of string", []string{"1", "2"}, []int64{1, 2})
	check("slice of any", []any{1, "2", float64(3)}, []int64{1, 2, 3})

	// Empty inputs
	check("empty slice of int64", []int64{}, []int64{})
	check("empty slice of string", []string{}, []int64{})

	// Failures come back as an EMPTY slice, because the wrapper drops the ok flag
	check("nil", nil, []int64{})
	check("bool is not convertible", true, []int64{})
	check("struct is not convertible", struct{ X int }{1}, []int64{})
}

// A TYPED nil slice reaches `case []int64` and comes back nil from Ok, which is the one
// path that exercises the wrapper's own empty-slice fallback.
func TestSliceOfInt64_TypedNilSlice(t *testing.T) {

	fromOk, ok := SliceOfInt64Ok([]int64(nil))
	require.True(t, ok)
	require.Nil(t, fromOk, "Ok passes the caller's nil slice straight through")

	result := SliceOfInt64([]int64(nil))
	require.NotNil(t, result, "the wrapper substitutes an empty slice for Ok's nil")
	require.Equal(t, []int64{}, result)
}

// The wrapper hides the difference between "converted successfully" and "gave up",
// because both report an empty slice. Pinned so the behavior is deliberate, not accidental.
func TestSliceOfInt64_SwallowsFailure(t *testing.T) {

	_, ok := SliceOfInt64Ok(true)
	require.False(t, ok, "Ok reports the failure")

	require.Equal(t, SliceOfInt64([]int64{}), SliceOfInt64(true),
		"the wrapper reports a failure and an empty input identically")
}

// Boundary values must survive the round trip without truncation or wraparound.
func TestSliceOfInt64_Boundaries(t *testing.T) {

	require.Equal(t, []int64{9223372036854775807}, SliceOfInt64(int64(9223372036854775807)))
	require.Equal(t, []int64{-9223372036854775808}, SliceOfInt64(int64(-9223372036854775808)))
	require.Equal(t, []int64{0}, SliceOfInt64(0))
	require.Equal(t, []int64{-1}, SliceOfInt64(-1))
	require.Equal(t, []int64{1, -1, 0}, SliceOfInt64([]int64{1, -1, 0}))
}
