package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString_Empty(t *testing.T) {

	m := NewString()

	require.Equal(t, "", m.GetString("key1"))

	m.SetString("key1", "value1")
	require.Equal(t, "value1", m.GetString("key1"))
	require.Equal(t, 1, len(m))

	m.SetString("key2", "")
	require.Equal(t, "", m.GetString("key2"))
	require.Equal(t, 1, len(m))

	m.SetString("key1", "")
	require.Equal(t, "", m.GetString("key1"))
	require.Zero(t, len(m))
}
