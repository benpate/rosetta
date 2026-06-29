package delta_test

import (
	"testing"

	"github.com/benpate/rosetta/delta"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

// Compile-time assertion that *delta.Slice satisfies schema.ArrayGetterSetter. delta.Slice is used
// as an Array property (e.g. Emissary's Stream.Syndication), and schema.validate_Array requires this
// interface. Losing it (e.g. a receiver flipped from pointer to value) would only surface at runtime
// during validation, so this pins it at compile time.
//
// RULE: assert on the POINTER (*Slice) — SetIndex has a pointer receiver.
var _ schema.ArrayGetterSetter = (*delta.Slice[string])(nil)

// TestSlice_ArrayGetterSetter exercises GetIndex/SetIndex and confirms a delta.Slice of scalars
// validates cleanly under the schema package.
func TestSlice_ArrayGetterSetter(t *testing.T) {

	value := delta.NewSlice("alpha", "bravo", "charlie")

	// GetIndex returns each value, and reports out-of-range correctly.
	got, ok := value.GetIndex(1)
	require.True(t, ok)
	require.Equal(t, "bravo", got)

	_, ok = value.GetIndex(99)
	require.False(t, ok)

	_, ok = value.GetIndex(-1)
	require.False(t, ok)

	// SetIndex grows the slice and records the added value.
	require.True(t, value.SetIndex(3, "delta"))
	require.Equal(t, 4, value.Length())
	got, _ = value.GetIndex(3)
	require.Equal(t, "delta", got)

	// SetIndex rejects a value of the wrong type.
	require.False(t, value.SetIndex(0, 12345))
}

// TestSlice_SchemaValidate confirms a schema with a delta.Slice array property validates without the
// "Value must implement ArrayGetterSetter interface" error that v0.28 produced before delta.Slice
// implemented the array interfaces.
func TestSlice_SchemaValidate(t *testing.T) {

	s := schema.New(schema.Array{Items: schema.String{}})

	value := delta.NewSlice("one", "two")

	changed, err := s.Validate(&value)
	require.Nil(t, err)
	require.False(t, changed)
}
