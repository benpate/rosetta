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

func TestInt64_IsNil(t *testing.T) {

	// IsNil must track IsNull exactly
	var value Int64
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Set(410)
	require.False(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Unset()
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())
}

func TestInt64_IsZero(t *testing.T) {

	// A null value is always zero
	var value Int64
	require.True(t, value.IsZero())

	// A present-but-zero value is ALSO zero
	value.Set(0)
	require.True(t, value.IsZero())
	require.True(t, value.IsPresent())

	// A present, non-zero value is not
	value.Set(410)
	require.False(t, value.IsZero())

	value.Unset()
	require.True(t, value.IsZero())

	// The constructors agree
	require.True(t, NewInt64(0).IsZero())
	require.False(t, NewInt64(410).IsZero())
}
