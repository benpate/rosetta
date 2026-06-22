package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDifference(t *testing.T) {

	a := []int{1, 2, 3, 4, 5}
	b := []int{2, 4}

	// Difference must contain only the elements of a that are not in b,
	// with no leading zero-values from the allocation.
	require.Equal(t, []int{1, 3, 5}, Difference(a, b))
}

func TestDifference_Empty(t *testing.T) {

	require.Equal(t, []int{}, Difference([]int{1, 2}, []int{1, 2}))
}
