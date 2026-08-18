package mapof

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlices_Add(t *testing.T) {

	s := make(Slices[string, int])

	s.Add("a", 1)
	s.Add("a", 2)
	s.Add("b", 3)

	require.Equal(t, []int{1, 2}, s["a"])
	require.Equal(t, []int{3}, s["b"])
}

func TestSlices_Flatten(t *testing.T) {

	s := make(Slices[string, int])
	s.Add("a", 1)
	s.Add("a", 2)
	s.Add("b", 3)

	flat := s.Flatten()
	sort.Ints(flat)

	require.Equal(t, []int{1, 2, 3}, flat)
}

func TestSlices_FlattenEmpty(t *testing.T) {
	s := make(Slices[string, int])
	require.Nil(t, s.Flatten())
}

func TestSlices_AddToNilMap(t *testing.T) {

	// A zero-value Slices is a nil map -- Add must allocate it rather than panic
	var s Slices[string, int]

	s.Add("first", 1)
	s.Add("first", 2)
	s.Add("second", 3)

	require.Equal(t, []int{1, 2}, s["first"])
	require.Equal(t, []int{3}, s["second"])
	require.Len(t, s, 2)
}

func TestSlices_FlattenNilMap(t *testing.T) {

	// Reading from a nil map is already safe, and must stay that way
	var s Slices[string, int]

	require.Empty(t, s.Flatten())
	require.Len(t, s, 0)
}
