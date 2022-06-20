package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	require.True(t, Contains([]string{"hello"}, "hello"))
	require.True(t, Contains([]string{"hello", "there", "general", "kenobi"}, "hello"))
	require.True(t, Contains([]string{"hello", "there", "general", "kenobi"}, "there"))
	require.True(t, Contains([]string{"hello", "there", "general", "kenobi"}, "general"))
	require.True(t, Contains([]string{"hello", "there", "general", "kenobi"}, "kenobi"))
	require.False(t, Contains([]string{"hello", "there", "general", "kenobi"}, "grievous"))
}

func TestContainsString(t *testing.T) {

	require.True(t, Contains("one", "on"))
	require.True(t, Contains("one", "ne"))
	require.True(t, Contains("three", "th"))
	require.True(t, Contains("three", "thr"))
	require.True(t, Contains("three", "hre"))
	require.True(t, Contains("three", "hree"))

	require.False(t, Contains("one", "four"))
}

func TestContainsSliceOfString(t *testing.T) {

	require.True(t, Contains([]string{"one", "two", "three"}, "one"))
	require.True(t, Contains([]string{"one", "two", "three"}, "two"))
	require.True(t, Contains([]string{"one", "two", "three"}, "three"))
	require.False(t, Contains([]string{"one", "two", "three"}, "four"))

	require.False(t, Contains([]string{}, "empty"))
}

func TestContainsSliceOfInt(t *testing.T) {

	require.True(t, Contains([]int{1, 2, 3}, 1))
	require.True(t, Contains([]int{1, 2, 3}, 2))
	require.True(t, Contains([]int{1, 2, 3}, 3))
	require.False(t, Contains([]int{1, 2, 3}, 4))

	require.False(t, Contains([]int{}, 5))
}

func TestContainsSliceOfFloat(t *testing.T) {

	require.True(t, Contains([]float64{1, 2, 3}, 1))
	require.True(t, Contains([]float64{1, 2, 3}, 2))
	require.True(t, Contains([]float64{1, 2, 3}, 3))
	require.False(t, Contains([]float64{1, 2, 3}, 4))

	require.False(t, Contains([]float64{}, 5))
}
