package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfMap_Nil(t *testing.T) {
	expected := []map[string]any{}

	actual, ok := SliceOfMapOk(nil)
	require.False(t, ok)
	require.Equal(t, expected, actual)
}

func TestSliceOfMap_SliceOfAny(t *testing.T) {

	input := []any{
		map[string]any{"a": 1},
		map[string]any{"b": 2},
	}

	expected := []map[string]any{
		{"a": 1},
		{"b": 2},
	}

	result, ok := SliceOfMapOk(input)
	require.True(t, ok)
	require.Equal(t, expected, result)
}

func TestSliceOfMap_SliceOfMap(t *testing.T) {

	input := []map[string]any{
		{"a": 1},
		{"b": 2},
	}

	expected := []map[string]any{
		{"a": 1},
		{"b": 2},
	}

	result, ok := SliceOfMapOk(input)
	require.True(t, ok)
	require.Equal(t, expected, result)
}

func TestSliceOfMap_SliceOfMapOfString(t *testing.T) {

	input := []map[string]string{
		{"a": "1"},
		{"b": "2"},
	}

	expected := []map[string]any{
		{"a": "1"},
		{"b": "2"},
	}

	result, ok := SliceOfMapOk(input)
	require.True(t, ok)
	require.Equal(t, expected, result)
}

func TestSliceOfMap_Other(t *testing.T) {
	input := 7

	expected := []map[string]any{}

	result, ok := SliceOfMapOk(input)
	require.False(t, ok)
	require.Equal(t, expected, result)
}
