package funcmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArraysFuncs(t *testing.T) {

	f := All()

	array := f["array"].(func(...any) []any)
	require.Equal(t, []any{1, "two", 3.0}, array(1, "two", 3.0))

	seq := f["seq"].(func(int) []int)
	require.Equal(t, []int{0, 1, 2, 3}, seq(4))
	require.Equal(t, []int{}, seq(0))

	first := f["first"].(func(...any) any)
	require.Equal(t, "found", first("", 0, "found", "other"))
	require.Nil(t, first("", 0, false))

	in := f["in"].(func(any, ...any) bool)
	require.True(t, in("b", "a", "b", "c"))
	require.False(t, in("z", "a", "b", "c"))
}
