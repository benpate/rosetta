package tests

import (
	"testing"

	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

type testTableItem struct {
	key   string
	value any
}

func testInline(t *testing.T, schema schema.Schema, object any, key string, value any) {
	err := schema.Set(object, key, value)
	require.Nil(t, err)

	result, err := schema.Get(object, key)
	require.Nil(t, err)
	require.Equal(t, value, result)
}
