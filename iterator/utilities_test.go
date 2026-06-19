package iterator

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

// testIterator is a simple Iterator implementation over a slice of ints,
// used to exercise the utility functions in this package.
type testIterator struct {
	values []int
	index  int
}

func newTestIterator(values ...int) *testIterator {
	return &testIterator{values: values}
}

// Next copies the current value into the provided pointer and advances the cursor.
func (it *testIterator) Next(value any) bool {

	if it.index >= len(it.values) {
		return false
	}

	if pointer, ok := value.(*int); ok {
		*pointer = it.values[it.index]
	}

	it.index++
	return true
}

func (it *testIterator) Count() int {
	return len(it.values)
}

func TestMap(t *testing.T) {

	it := newTestIterator(1, 2, 3)

	result := Map(it, func(value int) string {
		return string(rune('A' + value - 1))
	})

	require.Equal(t, []string{"A", "B", "C"}, result)
}

func TestMap_Empty(t *testing.T) {

	it := newTestIterator()

	result := Map(it, func(value int) int {
		return value * 2
	})

	require.NotNil(t, result)
	require.Equal(t, 0, len(result))
}

func TestSlice(t *testing.T) {

	it := newTestIterator(4, 5, 6)

	result := Slice(it, func() int {
		return 0
	})

	require.Equal(t, []int{4, 5, 6}, result)
}

func TestSlice_Empty(t *testing.T) {

	it := newTestIterator()

	result := Slice(it, func() int {
		return 0
	})

	require.NotNil(t, result)
	require.Equal(t, 0, len(result))
}

func TestChannel(t *testing.T) {

	it := newTestIterator(7, 8, 9)

	result := make([]int, 0)
	for value := range Channel(it, func() int { return 0 }) {
		result = append(result, value)
	}

	require.Equal(t, []int{7, 8, 9}, result)
}

func TestChannel_Empty(t *testing.T) {

	it := newTestIterator()

	result := make([]int, 0)
	for value := range Channel(it, func() int { return 0 }) {
		result = append(result, value)
	}

	require.Equal(t, 0, len(result))
}

func TestChannelWithCancel(t *testing.T) {

	it := newTestIterator(10, 20, 30)
	cancel := make(chan bool)

	result := make([]int, 0)
	for value := range ChannelWithCancel(it, func() int { return 0 }, cancel) {
		result = append(result, value)
	}

	sort.Ints(result)
	require.Equal(t, []int{10, 20, 30}, result)
}

func TestChannelWithCancel_Cancelled(t *testing.T) {

	it := newTestIterator(1, 2, 3, 4, 5)
	cancel := make(chan bool)

	// Cancel immediately so the goroutine stops at its next iteration.
	close(cancel)

	// Drain the channel. Depending on scheduling, zero or more items may be
	// delivered before the cancel signal is observed, but the channel must
	// always close cleanly without blocking.
	count := 0
	for range ChannelWithCancel(it, func() int { return 0 }, cancel) {
		count++
	}

	require.LessOrEqual(t, count, 5)
}
