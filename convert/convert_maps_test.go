package convert

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringOk_AllTypes(t *testing.T) {

	assert := func(value any, expected string) {
		result, _ := StringOk(value, "DEFAULT")
		require.Equal(t, expected, result, "input: %#v", value)
	}

	assert(nil, "DEFAULT")
	assert(true, "true")
	assert(false, "false")
	assert([]byte("bytes"), "bytes")
	assert(int(5), "5")
	assert(int8(5), "5")
	assert(int16(5), "5")
	assert(int32(5), "5")
	assert(int64(5), "5")
	assert(float32(5), "5")
	assert(float64(5), "5")
	assert("hello", "hello")
	assert([]string{}, "DEFAULT")
	assert([]string{"a"}, "a")
	assert([]string{"a", "b"}, "a")
	assert([]any{}, "DEFAULT")
	assert([]any{"x"}, "x")
	assert([]any{"x", "y"}, "x")
	assert(customInter(7), "7")
	assert(customFloater(7), "7")
	assert(customHexer("ff"), "ff")
	assert(customStringer("custom"), "custom")
}

func TestStringOk_Reader(t *testing.T) {
	result, ok := StringOk(strings.NewReader("from reader"), "")
	require.True(t, ok)
	require.Equal(t, "from reader", result)
}

func TestMapOfIntOk_AllTypes(t *testing.T) {

	{
		result, ok := MapOfIntOk(map[string]int{"a": 1})
		require.True(t, ok)
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		result, _ := MapOfIntOk(map[string]int8{"a": 1})
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		result, _ := MapOfIntOk(map[string]int16{"a": 1})
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		result, _ := MapOfIntOk(map[string]int32{"a": 1})
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		result, _ := MapOfIntOk(map[string]int64{"a": 1})
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		result, _ := MapOfIntOk(map[string]any{"a": 1})
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		result, _ := MapOfIntOk(map[string]string{"a": "1"})
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		// Reflection fallback for an unrecognized map type
		result, ok := MapOfIntOk(map[string]float64{"a": 1})
		require.True(t, ok)
		require.Equal(t, map[string]int{"a": 1}, result)
	}
	{
		// Non-map input fails
		_, ok := MapOfIntOk("not a map")
		require.False(t, ok)
	}
	{
		_, ok := MapOfIntOk(nil)
		require.False(t, ok)
	}
}

func TestMapOfInt32Ok_AllTypes(t *testing.T) {

	{
		result, ok := MapOfInt32Ok(map[string]int32{"a": 1})
		require.True(t, ok)
		require.Equal(t, map[string]int32{"a": 1}, result)
	}
	{
		result, _ := MapOfInt32Ok(map[string]int{"a": 1})
		require.Equal(t, map[string]int32{"a": 1}, result)
	}
	{
		result, _ := MapOfInt32Ok(map[string]string{"a": "1"})
		require.Equal(t, map[string]int32{"a": 1}, result)
	}
	{
		_, ok := MapOfInt32Ok("not a map")
		require.False(t, ok)
	}
}

func TestMapOfInt_Wrapper(t *testing.T) {
	require.Equal(t, map[string]int{"a": 1}, MapOfInt(map[string]int{"a": 1}))
	// Non-convertible input yields an empty map
	require.Equal(t, map[string]int{}, MapOfInt("nope"))
}
