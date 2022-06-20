package null

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {

	var b Bool

	require.True(t, b.IsNull())
	require.False(t, b.IsPresent())
	require.False(t, b.Bool())
	require.Equal(t, "", b.String())
	require.Nil(t, b.Interface())

	b.Set(false)
	require.False(t, b.IsNull())
	require.True(t, b.IsPresent())
	require.False(t, b.Bool())
	require.Equal(t, "false", b.String())
	require.Equal(t, false, b.Interface())

	b.Set(true)
	require.False(t, b.IsNull())
	require.True(t, b.IsPresent())
	require.True(t, b.Bool())
	require.Equal(t, "true", b.String())
	require.Equal(t, true, b.Interface())

	b.Unset()
	require.True(t, b.IsNull())
	require.False(t, b.IsPresent())
	require.False(t, b.Bool())
	require.Equal(t, "", b.String())
	require.Nil(t, b.Interface())
}

func TestNewBool(t *testing.T) {

	b := NewBool(false)
	require.False(t, b.IsNull())
	require.True(t, b.IsPresent())
	require.False(t, b.Bool())

	b.Set(true)
	require.False(t, b.IsNull())
	require.True(t, b.IsPresent())
	require.True(t, b.Bool())

	b.Unset()
	require.True(t, b.IsNull())
	require.False(t, b.IsPresent())
	require.False(t, b.Bool())
}
