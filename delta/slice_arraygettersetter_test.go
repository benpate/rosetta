package delta_test

import (
	"testing"

	"github.com/benpate/rosetta/delta"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// Compile-time assertion that *delta.Slice satisfies schema.ArrayGetterSetter. delta.Slice is used
// as an Array property (e.g. Emissary's Stream.Syndication), and schema.validate_Array requires this
// interface. Losing it (e.g. a receiver flipped from pointer to value) would only surface at runtime
// during validation, so this pins it at compile time.
//
// RULE: assert on the POINTER (*Slice) — SetIndex has a pointer receiver.
var _ schema.ArrayGetterSetter = (*delta.Slice[string])(nil)

// Compile-time assertion that *delta.Slice satisfies schema.ValueSetter. This is the seam
// that schema.SetProperty writes a whole array through, and it is preferred over every
// other setter, so losing it would silently reroute writes to reflection.
var _ schema.ValueSetter = (*delta.Slice[string])(nil)

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

// syndicationTarget stands in for Emissary's model.Stream: a delta.Slice reached through
// the PointerGetter interface, which is how schema.SetProperty descends into a struct.
type syndicationTarget struct {
	Syndication delta.Slice[string]
}

// GetPointer returns a pointer to the named property of this object
func (target *syndicationTarget) GetPointer(name string) (any, bool) {

	switch name {

	case "syndication":
		return &target.Syndication, true
	}

	return nil, false
}

// TestSlice_SchemaSet walks the FULL write path that a multiselect takes -- schema.Set,
// through validation, down PointerGetter, and into ValueSetter -- using the *sliceof.String
// that form.schemaSafeValue posts. Validating cleanly (TestSlice_SchemaValidate, above) was
// never enough on its own: the value passed validation and was then dropped on the floor.
func TestSlice_SchemaSet(t *testing.T) {

	s := schema.New(schema.Object{
		Properties: schema.ElementMap{
			"syndication": schema.Array{Items: schema.String{Required: true, Format: "token", MaxLength: 64}},
		},
	})

	{ // Selecting two targets writes both, and records both as added
		target := &syndicationTarget{Syndication: delta.NewSlice[string]()}
		posted := sliceof.String{"bandwagon", "spotify"}

		require.NoError(t, s.Set(target, "syndication", &posted))
		require.Equal(t, []string{"bandwagon", "spotify"}, target.Syndication.Values)
		require.Equal(t, []string{"bandwagon", "spotify"}, target.Syndication.Added)
		require.Equal(t, []string{}, target.Syndication.Deleted)
	}

	{ // Re-saving the same selection is not a change
		target := &syndicationTarget{Syndication: delta.NewSlice("bandwagon", "spotify")}
		posted := sliceof.String{"bandwagon", "spotify"}

		require.NoError(t, s.Set(target, "syndication", &posted))
		require.False(t, target.Syndication.IsChanged())
	}

	{ // Un-selecting everything clears the values and records the removals
		target := &syndicationTarget{Syndication: delta.NewSlice("bandwagon", "spotify")}
		posted := sliceof.String{}

		require.NoError(t, s.Set(target, "syndication", &posted))
		require.Equal(t, []string{}, target.Syndication.Values)
		require.Equal(t, []string{"bandwagon", "spotify"}, target.Syndication.Deleted)
	}

	{ // RULE: A value the schema rejects must not reach the object at all
		target := &syndicationTarget{Syndication: delta.NewSlice("bandwagon")}
		posted := sliceof.String{"not a token"}

		require.Error(t, s.Set(target, "syndication", &posted))
		require.Equal(t, []string{"bandwagon"}, target.Syndication.Values)
		require.False(t, target.Syndication.IsChanged())
	}
}
