package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt64(t *testing.T) {

	s := NewInt64()
	require.Zero(t, s.Length())
	require.True(t, s.SetInt64("0", 1))
	require.True(t, s.SetInt64("1", 2))
	require.True(t, s.SetInt64("2", 3))

	require.Equal(t, int64(1), s.GetInt64("0"))
	require.Equal(t, int64(2), s.GetInt64("1"))
	require.Equal(t, int64(3), s.GetInt64("2"))
	require.Equal(t, int64(0), s.GetInt64("3"))

	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.False(t, s.Remove("0"))
	require.Zero(t, s.Length())
}
