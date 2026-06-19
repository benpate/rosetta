package maps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {

	map1 := map[string]int{"a": 1, "b": 2, "c": 3}
	map2 := map[string]int{"a": 1, "b": 2, "c": 3}

	require.True(t, Equal(map1, map2))
	require.False(t, NotEqual(map1, map2))
}

func TestEqual_DifferentLengths(t *testing.T) {

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "b": 2, "c": 3}

	require.False(t, Equal(map1, map2))
	require.True(t, NotEqual(map1, map2))
}

func TestEqual_DifferentValues(t *testing.T) {

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "b": 99}

	require.False(t, Equal(map1, map2))
	require.True(t, NotEqual(map1, map2))
}

func TestEqual_DifferentKeys(t *testing.T) {

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "x": 2}

	require.False(t, Equal(map1, map2))
}

func TestEqual_Empty(t *testing.T) {

	require.True(t, Equal(map[string]int{}, map[string]int{}))
	require.True(t, Equal(map[string]string(nil), map[string]string(nil)))
}

func TestEqual_Strings(t *testing.T) {

	map1 := map[string]string{"first": "one", "second": "two"}
	map2 := map[string]string{"first": "one", "second": "two"}

	require.True(t, Equal(map1, map2))
}
