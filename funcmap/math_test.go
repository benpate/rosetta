package funcmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMathFuncs(t *testing.T) {

	f := All()

	require.Equal(t, 5, f["add"].(func(any, any) int)(2, 3))
	require.Equal(t, 1, f["subtract"].(func(any, any) int)(4, 3))
	require.Equal(t, int64(12), f["multiply"].(func(any, any) int64)(3, 4))
	require.Equal(t, int64(4), f["divide"].(func(any, any) int64)(12, 3))
	require.Equal(t, 6, f["inc"].(func(any) int)(5))

	min := f["min"].(func(...any) int)
	require.Equal(t, 2, min(5, 2, 8, 3))

	max := f["max"].(func(...any) int)
	require.Equal(t, 8, max(5, 2, 8, 3))

	require.Equal(t, 42, f["int"].(func(any) int)("42"))
	require.Equal(t, int64(42), f["int64"].(func(any) int64)("42"))
	require.Equal(t, 3.5, f["float"].(func(any) float64)("3.5"))
}
