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
