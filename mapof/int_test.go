package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt(t *testing.T) {

	m := NewInt()

	require.Equal(t, 0, m.Length())
	require.True(t, m.IsEmpty())
	require.False(t, m.NotEmpty())

	// NOTE: Int.SetInt stores on zero and deletes on non-zero, which is the
	// inverse of Int64/Float. These assertions document the current behavior.
	require.True(t, m.SetInt("zero", 0))
	require.Equal(t, 1, m.Length())

	require.True(t, m.SetInt("nonzero", 5))
	// "nonzero" was never stored (delete on non-zero), so length stays at 1
	require.Equal(t, 1, m.Length())

	value, ok := m.GetIntOK("zero")
	require.Equal(t, 0, value)
	require.True(t, ok)

	value, ok = m.GetIntOK("missing")
	require.Equal(t, 0, value)
	require.False(t, ok)

	require.Equal(t, 0, m.GetInt("zero"))
	require.Equal(t, []string{"zero"}, m.Keys())

	require.True(t, m.Remove("zero"))
	require.Equal(t, 0, m.Length())
}

func TestInt_Equal(t *testing.T) {

	a := Int{"x": 1, "y": 2}
	b := Int{"x": 1, "y": 2}
	c := Int{"x": 9}

	require.True(t, a.Equal(b))
	require.False(t, a.NotEqual(b))
	require.False(t, a.Equal(c))
	require.True(t, a.NotEqual(c))
}

func TestInt_NilMap(t *testing.T) {
	var m Int
	require.True(t, m.SetInt("key", 0))
	require.Equal(t, 1, m.Length())

	var n Int
	require.True(t, n.Remove("key"))
}
