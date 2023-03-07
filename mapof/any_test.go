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

	require.NotNil(t, s.Set(&v, "foo.bar", "baz"))
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
