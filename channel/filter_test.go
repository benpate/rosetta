package channel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {

	input := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	output := Filter(input, func(i int) bool {
		return i%2 == 0
	})

	require.Equal(t, 2, <-output)
	require.Equal(t, 4, <-output)
	require.Equal(t, 6, <-output)
	require.Equal(t, 8, <-output)
	require.Equal(t, 10, <-output)
}
