package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfInt_Nil(t *testing.T) {
	expected := []int{}

	actual, ok := SliceOfIntOk(nil)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_Float64(t *testing.T) {
	input := float64(3)
	expected := []int{3}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_Int(t *testing.T) {
	input := 42
	expected := []int{42}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_Int64(t *testing.T) {
	input := int64(42)
	expected := []int{42}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_String(t *testing.T) {
	input := "3"
	expected := []int{3}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_ReflectValue(t *testing.T) {
	input := ReflectValue(42)
	expected := []int{42}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_SliceOfAny(t *testing.T) {
	input := []any{42}
	expected := []int{42}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_SliceOfInt(t *testing.T) {
	input := []int{42}
	expected := []int{42}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_SliceOfInt64(t *testing.T) {
	input := []int64{42}
	expected := []int{42}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_SliceOfFloat64(t *testing.T) {
	input := []float64{3}
	expected := []int{3}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfInt_SliceOfString(t *testing.T) {
	input := []string{"3"}
	expected := []int{3}

	actual, ok := SliceOfIntOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}
