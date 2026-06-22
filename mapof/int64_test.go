package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt64(t *testing.T) {

	m := NewInt64()

	require.Equal(t, 0, m.Length())
	require.True(t, m.IsEmpty())
	require.False(t, m.NotEmpty())

	require.True(t, m.SetInt64("five", 5))
	require.Equal(t, 1, m.Length())
	require.True(t, m.NotEmpty())

	// Setting to zero deletes the key
	require.True(t, m.SetInt64("five", 0))
	require.Equal(t, 0, m.Length())

	m.SetInt64("a", 1)
	m.SetInt64("b", 2)
	require.Equal(t, []string{"a", "b"}, m.Keys())

	value, ok := m.GetInt64OK("a")
	require.Equal(t, int64(1), value)
	require.True(t, ok)

	value, ok = m.GetInt64OK("missing")
	require.Equal(t, int64(0), value)
	require.False(t, ok)

	require.Equal(t, int64(2), m.GetInt64("b"))

	require.True(t, m.Remove("a"))
	require.Equal(t, 1, m.Length())
}

func TestInt64_Equal(t *testing.T) {

	a := Int64{"x": 1, "y": 2}
	b := Int64{"x": 1, "y": 2}
	c := Int64{"x": 9}

	require.True(t, a.Equal(b))
	require.False(t, a.NotEqual(b))
	require.False(t, a.Equal(c))
	require.True(t, a.NotEqual(c))
}

func TestInt64_NilMap(t *testing.T) {
	var m Int64
	require.True(t, m.SetInt64("key", 1))
	require.Equal(t, 1, m.Length())

	var n Int64
	require.True(t, n.Remove("key"))
}
