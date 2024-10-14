package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfString_Nil(t *testing.T) {
	expected := []string{}

	actual, ok := SliceOfStringOk(nil)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_Float64(t *testing.T) {
	input := float64(3.14159)
	expected := []string{"3.14159"}

	actual, ok := SliceOfStringOk(input)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_Int(t *testing.T) {
	input := int(42)
	expected := []string{"42"}

	actual, ok := SliceOfStringOk(input)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_Int64(t *testing.T) {
	input := int64(42)
	expected := []string{"42"}

	actual, ok := SliceOfStringOk(input)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_String(t *testing.T) {
	input := "hello there"
	expected := []string{"hello there"}

	actual, ok := SliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_ReflectValue(t *testing.T) {
	input := ReflectValue("general kenobi")
	expected := []string{"general kenobi"}

	actual, ok := SliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_SliceOfAny(t *testing.T) {
	input := []any{"hello", "there"}
	expected := []string{"hello", "there"}

	actual, ok := SliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_SliceOfInt(t *testing.T) {
	input := []int{69, 420}
	expected := []string{"69", "420"}

	actual, ok := SliceOfStringOk(input)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_SliceOfInt64(t *testing.T) {
	input := []int64{69, 420}
	expected := []string{"69", "420"}

	actual, ok := SliceOfStringOk(input)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfString_SliceOfString(t *testing.T) {
	input := []string{"hello", "there"}
	expected := []string{"hello", "there"}

	actual, ok := SliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}
