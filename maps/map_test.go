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
	require.Equal(t, "true", m.GetString("boolValue"))
	require.Equal(t, "69", m.GetString("intValue"))
	require.Equal(t, "42.1", m.GetString("floatValue"))
	require.Equal(t, "there", m.GetString("hello"))
	require.Equal(t, "kenobi", m.GetString("general"))
	require.Equal(t, "", m.GetString("missing"))

	// boolValues
	require.True(t, m.GetBool("boolValue"))
	require.True(t, m.GetBool("intValue"))
	require.True(t, m.GetBool("floatValue"))
	require.False(t, m.GetBool("hello"))
	require.False(t, m.GetBool("missing"))

	// intValues
	require.Equal(t, 1, m.GetInt("boolValue"))
	require.Equal(t, 69, m.GetInt("intValue"))
	require.Equal(t, 42, m.GetInt("floatValue"))
	require.Zero(t, m.GetInt("hello"))
	require.Zero(t, m.GetInt("missing"))

	// floatValues
	require.Equal(t, float64(1), m.GetFloat("boolValue"))
	require.Equal(t, float64(69), m.GetFloat("intValue"))
	require.Equal(t, float64(42.1), m.GetFloat("floatValue"))
	require.Zero(t, m.GetInt("hello"))
	require.Zero(t, m.GetInt("missing"))

	// interfaceValues
	require.Equal(t, "there", m.GetInterface("hello"))
	require.Equal(t, 69, m.GetInterface("intValue"))
	require.Equal(t, float64(42.1), m.GetInterface("floatValue"))
	require.Equal(t, true, m.GetInterface("boolValue"))
	require.Nil(t, m.GetInterface("mising"))
}
