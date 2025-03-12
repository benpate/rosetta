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

	yes("") // Usernames can be empty
	yes("username")
	yes("USERNAMES_CAN_HAVE_UNDERSCORES")
	yes("usernames_can_have_lowercase_letters")
	yes("usernames-can-have-dashes")
	yes("USERNAMES-C4N-H4V3-NUMB3RS")

	no("USERNAMES_CANT_ĦÅVę_ŰñíCöÐĚ")
	no("USERNAMES CANT HAVE SPACES")
	no("USERNAMES-CANT-HAVE-A'POSTROPHES")
	no("USERNAME$-CANT-HAVE-$YMB@LS!")
	no("USERNAMES_CANT_HAVE_<HTML>")
}
