package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverse(t *testing.T) {

	check := func(name string, values []int, expected []int) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			require.Equal(t, expected, Reverse(values))
		})
	}

	check("odd length keeps the middle item in place", []int{1, 2, 3}, []int{3, 2, 1})
	check("even length", []int{1, 2, 3, 4}, []int{4, 3, 2, 1})
	check("two items", []int{1, 2}, []int{2, 1})
	check("one item", []int{1}, []int{1})
	check("empty slice", []int{}, []int{})
}

func TestReverse_NilSlice(t *testing.T) {
	require.Nil(t, Reverse[int](nil))
}

// Reverse works IN PLACE and returns the same slice it was given.
func TestReverse_IsInPlace(t *testing.T) {

	values := []string{"a", "b", "c"}
	result := Reverse(values)

	require.Equal(t, []string{"c", "b", "a"}, values, "the caller's slice is reversed in place")
	require.Equal(t, []string{"c", "b", "a"}, result)

	// Same backing array, so a write through one is visible through the other
	result[0] = "z"
	require.Equal(t, "z", values[0])
}

// Reversing twice restores the original order.
func TestReverse_IsItsOwnInverse(t *testing.T) {

	original := []int{5, 1, 4, 2, 3}
	values := append([]int{}, original...)

	require.Equal(t, original, Reverse(Reverse(values)))
}
