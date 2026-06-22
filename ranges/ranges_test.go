package ranges

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	require.Equal(t, []int{1, 2, 3}, Slice(Values(1, 2, 3)))
}

func TestValues_Empty(t *testing.T) {
	require.Equal(t, []int{}, Slice(Values[int]()))
}

func TestValues_EarlyTermination(t *testing.T) {

	// Break out of the loop after the first item to exercise the
	// `!yield(value)` path inside Values.
	count := 0
	for range Values(1, 2, 3, 4) {
		count++
		break
	}
	require.Equal(t, 1, count)
}

func TestEmpty(t *testing.T) {
	require.Equal(t, []int{}, Slice(Empty[int]()))
}

func TestSlice(t *testing.T) {
	require.Equal(t, []string{"a", "b"}, Slice(Values("a", "b")))
}

func TestFilter(t *testing.T) {

	result := Slice(Filter(Values(1, 2, 3, 4, 5, 6), func(v int) bool {
		return v%2 == 0
	}))

	require.Equal(t, []int{2, 4, 6}, result)
}

func TestFilter_EarlyTermination(t *testing.T) {

	count := 0
	for range Filter(Values(2, 4, 6, 8), func(v int) bool { return true }) {
		count++
		break
	}
	require.Equal(t, 1, count)
}

func TestFilterPointer(t *testing.T) {

	result := Slice(FilterPointer(Values(1, 2, 3, 4), func(v *int) bool {
		return *v > 2
	}))

	require.Equal(t, []int{3, 4}, result)
}

func TestFilterPointer_EarlyTermination(t *testing.T) {

	count := 0
	for range FilterPointer(Values(1, 2, 3), func(v *int) bool { return true }) {
		count++
		break
	}
	require.Equal(t, 1, count)
}

func TestJoin(t *testing.T) {

	result := Slice(Join(Values(1, 2), Values(3, 4), Values(5)))
	require.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestJoin_Empty(t *testing.T) {
	require.Equal(t, []int{}, Slice(Join[int]()))
}

func TestJoin_EarlyTermination(t *testing.T) {

	count := 0
	for range Join(Values(1, 2), Values(3, 4)) {
		count++
		break
	}
	require.Equal(t, 1, count)
}

func TestLimit(t *testing.T) {

	require.Equal(t, []int{1, 2, 3}, Slice(Limit(3, Values(1, 2, 3, 4, 5))))

	// Limit greater than length returns all items
	require.Equal(t, []int{1, 2}, Slice(Limit(10, Values(1, 2))))

	// Limit of zero returns nothing
	require.Equal(t, []int{}, Slice(Limit(0, Values(1, 2, 3))))
}

func TestLimit_EarlyTermination(t *testing.T) {

	count := 0
	for range Limit(5, Values(1, 2, 3, 4, 5)) {
		count++
		break
	}
	require.Equal(t, 1, count)
}

func TestMap(t *testing.T) {

	result := Slice(Map(Values(1, 2, 3), func(v int) string {
		return string(rune('A' + v - 1))
	}))

	require.Equal(t, []string{"A", "B", "C"}, result)
}

func TestMap_EarlyTermination(t *testing.T) {

	count := 0
	for range Map(Values(1, 2, 3), func(v int) int { return v * 2 }) {
		count++
		break
	}
	require.Equal(t, 1, count)
}

func TestUnique_EarlyTermination(t *testing.T) {

	count := 0
	for range Unique(Values(1, 2, 2, 3)) {
		count++
		break
	}
	require.Equal(t, 1, count)
}
