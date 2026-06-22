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

	// Setting a zero value deletes the key (keeping the map sparse), so the map stays empty.
	require.True(t, m.SetInt("zero", 0))
	require.Equal(t, 0, m.Length())

	// Setting a non-zero value stores it.
	require.True(t, m.SetInt("nonzero", 5))
	require.Equal(t, 1, m.Length())

	value, ok := m.GetIntOK("nonzero")
	require.Equal(t, 5, value)
	require.True(t, ok)

	value, ok = m.GetIntOK("missing")
	require.Equal(t, 0, value)
	require.False(t, ok)

	require.Equal(t, 5, m.GetInt("nonzero"))
	require.Equal(t, []string{"nonzero"}, m.Keys())

	require.True(t, m.Remove("nonzero"))
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
	// Setting a non-zero value on a nil map allocates it.
	var m Int
	require.True(t, m.SetInt("key", 5))
	require.Equal(t, 1, m.Length())

	// Remove on a nil map should not panic.
	var n Int
	require.True(t, n.Remove("key"))
}
