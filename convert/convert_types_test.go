package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// --- Helper types implementing the convert interfaces ---

type customInter int

func (c customInter) Int() int { return int(c) }

type customInt64er int64

func (c customInt64er) Int64() int64 { return int64(c) }

type customFloater float64

func (c customFloater) Float() float64 { return float64(c) }

type customStringer string

func (c customStringer) String() string { return string(c) }

type customHexer string

func (c customHexer) Hex() string { return string(c) }

type customSliceOfStringer []string

func (c customSliceOfStringer) SliceOfString() []string { return []string(c) }

// --- Int64Ok across every supported input type ---

func TestInt64Ok_AllTypes(t *testing.T) {

	assert := func(value any, expected int64) {
		result, _ := Int64Ok(value, -1)
		require.Equal(t, expected, result, "input: %#v", value)
	}

	assert(nil, -1)
	assert(true, 1)
	assert(false, 0)
	assert(int(5), 5)
	assert(int8(5), 5)
	assert(int16(5), 5)
	assert(int32(5), 5)
	assert(int64(5), 5)
	assert(float32(5), 5)
	assert(float64(5), 5)
	assert("5", 5)
	assert("garbage", -1)
	assert([]string{}, -1)
	assert([]string{"7"}, 7)
	assert([]string{"7", "8"}, 7)
	assert([]any{}, -1)
	assert([]any{int64(9)}, 9)
	assert([]any{int64(9), int64(10)}, 9)
	assert(customInter(11), 11)
	assert(customFloater(12), 12)
	assert(customStringer("13"), 13)
	assert(customHexer("ff"), 255)
	assert(customHexer("xyz"), 0) // invalid hex
}

// --- IntOk / Int32Ok mirror Int64Ok's structure ---

func TestIntOk_AllTypes(t *testing.T) {

	assert := func(value any, expected int) {
		result, _ := IntOk(value, -1)
		require.Equal(t, expected, result, "input: %#v", value)
	}

	assert(nil, -1)
	assert(true, 1)
	assert(int8(5), 5)
	assert(int64(5), 5)
	assert(float64(5), 5)
	assert("5", 5)
	assert("garbage", -1)
	assert([]string{"7"}, 7)
	assert([]any{int(9)}, 9)
	assert(customInter(11), 11)
	assert(customFloater(12), 12)
	assert(customStringer("13"), 13)
}

func TestInt32Ok_AllTypes(t *testing.T) {

	assert := func(value any, expected int32) {
		result, _ := Int32Ok(value, -1)
		require.Equal(t, expected, result, "input: %#v", value)
	}

	assert(nil, -1)
	assert(true, 1)
	assert(int(5), 5)
	assert(int64(5), 5)
	assert(float64(5), 5)
	assert("5", 5)
	assert("garbage", -1)
	assert(customInter(11), 11)
	assert(customStringer("13"), 13)
}

// --- FloatOk across supported input types ---

func TestFloatOk_AllTypes(t *testing.T) {

	assert := func(value any, expected float64) {
		result, _ := FloatOk(value, -1)
		require.Equal(t, expected, result, "input: %#v", value)
	}

	assert(nil, -1)
	assert(true, 1)
	assert(false, 0)
	assert(int(5), 5)
	assert(int64(5), 5)
	assert(float32(5), 5)
	assert(float64(5), 5)
	assert("5.5", 5.5)
	assert("garbage", -1)
	assert([]string{"7.5"}, 7.5)
	assert([]any{float64(9.5)}, 9.5)
	assert(customInter(11), 11)
	assert(customFloater(12.5), 12.5)
	assert(customStringer("13.5"), 13.5)
}

// --- BoolOk across supported input types ---

func TestBoolOk_AllTypes(t *testing.T) {

	assertTrue := func(value any) {
		result, _ := BoolOk(value, false)
		require.True(t, result, "input: %#v", value)
	}

	assertTrue(true)
	assertTrue(1)
	assertTrue(int8(1))
	assertTrue(int64(1))
	assertTrue(float64(1))
	assertTrue("true")
	assertTrue([]string{"true"})
	assertTrue([]any{"true"})

	// false / default cases
	result, ok := BoolOk(nil, false)
	require.False(t, result)
	require.False(t, ok)

	result, _ = BoolOk("false", true)
	require.False(t, result)

	result, _ = BoolOk(customStringer("true"), false)
	require.True(t, result)
}

// --- StringOk across supported input types ---

func TestSliceOfStringOk_AllTypes(t *testing.T) {

	assert := func(value any, expected []string) {
		result, _ := SliceOfStringOk(value)
		require.Equal(t, expected, result, "input: %#v", value)
	}

	assert(nil, []string{})
	assert(true, []string{"true"})
	assert(float64(5), []string{"5.00"}) // floats format to two decimals
	assert(int(5), []string{"5"})
	assert(int64(5), []string{"5"})
	assert("hello", []string{"hello"})
	assert([]string{"a", "b"}, []string{"a", "b"})
	assert([]int{1, 2}, []string{"1", "2"})
	assert([]int64{1, 2}, []string{"1", "2"})
	assert([]float64{1, 2}, []string{"1.00", "2.00"})
	assert([]bool{true, false}, []string{"true", "false"})
	assert([]any{"a", 1}, []string{"a", "1"})
	assert(customInter(5), []string{"5"})
	assert(customInt64er(5), []string{"5"})
	assert(customFloater(5), []string{"5.00"})
	assert(customHexer("ff"), []string{"ff"})
	assert(customStringer("x"), []string{"x"})
	assert(customSliceOfStringer{"p", "q"}, []string{"p", "q"})
}

func TestSliceOfIntOk_AllTypes(t *testing.T) {

	{
		result, _ := SliceOfIntOk([]int{1, 2, 3})
		require.Equal(t, []int{1, 2, 3}, result)
	}
	{
		result, _ := SliceOfIntOk([]string{"4", "5"})
		require.Equal(t, []int{4, 5}, result)
	}
	{
		result, _ := SliceOfIntOk(int(7))
		require.Equal(t, []int{7}, result)
	}
	{
		result, _ := SliceOfIntOk([]any{1, 2})
		require.Equal(t, []int{1, 2}, result)
	}
}

func TestSliceOfInt64Ok_AllTypes(t *testing.T) {

	{
		result, _ := SliceOfInt64Ok([]int64{1, 2, 3})
		require.Equal(t, []int64{1, 2, 3}, result)
	}
	{
		result, _ := SliceOfInt64Ok([]string{"4", "5"})
		require.Equal(t, []int64{4, 5}, result)
	}
	{
		result, _ := SliceOfInt64Ok(int64(7))
		require.Equal(t, []int64{7}, result)
	}
}

func TestSliceOfFloatOk_AllTypes(t *testing.T) {

	{
		result, _ := SliceOfFloatOk([]float64{1.5, 2.5})
		require.Equal(t, []float64{1.5, 2.5}, result)
	}
	{
		result, _ := SliceOfFloatOk([]string{"4.5", "5.5"})
		require.Equal(t, []float64{4.5, 5.5}, result)
	}
	{
		result, _ := SliceOfFloatOk(float64(7.5))
		require.Equal(t, []float64{7.5}, result)
	}
}

func TestSliceOfAnyOk_AllTypes(t *testing.T) {

	{
		result, ok := SliceOfAnyOk([]bool{true, false})
		require.True(t, ok)
		require.Equal(t, []any{true, false}, result)
	}
	{
		result, ok := SliceOfAnyOk([]int64{1, 2})
		require.True(t, ok)
		require.Equal(t, []any{int64(1), int64(2)}, result)
	}
	{
		result, ok := SliceOfAnyOk(customStringer("hello"))
		require.True(t, ok)
		require.Equal(t, []any{"hello"}, result)
	}
	{
		// Reflection-based fallback for a typed slice
		result, ok := SliceOfAnyOk([]uint{1, 2, 3})
		require.True(t, ok)
		require.Equal(t, 3, len(result))
	}
}

func TestBaseTypeOK(t *testing.T) {

	{
		value, ok := BaseTypeOK(42)
		require.True(t, ok)
		require.Equal(t, 42, value)
	}
	{
		value, ok := BaseTypeOK("hello")
		require.True(t, ok)
		require.Equal(t, "hello", value)
	}
	{
		value, ok := BaseTypeOK(true)
		require.True(t, ok)
		require.Equal(t, true, value)
	}
	{
		value, ok := BaseTypeOK([]int{1, 2})
		require.True(t, ok)
		require.Equal(t, []any{1, 2}, value)
	}
	{
		value, ok := BaseTypeOK(map[string]int{"a": 1})
		require.True(t, ok)
		require.Equal(t, map[string]any{"a": 1}, value)
	}
	{
		// An unsupported kind returns false
		_, ok := BaseTypeOK(make(chan int))
		require.False(t, ok)
	}
}
