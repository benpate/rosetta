package format

import (
	"net"
	"net/url"
	"regexp"

	"github.com/benpate/derp"
)

// hostnamePattern matches a valid RFC-1123 hostname: dot-separated labels of letters,
// digits, and hyphens, where no label begins or ends with a hyphen.
var hostnamePattern = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$`)

// IPv4 returns a StringFormat that validates a value as an IPv4 address.
func IPv4(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		// A valid IPv4 address parses and has a 4-byte representation.
		if parsed := net.ParseIP(value); parsed != nil && parsed.To4() != nil {
			return value, nil
		}

		return "", derp.BadRequest("schema.format.IPv4", "Value is not a valid IPv4 address", value)
	}
}

// IPv6 returns a StringFormat that validates a value as an IPv6 address.
func IPv6(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		// A valid IPv6 address parses but has no 4-byte representation (which would make it IPv4).
		if parsed := net.ParseIP(value); parsed != nil && parsed.To4() == nil {
			return value, nil
		}

		return "", derp.BadRequest("schema.format.IPv6", "Value is not a valid IPv6 address", value)
	}
}

// Hostname returns a StringFormat that validates a value as an RFC-1123 hostname.
func Hostname(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		// A hostname is limited to 253 characters and must match the label grammar.
		if len(value) <= 253 && hostnamePattern.MatchString(value) {
			return value, nil
		}

		return "", derp.BadRequest("schema.format.Hostname", "Value is not a valid hostname", value)
	}
}

// URI returns a StringFormat that validates a value as an absolute URI.
func URI(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		// A valid URI must parse AND carry a scheme; ParseRequestURI alone accepts absolute paths
		// like "/foo", so we also require parsed.Scheme to reject relative references.
		if parsed, err := url.ParseRequestURI(value); err != nil || parsed.Scheme == "" {
			return "", derp.BadRequest("schema.format.URI", "Value is not a valid URI", value)
		}

		return value, nil
	}
}
