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
