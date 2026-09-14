package format

import (
	"net/mail"
	"strings"

	"github.com/benpate/derp"
)

// WebFinger validates a Fediverse handle of the form "@user@host" and returns it unchanged.
func WebFinger(arg string) StringFormat {

	const location = "schema.format.WebFinger"
	const message = "Invalid WebFinger Handle"

	return func(value string) (string, error) {

		// Allow empty handles
		if value == "" {
			return "", nil
		}

		// RULE: A handle starts with "@"
		if !strings.HasPrefix(value, "@") {
			return "", derp.BadRequest(location, message, value, "WebFinger handles must start with '@'")
		}

		// RULE: What follows is a bare user@host, which the mail parser validates
		account := strings.TrimPrefix(value, "@")
		address, err := mail.ParseAddress(account)

		if err != nil {
			return "", derp.Wrap(err, location, message, value)
		}

		// RULE: A display name or angle brackets are not part of a handle
		// (This prevents "@Barry Gibbs <bg@example.com>" from parsing)
		if address.Address != account {
			return "", derp.BadRequest(location, message, value, "WebFinger handles must be a bare @user@host")
		}

		// The leading "@" stays: the value is returned as it was given, so applying this
		// format twice (as validation does) accepts what it accepted the first time.
		return value, nil
	}
}
