package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZero(t *testing.T) {
	require.True(t, IsZero(""))
	require.False(t, IsZero("howdy"))
	require.True(t, IsZero(0))
	require.False(t, IsZero(7))
	require.True(t, IsZero(nil))

	require.False(t, NotZero(""))
	require.True(t, NotZero("howdy"))
	require.False(t, NotZero(0))
	require.True(t, NotZero(7))
	require.False(t, NotZero(nil))
}
