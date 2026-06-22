package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {

	m := NewBool()

	require.Equal(t, 0, m.Length())
	require.True(t, m.IsEmpty())
	require.False(t, m.NotEmpty())

	require.True(t, m.SetBool("yes", true))
	require.True(t, m.SetBool("no", false))

	require.Equal(t, 2, m.Length())
	require.False(t, m.IsEmpty())
	require.True(t, m.NotEmpty())

	require.Equal(t, []string{"no", "yes"}, m.Keys())

	require.True(t, m.GetBool("yes"))
	require.False(t, m.GetBool("no"))

	// Missing key
	value, ok := m.GetBoolOK("missing")
	require.False(t, value)
	require.False(t, ok)

	// Existing key
	value, ok = m.GetBoolOK("yes")
	require.True(t, value)
	require.True(t, ok)

	require.True(t, m.Remove("yes"))
	require.Equal(t, 1, m.Length())
}

func TestBool_Equal(t *testing.T) {

	a := Bool{"x": true, "y": false}
	b := Bool{"x": true, "y": false}
	c := Bool{"x": false}

	require.True(t, a.Equal(b))
	require.False(t, a.NotEqual(b))
	require.False(t, a.Equal(c))
	require.True(t, a.NotEqual(c))
}

func TestBool_NilMap(t *testing.T) {

	// Setting on a nil map should initialize it
	var m Bool
	require.True(t, m.SetBool("key", true))
	require.Equal(t, 1, m.Length())

	// Remove on a nil map should not panic
	var n Bool
	require.True(t, n.Remove("key"))
}
