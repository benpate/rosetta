package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt(t *testing.T) {

	s := NewInt()
	require.Zero(t, s.Length())
	require.True(t, s.SetInt("0", 1))
	require.True(t, s.SetInt("1", 2))
	require.True(t, s.SetInt("2", 3))

	require.Equal(t, int(1), s.GetInt("0"))
	require.Equal(t, int(2), s.GetInt("1"))
	require.Equal(t, int(3), s.GetInt("2"))
	require.Equal(t, int(0), s.GetInt("3"))

	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.False(t, s.Remove("0"))
	require.Zero(t, s.Length())
}
