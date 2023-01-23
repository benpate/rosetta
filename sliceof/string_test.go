package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {

	s := NewString()
	require.Zero(t, s.Length())
	require.True(t, s.SetString("0", "one"))
	require.True(t, s.SetString("1", "two"))
	require.True(t, s.SetString("2", "three"))

	require.Equal(t, "one", s.GetString("0"))
	require.Equal(t, "two", s.GetString("1"))
	require.Equal(t, "three", s.GetString("2"))
	require.Equal(t, "", s.GetString("3"))

	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.False(t, s.Remove("0"))
	require.Zero(t, s.Length())
}
