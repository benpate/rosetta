package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapOfAny_FromTypedMaps(t *testing.T) {

	{
		result, ok := MapOfAny(Any{"a": 1})
		require.True(t, ok)
		require.Equal(t, 1, result["a"])
	}
	{
		input := Any{"a": 1}
		result, ok := MapOfAny(&input)
		require.True(t, ok)
		require.Equal(t, 1, result["a"])
	}
	{
		result, ok := MapOfAny(Bool{"a": true})
		require.True(t, ok)
		require.Equal(t, true, result["a"])
	}
	{
		result, ok := MapOfAny(Float{"a": 1.5})
		require.True(t, ok)
		require.Equal(t, 1.5, result["a"])
	}
	{
		result, ok := MapOfAny(Int{"a": 7})
		require.True(t, ok)
		require.Equal(t, 7, result["a"])
	}
	{
		result, ok := MapOfAny(Int64{"a": int64(7)})
		require.True(t, ok)
		require.Equal(t, int64(7), result["a"])
	}
	{
		result, ok := MapOfAny(String{"a": "hello"})
		require.True(t, ok)
		require.Equal(t, "hello", result["a"])
	}
}

func TestMapOfAny_FromPlainMaps(t *testing.T) {

	{
		result, ok := MapOfAny(map[string]any{"a": 1})
		require.True(t, ok)
		require.Equal(t, 1, result["a"])
	}
	{
		result, ok := MapOfAny(map[string]bool{"a": true})
		require.True(t, ok)
		require.Equal(t, true, result["a"])
	}
	{
		result, ok := MapOfAny(map[string]float64{"a": 1.5})
		require.True(t, ok)
		require.Equal(t, 1.5, result["a"])
	}
	{
		result, ok := MapOfAny(map[string]int{"a": 7})
		require.True(t, ok)
		require.Equal(t, 7, result["a"])
	}
	{
		result, ok := MapOfAny(map[string]int64{"a": int64(7)})
		require.True(t, ok)
		require.Equal(t, int64(7), result["a"])
	}
	{
		result, ok := MapOfAny(map[string]string{"a": "hello"})
		require.True(t, ok)
		require.Equal(t, "hello", result["a"])
	}
}

func TestMapOfAny_FromPointerMaps(t *testing.T) {

	{
		input := map[string]any{"a": 1}
		result, ok := MapOfAny(&input)
		require.True(t, ok)
		require.Equal(t, 1, result["a"])
	}
	{
		input := map[string]bool{"a": true}
		result, ok := MapOfAny(&input)
		require.True(t, ok)
		require.Equal(t, true, result["a"])
	}
	{
		input := map[string]float64{"a": 1.5}
		result, ok := MapOfAny(&input)
		require.True(t, ok)
		require.Equal(t, 1.5, result["a"])
	}
	{
		input := map[string]int{"a": 7}
		result, ok := MapOfAny(&input)
		require.True(t, ok)
		require.Equal(t, 7, result["a"])
	}
	{
		input := map[string]int64{"a": int64(7)}
		result, ok := MapOfAny(&input)
		require.True(t, ok)
		require.Equal(t, int64(7), result["a"])
	}
	{
		input := map[string]string{"a": "hello"}
		result, ok := MapOfAny(&input)
		require.True(t, ok)
		require.Equal(t, "hello", result["a"])
	}
}

func TestMapOfAny_Unsupported(t *testing.T) {
	result, ok := MapOfAny("not a map")
	require.False(t, ok)
	require.Nil(t, result)
}
