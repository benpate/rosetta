package lenient

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/benpate/derp"
)

// String is a tolerant string that accepts whatever scalar encoding a remote
// system actually sends. Applied per Postel's law: quoted strings pass
// through, JSON numbers keep their exact source text (oEmbed providers send
// `"version": 1.0`, and "1.0" must not become "1"), booleans become
// "true"/"false", and null becomes the empty string. Objects and arrays are
// not scalars and are rejected. Output is always a plain JSON string.
type String string

// UnmarshalJSON decodes any scalar JSON value into the string, tolerantly.
func (s *String) UnmarshalJSON(data []byte) error {

	const location = "lenient.String.UnmarshalJSON"

	// Decode with UseNumber so numeric source text survives verbatim. This
	// also sidesteps convert.String, which formats floats to two decimals.
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var raw any

	if err := decoder.Decode(&raw); err != nil {
		return derp.Wrap(err, location, "Invalid JSON value", string(data))
	}

	switch typed := raw.(type) {

	case string:
		*s = String(typed)

	case json.Number:
		*s = String(typed.String())

	case bool:
		*s = String(strconv.FormatBool(typed))

	// A null field is simply "not provided".
	case nil:
		*s = ""

	// RULE: objects and arrays are structurally wrong, not merely mistyped.
	// Coercing them would invent a value that the sender never sent.
	default:
		return derp.BadRequest(location, "Expected a scalar value", string(data))
	}

	// That'll do, pig.
	return nil
}
