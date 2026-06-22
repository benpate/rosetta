package funcmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogicFuncs(t *testing.T) {

	f := All()

	iif := f["iif"].(func(bool, any, any) any)
	require.Equal(t, "yes", iif(true, "yes", "no"))
	require.Equal(t, "no", iif(false, "yes", "no"))

	and := f["and"].(func(...bool) bool)
	require.True(t, and(true, true, true))
	require.False(t, and(true, false, true))
	require.True(t, and()) // empty is vacuously true

	or := f["or"].(func(...bool) bool)
	require.True(t, or(false, true))
	require.False(t, or(false, false))
	require.False(t, or()) // empty is false
}
