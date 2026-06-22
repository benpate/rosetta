package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaps_Equal(t *testing.T) {
	a := map[string]int{"x": 1, "y": 2}
	b := map[string]int{"x": 1, "y": 2}
	require.True(t, Maps(a, b))
}

func TestMaps_DifferentLength(t *testing.T) {
	a := map[string]int{"x": 1}
	b := map[string]int{"x": 1, "y": 2}
	require.False(t, Maps(a, b))
}

func TestMaps_MissingKey(t *testing.T) {
	a := map[string]int{"x": 1, "y": 2}
	b := map[string]int{"x": 1, "z": 2}
	require.False(t, Maps(a, b))
}

func TestMaps_DifferentValue(t *testing.T) {
	a := map[string]int{"x": 1, "y": 2}
	b := map[string]int{"x": 1, "y": 99}
	require.False(t, Maps(a, b))
}

func TestMaps_Empty(t *testing.T) {
	require.True(t, Maps(map[string]int{}, map[string]int{}))
}

func TestMaps_StringValues(t *testing.T) {
	a := map[string]string{"a": "one", "b": "two"}
	b := map[string]string{"a": "one", "b": "two"}
	require.True(t, Maps(a, b))
}
