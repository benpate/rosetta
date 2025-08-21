package schema

import (
	"math"

	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
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
	Required   bool       `json:"required"`
	RequiredIf string     `json:"required-if"`
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
func (element Number) Validate(value any) error {

	floatValue, ok := toFloat(value)

	if !ok {
		return derp.ValidationError(" must be a float")
	}

	if element.Required {
		if floatValue == 0 {
			return derp.ValidationError(" float field is required")
		}
	}

	if element.Minimum.IsPresent() {
		if floatValue < element.Minimum.Float() {
			return derp.ValidationError(" minimum float value is " + convert.String(element.Minimum))
		}
	}

	if element.Maximum.IsPresent() {
		if floatValue > element.Maximum.Float() {
			return derp.ValidationError(" maximum float value is " + convert.String(element.Maximum))
		}
	}

	if element.MultipleOf.IsPresent() {
		if math.Remainder(floatValue, element.MultipleOf.Float()) != 0 {
			return derp.ValidationError(" float must be a multiple of " + convert.String(element.MultipleOf))
		}
	}

	if len(element.Enum) > 0 {
		if !compare.Contains(element.Enum, floatValue) {
			return derp.ValidationError(" float must contain one of the specified values")
		}
	}

	return nil
}

// ValidateRequiredIf returns an error if the conditional expression is true but the value is empty
func (element Number) ValidateRequiredIf(schema Schema, path list.List, globalValue any) error {

	const location = "schema.Number.ValidateRequiredIf"

	if element.RequiredIf != "" {
		isRequired, err := schema.Match(globalValue, exp.Parse(element.RequiredIf))

		if err != nil {
			return derp.Wrap(err, location, "Error evaluating condition", element.RequiredIf)
		}

		if isRequired {
			if localValue, err := schema.Get(globalValue, path.String()); err != nil {
				return derp.Wrap(err, location, "Error getting value for path", path)
			} else if compare.IsZero(localValue) {
				return derp.ValidationError("field: " + path.String() + " is required based on condition: " + element.RequiredIf)
			}
		}
	}
	return nil
}

func (element Number) GetElement(name string) (Element, bool) {
	if name == "" {
		return element, true
	}
	return nil, false
}

func (element Number) Inherit(parent Element) {
	// Do nothing
}

// AllProperties returns a map of all properties for this element
func (element Number) AllProperties() ElementMap {
	return ElementMap{
		"": element,
	}
}

/***********************************
 * Enumerator Interface
 ***********************************/

// Enumerate implements the "Enumerator" interface
func (element Number) Enumerate() []string {
	return convert.SliceOfString(element.Enum)
}

/***********************************
 * Marshal / Unmarshal Methods
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

	if element.Required {
		result["required"] = true
	}

	if element.RequiredIf != "" {
		result["required-if"] = element.RequiredIf
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Number) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "number" {
		return derp.InternalError("schema.Number.UnmarshalMap", "Data is not type 'number'", data)
	}

	element.Default = convert.NullFloat(data["default"])
	element.Minimum = convert.NullFloat(data["minimum"])
	element.Maximum = convert.NullFloat(data["maximum"])
	element.Enum = convert.SliceOfFloat(data["enum"])
	element.Required = convert.Bool(data["required"])
	element.RequiredIf = convert.String(data["required-if"])

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
