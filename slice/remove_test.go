package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveAt(t *testing.T) {

	// Test Integers
	doInt := func(slice []int, item int, result []int) {
		require.Equal(t, result, RemoveAt(slice, item))
	}

	doInt([]int{1, 2, 3}, -1, []int{1, 2, 3})
	doInt([]int{1, 2, 3}, 0, []int{2, 3})
	doInt([]int{1, 2, 3}, 1, []int{1, 3})
	doInt([]int{1, 2, 3}, 2, []int{1, 2})
	doInt([]int{1, 2, 3}, 3, []int{1, 2, 3})

	// Test Strings
	doString := func(slice []string, index int, result []string) {
		require.Equal(t, result, RemoveAt(slice, index))
	}

	doString([]string{"a", "b", "c"}, -1, []string{"a", "b", "c"})
	doString([]string{"a", "b", "c"}, 0, []string{"b", "c"})
	doString([]string{"a", "b", "c"}, 1, []string{"a", "c"})
	doString([]string{"a", "b", "c"}, 2, []string{"a", "b"})
	doString([]string{"a", "b", "c"}, 3, []string{"a", "b", "c"})
}
