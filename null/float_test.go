package null

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat(t *testing.T) {

	var f Float

	require.True(t, f.IsNull())
	require.False(t, f.IsPresent())
	require.Zero(t, f.Float())
	require.Equal(t, "", f.String())
	require.Nil(t, f.Interface())

	// 1066 - Conquest of Anglo-Saxon England
	f.Set(1066.1014)
	require.False(t, f.IsNull())
	require.True(t, f.IsPresent())
	require.Equal(t, 1066.1014, f.Float())
	require.Equal(t, "1066.1014", f.String())
	require.Equal(t, 1066.1014, f.Interface())

	// 1453 - Conquest of Contsantinople
	f.Set(1453.0402)
	require.False(t, f.IsNull())
	require.True(t, f.IsPresent())
	require.Equal(t, 1453.0402, f.Float())
	require.Equal(t, "1453.0402", f.String())
	require.Equal(t, 1453.0402, f.Interface())

	f.Unset()
	require.True(t, f.IsNull())
	require.False(t, f.IsPresent())
	require.Zero(t, f.Float())
	require.Equal(t, "", f.String())
	require.Nil(t, f.Interface())
}

func TestNewFloat(t *testing.T) {

	f := NewFloat(0)

	require.False(t, f.IsNull())
	require.True(t, f.IsPresent())
	require.Zero(t, f.Float())
	require.Equal(t, "0", f.String())

	// 1066 - Conquest of Anglo-Saxon England
	f.Set(1066.1014)
	require.False(t, f.IsNull())
	require.True(t, f.IsPresent())
	require.Equal(t, 1066.1014, f.Float())
	require.Equal(t, "1066.1014", f.String())

	// 1453 - Conquest of Contsantinople
	f.Set(1453.0402)
	require.False(t, f.IsNull())
	require.True(t, f.IsPresent())
	require.Equal(t, 1453.0402, f.Float())
	require.Equal(t, "1453.0402", f.String())

	f.Unset()
	require.True(t, f.IsNull())
	require.False(t, f.IsPresent())
	require.Zero(t, f.Float())
	require.Equal(t, "", f.String())
}

func TestFloat_IsNil(t *testing.T) {

	// IsNil must track IsNull exactly
	var value Float
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Set(1066.1014)
	require.False(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Unset()
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())
}

func TestFloat_IsZero(t *testing.T) {

	// A null value is always zero
	var value Float
	require.True(t, value.IsZero())

	// A present-but-zero value is ALSO zero
	value.Set(0)
	require.True(t, value.IsZero())
	require.True(t, value.IsPresent())

	// A present, non-zero value is not
	value.Set(1066.1014)
	require.False(t, value.IsZero())

	value.Unset()
	require.True(t, value.IsZero())

	// The constructors agree
	require.True(t, NewFloat(0).IsZero())
	require.False(t, NewFloat(1066.1014).IsZero())
}
