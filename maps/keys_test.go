package maps

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {

	value := map[string]int{"a": 1, "b": 2, "c": 3}

	keys := Keys(value)
	sort.Strings(keys)

	require.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestKeys_Empty(t *testing.T) {

	keys := Keys(map[string]int{})

	require.NotNil(t, keys)
	require.Equal(t, 0, len(keys))
}

func TestKeys_IntKeys(t *testing.T) {

	value := map[int]string{3: "three", 1: "one", 2: "two"}

	keys := Keys(value)
	sort.Ints(keys)

	require.Equal(t, []int{1, 2, 3}, keys)
}

func TestKeysSorted(t *testing.T) {

	value := map[string]int{"charlie": 3, "alpha": 1, "bravo": 2}

	require.Equal(t, []string{"alpha", "bravo", "charlie"}, KeysSorted(value))
}

func TestKeysSorted_Ints(t *testing.T) {

	value := map[int]string{30: "c", 10: "a", 20: "b"}

	require.Equal(t, []int{10, 20, 30}, KeysSorted(value))
}

func TestKeysSorted_Empty(t *testing.T) {

	keys := KeysSorted(map[string]int{})

	require.NotNil(t, keys)
	require.Equal(t, 0, len(keys))
}
