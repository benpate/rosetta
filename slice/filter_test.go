package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {

	value := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	value = Filter(value, func(value int) bool {
		return value%2 == 1
	})

	require.Equal(t, []int{1, 3, 5, 7, 9}, value)
}
