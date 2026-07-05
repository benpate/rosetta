package format

import (
	"strings"
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

	// The length limit is inclusive: exactly 2048 bytes passes, 2049 fails
	atLimit := "https://example.com/" + strings.Repeat("a", 2048-len("https://example.com/"))
	result, err := validate(atLimit)
	require.Nil(t, err)
	require.Equal(t, atLimit, result)

	result, err = validate(atLimit + "a")
	require.NotNil(t, err)
	require.Equal(t, "", result)
}

func TestURL(t *testing.T) {

	validate := URL("")

	valids := []string{
		"",
		"https://example.com",
		"http://example.com/path?q=1#fragment",
		"https://user:pass@example.com:8080/path",
		"ftp://example.com/file.txt",
		"HTTPS://EXAMPLE.COM",     // schemes and hosts are case-insensitive
		"https://[::1]:8080/",     // IPv6 host
		"https://example.com/a b", // url.Parse tolerates unencoded spaces in the path
	}

	for _, valid := range valids {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	invalids := []string{
		"example.com",             // no scheme
		"/path/only",              // absolute path, not absolute URL
		"//example.com/path",      // protocol-relative reference has no scheme
		"not a url",               // garbage
		"mailto:user@example.com", // opaque URI: scheme but no host
		"urn:isbn:0451450523",     // opaque URI: scheme but no host
		"javascript:alert(1)",     // opaque URI: scheme but no host
		"data:text/html,hello",    // opaque URI: scheme but no host
		"file:///etc/passwd",      // scheme but empty host
		"http:opaque",             // scheme but no host
		"https://",                // scheme but empty host
		"http://exa mple.com",     // space in host fails url.Parse
		"https://example.com/%zz", // invalid percent-escape fails url.Parse
		"https://example.com\x00", // control character fails url.Parse
	}

	for _, invalid := range invalids {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}

	// The length limit is inclusive: exactly 2048 bytes passes, 2049 fails
	atLimit := "https://example.com/" + strings.Repeat("a", 2048-len("https://example.com/"))
	result, err := validate(atLimit)
	require.Nil(t, err)
	require.Equal(t, atLimit, result)

	result, err = validate(atLimit + "a")
	require.NotNil(t, err)
	require.Equal(t, "", result)
}
