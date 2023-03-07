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

	// Validate and Clean
	require.Nil(t, schema.Validate(data))
	require.Nil(t, schema.Clean(data))

	require.Equal(t, "bar", data["foo"])
	require.Equal(t, 123, data["baz"])
}
