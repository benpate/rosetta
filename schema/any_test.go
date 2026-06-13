package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {

	schema := New(Object{
		Wildcard: Any{},
	})

	data := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	// Validate the data
	value, changed, err := Validate(schema, data)

	require.Nil(t, err)
	require.False(t, changed)
	require.Equal(t, data, value)

	require.Equal(t, "bar", data["foo"])
	require.Equal(t, 123, data["baz"])
}
