package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareBool(t *testing.T) {

	// FALSE is ordered before TRUE: equal values return 0, (false, true)
	// returns -1, and (true, false) returns 1.
	require.Equal(t, 0, Bool(true, true))
	require.Equal(t, 0, Bool(false, false))
	require.Equal(t, -1, Bool(false, true))
	require.Equal(t, 1, Bool(true, false))
}
