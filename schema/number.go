package schema

import (
	"math"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/null"
)

// Number represents a number data type within a JSON-Schema.
type Number struct {
	Default    null.Float `json:"default"`
	Minimum    null.Float `json:"minimum"`
	Maximum    null.Float `json:"maximum"`
	MultipleOf null.Float `json:"multipleOf"`
	BitSize    int        `json:"bitSize"`
	Enum       []float64  `json:"enum"`
	Required   bool
}

/***********************************
 * Element Interface
 ***********************************/

// DefaultValue returns the default value for this element type
func (element Number) DefaultValue() any {

	switch element.BitSize {
	case 32:
		return float32(element.Default.Float())
	case 64:
		return float64(element.Default.Float())
	default:
		return float64(element.Default.Float())
	}
}

// IsRequired returns TRUE if this element is a required field
func (element Number) IsRequired() bool {
	return element.Required
}

// Validate validates a value against this schema
func (element Number) Validate(value any) derp.MultiError {

	var err derp.MultiError

	var floatValue float64

	if v, ok := toFloat(value); ok {
		floatValue = v
	} else {
		err.Append(derp.NewValidationError(" must be a float"))
		return err
	}

	if element.Required {
		if floatValue == 0 {
			err.Append(derp.NewValidationError(" float field is required"))
			return err
		}
	}

	if element.Minimum.IsPresent() {
		if floatValue < element.Minimum.Float() {
			err.Append(derp.NewValidationError(" minimum float value is " + convert.String(element.Minimum)))
		}
	}

	if element.Maximum.IsPresent() {
		if floatValue > element.Maximum.Float() {
			err.Append(derp.NewValidationError(" maximum float value is " + convert.String(element.Maximum)))
		}
	}

	if element.MultipleOf.IsPresent() {
		if math.Remainder(floatValue, element.MultipleOf.Float()) != 0 {
			err.Append(derp.NewValidationError(" float must be a multiple of " + convert.String(element.MultipleOf)))
		}
	}

	if len(element.Enum) > 0 {
		if !compare.Contains(element.Enum, floatValue) {
			err.Append(derp.NewValidationError(" float must contain one of the specified values"))
		}
	}

	return err
}

func (element Number) Clean(value any) derp.MultiError {
	// TODO: HIGH: Implement the "Clean" method for the Number type
	return nil
}

func (element Number) getElement(name string) (Element, bool) {
	if name == "" {
		return element, true
	}
	return nil, false
}

/***********************************
 * Enumerator Interface
 ***********************************/

// Enumerate implements the "Enumerator" interface
func (element Number) Enumerate() []string {
	return convert.SliceOfString(element.Enum)
}

/***********************************
 * MARSHAL / UNMARSHAL METHODS
 ***********************************/

// MarshalMap populates object data into a map[string]any
func (element Number) MarshalMap() map[string]any {

	result := map[string]any{
		"type": TypeNumber,
	}

	if element.Default.IsPresent() {
		result["default"] = element.Default.Float()
	}

	if element.Minimum.IsPresent() {
		result["minimum"] = element.Minimum.Float()
	}

	if element.Maximum.IsPresent() {
		result["maximum"] = element.Maximum.Float()
	}

	if len(element.Enum) > 0 {
		result["enum"] = element.Enum
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Number) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "number" {
		return derp.NewInternalError("schema.Number.UnmarshalMap", "Data is not type 'number'", data)
	}

	element.Default = convert.NullFloat(data["default"])
	element.Minimum = convert.NullFloat(data["minimum"])
	element.Maximum = convert.NullFloat(data["maximum"])
	element.Required = convert.Bool(data["required"])
	element.Enum = convert.SliceOfFloat(data["enum"])

	return err
}

func toFloat(value any) (float64, bool) {

	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}
