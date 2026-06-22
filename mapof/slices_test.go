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
