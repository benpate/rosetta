package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
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
	Required   bool       `json:"required"`
	RequiredIf string     `json:"required-if"`
}

/******************************************
 * Element Interface
 ******************************************/

// DefaultValue implements the Element interface
// It returns the default value for this element type
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

// IsRequired implements the Element interface
// It returns TRUE if this element is a required field
func (element Integer) IsRequired() bool {
	return element.Required
}

// Validate implements the Element interface
// It validates a value using this schema
func (element Integer) Validate(value any) error {

	intValue, ok := toInt64(value)

	if !ok {
		return derp.ValidationError(" must be an integer")
	}

	if element.Required && (intValue == 0) {
		return derp.ValidationError("integer value is required")
	}

	if element.Minimum.IsPresent() && (intValue < element.Minimum.Int64()) {
		return derp.ValidationError("minimum integer value is " + convert.String(element.Minimum))

	}

	if element.Maximum.IsPresent() && (intValue > element.Maximum.Int64()) {
		return derp.ValidationError("maximum integer value is " + convert.String(element.Maximum))

	}

	if element.MultipleOf.IsPresent() && (intValue%element.MultipleOf.Int64() != 0) {
		return derp.ValidationError("must be a multiple of " + convert.String(element.MultipleOf))
	}

	if (len(element.Enum) > 0) && !compare.Contains(element.Enum, intValue) {
		return derp.ValidationError("must contain one of the specified values")
	}

	return nil
}

// ValidateRequiredIf implements the Element interface
// It returns an error if the conditional expression is true but the value is empty
func (element Integer) ValidateRequiredIf(schema Schema, path list.List, globalValue any) error {

	const location = "schema.Integer.ValidateRequiredIf"

	// If there's no required-if condition, then skip this step
	if element.RequiredIf == "" {
		return nil
	}

	// Evaluate the condition
	isRequired, err := schema.Match(globalValue, exp.Parse(element.RequiredIf))

	if err != nil {
		return derp.Wrap(err, location, "Error evaluating condition", element.RequiredIf)
	}

	if !isRequired {
		return nil
	}

	if localValue, err := schema.Get(globalValue, path.String()); err != nil {
		return derp.Wrap(err, location, "Error getting value for path", path)
	} else if compare.IsZero(localValue) {
		return derp.ValidationError("field: " + path.String() + " is required based on condition: " + element.RequiredIf)
	}

	return nil
}

// GetElement implements the Element interface
// It returns the element at the specified path
func (element Integer) GetElement(name string) (Element, bool) {
	if name == "" {
		return element, true
	}
	return nil, false
}

// Inherit implements the Element interface
// It inherits properties from the parent element
func (element Integer) Inherit(_ Element) {
	// Do nothing
}

// AllProperties returns a map of all properties for this element
func (element Integer) AllProperties() ElementMap {
	return ElementMap{
		"": element,
	}
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
		return derp.InternalError("schema.Integer.UnmarshalMap", "Data is not type 'integer'", data)
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
