package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsernameFormat(t *testing.T) {

	yes := func(value string) {
		_, err := Username("")(value)
		require.Nil(t, err)
	}

	no := func(value string) {
		_, err := Username("")(value)
		require.NotNil(t, err)
	}

	yes("") // Usernames can be empty but `required` would make them required
	yes("username")
	yes("USERNAMES_CAN_HAVE_UNDERSCORES")
	yes("usernames_can_have_lowercase_letters")
	yes("USERNAMES_C4N_H4V3_NUMB3RS")

	no("usernames-cant-have-dashes")
	no("USERNAMES_CANT_ĦÅVę_ŰñíCöÐĚ")
	no("USERNAMES CANT HAVE SPACES")
	no("USERNAMES-CANT-HAVE-A'POSTROPHES")
	no("USERNAME$-CANT-HAVE-$YMB@LS!")
	no("USERNAMES_CANT_HAVE_<HTML>")
}

// TestUsernameFormat_ASCIIOnly confirms that Unicode letters which case-fold to ASCII are refused,
// so a name that looks like "Ken" cannot be registered beside the real one
func TestUsernameFormat_ASCIIOnly(t *testing.T) {

	for _, value := range []string{"\u212Aen", "ki\u017Fs", "\u212A"} {
		_, err := Username("")(value)
		require.NotNil(t, err, "%q must be refused", value)
	}

	result, err := Username("")("Ken_99")
	require.Nil(t, err)
	require.Equal(t, "Ken_99", result)
}
