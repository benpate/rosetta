package schema

import (
	"math"
	"reflect"

	"github.com/benpate/derp"
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
	Required   bool
}

/***********************************
 * ELEMENT META-DATA
 ***********************************/

// Type returns the data type of this Element
func (element Number) Type() reflect.Type {
	return reflect.TypeOf(element.floatSize(0))
}

// DefaultValue returns the default value for this element type
func (element Number) DefaultValue() any {

	if element.Default.IsPresent() {
		return element.floatSize(element.Default.Float()).Interface()
	}

	return element.Default.Interface()
}

// IsRequired returns TRUE if this element is a required field
func (element Number) IsRequired() bool {
	return element.Required
}

/***********************************
 * PRIMARY INTERFACE METHODS
 ***********************************/

func (element Number) Get(object reflect.Value, path list.List) (reflect.Value, error) {

	// RULE: Cannot get sub-properties on a number
	if !path.IsEmpty() {
		return reflect.ValueOf(nil), derp.NewInternalError("schema.Number.Find", "Can't find sub-properties on a 'number' type", path)
	}

	// Try to convert and return the value
	if intValue, ok := convert.FloatOk(object, 0); ok {
		return element.floatSize(intValue), nil
	}

	// Try to use the default value
	if element.Default.IsPresent() {
		defaultValue := convert.FloatDefault(object, element.Default.Float())
		return element.floatSize(defaultValue), nil
	}

	// Return nil if no value is present
	return reflect.ValueOf(nil), nil
}

// GetElement returns a sub-element of this schema
func (element Number) GetElement(path list.List) (Element, error) {

	if path.IsEmpty() {
		return element, nil
	}

	return nil, derp.NewInternalError("schema.Number.GetElement", "Can't find sub-properties on an 'number' type", path)
}

// Set formats a value and applies it to the provided object/path
func (element Number) Set(object reflect.Value, path list.List, value any) (reflect.Value, error) {

	// RULE: Cannot set sub-properties on a number
	if !path.IsEmpty() {
		return reflect.ValueOf(nil), derp.NewInternalError("schema.Number.Set", "Can't set sub-properties on a number", path, value)
	}

	// Convert and return the new value
	floatValue, ok := convert.FloatOk(value, element.Default.Float())

	if !ok {
		return reflect.ValueOf(nil), derp.NewBadRequestError("schema.Number.Set", "Value must be convertable to a number", value)
	}

	return element.floatSize(floatValue), nil
}

// Remove removes a value from the provided object/path.  In the case of numbers, this is a no-op.
func (element Number) Remove(_ reflect.Value, _ list.List) (reflect.Value, error) {
	return reflect.ValueOf(nil), derp.NewInternalError("schema.Number.Remove", "Can't remove properties from a number.  This should never happen.")
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

func (element Number) Clean(value any) error {
	// TODO: HIGH: Implement this
	return nil
}

func (element Number) floatSize(value float64) reflect.Value {
	switch element.BitSize {
	case 32:
		return reflect.ValueOf(float32(value))
	case 64:
		return reflect.ValueOf(value)
	default:
		return reflect.ValueOf(value)
	}
}

/***********************************
 * ENUMERATOR INTERFACE
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
