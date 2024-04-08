package mapof

import (
	"testing"

	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {

	s := schema.New(schema.Object{
		Wildcard: schema.Object{
			Wildcard: schema.String{},
		},
	})

	v := NewAny()

	require.Nil(t, s.Set(&v, "foo.bar", "baz"))

	value, err := s.Get(&v, "foo.bar")
	require.Nil(t, err)
	require.Equal(t, "baz", value)
}

func TestAny_Any(t *testing.T) {

	s := schema.New(schema.Object{
		Wildcard: schema.Any{},
	})

	v := NewAny()

	require.Nil(t, testTable(s, &v, []testTableItem{
		{"foo.bar", "baz"},
		{"itsy.bitsy.spider.went", "up the water spout"},
		{"string", "value"},
		{"int", int(42)},
		{"int64", int64(42)},
		{"float64", float64(42)},
		{"bool", true},
	}))
}

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
	require.Equal(t, int64(69), m.GetAny("key"))
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
