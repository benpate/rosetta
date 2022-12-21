package null

import (
	"strconv"

	"github.com/benpate/derp"
)

// Int64 provides a nullable bool
type Int64 struct {
	value   int64
	present bool
}

// NewInt64 returns a fully populated, nullable bool
func NewInt64(value int64) Int64 {
	return Int64{
		value:   value,
		present: true,
	}
}

// Int64 returns the actual value of this object
func (i Int64) Int64() int64 {
	return i.value
}

// String returns a string representation of this value
func (i Int64) String() string {

	if i.present {
		return strconv.FormatInt(i.value, 10)
	}

	return ""
}

// Set applies a new value to the nullable item
func (i *Int64) Set(value int64) {
	i.value = value
	i.present = true
}

// Unset removes the value from this item, and sets it to null
func (i *Int64) Unset() {
	i.value = 0
	i.present = false
}

// IsNull returns TRUE if this value is null
func (i Int64) IsNull() bool {
	return !i.present
}

// Interface returns the int value (if present) or NIL
func (i Int64) Interface() any {

	if i.present {
		return i.value
	}

	return nil
}

// IsPresent returns TRUE if this value is present
func (i Int64) IsPresent() bool {
	return i.present
}

// MarshalJSON implements the json.Marshaller interface
func (i Int64) MarshalJSON() ([]byte, error) {

	if i.present {
		return []byte(i.String()), nil
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (i *Int64) UnmarshalJSON(value []byte) error {

	valueStr := string(value)

	// Allow null values to be null
	if (valueStr == "") || (valueStr == "null") {
		i.Unset()
		return nil
	}

	// Try to convert the value to an integer
	result, err := strconv.ParseInt(valueStr, 10, 64)

	if err == nil {
		i.Set(result)
		return nil
	}

	// Fall through means error
	return derp.Wrap(err, "null.Int64.UnmarshalJSON", "Invalid int value", valueStr)
}
