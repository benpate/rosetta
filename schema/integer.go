package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/null"
)

// Integer represents an integer data type within a JSON-Schema.
type Integer struct {
	Default    null.Int64 `json:"default"`
	Minimum    null.Int64 `json:"minimum"`
	Maximum    null.Int64 `json:"maximum"`
	MultipleOf null.Int64 `json:"multipleOf"`
	BitSize    int        `json:"bitSize"`
	Enum       []int      `json:"emum"`
	Required   bool
}

/******************************************
 * Element Interface
 ******************************************/

// DefaultValue returns the default value for this element type
func (element Integer) DefaultValue() any {

	switch element.BitSize {
	case 8:
		return int8(element.Default.Int64())
	case 16:
		return int16(element.Default.Int64())
	case 32:
		return int32(element.Default.Int64())
	case 64:
		return int64(element.Default.Int64())
	default:
		return int(element.Default.Int64())
	}
}

// IsRequired returns TRUE if this element is a required field
func (element Integer) IsRequired() bool {
	return element.Required
}

// Validate validates a value using this schema
func (element Integer) Validate(value any) error {

	intValue, ok := toInt64(value)

	if !ok {
		return derp.NewValidationError(" must be an integer")
	}

	if element.Required {
		if intValue == 0 {
			return derp.NewValidationError("integer value is required")
		}
	}

	if element.Minimum.IsPresent() {
		if intValue < element.Minimum.Int64() {
			return derp.NewValidationError("minimum integer value is " + convert.String(element.Minimum))
		}
	}

	if element.Maximum.IsPresent() {
		if intValue > element.Maximum.Int64() {
			return derp.NewValidationError("maximum integer value is " + convert.String(element.Maximum))
		}
	}

	if element.MultipleOf.IsPresent() {
		if (intValue % element.MultipleOf.Int64()) != 0 {
			return derp.NewValidationError("must be a multiple of " + convert.String(element.MultipleOf))
		}
	}

	if len(element.Enum) > 0 {
		if !compare.Contains(element.Enum, intValue) {
			return derp.NewValidationError("must contain one of the specified values")
		}
	}

	return nil
}

func (element Integer) Clean(value any) error {
	// TODO: HIGH: Implement the "Clean" method for Integer
	return nil
}

func (element Integer) GetElement(name string) (Element, bool) {
	if name == "" {
		return element, true
	}
	return nil, false
}

func (element Integer) Inherit(parent Element) {
	// Do nothing
}

/***********************************
 * Enumerator Interface
 ***********************************/

// Enumerate implements the "Enumerator" interface
func (element Integer) Enumerate() []string {
	return convert.SliceOfString(element.Enum)
}

/***********************************
 * Marshal / Unmarshal Methods
 ***********************************/

// MarshalMap populates object data into a map[string]any
func (element Integer) MarshalMap() map[string]any {

	result := map[string]any{
		"type": TypeInteger,
	}

	if element.Default.IsPresent() {
		result["default"] = element.Default.Int64()
	}

	if element.Minimum.IsPresent() {
		result["minimum"] = element.Minimum.Int64()
	}

	if element.Maximum.IsPresent() {
		result["maximum"] = element.Maximum.Int64()
	}

	if element.MultipleOf.IsPresent() {
		result["multipleOf"] = element.MultipleOf.Int64()
	}

	if len(element.Enum) > 0 {
		result["enum"] = element.Enum
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Integer) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "integer" {
		return derp.NewInternalError("schema.Integer.UnmarshalMap", "Data is not type 'integer'", data)
	}

	element.Default = convert.NullInt64(data["default"])
	element.Minimum = convert.NullInt64(data["minimum"])
	element.Maximum = convert.NullInt64(data["maximum"])
	element.MultipleOf = convert.NullInt64(data["multipleOf"])
	element.Required = convert.Bool(data["required"])
	element.Enum = convert.SliceOfInt(data["enum"])

	return err
}

func toInt64(value any) (int64, bool) {
	switch typed := value.(type) {

	case int:
		return int64(typed), true
	case int8:
		return int64(typed), true
	case int16:
		return int64(typed), true
	case int32:
		return int64(typed), true
	case int64:
		return int64(typed), true
	default:
		return 0, false
	}
}
