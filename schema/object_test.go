package schema

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestObjectA_Get(t *testing.T) {

	object := newTestStructA()
	schema := New(testStructA_Schema())

	{
		value, err := schema.Get(object, "name")
		require.Nil(t, err)
		require.Equal(t, "John Connor", value)
	}

	{
		value, err := schema.Get(object, "latitude")
		require.Nil(t, err)
		require.Equal(t, 45.123456, value)
	}

	{
		value, err := schema.Get(object, "active")
		require.Nil(t, err)
		require.Equal(t, true, value)
	}

	{
		_, err := schema.Get(object, "age")
		require.NotNil(t, err)
	}

	{
		_, err := schema.Get(object, "published")
		require.NotNil(t, err)
	}
}

func TestObjectB_Get(t *testing.T) {

	object := newTestStructB()
	schema := New(testStructB_Schema())

	{
		value, err := schema.Get(object, "name")
		require.Nil(t, err)
		require.Equal(t, "John Connor", value)
	}

	{
		value, err := schema.Get(object, "age")
		require.Nil(t, err)
		require.Equal(t, 42, value)
	}

	{
		value, err := schema.Get(object, "published")
		require.Nil(t, err)
		require.Equal(t, int64(1234567890), value)
	}

	{
		_, err := schema.Get(object, "latitude")
		require.NotNil(t, err)
	}

	{
		_, err := schema.Get(object, "active")
		require.NotNil(t, err)
	}
}

func TestObject_Set(t *testing.T) {

	object := newTestStructA()
	schema := New(testStructA_Schema())

	{
		require.Nil(t, schema.Set(&object, "name", "Sarah Connor"))
		value, err := schema.Get(object, "name")
		require.Nil(t, err)
		require.Equal(t, "Sarah Connor", value)
	}

	{
		require.Nil(t, schema.Set(&object, "latitude", 42.424242))
		value, err := schema.Get(object, "latitude")
		require.Nil(t, err)
		require.Equal(t, 42.424242, value)
	}

	{
		require.Nil(t, schema.Set(&object, "active", false))
		value, err := schema.Get(object, "active")
		require.Nil(t, err)
		require.Equal(t, false, value)
	}
}

func TestObject_Create(t *testing.T) {

	object := newTestStructA()
	schema := New(testStructA_Schema())

	require.Nil(t, schema.Set(&object, "array.0", "Sarah Connor"))

	result, err := schema.Get(&object, "array.0")
	require.Nil(t, err)
	require.Equal(t, "Sarah Connor", result)
}

func TestObject_Wildcard(t *testing.T) {

	schema := New(Object{
		Wildcard: String{Format: "email", MinLength: 42},
	})

	element, ok := schema.GetElement("does-not-exist")

	require.True(t, ok)
	require.Equal(t, String{Format: "email", MinLength: 42}, element)
}

func TestObject_Validate(t *testing.T) {

	object := newTestStructA()
	schema := New(testStructA_Schema())

	{
		err := schema.Validate(&object)
		spew.Config.DisableMethods = true
		require.Nil(t, err)
	}
}
