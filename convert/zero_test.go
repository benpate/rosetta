package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZero(t *testing.T) {
	require.True(t, IsZeroValue(""))
	require.False(t, IsZeroValue("howdy"))
	require.True(t, IsZeroValue(0))
	require.False(t, IsZeroValue(7))
	require.True(t, IsZeroValue(nil))
}
