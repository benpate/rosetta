package channel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	result := Slice(FromSlice([]string{"a", "b", "c"}))
	require.Equal(t, []string{"a", "b", "c"}, result)
}

func TestSlice_Empty(t *testing.T) {
	result := Slice(FromSlice([]int{}))
	require.Equal(t, []int{}, result)
}

func TestFromSlice(t *testing.T) {

	collected := make([]int, 0)
	for value := range FromSlice([]int{10, 20, 30}) {
		collected = append(collected, value)
	}

	require.Equal(t, []int{10, 20, 30}, collected)
}

func TestReverse(t *testing.T) {
	result := Slice(Reverse(FromSlice([]int{1, 2, 3, 4})))
	require.Equal(t, []int{4, 3, 2, 1}, result)
}

func TestReverse_Empty(t *testing.T) {
	result := Slice(Reverse(FromSlice([]int{})))
	require.Equal(t, []int{}, result)
}

func TestClosed(t *testing.T) {

	// An open channel with no data is not "closed"
	open := make(chan int)
	require.False(t, Closed(open))

	// A closed channel reports as closed
	closed := make(chan int)
	close(closed)
	require.True(t, Closed(closed))
}
