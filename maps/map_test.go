package maps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	m := New()
	require.Equal(t, map[string]any{}, m.AsMapOfInterface())
}

func TestMap_Get(t *testing.T) {

	m := Map{"hello": "there", "general": "kenobi", "intValue": 69, "boolValue": true, "floatValue": 42.1}

	// stringValues
	{
		result, ok := m.GetString("boolValue")
		require.True(t, ok)
		require.Equal(t, "true", result)

		result, ok = m.GetString("intValue")
		require.True(t, ok)
		require.Equal(t, "69", result)

		result, ok = m.GetString("floatValue")
		require.True(t, ok)
		require.Equal(t, "42.1", result)

		result, ok = m.GetString("hello")
		require.True(t, ok)
		require.Equal(t, "there", result)

		result, ok = m.GetString("general")
		require.True(t, ok)
		require.Equal(t, "kenobi", result)

		result, ok = m.GetString("missing")
		require.False(t, ok)
		require.Equal(t, "", result)
	}

	// boolValues
	{
		result, ok := m.GetBool("boolValue")
		require.True(t, ok)
		require.True(t, result)

		result, ok = m.GetBool("intValue")
		require.True(t, ok)
		require.True(t, result)

		result, ok = m.GetBool("floatValue")
		require.True(t, ok)
		require.True(t, result)

		result, ok = m.GetBool("hello")
		require.True(t, ok)
		require.False(t, result)

		result, ok = m.GetBool("missing")
		require.False(t, ok)
		require.False(t, result)
	}

	// intValues
	{
		result, ok := m.GetInt("boolValue")
		require.True(t, ok)
		require.Equal(t, 1, result)

		result, ok = m.GetInt("intValue")
		require.True(t, ok)
		require.Equal(t, 69, result)

		result, ok = m.GetInt("floatValue")
		require.True(t, ok)
		require.Equal(t, 42, result)

		result, ok = m.GetInt("hello")
		require.True(t, ok)
		require.Zero(t, result)

		result, ok = m.GetInt("missing")
		require.False(t, ok)
		require.Zero(t, result)
	}

	// floatValues
	{
		result, ok := m.GetFloat("boolValue")
		require.True(t, ok)
		require.Equal(t, float64(1), result)

		result, ok = m.GetFloat("intValue")
		require.True(t, ok)
		require.Equal(t, float64(69), result)

		result, ok = m.GetFloat("floatValue")
		require.True(t, ok)
		require.Equal(t, float64(42.1), result)

		result, ok = m.GetFloat("hello")
		require.False(t, ok)
		require.Zero(t, result)

		result, ok = m.GetFloat("missing")
		require.False(t, ok)
		require.Zero(t, result)
	}

	// interfaceValues
	{
		result, ok := m.GetInterface("hello")
		require.True(t, ok)
		require.Equal(t, "there", result)

		result, ok = m.GetInterface("intValue")
		require.True(t, ok)
		require.Equal(t, 69, result)

		result, ok = m.GetInterface("floatValue")
		require.True(t, ok)
		require.Equal(t, float64(42.1), result)

		result, ok = m.GetInterface("boolValue")
		require.True(t, ok)
		require.Equal(t, true, result)

		result, ok = m.GetInterface("mising")
		require.False(t, ok)
		require.Nil(t, result)
	}
}
