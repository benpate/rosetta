package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	require.True(t, IsNil(nil))
	require.True(t, IsNil((*int)(nil)))
	require.True(t, IsNil([]int(nil)))
	require.True(t, IsNil(map[string]int(nil)))

	require.False(t, IsNil(0))
	require.False(t, IsNil(1))
	require.False(t, IsNil(""))
	require.False(t, IsNil([]int{}))
	require.False(t, IsNil(map[string]int{}))
}

// TestNil_ArraysAndFuncs confirms that an array is never nil, and a nil func is
func TestNil_ArraysAndFuncs(t *testing.T) {
	require.False(t, IsNil([12]byte{}))
	require.False(t, IsNil([3]int{1, 2, 3}))
	require.True(t, IsNil((func())(nil)))
	require.False(t, IsNil(func() {}))
}
