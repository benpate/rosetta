package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat(t *testing.T) {

	m := NewFloat()

	require.Equal(t, 0, m.Length())
	require.True(t, m.IsEmpty())
	require.False(t, m.NotEmpty())

	require.True(t, m.SetFloat("pi", 3.14))
	require.Equal(t, 1, m.Length())
	require.True(t, m.NotEmpty())

	// Setting to zero deletes the key
	require.True(t, m.SetFloat("pi", 0))
	require.Equal(t, 0, m.Length())

	m.SetFloat("a", 1.5)
	m.SetFloat("b", 2.5)
	require.Equal(t, []string{"a", "b"}, m.Keys())

	value, ok := m.GetFloatOK("a")
	require.Equal(t, 1.5, value)
	require.True(t, ok)

	value, ok = m.GetFloatOK("missing")
	require.Equal(t, float64(0), value)
	require.False(t, ok)

	require.Equal(t, 2.5, m.GetFloat("b"))

	require.True(t, m.Remove("a"))
	require.Equal(t, 1, m.Length())
}

func TestFloat_Equal(t *testing.T) {

	a := Float{"x": 1.1, "y": 2.2}
	b := Float{"x": 1.1, "y": 2.2}
	c := Float{"x": 9.9}

	require.True(t, a.Equal(b))
	require.False(t, a.NotEqual(b))
	require.False(t, a.Equal(c))
	require.True(t, a.NotEqual(c))
}

func TestFloat_NilMap(t *testing.T) {
	var m Float
	require.True(t, m.SetFloat("key", 1.5))
	require.Equal(t, 1, m.Length())

	var n Float
	require.True(t, n.Remove("key"))
}
