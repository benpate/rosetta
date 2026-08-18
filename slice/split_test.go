package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {

	check := func(name string, values []int, wantHead int, wantTail []int) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			head, tail := Split(values)
			require.Equal(t, wantHead, head)
			require.Equal(t, wantTail, tail)
		})
	}

	check("many items", []int{1, 2, 3}, 1, []int{2, 3})
	check("two items", []int{1, 2}, 1, []int{2})
	check("one item yields an empty (not nil) tail", []int{1}, 1, []int{})
	check("empty slice yields the zero value and the original slice", []int{}, 0, []int{})
}

// A nil slice returns the zero value and the original nil slice -- NOT an empty slice.
func TestSplit_NilSlice(t *testing.T) {

	head, tail := Split[int](nil)

	require.Equal(t, 0, head)
	require.Nil(t, tail)
}

// For len > 1 the tail SHARES the caller's backing array, so writing through it is visible
// to the caller. The single-item case returns a fresh empty slice instead.
func TestSplit_TailAliasesBackingArray(t *testing.T) {

	values := []int{1, 2, 3}
	_, tail := Split(values)

	tail[0] = 99
	require.Equal(t, []int{1, 99, 3}, values, "the tail shares the caller's backing array")
}

func TestSplit_StringsAndStructs(t *testing.T) {

	head, tail := Split([]string{"a", "b"})
	require.Equal(t, "a", head)
	require.Equal(t, []string{"b"}, tail)

	type point struct{ X int }
	emptyHead, emptyTail := Split([]point{})
	require.Equal(t, point{}, emptyHead)
	require.Empty(t, emptyTail)
}
