package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArray_GetSet(t *testing.T) {
	value := make(testArrayA, 0)
	schema := New(testArrayA_Schema())

	{
		result, err := schema.Get(value, "0")
		require.NotNil(t, err)
		require.Equal(t, nil, result)
	}

	{
		require.Nil(t, schema.Set(&value, "0", "one"))
		result, err := schema.Get(value, "0")
		require.Nil(t, err)
		require.Equal(t, "one", result)
	}

	{
		require.Nil(t, schema.Set(&value, "1", "two"))
		result, err := schema.Get(value, "1")
		require.Nil(t, err)
		require.Equal(t, "two", result)
	}

	{
		require.NotNil(t, schema.Set(&value, "10", "out of bounds"))
	}
}

func TestArray_Validation(t *testing.T) {

	s := New(Array{
		Items: String{MaxLength: 10},
	})

	{
		v := testArrayA{"one", "two", "three", "valid"}
		require.Nil(t, s.Validate(&v))
	}

	{
		v := testArrayA{"one", "two", "three", "invalid because its way too long"}
		require.NotNil(t, s.Validate(&v))
	}

	{
		err := s.Validate(17)
		require.NotNil(t, err)
	}
}
