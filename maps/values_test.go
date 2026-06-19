package maps

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {

	value := map[string]int{"a": 1, "b": 2, "c": 3}

	values := Values(value)
	sort.Ints(values)

	require.Equal(t, []int{1, 2, 3}, values)
}

func TestValues_Empty(t *testing.T) {

	values := Values(map[string]int{})

	require.NotNil(t, values)
	require.Equal(t, 0, len(values))
}

func TestValues_StringValues(t *testing.T) {

	value := map[int]string{1: "one", 2: "two", 3: "three"}

	values := Values(value)
	sort.Strings(values)

	require.Equal(t, []string{"one", "three", "two"}, values)
}
