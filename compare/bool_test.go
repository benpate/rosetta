package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareBool(t *testing.T) {

	// Matches the implementation: equal values return 0, otherwise
	// (true, false) returns -1 and (false, true) returns 1.
	require.Equal(t, 0, Bool(true, true))
	require.Equal(t, 0, Bool(false, false))
	require.Equal(t, -1, Bool(true, false))
	require.Equal(t, 1, Bool(false, true))
}
