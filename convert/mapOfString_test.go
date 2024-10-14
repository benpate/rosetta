package convert

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapOfString_MapOfAny(t *testing.T) {

	input := map[string]any{
		"one": 1,
	}

	expected := map[string]string{
		"one": "1",
	}

	actual, ok := MapOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfString_MapOfString(t *testing.T) {

	input := map[string]string{
		"one": "1",
	}

	expected := map[string]string{
		"one": "1",
	}

	actual, ok := MapOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfString_URLValues(t *testing.T) {

	input := url.Values{}
	input.Set("one", "1")
	input.Set("two", "2")

	expected := map[string]string{
		"one": "1",
		"two": "2",
	}

	actual, ok := MapOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfString_URLValues2s(t *testing.T) {

	input := url.Values{}
	input.Add("one", "1")
	input.Add("one", "2")

	expected := map[string]string{
		"one": "1",
	}

	actual, ok := MapOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfString_HTTPHeader(t *testing.T) {

	input := http.Header{}
	input.Set("Accept", "application/json")
	input.Set("Etag", "howdy")

	expected := map[string]string{
		"Accept": "application/json",
		"Etag":   "howdy",
	}

	actual, ok := MapOfStringOk(input)
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfString_Reflect(t *testing.T) {

	input := map[string]any{
		"one": 1,
	}

	expected := map[string]string{
		"one": "1",
	}

	actual, ok := MapOfStringOk(ReflectValue(input))
	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestMapOfString_Other(t *testing.T) {
	input := 7

	expected := map[string]string{}

	result, ok := MapOfStringOk(input)
	require.False(t, ok)
	require.Equal(t, expected, result)
}
