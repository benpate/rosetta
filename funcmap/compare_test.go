package funcmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareFuncs(t *testing.T) {

	f := All()

	isZero := f["isZero"].(func(any) bool)
	require.True(t, isZero(0))
	require.True(t, isZero(""))
	require.False(t, isZero("hello"))

	notZero := f["notZero"].(func(any) bool)
	require.True(t, notZero("hello"))
	require.False(t, notZero(0))
}
