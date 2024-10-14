package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfAny_Nil(t *testing.T) {
	expected := []any{}

	actual, ok := SliceOfAnyOk(nil)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_Float64(t *testing.T) {
	input := float64(3.14159)
	expected := []any{3.14159}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_Int(t *testing.T) {
	input := 42
	expected := []any{42}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_Int64(t *testing.T) {
	input := int64(42)
	expected := []any{int64(42)}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_String(t *testing.T) {
	input := "hello"
	expected := []any{"hello"}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_ReflectValue(t *testing.T) {
	input := ReflectValue(42)
	expected := []any{42}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_SliceOfAny(t *testing.T) {
	input := []any{42}
	expected := []any{42}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_SliceOfInt(t *testing.T) {
	input := []int{42}
	expected := []any{42}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_SliceOfInt64(t *testing.T) {
	input := []int64{42}
	expected := []any{int64(42)}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_SliceOfFloat64(t *testing.T) {
	input := []float64{3.14159}
	expected := []any{3.14159}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfAny_SliceOfString(t *testing.T) {
	input := []string{"hello"}
	expected := []any{"hello"}

	actual, ok := SliceOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}
