package slice

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

// Shuffle is random, so it is tested by its invariants rather than an exact result:
// it must be a permutation of the input, in place, and never lose or duplicate items.
func TestShuffle_IsAPermutation(t *testing.T) {

	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	values := append([]int{}, original...)

	result := Shuffle(values)

	sorted := append([]int{}, result...)
	sort.Ints(sorted)
	require.Equal(t, original, sorted, "shuffling must not add, drop, or duplicate items")
}

func TestShuffle_EdgeCases(t *testing.T) {

	require.Equal(t, []int{}, Shuffle([]int{}))
	require.Nil(t, Shuffle[int](nil))
	require.Equal(t, []int{7}, Shuffle([]int{7}))

	// Two items: either order is valid, but both items must survive
	pair := Shuffle([]int{1, 2})
	require.Len(t, pair, 2)
	require.ElementsMatch(t, []int{1, 2}, pair)
}

// Shuffle works IN PLACE and returns the same slice it was given.
func TestShuffle_IsInPlace(t *testing.T) {

	values := []int{1, 2, 3, 4, 5}
	result := Shuffle(values)

	require.ElementsMatch(t, values, result)

	result[0] = 99
	require.Equal(t, 99, values[0], "the returned slice shares the caller's backing array")
}

// Over many runs a 10-item slice must not come back in its original order every time.
// The chance of a false failure is (1/10!)^20, which is effectively zero.
func TestShuffle_ActuallyReorders(t *testing.T) {

	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	reordered := false

	for range 20 {
		if values := append([]int{}, original...); !Equal(Shuffle(values), original) {
			reordered = true
			break
		}
	}

	require.True(t, reordered, "20 shuffles of 10 items never changed the order")
}
