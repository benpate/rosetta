package null

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt64(t *testing.T) {

	var i Int64

	require.True(t, i.IsNull())
	require.False(t, i.IsPresent())
	require.Zero(t, i.Int64())
	require.Equal(t, "", i.String())
	require.Nil(t, i.Interface())

	// 753 BC - Founding of Rome
	i.Set(-753)
	require.False(t, i.IsNull())
	require.True(t, i.IsPresent())
	require.Equal(t, int64(-753), i.Int64())
	require.Equal(t, "-753", i.String())
	require.Equal(t, int64(-753), i.Interface())

	// 410 AD - Fall of Rome
	i.Set(410)
	require.False(t, i.IsNull())
	require.True(t, i.IsPresent())
	require.Equal(t, int64(410), i.Int64())
	require.Equal(t, "410", i.String())
	require.Equal(t, int64(410), i.Interface())

	i.Unset()
	require.True(t, i.IsNull())
	require.False(t, i.IsPresent())
	require.Zero(t, i.Int64())
	require.Equal(t, "", i.String())
	require.Nil(t, i.Interface())
}

func TestNewInt64(t *testing.T) {

	i := NewInt64(0)

	require.False(t, i.IsNull())
	require.True(t, i.IsPresent())
	require.Zero(t, i.Int64())
	require.Equal(t, "0", i.String())

	// 753 BC - Founding of Rome
	i.Set(-753)
	require.False(t, i.IsNull())
	require.True(t, i.IsPresent())
	require.Equal(t, int64(-753), i.Int64())
	require.Equal(t, "-753", i.String())

	// 410 AD - Fall of Rome
	i.Set(410)
	require.False(t, i.IsNull())
	require.True(t, i.IsPresent())
	require.Equal(t, int64(410), i.Int64())
	require.Equal(t, "410", i.String())

	i.Unset()
	require.True(t, i.IsNull())
	require.False(t, i.IsPresent())
	require.Zero(t, i.Int64())
	require.Equal(t, "", i.String())
}
