package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	require.True(t, Contains([]string{"a", "b", "c"}, "a"))
	require.True(t, Contains([]string{"a", "b", "c"}, "b"))
	require.True(t, Contains([]string{"a", "b", "c"}, "c"))
	require.False(t, Contains([]string{"a", "b", "c"}, "d"))
}

func TestContainsAny(t *testing.T) {
	require.False(t, ContainsAny([]string{"a", "b", "c"}))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "a"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "b"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "c"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "a", "b"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "b", "a"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "c", "b"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "a", "1", "2"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "b", "1", "2"))
	require.True(t, ContainsAny([]string{"a", "b", "c"}, "c", "1", "2"))
	require.False(t, ContainsAny([]string{"a", "b", "c"}, "1", "2"))
}

func TestContainsAll(t *testing.T) {
	require.True(t, ContainsAll([]string{"a", "b", "c"}))
	require.True(t, ContainsAll([]string{"a", "b", "c"}, "a", "b", "c"))
	require.True(t, ContainsAll([]string{"a", "b", "c"}, "b", "c", "a"))
	require.True(t, ContainsAll([]string{"a", "b", "c"}, "c", "a", "b"))
	require.True(t, ContainsAll([]string{"a", "b", "c"}, "a", "b"))
	require.True(t, ContainsAll([]string{"a", "b", "c"}, "b", "a"))
	require.True(t, ContainsAll([]string{"a", "b", "c"}, "c", "b"))

	require.False(t, ContainsAll([]string{"a", "b", "c"}, "a", "b", "c", "d"))
	require.False(t, ContainsAll([]string{"a", "b", "c"}, "b", "c", "d"))
	require.False(t, ContainsAll([]string{"a", "b", "c"}, "1", "2"))
}
