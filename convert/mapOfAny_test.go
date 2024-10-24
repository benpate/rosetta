package convert

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapOfAny_MapOfAny(t *testing.T) {

	input := map[string]any{
		"one": 1,
	}

	actual, ok := MapOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, input, actual)
}

func TestMapOfAny_MapOfString(t *testing.T) {

	input := map[string]string{
		"one": "1",
	}

	expected := map[string]any{
		"one": "1",
	}

	actual, ok := MapOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfAny_URLValues(t *testing.T) {

	input := url.Values{}
	input.Set("one", "1")
	input.Set("two", "2")

	expected := map[string]any{
		"one": "1",
		"two": "2",
	}

	actual, ok := MapOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfAny_URLValues2s(t *testing.T) {

	input := url.Values{}
	input.Add("one", "1")
	input.Add("one", "2")

	expected := map[string]any{
		"one": []string{"1", "2"},
	}

	actual, ok := MapOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfAny_HTTPHeader(t *testing.T) {

	input := http.Header{}
	input.Set("Accept", "application/json")
	input.Set("Etag", "howdy")

	expected := map[string]any{
		"Accept": "application/json",
		"Etag":   "howdy",
	}

	actual, ok := MapOfAnyOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfAny_Reflect(t *testing.T) {

	input := map[string]any{
		"one": 1,
	}

	actual, ok := MapOfAnyOk(ReflectValue(input))
	require.True(t, ok)
	require.Equal(t, input, actual)
}

func TestMapOfAny_Other(t *testing.T) {
	input := 7

	expected := map[string]any{}

	result, ok := MapOfAnyOk(input)
	require.False(t, ok)
	require.Equal(t, expected, result)
}