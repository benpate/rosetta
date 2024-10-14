package convert

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapOfSliceOfString_MapOfAny(t *testing.T) {
	input := map[string]any{
		"one": 1,
	}

	expected := map[string][]string{
		"one": {"1"},
	}

	actual, ok := MapOfSliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfSliceOfString_MapOfString(t *testing.T) {
	input := map[string]string{
		"one": "1",
	}

	expected := map[string][]string{
		"one": {"1"},
	}

	actual, ok := MapOfSliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfSliceOfString_MapOfSliceOfString(t *testing.T) {
	input := map[string][]string{
		"one": {"1"},
	}

	expected := map[string][]string{
		"one": {"1"},
	}

	actual, ok := MapOfSliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfSliceOfString_URLValues(t *testing.T) {
	input := url.Values{}
	input.Add("name1", "value1")
	input.Add("name2", "value2")
	input.Add("name3", "value3")
	input.Add("name3", "another value3")

	expected := map[string][]string{
		"name1": {"value1"},
		"name2": {"value2"},
		"name3": {"value3", "another value3"},
	}

	actual, ok := MapOfSliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfSliceOfString_HTTPHeader(t *testing.T) {
	input := http.Header{}
	input.Add("name1", "value1")
	input.Add("name2", "value2")
	input.Add("name3", "value3")
	input.Add("name3", "another value3")

	expected := map[string][]string{
		"Name1": {"value1"},
		"Name2": {"value2"},
		"Name3": {"value3", "another value3"},
	}

	actual, ok := MapOfSliceOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfSliceOfString_ReflectValue(t *testing.T) {
	input := map[string]any{
		"one": 1,
	}

	expected := map[string][]string{
		"one": {"1"},
	}

	actual, ok := MapOfSliceOfStringOk(ReflectValue(input))
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfSliceOfString_Other(t *testing.T) {
	input := 7

	expected := map[string][]string{}

	result, ok := MapOfSliceOfStringOk(input)
	require.False(t, ok)
	require.Equal(t, expected, result)
}
