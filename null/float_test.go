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
