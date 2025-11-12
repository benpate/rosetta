package tests

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {

	s := schema.New(schema.Object{
		Wildcard: schema.Object{
			Wildcard: schema.String{},
		},
	})

	v := mapof.NewAny()

	require.Nil(t, s.Set(&v, "foo.bar", "baz"))

	value, err := s.Get(&v, "foo.bar")
	require.Nil(t, err)
	require.Equal(t, "baz", value)
}

func TestAny_Any(t *testing.T) {

	schema := schema.New(schema.Object{
		Wildcard: schema.Any{},
	})

	object := mapof.NewAny()

	run := func(key string, value any) {
		testInline(t, schema, &object, key, value)
	}

	run("foo.bar", "baz")
	run("itsy.bitsy.spider.went", "up the water spout")
	run("string", "value")
	run("int", int(42))
	run("int64", int64(42))
	run("float64", float64(42))
	run("bool", true)
}
