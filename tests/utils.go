package tests

import (
	"testing"

	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

// testInline asserts that a value written to the named path reads back unchanged.
func testInline(t *testing.T, schema schema.Schema, object any, key string, value any) {
	err := schema.Set(object, key, value)
	require.Nil(t, err)

	result, err := schema.Get(object, key)
	require.Nil(t, err)
	require.Equal(t, value, result)
}
