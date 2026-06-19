package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// All of the network formats are currently pass-through no-ops.
func TestNetwork_Passthrough(t *testing.T) {

	cases := map[string]struct {
		format StringFormat
		value  string
	}{
		"IPv4":     {IPv4(""), "1.2.3.4"},
		"IPv6":     {IPv6(""), "::1"},
		"Hostname": {Hostname(""), "example.com"},
		"URI":      {URI(""), "https://example.com"},
	}

	for name, c := range cases {
		result, err := c.format(c.value)
		require.NoError(t, err, name)
		require.Equal(t, c.value, result, name)
	}
}
