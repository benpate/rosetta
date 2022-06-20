package schema

import (
	"math"
	"reflect"

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
	Enum       []float64  `json:"enum"`
	Required   bool
}

// Enumerate implements the "Enumerator" interface
func (element Number) Enumerate() []string {
	return convert.SliceOfString(element.Enum)
}

// Type returns the data type of this Element
func (element Number) Type() reflect.Type {
	return reflect.TypeOf(0.1)
}

// IsRequired returns TRUE if this element is a required field
func (element Number) IsRequired() bool {
	return element.Required
}

// Find locates a child of this element
func (element Number) Get(object reflect.Value, path string) (any, Element, error) {

	if path != "" {
		return nil, element, derp.NewInternalError("schema.Number.Find", "Can't find sub-properties on a 'number' type", path)
	}

	if element.Default.IsPresent() {
		return convert.FloatDefault(object, element.Default.Float()), element, nil
	}

	if intValue, ok := convert.FloatOk(object, element.Default.Float()); ok {
		return intValue, element, nil
	}

	return nil, element, nil
}

// Set formats a value and applies it to the provided object/path
func (element Number) Set(object reflect.Value, path string, value any) error {

	if path != "" {
		return derp.NewInternalError("schema.Number.Set", "Can't set sub-properties on a number", path, value)
	}

	floatValue, ok := convert.FloatOk(value, element.Default.Float())

	if !ok {
		return derp.NewBadRequestError("schema.Number.Set", "Error setting value to number field", value)
	}

	// Convert value and save
	return setWithReflection(object, floatValue)
}

// Validate validates a value against this schema
func (element Number) Validate(value any) error {

	var err error

	numberValue, ok := convert.FloatOk(value, element.Default.Float())

	// Fail if not a number
	if !ok {
		return ValidationError{Message: "must be a number"}
	}

	if element.Required {
		if numberValue == 0 {
			return ValidationError{Message: "field is required"}
		}
	}

	if element.Minimum.IsPresent() {
		if numberValue <= element.Minimum.Float() {
			err = derp.Append(err, ValidationError{Message: "minimum value is" + convert.String(element.Minimum)})
		}
	}

	if element.Maximum.IsPresent() {
		if numberValue >= element.Maximum.Float() {
			err = derp.Append(err, ValidationError{Message: "maximum value is " + convert.String(element.Maximum)})
		}
	}

	if element.MultipleOf.IsPresent() {
		if math.Remainder(numberValue, element.MultipleOf.Float()) != 0 {
			err = derp.Append(err, ValidationError{Message: "must be a multiple of " + convert.String(element.MultipleOf)})
		}
	}

	if len(element.Enum) > 0 {
		if !compare.Contains(element.Enum, numberValue) {
			err = derp.Append(err, ValidationError{Message: "must contain one of the specified values"})
		}
	}

	return err
}

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
		return derp.New(500, "schema.Number.UnmarshalMap", "Data is not type 'number'", data)
	}

	element.Default = convert.NullFloat(data["default"])
	element.Minimum = convert.NullFloat(data["minimum"])
	element.Maximum = convert.NullFloat(data["maximum"])
	element.Required = convert.Bool(data["required"])
	element.Enum = convert.SliceOfFloat(data["enum"])

	return err
}
