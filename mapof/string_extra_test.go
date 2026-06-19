package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString_Manipulations(t *testing.T) {

	m := NewString()
	require.Equal(t, 0, m.Length())
	require.True(t, m.IsEmpty())
	require.False(t, m.NotEmpty())

	m.SetString("b", "two")
	m.SetString("a", "one")

	require.Equal(t, 2, m.Length())
	require.False(t, m.IsEmpty())
	require.True(t, m.NotEmpty())
	require.Equal(t, []string{"a", "b"}, m.Keys())

	require.True(t, m.Remove("a"))
	require.Equal(t, 1, m.Length())
}

func TestString_Equal(t *testing.T) {

	a := String{"x": "1", "y": "2"}

	require.True(t, a.Equal(String{"x": "1", "y": "2"}))
	require.False(t, a.NotEqual(String{"x": "1", "y": "2"}))

	// Different length
	require.False(t, a.Equal(String{"x": "1"}))
	require.True(t, a.NotEqual(String{"x": "1"}))

	// Different value
	require.False(t, a.Equal(String{"x": "1", "y": "9"}))
}

func TestString_GetStringOK(t *testing.T) {

	m := String{"key": "value"}

	value, ok := m.GetStringOK("key")
	require.Equal(t, "value", value)
	require.True(t, ok)

	value, ok = m.GetStringOK("missing")
	require.Equal(t, "", value)
	require.False(t, ok)
}

func TestString_MapConversions(t *testing.T) {

	m := String{"a": "one", "b": "two"}

	asAny := m.MapOfAny()
	require.Equal(t, "one", asAny["a"])

	asString := m.MapOfString()
	require.Equal(t, "two", asString["b"])
}

func TestString_NilMap(t *testing.T) {
	var m String
	require.True(t, m.SetString("key", "value"))
	require.Equal(t, 1, m.Length())

	var n String
	require.True(t, n.Remove("key"))
}
