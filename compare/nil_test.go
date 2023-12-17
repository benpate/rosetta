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
