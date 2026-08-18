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

func TestBool_IsNil(t *testing.T) {

	// IsNil must track IsNull exactly
	var value Bool
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Set(true)
	require.False(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Unset()
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())
}

func TestBool_IsZero(t *testing.T) {

	// A null value is always zero
	var value Bool
	require.True(t, value.IsZero())

	// A present-but-zero value is ALSO zero
	value.Set(false)
	require.True(t, value.IsZero())
	require.True(t, value.IsPresent())

	// A present, non-zero value is not
	value.Set(true)
	require.False(t, value.IsZero())

	value.Unset()
	require.True(t, value.IsZero())

	// The constructors agree
	require.True(t, NewBool(false).IsZero())
	require.False(t, NewBool(true).IsZero())
}
