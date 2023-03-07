package mapof

import (
	"testing"

	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {

	value := NewObject[String]()

	s := schema.New(schema.Object{
		Properties: schema.ElementMap{
			"key1": schema.Object{
				Properties: schema.ElementMap{
					"subkey1": schema.String{},
					"subkey2": schema.String{},
				},
			},
			"key2": schema.Object{
				Properties: schema.ElementMap{
					"subkey1": schema.String{},
					"subkey2": schema.String{},
				},
			},
		},
	})

	require.Nil(t, s.Set(&value, "key1.subkey1", "subvalue.1.1"))
	require.Nil(t, s.Set(&value, "key1.subkey2", "subvalue.1.2"))
	require.Nil(t, s.Set(&value, "key2.subkey1", "subvalue.2.1"))
	require.Nil(t, s.Set(&value, "key2.subkey2", "subvalue.2.2"))

	{
		value, err := s.Get(&value, "key1.subkey1")
		require.Nil(t, err)
		require.Equal(t, "subvalue.1.1", value)
	}

	{
		value, err := s.Get(&value, "key1.subkey2")
		require.Nil(t, err)
		require.Equal(t, "subvalue.1.2", value)
	}

	{
		value, err := s.Get(&value, "key2.subkey1")
		require.Nil(t, err)
		require.Equal(t, "subvalue.2.1", value)
	}

	{
		value, err := s.Get(&value, "key2.subkey2")
		require.Nil(t, err)
		require.Equal(t, "subvalue.2.2", value)
	}
}

func TestObjectConversion(t *testing.T) {
	var value map[string][]string = make(Object[[]string])
	t.Log(value)
}
