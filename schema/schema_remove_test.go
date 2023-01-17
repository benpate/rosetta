package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemove_A(t *testing.T) {

	array := newTestArrayA()
	schema := New(testArrayA_Schema())

	require.Equal(t, 3, array.Length())
	require.Equal(t, "one", array[0])
	require.Equal(t, "two", array[1])
	require.Equal(t, "three", array[2])

	require.True(t, schema.Remove(&array, "1"))
	require.Equal(t, 2, array.Length())
	require.Equal(t, "one", array[0])
	require.Equal(t, "three", array[1])
}

func TestRemove_B(t *testing.T) {

	value := newTestStructA()
	schema := New(testStructA_Schema())

	require.Equal(t, 3, value.Array.Length())
	require.Equal(t, "one", value.Array[0])
	require.Equal(t, "two", value.Array[1])
	require.Equal(t, "three", value.Array[2])

	require.True(t, schema.Remove(&value, "array.1"))
	require.Equal(t, 2, value.Array.Length())
	require.Equal(t, "one", value.Array[0])
	require.Equal(t, "three", value.Array[1])
}

func TestRemove_Errors(t *testing.T) {

	value := newTestStructA()
	schema := New(testStructA_Schema())

	require.False(t, schema.Remove(&value, "array.other.wrong"))
	require.False(t, schema.Remove(&value, "still-wrong"))
}
