package delta

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice1(t *testing.T) {

	s := NewSlice[int](1, 2, 3, 4)

	err := s.SetValue([]int{1, 2, 3, 4})
	require.Nil(t, err)

	require.Equal(t, []int{1, 2, 3, 4}, s.Values)
	require.Equal(t, []int{}, s.Added)
	require.Equal(t, []int{}, s.Deleted)
}

func TestSlice2(t *testing.T) {

	s := NewSlice[int]()

	err := s.SetValue([]int{1, 2, 3, 4})
	require.Nil(t, err)

	require.Equal(t, []int{1, 2, 3, 4}, s.Values)
	require.Equal(t, []int{1, 2, 3, 4}, s.Added)
	require.Equal(t, []int{}, s.Deleted)
}

func TestSlice3(t *testing.T) {

	s := NewSlice[int](1, 2, 3, 4)

	err := s.SetValue([]int{1, 3, 5})
	require.Nil(t, err)

	require.Equal(t, []int{1, 3, 5}, s.Values)
	require.Equal(t, []int{5}, s.Added)
	require.Equal(t, []int{2, 4}, s.Deleted)
}

func TestSlice_SetValueDoesNotAliasCaller(t *testing.T) {

	// SetValue compacts the "deleted" list in place, so it must be working on a
	// copy -- otherwise it writes through to the slice the caller handed us.
	existing := []string{"a", "b", "c"}
	s := NewSlice(existing...)

	require.Nil(t, s.SetValue([]string{"b"}))

	require.Equal(t, []string{"a", "b", "c"}, existing, "the caller's slice must not be modified")
	require.Equal(t, []string{"a", "c"}, s.Deleted)
	require.Equal(t, []string{"b"}, s.Values)
	require.Equal(t, []string{}, s.Added)
}

func TestSlice_SetValueEmpty(t *testing.T) {

	// Every previous value is deleted when the new list is empty
	existing := []string{"a", "b"}
	s := NewSlice(existing...)

	require.Nil(t, s.SetValue([]string{}))

	require.Equal(t, []string{"a", "b"}, existing)
	require.Equal(t, []string{"a", "b"}, s.Deleted)
	require.Empty(t, s.Values)
}

func TestSlice_SetValueDuplicates(t *testing.T) {

	// A repeated value removes exactly ONE occurrence from the deleted list
	s := NewSlice("a", "a", "b")

	require.Nil(t, s.SetValue([]string{"a"}))

	require.Equal(t, []string{"a", "b"}, s.Deleted)
	require.Equal(t, []string{"a"}, s.Values)
	require.Empty(t, s.Added)
}

func TestSlice_SetValueZeroValue(t *testing.T) {

	// A zero-value Slice has a nil Values list, which must not panic
	var s Slice[string]

	require.Nil(t, s.SetValue([]string{"a"}))

	require.Equal(t, []string{"a"}, s.Values)
	require.Equal(t, []string{"a"}, s.Added)
	require.Empty(t, s.Deleted)
}
