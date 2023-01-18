package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObjectReverseOdd(t *testing.T) {

	x := Object[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
	require.True(t, x.IsLength(9))

	x.Reverse()

	require.Equal(t, 9, x.Length())
	require.Equal(t, 9, x[0])
	require.Equal(t, 8, x[1])
	require.Equal(t, 7, x[2])
	require.Equal(t, 6, x[3])
	require.Equal(t, 5, x[4])
	require.Equal(t, 4, x[5])
	require.Equal(t, 3, x[6])
	require.Equal(t, 2, x[7])
	require.Equal(t, 1, x[8])
}

func TestObjectReverseEven(t *testing.T) {

	x := Object[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	require.True(t, x.IsLength(10))

	x.Reverse()

	require.Equal(t, 10, x.Length())
	require.Equal(t, 10, x[0])
	require.Equal(t, 9, x[1])
	require.Equal(t, 8, x[2])
	require.Equal(t, 7, x[3])
	require.Equal(t, 6, x[4])
	require.Equal(t, 5, x[5])
	require.Equal(t, 4, x[6])
	require.Equal(t, 3, x[7])
	require.Equal(t, 2, x[8])
	require.Equal(t, 1, x[9])
}
