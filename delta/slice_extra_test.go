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

func TestSlice_SetValue_WrongType(t *testing.T) {

	s := NewSlice[int](1, 2, 3)

	// A value of the wrong type is treated as an empty slice
	require.NoError(t, s.SetValue("not a slice"))
	require.Equal(t, []int{}, s.Values)
	require.Equal(t, []int{1, 2, 3}, s.Deleted)
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
