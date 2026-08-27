package delta

import (
	"testing"

	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestSlice_SetValue_PointerToNamedSlice covers the exact shape that multi-value form
// widgets post.  form.schemaSafeValue wraps every Array-typed path in a *sliceof.String,
// because that is the only shape schema.validate_Array accepts, so this is what SetValue
// actually receives in production.
func TestSlice_SetValue_PointerToNamedSlice(t *testing.T) {

	s := NewSlice[string]()

	require.NoError(t, s.SetValue(&sliceof.String{"bandwagon", "spotify"}))

	require.Equal(t, []string{"bandwagon", "spotify"}, s.Values)
	require.Equal(t, []string{"bandwagon", "spotify"}, s.Added)
	require.Equal(t, []string{}, s.Deleted)
	require.True(t, s.IsChanged())
}

// TestSlice_SetValue_ClearsWhenEmpty confirms that un-checking every option still clears
// the slice, and records the departing values so that a caller can act on the removals.
func TestSlice_SetValue_ClearsWhenEmpty(t *testing.T) {

	s := NewSlice("bandwagon", "spotify")

	require.NoError(t, s.SetValue(&sliceof.String{}))

	require.Equal(t, []string{}, s.Values)
	require.Equal(t, []string{}, s.Added)
	require.Equal(t, []string{"bandwagon", "spotify"}, s.Deleted)
	require.True(t, s.IsChanged())
}

// TestSlice_SetValue_UnchangedIsNotAChange re-saves a Slice with exactly the values it
// already holds.  This is the assertion that guards against phantom diffs: a caller that
// broadcasts Deleted (Emissary sends syndication messages this way) must not be told that
// everything was removed simply because the object was saved again.
func TestSlice_SetValue_UnchangedIsNotAChange(t *testing.T) {

	s := NewSlice("bandwagon", "spotify")

	require.NoError(t, s.SetValue(&sliceof.String{"bandwagon", "spotify"}))

	require.Equal(t, []string{"bandwagon", "spotify"}, s.Values)
	require.Equal(t, []string{}, s.Added)
	require.Equal(t, []string{}, s.Deleted)
	require.False(t, s.IsChanged())
}

// TestSlice_SetValue_PartialChange confirms the added/deleted bookkeeping when a save both
// adds and removes values at once.
func TestSlice_SetValue_PartialChange(t *testing.T) {

	s := NewSlice("bandwagon", "spotify")

	require.NoError(t, s.SetValue(&sliceof.String{"spotify", "tidal"}))

	require.Equal(t, []string{"spotify", "tidal"}, s.Values)
	require.Equal(t, []string{"tidal"}, s.Added)
	require.Equal(t, []string{"bandwagon"}, s.Deleted)
	require.Equal(t, []string{"spotify"}, s.Unchanged())
}

// TestSlice_SetValue_DoesNotAliasInput confirms that SetValue copies the incoming values.
// They arrive straight from a url.Values entry, which the caller still owns and may reuse.
func TestSlice_SetValue_DoesNotAliasInput(t *testing.T) {

	input := sliceof.String{"bandwagon", "spotify"}

	s := NewSlice[string]()
	require.NoError(t, s.SetValue(&input))

	// Mutating the caller's slice must not reach into the Slice
	input[0] = "CLOBBERED"
	require.Equal(t, []string{"bandwagon", "spotify"}, s.Values)

	// And growing the Slice must not reach back into the caller's slice
	require.True(t, s.SetIndex(0, "changed"))
	require.Equal(t, sliceof.String{"CLOBBERED", "spotify"}, input)
}

// TestSlice_SetValue_ErrorLeavesValuesIntact confirms that a value SetValue cannot read is
// reported without disturbing the Slice.  A failed write that still cleared the values
// would hand the caller an empty list AND a full Deleted list -- the worst of both.
func TestSlice_SetValue_ErrorLeavesValuesIntact(t *testing.T) {

	s := NewSlice("bandwagon", "spotify")
	s.Added = []string{"spotify"}

	// A struct is not a collection, and cannot be rendered as one
	require.Error(t, s.SetValue(struct{ Name string }{Name: "nope"}))

	require.Equal(t, []string{"bandwagon", "spotify"}, s.Values)
	require.Equal(t, []string{"spotify"}, s.Added)
	require.Equal(t, []string{}, s.Deleted)
}

// TestSlice_SetValue_CoercesScalars documents the deliberate looseness inherited from the
// convert package: for an element type convert knows how to produce, values are rendered
// into it rather than rejected.  Only a value that is not a collection at all -- and cannot
// be rendered as one -- fails.
func TestSlice_SetValue_CoercesScalars(t *testing.T) {

	{ // A bare scalar becomes a one-item slice
		s := NewSlice[string]()
		require.NoError(t, s.SetValue(12345))
		require.Equal(t, []string{"12345"}, s.Values)
	}

	{ // A slice of another scalar type is rendered element by element
		s := NewSlice[string]()
		require.NoError(t, s.SetValue([]int{1, 2}))
		require.Equal(t, []string{"1", "2"}, s.Values)
	}

	// A scalar convert can only render approximately is still stored. For a Slice[int] that
	// means a non-numeric string lands as the zero value -- the cost of allowing lossy
	// conversions, and the reason the error path is now reachable only for non-collections.
	{
		s := NewSlice[int]()
		require.NoError(t, s.SetValue("not a number"))
		require.Equal(t, []int{0}, s.Values)
	}

	{ // RULE: Coercion still has a floor. A value that is not a collection is an error.
		s := NewSlice[int]()
		require.Error(t, s.SetValue(map[string]any{"a": 1}))
		require.Empty(t, s.Values)
	}
}

// TestSlice_SetValue_NilClears confirms that a nil value empties the Slice, which is how an
// un-checked multi-value widget reports that nothing is selected any more.
func TestSlice_SetValue_NilClears(t *testing.T) {

	s := NewSlice("bandwagon")

	require.NoError(t, s.SetValue(nil))
	require.Equal(t, []string{}, s.Values)
	require.Equal(t, []string{"bandwagon"}, s.Deleted)
}

// TestSlice_SetValue_TypedNilPointer pins the treatment of a typed nil pointer, which is NOT
// the same as a nil value: convert reports it as unconvertible, so it is an error that leaves
// the Slice alone.  A pointer to an empty slice -- what form.schemaSafeValue actually posts
// when nothing is checked -- is a different thing entirely, and does clear.
func TestSlice_SetValue_TypedNilPointer(t *testing.T) {

	{ // A nil POINTER is unreadable, and reported as such
		s := NewSlice("bandwagon")
		require.Error(t, s.SetValue((*sliceof.String)(nil)))
		require.Equal(t, []string{"bandwagon"}, s.Values)
		require.False(t, s.IsChanged())
	}

	{ // A pointer to a nil SLICE is an empty collection, and clears
		s := NewSlice("bandwagon")
		require.NoError(t, s.SetValue(&sliceof.String{}))
		require.Equal(t, []string{}, s.Values)
		require.Equal(t, []string{"bandwagon"}, s.Deleted)
	}
}

// TestSlice_SetValue_AcceptsLossyConversion confirms that a value convert can render only
// approximately is still stored.  SetValue reads convert's `converted` flag, not `lossless`:
// a value convert can render is a value this Slice will hold.
func TestSlice_SetValue_AcceptsLossyConversion(t *testing.T) {

	{ // The fixed two-decimal float rendering does not round-trip, and is kept anyway
		s := NewSlice[string]()
		require.NoError(t, s.SetValue(3.14159))
		require.Equal(t, []string{"3.14"}, s.Values)
	}

	{ // RULE: The floor is "could convert render this at all", not "did it round-trip".
		// A struct is not a collection, so it is still an error.
		s := NewSlice("bandwagon")
		require.Error(t, s.SetValue(struct{ Name string }{Name: "nope"}))
		require.Equal(t, []string{"bandwagon"}, s.Values)
		require.False(t, s.IsChanged())
	}
}

// TestSlice_SetValue_ValuesStayNotNil confirms that Values is never left nil, matching the
// guarantee that NewSlice and Reset already make for Added and Deleted.
func TestSlice_SetValue_ValuesStayNotNil(t *testing.T) {

	var s Slice[string]

	require.NoError(t, s.SetValue([]string(nil)))

	require.NotNil(t, s.Values)
	require.NotNil(t, s.Added)
	require.NotNil(t, s.Deleted)
	require.Zero(t, s.Length())
}
