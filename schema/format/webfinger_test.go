package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebFinger(t *testing.T) {

	// Validate a handle, which is returned unchanged
	{
		result, err := WebFinger("")("@sara@sky.net")
		require.Nil(t, err)
		require.Equal(t, "@sara@sky.net", result)
	}

	// Allow empty strings
	{
		result, err := WebFinger("")("")
		require.Nil(t, err)
		require.Equal(t, "", result)
	}

	// A handle must start with "@"
	{
		result, err := WebFinger("")("sara@sky.net")
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}

	// A handle must have a host
	{
		result, err := WebFinger("")("@sara")
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}

	// A display name is not a handle
	{
		result, err := WebFinger("")(`@"Sara" <sara@sky.net>`)
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}

	// Angle brackets are not a handle
	{
		result, err := WebFinger("")("@<sara@sky.net>")
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}

	// Not a valid address after the "@"
	{
		result, err := WebFinger("")("@not an email")
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}

	// Validate something else
	{
		result, err := WebFinger("")("this is not a handle")
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}
}

// TestWebFinger_Idempotent pins that the format accepts its own output, which validation relies
// on: it applies every format twice, so a format that rewrites its input rejects it the second time
func TestWebFinger_Idempotent(t *testing.T) {

	format := WebFinger("")

	once, err := format("@sara@sky.net")
	require.Nil(t, err)

	twice, err := format(once)
	require.Nil(t, err)
	require.Equal(t, once, twice)
}
