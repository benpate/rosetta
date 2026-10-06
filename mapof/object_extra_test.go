package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObject_Manipulations(t *testing.T) {

	m := NewObject[string]()
	require.Equal(t, 0, m.Length())
	require.True(t, m.IsEmpty())
	require.False(t, m.NotEmpty())

	m["b"] = "two"
	m["a"] = "one"

	require.Equal(t, 2, m.Length())
	require.False(t, m.IsEmpty())
	require.True(t, m.NotEmpty())
	require.Equal(t, []string{"a", "b"}, m.Keys())
}

func TestObject_GetPointer(t *testing.T) {

	m := Object[string]{"key": "value"}

	value, ok := m.GetPointer("key")
	require.True(t, ok)
	require.Equal(t, "value", value)

	_, ok = m.GetPointer("missing")
	require.False(t, ok)
}

func TestObject_Remove(t *testing.T) {

	m := Object[string]{"a": "one", "b": "two"}

	require.True(t, m.Remove("a"))
	require.Equal(t, 1, m.Length())

	// Remove on a nil map should not panic
	var n Object[string]
	require.True(t, n.Remove("key"))
	require.NotNil(t, n)
}

// TestObject_SetKey requires that SetKey stores only a value of type T
func TestObject_SetKey(t *testing.T) {

	t.Run("a value of type T is stored", func(t *testing.T) {
		var object Object[Any]
		require.NoError(t, object.SetKey("child", Any{"name": "v"}))
		require.Equal(t, Object[Any]{"child": {"name": "v"}}, object)
	})

	t.Run("a value of another type is an error, and nothing is stored", func(t *testing.T) {
		object := Object[Any]{"child": {"keep": "k"}}
		require.Error(t, object.SetKey("child", map[string]any{"name": "v"}))
		require.Error(t, object.SetKey("other", "v"))
		require.Error(t, object.SetKey("other", nil))
		require.Equal(t, Object[Any]{"child": {"keep": "k"}}, object)
	})
}
