package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfInt64_Nil(t *testing.T) {
	expected := []int64{}

	actual, ok := SliceOfInt64Ok(nil)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_Float64(t *testing.T) {
	input := float64(3)
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_Int(t *testing.T) {
	input := 42
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_Int64(t *testing.T) {
	input := int64(42)
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_String(t *testing.T) {
	input := "3"
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_StringDelimited(t *testing.T) {
	input := "3,4,5,6"
	expected := []int64{3, 4, 5, 6}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_ReflectValue(t *testing.T) {
	input := ReflectValue(42)
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfAny(t *testing.T) {
	input := []any{42}
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfInt(t *testing.T) {
	input := []int{42}
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfInt64(t *testing.T) {
	input := []int64{42}
	expected := []int64{42}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfFloat64(t *testing.T) {
	input := []float64{3}
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt64_SliceOfString(t *testing.T) {
	input := []string{"3"}
	expected := []int64{3}

	actual, ok := SliceOfInt64Ok(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}
