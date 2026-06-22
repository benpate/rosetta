package mapof

import (
	"testing"
	"time"

	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestAny_Manipulations(t *testing.T) {

	m := NewAny()
	require.Equal(t, 0, m.Length())
	require.True(t, m.IsEmpty())
	require.False(t, m.NotEmpty())

	m.SetString("b", "two")
	m.SetString("a", "one")

	require.Equal(t, 2, m.Length())
	require.False(t, m.IsEmpty())
	require.True(t, m.NotEmpty())
	require.Equal(t, []string{"a", "b"}, m.Keys())
}

func TestAny_NotEqual(t *testing.T) {
	a := Any{"x": 1}
	b := Any{"x": 2}
	require.True(t, a.NotEqual(b))
	require.False(t, a.NotEqual(Any{"x": 1}))
}

func TestAny_Getters(t *testing.T) {

	m := Any{
		"flag":   true,
		"pi":     3.14,
		"count":  42,
		"big":    int64(99),
		"name":   "rosetta",
		"when":   "2026-01-02T00:00:00Z",
		"absent": nil,
	}

	require.True(t, m.GetBool("flag"))
	require.Equal(t, 3.14, m.GetFloat("pi"))
	require.Equal(t, 42, m.GetInt("count"))
	require.Equal(t, int64(99), m.GetInt64("big"))
	require.Equal(t, "rosetta", m.GetString("name"))

	// OK variants for existing and missing keys
	b, ok := m.GetBoolOK("flag")
	require.True(t, b)
	require.True(t, ok)

	_, ok = m.GetBoolOK("missing")
	require.False(t, ok)

	f, ok := m.GetFloatOK("pi")
	require.Equal(t, 3.14, f)
	require.True(t, ok)

	_, ok = m.GetFloatOK("missing")
	require.False(t, ok)

	when, ok := m.GetTimeOK("when")
	require.True(t, ok)
	require.Equal(t, 2026, when.Year())

	require.Equal(t, when, m.GetTime("when"))

	_, ok = m.GetTimeOK("missing")
	require.False(t, ok)
}

func TestAny_SetBoolFloat(t *testing.T) {

	m := NewAny()

	require.True(t, m.SetBool("flag", true))
	require.Equal(t, true, m.GetAny("flag"))

	// SetBool stores even false values
	require.True(t, m.SetBool("flag", false))
	require.Equal(t, false, m.GetAny("flag"))

	require.True(t, m.SetFloat("pi", 3.14))
	require.Equal(t, 3.14, m.GetFloat("pi"))

	// SetFloat deletes on zero
	require.True(t, m.SetFloat("pi", 0))
	_, ok := m.GetAnyOK("pi")
	require.False(t, ok)
}

func TestAny_Append(t *testing.T) {

	m := NewAny()

	m.Append("list", "one")
	require.Equal(t, []any{"one"}, m.GetAny("list"))

	m.Append("list", "two")
	require.Equal(t, []any{"one", "two"}, m.GetAny("list"))

	// Appending to a non-slice value wraps both into a slice
	m.SetString("scalar", "x")
	m.Append("scalar", "y")
	require.Equal(t, []any{"x", "y"}, m.GetAny("scalar"))
}

func TestAny_GetPointerAndRemove(t *testing.T) {

	m := Any{"key": "value"}

	value, ok := m.GetPointer("key")
	require.True(t, ok)
	require.Equal(t, "value", value)

	_, ok = m.GetPointer("missing")
	require.False(t, ok)

	require.True(t, m.Remove("key"))
	require.Equal(t, 0, m.Length())
}

func TestAny_SetValue(t *testing.T) {

	m := NewAny()
	require.NoError(t, m.SetValue(map[string]any{"a": 1}))
	require.Equal(t, 1, m.GetInt("a"))

	// Unsupported value returns an error
	require.Error(t, m.SetValue("not a map"))
}

func TestAny_Slices(t *testing.T) {

	m := Any{
		"anys":    []any{"a", "b"},
		"strings": []string{"x", "y"},
		"ints":    []int{1, 2},
		"floats":  []float64{1.5, 2.5},
	}

	require.Equal(t, []any{"a", "b"}, m.GetSliceOfAny("anys"))
	require.Equal(t, []string{"x", "y"}, m.GetSliceOfString("strings"))
	require.Equal(t, []int{1, 2}, m.GetSliceOfInt("ints"))
	require.Equal(t, []float64{1.5, 2.5}, m.GetSliceOfFloat("floats"))
}

func TestAny_GetMap(t *testing.T) {

	m := Any{
		"nested":      Any{"a": 1},
		"plain":       map[string]any{"b": 2},
		"strings":     String{"s": "value"},
		"plainString": map[string]string{"t": "value"},
	}

	require.Equal(t, 1, m.GetMap("nested").GetInt("a"))
	require.Equal(t, 2, m.GetMapOfAny("plain").GetInt("b"))
	require.Equal(t, "value", m.GetMapOfString("strings").GetString("s"))
	require.Equal(t, "value", m.GetMapOfString("plainString").GetString("t"))

	// Missing/invalid keys return empty maps
	require.Equal(t, 0, m.GetMapOfAny("missing").Length())
	require.Equal(t, 0, m.GetMapOfString("missing").Length())
}

func TestAny_GetSliceOfMap(t *testing.T) {

	{
		m := Any{"value": []Any{{"k": "v"}}}
		result := m.GetSliceOfMap("value")
		require.Equal(t, 1, len(result))
		require.Equal(t, "v", result[0].GetString("k"))
	}
	{
		single := Any{"k": "v"}
		m := Any{"value": single}
		result := m.GetSliceOfMap("value")
		require.Equal(t, 1, len(result))
	}
	{
		m := Any{}
		require.Equal(t, 0, len(m.GetSliceOfMap("missing")))
	}
}

func TestAny_MapOfAnyAndString(t *testing.T) {

	m := Any{"a": "one", "b": "two"}

	require.Equal(t, map[string]any(m), m.MapOfAny())

	asStrings := m.MapOfString()
	require.Equal(t, "one", asStrings["a"])
	require.Equal(t, "two", asStrings["b"])
}

func TestAny_SetObject(t *testing.T) {

	s := schema.Object{
		Properties: schema.ElementMap{
			"outer": schema.Object{
				Properties: schema.ElementMap{
					"inner": schema.String{},
				},
			},
		},
	}

	m := NewAny()

	// Empty path is an error
	require.Error(t, m.SetObject(s, list.ByDot(""), "value"))

	// Single-segment path sets directly
	require.NoError(t, m.SetObject(s, list.ByDot("outer"), Any{}))

	// Nested path sets through a child map
	require.NoError(t, m.SetObject(s, list.ByDot("outer.inner"), "value"))
	require.Equal(t, "value", m.GetMapOfAny("outer").GetString("inner"))

	// Unknown property is an error
	require.Error(t, m.SetObject(s, list.ByDot("bogus.inner"), "value"))
}

func TestAny_GetTimeType(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	m := Any{"now": now}
	require.Equal(t, now, m.GetTime("now"))
}
