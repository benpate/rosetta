package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// NOTE: Schema.Append's success path is not exercised here. Append grows the
// array correctly, but its final step re-validates the value via Schema.Set,
// which passes the dereferenced (non-pointer) array to validate_Array. Because
// array SetIndex methods are pointer-receivers, the non-pointer value does not
// satisfy ArrayGetterSetter, so the re-validation fails. Only Append's error
// branches are reachable. See the summary notes for details.

func TestSchema_Append_NotAnArray(t *testing.T) {
	// Appending to a non-array element returns an error
	schema := New(Object{Properties: ElementMap{"name": String{}}})

	err := schema.Append(&testStructA{}, "name", "value")
	require.Error(t, err)
}

func TestSchema_Append_GetError(t *testing.T) {
	// The path resolves to an array element, but the object can't provide it
	schema := New(Object{Properties: ElementMap{"list": Array{Items: String{}}}})

	err := schema.Append(&testStructA{}, "list", "value")
	require.Error(t, err)
}

func TestSchema_Remove(t *testing.T) {
	value := newTestArrayA()
	schema := New(testArrayA_Schema())

	// Remove the item at index 1
	require.True(t, schema.Remove(&value, "1"))
	require.Equal(t, testArrayA{"one", "three"}, value)
}

func TestSchema_Remove_NotARemover(t *testing.T) {
	// A value that does not implement Remover cannot be removed from
	require.False(t, New(testStructA_Schema()).Remove(&testStructA{}, "name"))
}

func TestSchema_Remove_Nested(t *testing.T) {
	value := newTestStructA()
	schema := New(testStructA_Schema())

	// Dig into the "array" property, then remove index 1 from it
	require.True(t, schema.Remove(&value, "array.1"))
	require.Equal(t, testArrayA{"one", "three"}, value.Array)
}

func TestSchema_Remove_NestedMissing(t *testing.T) {
	// An intermediate path segment that does not exist returns false
	value := newTestStructA()
	schema := New(testStructA_Schema())

	require.False(t, schema.Remove(&value, "missing.0"))
}
