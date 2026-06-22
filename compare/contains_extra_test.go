package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains_SliceOfInt64(t *testing.T) {
	require.True(t, Contains([]int64{1, 2, 3}, 2))
	require.False(t, Contains([]int64{1, 2, 3}, 4))
}

func TestContains_TypeCoercion(t *testing.T) {
	// value2 is coerced to match the slice element type
	require.True(t, Contains([]int{1, 2, 3}, "2"))
	require.True(t, Contains([]float64{1, 2, 3}, "3"))
	require.True(t, Contains([]int64{1, 2, 3}, int64(1)))
}

// containsType implements the ContainsInterfacer interface.
type containsType []string

func (c containsType) ContainsInterface(value any) bool {
	target, ok := value.(string)
	if !ok {
		return false
	}
	for _, item := range c {
		if item == target {
			return true
		}
	}
	return false
}

func TestContains_ContainsInterfacer(t *testing.T) {
	value := containsType{"alpha", "beta"}
	require.True(t, Contains(value, "alpha"))
	require.False(t, Contains(value, "gamma"))
}

func TestContains_UnsupportedType(t *testing.T) {
	// A type that is not a string, slice, or ContainsInterfacer returns FALSE
	require.False(t, Contains(42, 4))
}
