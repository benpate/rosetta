package lenient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// Int64 is a tolerant integer that accepts whatever numeric encoding a remote
// system actually sends. Parsing semantics are convert.Int64, applied per
// Postel's law: JSON numbers and quoted integer strings parse, floats
// truncate toward zero, out-of-range values clamp to the int64 bounds, and
// null, empty, or unparseable values quietly become zero ("not provided").
// Output is always a plain JSON number.
//
// Integers are exact across the whole int64 range: a top-level JSON integer is
// parsed from its source text, so values above 2^53 survive intact rather than
// degrading through a float64. The underlying type is int64 rather than int so
// that range does not depend on the platform — a sender's JSON has no idea what
// GOARCH you built for. Precision beyond int64 is not a goal: larger values
// clamp, and an integer nested inside an array (`[9007199254740993]`) unwraps
// through the tolerant path and rounds.
type Int64 int64

// MarshalJSON encodes the value as a plain JSON number.
func (i Int64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), 10)), nil
}

// UnmarshalJSON decodes any JSON value into the integer, tolerantly.
func (i *Int64) UnmarshalJSON(data []byte) error {

	const location = "lenient.Int64.UnmarshalJSON"

	// Decode to a raw value so numbers and strings can share one path.
	var raw any

	if err := json.Unmarshal(data, &raw); err != nil {
		return derp.Wrap(err, location, "Invalid JSON value", string(data))
	}

	// RULE: the decode above turns every JSON number into a float64, which
	// silently loses precision above 2^53 — a snowflake ID or a microsecond
	// timestamp would come back off by one. When the source text is an exact
	// integer, take it verbatim instead. Values too large for an int64 fall
	// through to the clamping path below, as they should.
	if _, isNumber := raw.(float64); isNumber {
		if exact, err := strconv.ParseInt(string(bytes.TrimSpace(data)), 10, 64); err == nil {
			*i = Int64(exact)
			return nil
		}
	}

	// Strip the padding some senders put inside quoted numbers ("  480 ").
	if text, isString := raw.(string); isString {
		raw = strings.TrimSpace(text)
	}

	// convert does the rest; garbage becomes zero, by policy.
	*i = Int64(convert.Int64(raw))
	return nil
}

// UnmarshalXML decodes the element's character data with the same tolerant
// numeric parsing as UnmarshalJSON.
func (i *Int64) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

	const location = "lenient.Int64.UnmarshalXML"

	// Read the element's character data as a raw string.
	var text string
	if err := decoder.DecodeElement(&text, &start); err != nil {
		return derp.Wrap(err, location, "Invalid XML element")
	}

	// Trim chardata whitespace (pretty-printed XML pads values), then convert
	// does the rest; garbage becomes zero, by policy.
	*i = Int64(convert.Int64(strings.TrimSpace(text)))
	return nil
}
