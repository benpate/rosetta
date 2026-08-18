package slice

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {

	check := func(name string, values []string, fn func(string) bool, wantValue string, wantFound bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			value, found := Find(values, fn)
			require.Equal(t, wantFound, found)
			require.Equal(t, wantValue, value)
		})
	}

	hasB := func(value string) bool { return strings.Contains(value, "b") }
	never := func(string) bool { return false }
	always := func(string) bool { return true }

	check("finds the only match", []string{"a", "b", "c"}, hasB, "b", true)
	check("returns the FIRST match, not the last", []string{"ab", "bc"}, hasB, "ab", true)
	check("no match returns the zero value", []string{"a", "c"}, hasB, "", false)
	check("predicate never matches", []string{"a", "b"}, never, "", false)
	check("predicate always matches, so the first item wins", []string{"a", "b"}, always, "a", true)
	check("empty slice", []string{}, always, "", false)
	check("nil slice", nil, always, "", false)
}

// Find must stop calling the predicate as soon as it matches.
func TestFind_StopsAtFirstMatch(t *testing.T) {

	calls := 0
	value, found := Find([]int{1, 2, 3, 4, 5}, func(item int) bool {
		calls++
		return item == 2
	})

	require.True(t, found)
	require.Equal(t, 2, value)
	require.Equal(t, 2, calls, "the predicate must not be called after it matches")
}

// The zero value returned on a miss is the zero value of T, not of string.
func TestFind_ZeroValueOfStructType(t *testing.T) {

	type point struct{ X, Y int }

	value, found := Find([]point{{1, 2}}, func(point) bool { return false })

	require.False(t, found)
	require.Equal(t, point{}, value)
}
