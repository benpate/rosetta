package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	require.True(t, Equal([]int{1, 2, 3}, []int{1, 2, 3}))
	require.False(t, Equal([]int{1, 2, 3}, []int{1, 2, 3, 4}))
	require.False(t, Equal([]int{1, 2, 3, 4}, []int{1, 2, 3}))
	require.False(t, Equal([]int{1, 2, 3}, []int{1, 2, 4}))
}

// NotEqual is the exact inverse of Equal, for every input.
func TestNotEqual(t *testing.T) {

	require.False(t, NotEqual([]string{"a", "b"}, []string{"a", "b"}))
	require.True(t, NotEqual([]string{"a", "b"}, []string{"b", "a"}), "order matters")
	require.True(t, NotEqual([]string{"a"}, []string{"a", "b"}), "length matters")
	require.True(t, NotEqual([]string{"a", "b"}, []string{"a"}))
	require.False(t, NotEqual([]string{}, []string{}))
	require.False(t, NotEqual[string](nil, nil))
	require.False(t, NotEqual([]string{}, nil), "a nil slice equals an empty slice")

	inputs := [][]int{nil, {}, {0}, {1}, {1, 2}, {2, 1}, {1, 2, 3}}
	for _, left := range inputs {
		for _, right := range inputs {
			require.Equal(t, !Equal(left, right), NotEqual(left, right),
				"NotEqual must invert Equal for %v vs %v", left, right)
		}
	}
}
