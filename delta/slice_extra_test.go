package delta

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice_NewSlice(t *testing.T) {

	// Constructor with values
	s := NewSlice[int](1, 2, 3)
	require.Equal(t, []int{1, 2, 3}, s.Values)
	require.Equal(t, []int{}, s.Added)
	require.Equal(t, []int{}, s.Deleted)

	// Constructor with no values is still initialized
	empty := NewSlice[string]()
	require.Equal(t, []string{}, empty.Values)
	require.Equal(t, 0, empty.Length())
}

func TestSlice_LengthAndGetValue(t *testing.T) {

	s := NewSlice[int](1, 2, 3, 4)
	require.Equal(t, 4, s.Length())
	require.Equal(t, []int{1, 2, 3, 4}, s.GetValue())
}

func TestSlice_IsChangedAndUnchanged(t *testing.T) {

	s := NewSlice[int](1, 2, 3, 4)
	require.False(t, s.IsChanged())

	require.NoError(t, s.SetValue([]int{1, 3, 5}))
	require.True(t, s.IsChanged())

	// Unchanged returns values present in both old and new
	require.Equal(t, []int{1, 3}, s.Unchanged())
}

func TestSlice_IsChanged_OnlyDeleted(t *testing.T) {

	s := NewSlice[int](1, 2, 3)
	require.NoError(t, s.SetValue([]int{1, 2}))

	require.True(t, s.IsChanged())
	require.Equal(t, []int{3}, s.Deleted)
	require.Equal(t, []int{}, s.Added)
}

func TestSlice_Reset(t *testing.T) {

	s := NewSlice[int](1, 2, 3)
	require.NoError(t, s.SetValue([]int{1, 4}))
	require.True(t, s.IsChanged())

	s.Reset()
	require.False(t, s.IsChanged())
	require.Equal(t, []int{}, s.Added)
	require.Equal(t, []int{}, s.Deleted)
}

// TestSlice_SetValue_WrongType confirms that a value that cannot be read as a collection of
// T is REPORTED, and that the Slice is left exactly as it was.  Quietly swallowing the
// mismatch (and clearing the values) is what let a broken multiselect wipe its field, and
// populate Deleted, without anything upstream noticing.
func TestSlice_SetValue_WrongType(t *testing.T) {

	s := NewSlice[int](1, 2, 3)

	// A channel is not a collection, and convert cannot render it as one
	require.Error(t, s.SetValue(make(chan int)))

	// Nothing moved: the values survive, and no phantom diff was recorded
	require.Equal(t, []int{1, 2, 3}, s.Values)
	require.Equal(t, []int{}, s.Added)
	require.Equal(t, []int{}, s.Deleted)
	require.False(t, s.IsChanged())
}

func TestSlice_MarshalJSON(t *testing.T) {

	s := NewSlice[int](1, 2, 3)

	data, err := json.Marshal(s)
	require.NoError(t, err)
	require.Equal(t, "[1,2,3]", string(data))
}

func TestSlice_UnmarshalJSON(t *testing.T) {

	var s Slice[int]
	require.NoError(t, json.Unmarshal([]byte("[4,5,6]"), &s))

	require.Equal(t, []int{4, 5, 6}, s.Values)
	require.Equal(t, []int{}, s.Added)
	require.Equal(t, []int{}, s.Deleted)
	require.False(t, s.IsChanged())
}

func TestSlice_UnmarshalJSON_Error(t *testing.T) {

	var s Slice[int]
	require.Error(t, json.Unmarshal([]byte("not json"), &s))
}

func TestSlice_RoundTrip(t *testing.T) {

	type wrapper struct {
		Tags Slice[string] `json:"tags"`
	}

	original := wrapper{Tags: NewSlice("a", "b", "c")}

	data, err := json.Marshal(original)
	require.NoError(t, err)
	require.JSONEq(t, `{"tags":["a","b","c"]}`, string(data))

	var decoded wrapper
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, []string{"a", "b", "c"}, decoded.Tags.Values)
}

// REGRESSION: NewSlice and Reset both guarantee that Added and Deleted are non-nil.
// SetValue used to clone Values into Deleted, and slices.Clone returns nil for a nil
// input, so calling SetValue on a zero-value Slice left Deleted nil.
func TestSlice_SetValue_ListsStayNotNil(t *testing.T) {

	var s Slice[string] // zero value: Values is nil

	require.NoError(t, s.SetValue([]string{"a"}))

	require.NotNil(t, s.Added)
	require.NotNil(t, s.Deleted)
	require.Equal(t, []string{"a"}, s.Added)
	require.Empty(t, s.Deleted)
}
