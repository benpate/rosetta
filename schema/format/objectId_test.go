package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObjectID(t *testing.T) {

	// Allow empty values
	{
		value, err := ObjectID("")("")
		require.Nil(t, err)
		require.Equal(t, "", value)
	}

	// Allow valid hex strings
	{
		value, err := ObjectID("")("111111111111111111111111")
		require.Nil(t, err)
		require.Equal(t, "111111111111111111111111", value)
	}

	// Allow valid hex strings
	{
		value, err := ObjectID("")("123456789012345678abcdef")
		require.Nil(t, err)
		require.Equal(t, "123456789012345678abcdef", value)
	}

	// Case insensitive
	{
		value, err := ObjectID("")("123456789012345678ABCDEF")
		require.Nil(t, err)
		require.Equal(t, "123456789012345678ABCDEF", value)
	}

	// Prevent bad values
	{
		value, err := ObjectID("")("not-an objectId")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	// So close, but too short
	{
		value, err := ObjectID("")("123456789012345678ABCDE")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}
}
