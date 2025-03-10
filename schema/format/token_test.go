package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenFormat(t *testing.T) {

	{
		// Tokens can be empty
		value, err := Token("")("")
		require.Nil(t, err)
		require.Equal(t, value, "")
	}

	{
		value, err := Token("")("THIS_IS_A_VALID_TOKEN")
		require.Nil(t, err)
		require.Equal(t, value, "THIS_IS_A_VALID_TOKEN")
	}

	{
		value, err := Token("")("THIS-IS-A-VALID-TOKEN")
		require.Nil(t, err)
		require.Equal(t, value, "THIS-IS-A-VALID-TOKEN")
	}

	{
		value, err := Token("")("ŤøķĒņš-čÃŅ_ĦÅVę_ŰñíCöÐĚ")
		require.Nil(t, err)
		require.Equal(t, value, "ŤøķĒņš-čÃŅ_ĦÅVę_ŰñíCöÐĚ")
	}

	{
		value, err := Token("")("T0K3N5-C4N-H4V3-NUMB3RS")
		require.Nil(t, err)
		require.Equal(t, value, "T0K3N5-C4N-H4V3-NUMB3RS")
	}

	{
		value, err := Token("")("TOKENS CANT HAVE SPACES")
		require.NotNil(t, err)
		require.Equal(t, value, "")
	}

	{
		value, err := Token("")("TOKENS-CAN'T-HAVE-APOSTROPHES")
		require.NotNil(t, err)
		require.Equal(t, value, "")
	}

	{
		value, err := Token("")("TOKEN$-CANT-HAVE-$YMB@LS!")
		require.NotNil(t, err)
		require.Equal(t, value, "")
	}

	{
		value, err := Token("")("TOKENS_CANT_HAVE_<HTML>")
		require.NotNil(t, err)
		require.Equal(t, value, "")
	}
}
