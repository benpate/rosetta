package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRange(t *testing.T) {

	check := func(name string, values []string, wantIndexes []int, wantValues []string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {

			gotIndexes := make([]int, 0)
			gotValues := make([]string, 0)

			for index, value := range Range(values) {
				gotIndexes = append(gotIndexes, index)
				gotValues = append(gotValues, value)
			}

			require.Equal(t, wantIndexes, gotIndexes)
			require.Equal(t, wantValues, gotValues)
		})
	}

	check("many items", []string{"a", "b", "c"}, []int{0, 1, 2}, []string{"a", "b", "c"})
	check("one item", []string{"a"}, []int{0}, []string{"a"})
	check("empty slice yields nothing", []string{}, []int{}, []string{})
	check("nil slice yields nothing", nil, []int{}, []string{})
}

// Breaking out of the loop must stop the iterator, not keep walking the slice.
func TestRange_BreakStopsIteration(t *testing.T) {

	seen := make([]int, 0)

	for index := range Range([]string{"a", "b", "c", "d"}) {
		seen = append(seen, index)
		if index == 1 {
			break
		}
	}

	require.Equal(t, []int{0, 1}, seen, "iteration must stop when the caller breaks")
}

// The sequence is re-runnable: iterating the same iter.Seq2 twice yields the same values.
func TestRange_IsReusable(t *testing.T) {

	sequence := Range([]int{10, 20})

	for range 2 {
		total := 0
		count := 0
		for index, value := range sequence {
			total = total + index + value
			count++
		}
		require.Equal(t, 2, count)
		require.Equal(t, 31, total)
	}
}

func TestRange_StructValues(t *testing.T) {

	type point struct{ X int }

	got := make([]point, 0)
	for _, value := range Range([]point{{1}, {2}}) {
		got = append(got, value)
	}

	require.Equal(t, []point{{1}, {2}}, got)
}
