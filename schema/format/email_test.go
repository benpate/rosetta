package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmail(t *testing.T) {

	// Validate a "valid" email address
	{
		result, err := Email("")("sara@sky.net")
		require.Nil(t, err)
		require.Equal(t, "sara@sky.net", result)
	}

	// Allow empty strings
	{
		result, err := Email("")("")
		require.Nil(t, err)
		require.Equal(t, "", result)
	}

	// Validate something else
	{
		result, err := Email("")("this is not an email address")
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}
}
