package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSliceOf_ScalarTypes confirms that SliceOf applies the same coercion rules as the
// concrete SliceOfXxx function matching its element type.
func TestSliceOf_ScalarTypes(t *testing.T) {

	{ // Strings coerce from anything the concrete converter accepts
		require.Equal(t, []string{"a", "b"}, SliceOf[string]([]string{"a", "b"}))
		require.Equal(t, []string{"1", "2"}, SliceOf[string]([]int{1, 2}))
		require.Equal(t, []string{"12345"}, SliceOf[string](12345))
		require.Equal(t, []string{"a"}, SliceOf[string]("a"))
	}

	{ // So do ints
		require.Equal(t, []int{1, 2}, SliceOf[int]([]int{1, 2}))
		require.Equal(t, []int{1, 2}, SliceOf[int]([]string{"1", "2"}))
		require.Equal(t, []int{12345}, SliceOf[int](12345))
	}

	// ...and int64s, float64s, and maps.  Each input is deliberately NOT already a []T, so
	// that the conversion is exercised instead of the pass-through fast path.
	{
		require.Equal(t, []int64{1, 2}, SliceOf[int64]([]string{"1", "2"}))
		require.Equal(t, []float64{1.5, 2.5}, SliceOf[float64]([]string{"1.5", "2.5"}))
		require.Equal(t, []float64{3}, SliceOf[float64](3))
		require.Equal(t, []map[string]any{{"a": 1}}, SliceOf[map[string]any]([]any{map[string]any{"a": 1}}))
	}

	{ // An interface element type is recognized, and does not vanish into a nil interface
		require.Equal(t, []any{"a", "b"}, SliceOf[any]([]string{"a", "b"}))
		require.Equal(t, []any{1, 2}, SliceOf[any]([]int{1, 2}))
	}
}

// namedString is a named type with no concrete converter, standing in for the element types
// that SliceOf must reshape structurally.
type namedString string

// TestSliceOf_StructuralTypes confirms that an element type this package cannot convert is
// reshaped structurally: pointers and named slice types are unwrapped, but each item must
// already be a T.
func TestSliceOf_StructuralTypes(t *testing.T) {

	{ // Items that are already the right type pass through
		result, ok := SliceOfOk[namedString]([]namedString{"a", "b"})
		require.True(t, ok)
		require.Equal(t, []namedString{"a", "b"}, result)
	}

	{ // A []any holding the right type is unwrapped item by item
		result, ok := SliceOfOk[namedString]([]any{namedString("a"), namedString("b")})
		require.True(t, ok)
		require.Equal(t, []namedString{"a", "b"}, result)
	}

	{ // RULE: There is no value conversion for an unknown element type
		_, ok := SliceOfOk[namedString]([]string{"a", "b"})
		require.False(t, ok)

		_, ok = SliceOfOk[namedString](12345)
		require.False(t, ok)
	}
}

// TestSliceOf_RejectsNonCollections confirms that a value which is not a collection at all is
// reported as such, for every element type.
func TestSliceOf_RejectsNonCollections(t *testing.T) {

	for _, value := range []any{
		struct{ Name string }{Name: "a"},
		map[string]any{"a": 1},
		make(chan int),
		func() {},
	} {
		_, ok := SliceOfOk[string](value)
		require.False(t, ok, "value %T must not convert", value)

		_, ok = SliceOfOk[int](value)
		require.False(t, ok, "value %T must not convert", value)

		// The structural path rejects it too, before ever reaching the per-item check
		_, ok = SliceOfOk[namedString](value)
		require.False(t, ok, "value %T must not convert", value)
	}
}

// TestSliceOf_PassthroughDoesNotCopy confirms the fast path: a value that is already a []T is
// returned as-is, without an intermediate slice or per-item boxing.
func TestSliceOf_PassthroughDoesNotCopy(t *testing.T) {

	original := []string{"a", "b"}
	result := SliceOf[string](original)

	// Writing through the result reaches the original, proving no copy was made
	result[0] = "changed"
	require.Equal(t, "changed", original[0])
}

// TestSliceOf_NilInputs confirms that every flavor of nil is reported rather than panicking.
func TestSliceOf_NilInputs(t *testing.T) {

	var nilSlicePtr *[]string

	for _, value := range []any{nil, nilSlicePtr, []string(nil)} {
		require.NotPanics(t, func() { SliceOf[string](value) })
		require.NotPanics(t, func() { SliceOf[int](value) })
		require.NotPanics(t, func() { SliceOf[any](value) })
	}

	// A nil slice IS a collection -- an empty one -- so it converts successfully
	result, ok := SliceOfOk[string]([]string(nil))
	require.True(t, ok)
	require.Empty(t, result)
}

// TestSliceOf_LossyVersusImpossible pins what the single flag does and does not tell a
// caller.  FALSE covers two different situations -- a conversion that worked but did not
// round-trip, and a value that could not be converted at all -- which are told apart by the
// returned slice, not by the flag.
func TestSliceOf_LossyVersusImpossible(t *testing.T) {

	{ // Lossy but usable: the fixed two-decimal float rendering does not round-trip
		result, lossless := SliceOfOk[string](3.14159)
		require.False(t, lossless)
		require.Equal(t, []string{"3.14"}, result)
	}

	{ // Lossy but usable: a non-numeric string renders as the zero int
		result, lossless := SliceOfOk[int]("not a number")
		require.False(t, lossless)
		require.Equal(t, []int{0}, result)
	}

	{ // Lossless
		result, lossless := SliceOfOk[string]([]string{"a"})
		require.True(t, lossless)
		require.Equal(t, []string{"a"}, result)
	}

	{ // Impossible: a struct is not a collection, so it converts to an empty slice
		result, lossless := SliceOfOk[string](struct{ Name string }{Name: "a"})
		require.False(t, lossless)
		require.Empty(t, result)
	}

	{ // Impossible: an element type with no concrete converter, whose items are not Ts
		result, lossless := SliceOfOk[namedString]([]string{"a"})
		require.False(t, lossless)
		require.Empty(t, result)
	}

	{ // RULE: An empty collection is LOSSLESS, not a failure -- it just has nothing in it
		result, lossless := SliceOfOk[string]([]string{})
		require.True(t, lossless)
		require.Empty(t, result)
	}
}
