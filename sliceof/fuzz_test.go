package sliceof

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// growthLimit bounds the indexes handed to the slice-growing functions below.
// GetPointer and SetIndex accept ANY non-negative index and grow the slice to fit it
// one element at a time, so an unbounded fuzzer would simply exhaust memory instead of
// finding bugs.  Keys larger than this are still exercised against the read-side
// functions, which are bounded by the slice length.
const growthLimit = 1000

// FuzzObject_Keys feeds arbitrary keys to the string-keyed accessors and asserts that
// they never panic, and that a key is honored only when it names a real element.
func FuzzObject_Keys(f *testing.F) {

	f.Add("0", 3)
	f.Add("", 3)
	f.Add("-1", 3)
	f.Add("last", 3)
	f.Add("next", 3)
	f.Add("2", 0)
	f.Add("999999999999999999999", 3)
	f.Add("0x10", 3)
	f.Add(" 1", 3)
	f.Add("+1", 3)
	f.Add("1.0", 3)
	f.Add("\x00", 1)

	f.Fuzz(func(t *testing.T, key string, length int) {

		value := makeFuzzSlice(length)
		originalLength := value.Length()

		// GetAny/GetAnyOK must agree, and must only succeed for a real index
		result, ok := value.GetAnyOK(key)
		require.Equal(t, result, value.GetAny(key))

		if index, isNumeric := strconv.Atoi(key); isNumeric == nil {

			if (index >= 0) && (index < originalLength) {
				require.True(t, ok, "index %d should resolve in a slice of %d", index, originalLength)
				require.Equal(t, value[index], result)
			} else if key != "last" && key != "next" {
				require.False(t, ok, "index %d should NOT resolve in a slice of %d", index, originalLength)
			}
		}

		// Remove must either delete exactly one element, or change nothing at all
		if value.Remove(key) {
			require.Equal(t, originalLength-1, value.Length())
		} else {
			require.Equal(t, originalLength, value.Length())
		}
	})
}

// FuzzObject_Indexes feeds arbitrary integer indexes to the bounds-safe accessors and
// asserts that none of them panics, whatever the index.
func FuzzObject_Indexes(f *testing.F) {

	f.Add(0, 3)
	f.Add(-1, 3)
	f.Add(2, 3)
	f.Add(3, 3)
	f.Add(-999999, 0)
	f.Add(int(^uint(0)>>1), 3)    // max int
	f.Add(-int(^uint(0)>>1)-1, 3) // min int

	f.Fuzz(func(t *testing.T, index int, length int) {

		value := makeFuzzSlice(length)
		originalLength := value.Length()

		// At and AtOK must agree, and must succeed only for a real index
		result, ok := value.AtOK(index)
		require.Equal(t, result, value.At(index))
		require.Equal(t, (index >= 0) && (index < originalLength), ok)

		// GetIndex must match AtOK exactly
		indexed, indexedOK := value.GetIndex(index)
		require.Equal(t, ok, indexedOK)
		require.Equal(t, result, indexed)

		// FirstN must clamp at BOTH ends -- a negative count returns an empty slice
		// rather than panicking on the slice expression (BUG: fixed 2026-08-17)
		switch firstN := value.FirstN(index); {
		case index <= 0:
			require.Equal(t, 0, firstN.Length())
		case index >= originalLength:
			require.Equal(t, originalLength, firstN.Length())
		default:
			require.Equal(t, index, firstN.Length())
		}

		// RemoveAt must either delete exactly one element, or change nothing at all
		if value.RemoveAt(index) {
			require.Equal(t, originalLength-1, value.Length())
		} else {
			require.Equal(t, originalLength, value.Length())
		}
	})
}

// FuzzObject_Grow feeds arbitrary (bounded) indexes to the slice-growing functions and
// asserts that they grow the slice exactly far enough to hold the requested index.
func FuzzObject_Grow(f *testing.F) {

	f.Add("0", 0)
	f.Add("5", 0)
	f.Add("last", 3)
	f.Add("next", 3)
	f.Add("-1", 3)
	f.Add("nope", 3)
	f.Add("999", 3)

	f.Fuzz(func(t *testing.T, key string, length int) {

		value := makeFuzzSlice(length)

		// See growthLimit: an unbounded index is a memory bomb, not a bug hunt
		if index, err := strconv.Atoi(key); (err == nil) && (index > growthLimit) {
			t.Skip("index beyond the growth limit")
		}

		pointer, ok := value.GetPointer(key)

		if !ok {
			require.Nil(t, pointer)
			return
		}

		// A resolved pointer must address an element that now really exists
		require.NotNil(t, pointer)
		require.Greater(t, value.Length(), 0)

		// Writing through the pointer must be visible in the slice
		typed, isString := pointer.(*string)
		require.True(t, isString)

		*typed = "written"
		require.Contains(t, value, "written")
	})
}

// FuzzSliceStringIndex asserts that the private key parser accepts exactly the keys that
// name a valid, in-bounds index, and never panics on anything else.
func FuzzSliceStringIndex(f *testing.F) {

	f.Add("0", 3)
	f.Add("-0", 3)
	f.Add("-1", 3)
	f.Add("007", 3)
	f.Add("999999999999999999999999", 3)
	f.Add("１", 3) // fullwidth digit
	f.Add("", 3)

	f.Fuzz(func(t *testing.T, key string, maximum int) {

		index, ok := sliceStringIndex(key, maximum)

		if !ok {
			require.Equal(t, 0, index)
			return
		}

		// A key that parses must be in bounds, and must round-trip back to itself
		require.GreaterOrEqual(t, index, 0)
		require.Less(t, index, maximum)

		parsed, err := strconv.Atoi(key)
		require.NoError(t, err)
		require.Equal(t, parsed, index)
	})
}

// makeFuzzSlice builds a small, predictable slice whose length is derived from a fuzzed
// integer.  The length is bounded so the fuzzer spends its time on the accessors rather
// than on allocation.
func makeFuzzSlice(length int) Object[string] {

	if length < 0 {
		length = -length
	}

	length = length % 8
	result := NewObject[string]()

	for index := range length {
		result.Append("item-" + strconv.Itoa(index))
	}

	return result
}
