package funcmap

import (
	"testing"

	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

func TestStringFuncs(t *testing.T) {

	f := All()

	require.Equal(t, "42", f["string"].(func(any) string)(42))

	split := f["split"].(func(string, string) sliceof.String)
	require.Equal(t, sliceof.String{"a", "b", "c"}, split("a,b,c", ","))
	require.Equal(t, sliceof.String{}, split("", ","))

	require.Equal(t, "abc", f["join"].(func(...string) string)("a", "b", "c"))

	appendFn := f["append"].(func([]string, []string) sliceof.String)
	require.Equal(t, sliceof.String{"a", "b", "c"}, appendFn([]string{"a"}, []string{"b", "c"}))

	pluralize := f["pluralize"].(func(any, string, string) string)
	require.Equal(t, "item", pluralize(1, "item", "items"))
	require.Equal(t, "items", pluralize(2, "item", "items"))

	require.Equal(t, "hello", f["lowerCase"].(func(any) string)("HELLO"))
	require.Equal(t, "trimmed", f["trim"].(func(string) string)("  trimmed  "))

	hasPrefix := f["hasPrefix"].(func(string, string) bool)
	require.True(t, hasPrefix("hello world", "hello"))
	require.False(t, hasPrefix("hello world", "world"))

	require.Equal(t, "a1b2", f["concat"].(func(...any) string)("a", 1, "b", 2))
}
