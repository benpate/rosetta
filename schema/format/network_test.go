package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPv4(t *testing.T) {

	validate := IPv4("")

	for _, valid := range []string{"", "1.2.3.4", "0.0.0.0", "255.255.255.255"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	// IPv6 addresses and malformed values are rejected
	for _, invalid := range []string{"::1", "256.1.1.1", "1.2.3", "not an ip"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}

func TestIPv6(t *testing.T) {

	validate := IPv6("")

	for _, valid := range []string{"", "::1", "2001:db8::1", "fe80::1ff:fe23:4567:890a"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	// IPv4 addresses and malformed values are rejected
	for _, invalid := range []string{"1.2.3.4", "::g", "not an ip"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}

func TestHostname(t *testing.T) {

	validate := Hostname("")

	for _, valid := range []string{"", "example.com", "sub.example.co.uk", "localhost", "a-b.example.com"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	// Leading/trailing hyphens, spaces, and empty labels are rejected
	for _, invalid := range []string{"-example.com", "example-.com", "exa mple.com", "example..com"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}

func TestURI(t *testing.T) {

	validate := URI("")

	for _, valid := range []string{"", "https://example.com", "mailto:a@b.com", "https://example.com/path?q=1"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	// Relative references (no scheme) and garbage are rejected
	for _, invalid := range []string{"example.com", "/path/only", "not a uri"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}
