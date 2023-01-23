package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat(t *testing.T) {

	s := NewFloat()
	require.Zero(t, s.Length())
	require.True(t, s.SetFloat("0", 1))
	require.True(t, s.SetFloat("1", 2))
	require.True(t, s.SetFloat("2", 3))

	require.Equal(t, float64(1), s.GetFloat("0"))
	require.Equal(t, float64(2), s.GetFloat("1"))
	require.Equal(t, float64(3), s.GetFloat("2"))
	require.Equal(t, float64(0), s.GetFloat("3"))

	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.False(t, s.Remove("0"))
	require.Zero(t, s.Length())
}
