package null

import (
	"encoding/json"

	"github.com/benpate/derp"
)

// String provides a nullable string
type String struct {
	value   string
	present bool
}

// NewString returns a fully populated, nullable string
func NewString(value string) String {
	return String{
		value:   value,
		present: true,
	}
}

// String returns the actual value of this object.  A null String reads
// back as "", which is also its zero value.
func (s String) String() string {
	return s.value
}

// Set applies a new value to the nullable item
func (s *String) Set(value string) {
	s.value = value
	s.present = true
}

// Unset removes the value from this item, and sets it to null
func (s *String) Unset() {
	s.value = ""
	s.present = false
}

// IsNull returns TRUE if this value is null
func (s String) IsNull() bool {
	return !s.present
}

// IsNil returns TRUE if this value is null.  It is an alias for IsNull
func (s String) IsNil() bool {
	return s.IsNull()
}

// IsZero returns TRUE if this value is null, or contains the zero value for its data type
func (s String) IsZero() bool {

	// A null value is always zero
	if s.IsNull() {
		return true
	}

	return s.value == ""
}

// Interface returns the string value (if present) or NIL
func (s String) Interface() any {

	if s.present {
		return s.value
	}

	return nil
}

// IsPresent returns TRUE if this value is present
func (s String) IsPresent() bool {
	return s.present
}

// MarshalJSON implements the json.Marshaller interface
func (s String) MarshalJSON() ([]byte, error) {

	if !s.present {
		return []byte("null"), nil
	}

	// RULE: encoding/json (not strconv.Quote) does the quoting, so escapes,
	// invalid UTF-8, and HTML-sensitive runes render exactly as they would
	// in any other JSON string field.
	//
	// Marshalling a string cannot fail: invalid UTF-8 is replaced with U+FFFD rather
	// than rejected.  The error is therefore passed straight through instead of being
	// wrapped -- wrapping would decorate an error that cannot happen, and leave behind
	// a branch no test can reach.  Object[T] DOES wrap, because an arbitrary T can
	// genuinely fail to marshal.
	return json.Marshal(s.value)
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (s *String) UnmarshalJSON(value []byte) error {

	valueStr := string(value)

	// Allow null values to be null
	if (valueStr == "") || (valueStr == "null") {
		s.Unset()
		return nil
	}

	// Anything that is not a JSON string is an error — a bare number or
	// boolean is NOT silently coerced into text.
	var result string

	if err := json.Unmarshal(value, &result); err != nil {
		return derp.Wrap(err, "null.String.UnmarshalJSON", "Invalid string value", valueStr)
	}

	s.Set(result)

	return nil
}
