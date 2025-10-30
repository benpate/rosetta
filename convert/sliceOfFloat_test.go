package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfFloat_Nil(t *testing.T) {
	expected := []float64{}

	actual, ok := SliceOfFloatOk(nil)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_Float64(t *testing.T) {
	input := float64(3.14159)
	expected := []float64{3.14159}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_Int(t *testing.T) {
	input := 42
	expected := []float64{42}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_Int64(t *testing.T) {
	input := int64(42)
	expected := []float64{42}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_String(t *testing.T) {
	input := "3.14159"
	expected := []float64{3.14159}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_StringDelimited(t *testing.T) {

	input := "1.0,2.0,3.0,4.0"
	expected := []float64{1.0, 2.0, 3.0, 4.0}
	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_ReflectValue(t *testing.T) {
	input := ReflectValue(42)
	expected := []float64{42}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_SliceOfAny(t *testing.T) {
	input := []any{42}
	expected := []float64{42}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_SliceOfInt(t *testing.T) {
	input := []int{42}
	expected := []float64{42}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_SliceOfInt64(t *testing.T) {
	input := []int64{42}
	expected := []float64{42}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_SliceOfFloat64(t *testing.T) {
	input := []float64{3.14159}
	expected := []float64{3.14159}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfFloat_SliceOfString(t *testing.T) {
	input := []string{"3.14159"}
	expected := []float64{3.14159}

	actual, ok := SliceOfFloatOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}
