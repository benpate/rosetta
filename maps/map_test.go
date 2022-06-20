package maps

import (
	"testing"

	"github.com/benpate/rosetta/path"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	m := New()
	require.Equal(t, map[string]any{}, m.AsMapOfInterface())
}

func TestMapGet(t *testing.T) {

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

func TestMapSet(t *testing.T) {

	m := Map{}

	m.SetBool("bool", true)
	require.True(t, m["bool"].(bool))
	require.True(t, m.GetBool("bool"))

	m.SetInt("int", 42)
	require.Equal(t, 42, m.GetInt("int"))

	m.SetFloat("float", 42.69)
	require.Equal(t, 42.69, m.GetFloat("float"))

	m.SetString("string", "John Doe")
	require.Equal(t, "John Doe", m.GetString("string"))
}

func TestMapSetPath(t *testing.T) {

	m := Map{"stringValue": "kenobi", "intValue": 69, "boolValue": true}

	path.Set(m, "stringValue", "General Kenobi")
	require.Equal(t, "General Kenobi", m["stringValue"])
	require.Equal(t, "General Kenobi", m.GetString("stringValue"))

	path.Set(m, "intValue", 420)
	require.Equal(t, 420, m["intValue"])
	require.Equal(t, 420, m.GetInt("intValue"))

	// Add a slice of strings to the map
	path.Set(m, "slice", []string{"hello", "there", "general", "kenobi"})

	{
		value, ok := path.GetOK(m, "slice.0")
		require.True(t, ok)
		require.Equal(t, "hello", value)
	}

	{
		value, ok := path.GetOK(m, "slice.1")
		require.True(t, ok)
		require.Equal(t, "there", value)
	}

	{
		value, ok := path.GetOK(m, "slice.2")
		require.True(t, ok)
		require.Equal(t, "general", value)
	}

	{
		value, ok := path.GetOK(m, "slice.3")
		require.True(t, ok)
		require.Equal(t, "kenobi", value)
	}
}

func TestMapDelete(t *testing.T) {

	{
		m := Map{"hello": Map{"there": Map{"general": "kenobi"}}}
		path.Delete(m, "hello.there.general")

		hello := m["hello"].(Map)
		there := hello["there"].(Map)
		require.Nil(t, there["general"])
	}
	{
		m := Map{"hello": map[string]interface{}{
			"the":    1337,
			"answer": 69,
			"is":     42,
		}}

		path.Delete(m, "hello.answer")

		hello := m["hello"].(map[string]interface{})
		require.Equal(t, 1337, hello["the"])
		require.Equal(t, 42, hello["is"])
		require.Nil(t, hello["answer"])
	}
}
