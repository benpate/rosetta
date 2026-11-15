package deep

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestCopy_NullInsideContainer proves that a nil value inside a map or a slice is copied as nil
func TestCopy_NullInsideContainer(t *testing.T) {

	original := map[string]any{
		"name": nil,
		"list": []any{nil, "a"},
	}

	result, err := Copy(original)

	require.NoError(t, err)
	require.Equal(t, original, result)
}

// TestCopy_SliceOfNamedMaps proves that the sliceof map types copy, and share nothing
func TestCopy_SliceOfNamedMaps(t *testing.T) {

	t.Run("MapOfAny", func(t *testing.T) {

		original := sliceof.MapOfAny{{"a": map[string]any{"b": 1}}}

		result, err := Copy(original)
		require.NoError(t, err)

		typed, isCorrectType := result.(sliceof.MapOfAny)
		require.True(t, isCorrectType)
		require.Equal(t, original, typed)

		// Changing the copy leaves the original alone
		typed[0]["a"].(map[string]any)["b"] = 2
		require.Equal(t, 1, original[0]["a"].(map[string]any)["b"])
	})

	t.Run("MapOfString", func(t *testing.T) {

		original := sliceof.MapOfString{{"a": "b"}}

		result, err := Copy(original)
		require.NoError(t, err)

		typed, isCorrectType := result.(sliceof.MapOfString)
		require.True(t, isCorrectType)
		require.Equal(t, original, typed)

		// Changing the copy leaves the original alone
		typed[0]["a"] = "c"
		require.Equal(t, "b", original[0]["a"])
	})
}

// TestCopy_NamedMapsKeepTheirType proves that mapof.Any and mapof.String copy to their own types
func TestCopy_NamedMapsKeepTheirType(t *testing.T) {

	anyResult, err := Copy(mapof.Any{"a": 1})
	require.NoError(t, err)
	require.Equal(t, mapof.Any{"a": 1}, anyResult)

	stringResult, err := Copy(mapof.String{"a": "b"})
	require.NoError(t, err)
	require.Equal(t, mapof.String{"a": "b"}, stringResult)
}

// TestCopy_ErrorNamesTypeNotValue proves that an unsupported value is never written into the error
func TestCopy_ErrorNamesTypeNotValue(t *testing.T) {

	type unsupported struct{ Secret string }

	_, err := Copy(map[string]any{"bcc": unsupported{Secret: "hunter2-recipient"}})

	require.Error(t, err)
	require.NotContains(t, derpText(err), "hunter2-recipient")
	require.Contains(t, derpText(err), "deep.unsupported")
}

// derpText returns an error and everything it wraps as JSON, details included
func derpText(err error) string {

	result, marshalErr := json.Marshal(err)

	if marshalErr != nil {
		return ""
	}

	return string(result)
}

// TestCopy_SliceOfString proves that a plain []string copies, and shares nothing
func TestCopy_SliceOfString(t *testing.T) {

	original := map[string]any{"to": []string{"a", "b"}}

	result, err := Copy(original)
	require.NoError(t, err)
	require.Equal(t, original, result)

	// Changing the copy leaves the original alone
	result.(map[string]any)["to"].([]string)[0] = "changed"
	require.Equal(t, "a", original["to"].([]string)[0])
}
