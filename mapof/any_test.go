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
