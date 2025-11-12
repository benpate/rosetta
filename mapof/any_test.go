package mapof

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAny_Zero(t *testing.T) {
	var value = NewAny()
	require.Equal(t, nil, value["this-is-zero"])
	require.True(t, value.IsZeroValue("this-is-zero"))

	value["exists"] = "true"
	require.False(t, value.IsZeroValue("exists"))
}

func TestAny_EmptyString(t *testing.T) {

	m := NewAny()

	m.SetAny("key", "any")
	require.Equal(t, "any", m.GetString("key"))
	require.Equal(t, "any", m.GetAny("key"))
	require.Equal(t, 1, len(m))

	m.SetAny("key", "")
	require.Equal(t, "", m.GetString("key"))
	require.Equal(t, nil, m.GetAny("key"))
	require.Zero(t, len(m))

	m.SetString("key", "string")
	require.Equal(t, "string", m.GetString("key"))
	require.Equal(t, "string", m.GetAny("key"))
	require.Equal(t, 1, len(m))

	m.SetString("key", "")
	require.Equal(t, "", m.GetString("key"))
	require.Equal(t, nil, m.GetAny("key"))
	require.Zero(t, len(m))
}

func TestAny_EmptyInt(t *testing.T) {

	m := NewAny()

	m.SetAny("key", 69)
	require.Equal(t, 69, m.GetInt("key"))
	require.Equal(t, 69, m.GetAny("key"))
	require.Equal(t, 1, len(m))

	m.SetAny("key", 0)
	require.Equal(t, int64(0), m.GetInt64("key"))
	require.Equal(t, nil, m.GetAny("key"))
	require.Zero(t, len(m))

	m.SetInt("key", 42)
	require.Equal(t, 42, m.GetInt("key"))
	require.Equal(t, 42, m.GetAny("key"))
	require.Equal(t, 1, len(m))

	m.SetInt("key", 0)
	require.Equal(t, 0, m.GetInt("key"))
	require.Equal(t, nil, m.GetAny("key"))
	require.Zero(t, len(m))
}

func TestAny_EmptyInt64(t *testing.T) {

	m := NewAny()

	m.SetAny("key", 69)
	require.Equal(t, int64(69), m.GetInt64("key"))
	require.Equal(t, 69, m.GetAny("key"))
	require.Equal(t, 1, len(m))

	m.SetAny("key", int64(0))
	require.Equal(t, int64(0), m.GetInt64("key"))
	require.Equal(t, nil, m.GetAny("key"))
	require.Zero(t, len(m))

	m.SetInt64("key", 42)
	require.Equal(t, int64(42), m.GetInt64("key"))
	require.Equal(t, int64(42), m.GetAny("key"))
	require.Equal(t, 1, len(m))

	m.SetInt64("key", 0)
	require.Equal(t, int64(0), m.GetInt64("key"))
	require.Equal(t, nil, m.GetAny("key"))
	require.Zero(t, len(m))
}

func TestSliceOfPlainMap(t *testing.T) {

	data := Any{
		"value": []Any{
			{"key": "valueA"},
			{"key": "valueB"},
			{"key": "valueC"},
		},
	}

	plainMap := data.GetSliceOfPlainMap("value")

	require.Equal(t, 3, len(plainMap))
	require.Equal(t, "valueA", plainMap[0]["key"])
	require.Equal(t, "valueB", plainMap[1]["key"])
	require.Equal(t, "valueC", plainMap[2]["key"])
}

func TestEqual(t *testing.T) {

	a := Any{
		"a": "value",
		"b": []string{"one", "two", "three"},
		"c": Any{
			"deeply": "valued",
		},
	}

	b := Any{
		"a": "value",
		"b": []string{"one", "two", "three"},
		"c": Any{
			"deeply": "valued",
		},
	}

	require.True(t, reflect.DeepEqual(a, b))
	require.True(t, a.Equal(b))
}
