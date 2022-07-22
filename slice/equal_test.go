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
