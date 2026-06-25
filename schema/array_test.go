package schema

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestArray_GetSet(t *testing.T) {
	value := make(testArrayA, 0)
	schema := New(testArrayA_Schema())

	{
		result, err := schema.Get(value, "0")
		require.Error(t, err)
		require.Equal(t, nil, result)
	}

	{
		require.NoError(t, schema.Set(&value, "0", "one"))
		result, err := schema.Get(value, "0")
		require.NoError(t, err)
		require.Equal(t, "one", result)
	}

	{
		require.NoError(t, schema.Set(&value, "1", "two"))
		result, err := schema.Get(value, "1")
		require.NoError(t, err)
		require.Equal(t, "two", result)
	}

	{
		require.Error(t, schema.Set(&value, "10", "out of bounds"))
	}
}

func TestArray_Validation(t *testing.T) {

	spew.Config.DisableMethods = true

	schema := New(Array{
		Items: String{MaxLength: 10},
	})

	{
		v := testArrayA{"one", "two", "three", "valid"}
		value, changed, err := Validate(schema, &v)
		require.NoError(t, err)
		require.False(t, changed)
		require.Equal(t, &testArrayA{"one", "two", "three", "valid"}, value)
	}

	{
		// An item that exceeds MaxLength would require rewriting, so Validate now
		// rejects it rather than truncating (Set still clamps the value in place).
		v := testArrayA{"one", "two", "three", "invalid because its way too long"}
		_, _, err := Validate(schema, &v)
		require.Error(t, err)
	}
}

func TestArray_Validation_Success(t *testing.T) {

	schema := New(Array{
		Items: String{MaxLength: 10},
	})

	newValue, changed, err := Validate(schema, &testArrayA{"one", "two", "three", "valid"})
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, &testArrayA{"one", "two", "three", "valid"}, newValue)
}

func TestArray_Validation_Fail1(t *testing.T) {

	schema := New(Array{
		Items: String{MaxLength: 10},
	})

	// This should fail because &testArrayA{} is not addressable
	_, _, err := Validate(schema, testArrayA{"one", "two", "three", "valid"})
	require.Error(t, err)
}
