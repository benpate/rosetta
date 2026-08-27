package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// testArray presents itself as an array through the ArrayGetter interface without being a Go
// slice, the way rosetta's change-tracking and schema containers do.
type testArray struct {
	items []any
}

// Length returns the number of items in the array
func (t testArray) Length() int {
	return len(t.items)
}

// GetIndex returns the value at the specified index, and a boolean indicating success
func (t testArray) GetIndex(index int) (any, bool) {

	if (index < 0) || (index >= len(t.items)) {
		return nil, false
	}

	return t.items[index], true
}

// brokenArray claims a length it cannot serve, standing in for an implementation that does not
// hold up its end of the interface.
type brokenArray struct {
	length int
}

// Length returns the number of items this array claims to hold
func (b brokenArray) Length() int {
	return b.length
}

// GetIndex always reports failure
func (brokenArray) GetIndex(int) (any, bool) {
	return nil, false
}

// TestArrayGetter_Conversions confirms that a value which is array-shaped but not a Go slice can
// still be converted.  Without this, a container like delta.Slice converted to nothing at all.
func TestArrayGetter_Conversions(t *testing.T) {

	strings := testArray{items: []any{"a", "b"}}
	numbers := testArray{items: []any{1, 2}}

	require.Equal(t, []any{"a", "b"}, SliceOfAny(strings))
	require.Equal(t, []string{"a", "b"}, SliceOfString(strings))
	require.Equal(t, []int{1, 2}, SliceOfInt(numbers))
	require.Equal(t, []int64{1, 2}, SliceOfInt64(numbers))
	require.Equal(t, []float64{1, 2}, SliceOfFloat(numbers))

	maps := testArray{items: []any{map[string]any{"a": 1}}}
	require.Equal(t, []map[string]any{{"a": 1}}, SliceOfMap(maps))

	// The generic converter reaches it too, through both the scalar and structural paths
	require.Equal(t, []string{"a", "b"}, SliceOf[string](strings))
	require.Equal(t, []any{"a", "b"}, SliceOf[any](strings))

	// And so does a POINTER to one
	require.Equal(t, []string{"a", "b"}, SliceOfString(&strings))
}

// TestArrayGetter_Empty confirms that an empty array converts successfully rather than being
// mistaken for a failure.
func TestArrayGetter_Empty(t *testing.T) {

	result, ok := SliceOfAnyOk(testArray{items: []any{}})
	require.True(t, ok)
	require.Empty(t, result)
}

// TestArrayGetter_Broken confirms that an implementation which cannot serve the length it
// reports is rejected, rather than yielding a half-filled slice or panicking inside make().
func TestArrayGetter_Broken(t *testing.T) {

	{ // Promises items it cannot produce
		_, ok := SliceOfAnyOk(brokenArray{length: 4})
		require.False(t, ok)

		_, ok = SliceOfStringOk(brokenArray{length: 4})
		require.False(t, ok)
	}

	{ // RULE: A negative length must not reach make()
		require.NotPanics(t, func() { SliceOfAny(brokenArray{length: -1}) })

		_, ok := SliceOfAnyOk(brokenArray{length: -1})
		require.False(t, ok)
	}
}

// TestArrayGetter_DoesNotShadowSlices confirms that a real Go slice still takes the reflection
// fast path.  A named slice type satisfies neither Length nor GetIndex, but this pins that the
// interface check cannot start intercepting ordinary slices.
func TestArrayGetter_DoesNotShadowSlices(t *testing.T) {

	require.Equal(t, []string{"a", "b"}, SliceOfString([]string{"a", "b"}))
	require.Equal(t, []any{"a", "b"}, SliceOfAny([]any{"a", "b"}))
}
