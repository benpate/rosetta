package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {

	s := NewBool()
	require.Zero(t, s.Length())
	require.True(t, s.SetBool("0", true))
	require.True(t, s.SetBool("1", true))
	require.True(t, s.SetBool("2", true))

	require.True(t, s.GetBool("0"))
	require.True(t, s.GetBool("1"))
	require.True(t, s.GetBool("2"))
	require.False(t, s.GetBool("3"))

	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.True(t, s.Remove("0"))
	require.False(t, s.Remove("0"))
	require.Zero(t, s.Length())
}
