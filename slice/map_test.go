package slice

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {

	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	result := Map(original, func(v int) string {
		return strconv.Itoa(v * 11)
	})

	require.Equal(t, []string{"11", "22", "33", "44", "55", "66", "77", "88", "99"}, result)
}
