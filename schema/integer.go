package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/null"
)

// Integer represents an integer data type within a JSON-Schema.
type Integer struct {
	Default    null.Int64 `json:"default"` // TODO: Int64??
	Minimum    null.Int64 `json:"minimum"`
	Maximum    null.Int64 `json:"maximum"`
	MultipleOf null.Int64 `json:"multipleOf"`
	Enum       []int      `json:"emum"`
	Required   bool
}

// Enumerate implements the "Enumerator" interface
func (element Integer) Enumerate() []string {
	return convert.SliceOfString(element.Enum)
}

// Type returns the data type of this Schema
func (element Integer) Type() reflect.Type {
	return reflect.TypeOf(0)
}

// IsRequired returns TRUE if this element is a required field
func (element Integer) IsRequired() bool {
	return element.Required
}

func (element Integer) Get(object reflect.Value, path string) (any, Element, error) {

	if path != "" {
		return nil, element, derp.NewInternalError("schema.Integer.Find", "Can't find sub-properties on an 'integer' type", path)
	}

	if element.Default.IsPresent() {
		return convert.Int64Default(object, element.Default.Int64()), element, nil
	}

	if intValue, ok := convert.Int64Ok(object, element.Default.Int64()); ok {
		return intValue, element, nil
	}

	return nil, element, nil
}

// Set formats a value and applies it to the provided object/path
func (element Integer) Set(object reflect.Value, path string, value any) error {

	// Cannot set sub-properties of an Integer
	if path != "" {
		return derp.NewInternalError("schema.Integer.Set", "Can't set sub-properties on an integer", path, value)
	}

	// Convert and set value
	intValue, ok := convert.Int64Ok(value, element.Default.Int64())

	if !ok {
		return derp.NewBadRequestError("schema.Integer.Set", "Value must be convertable to an integer", value)
	}

	return setWithReflection(object, intValue)
}

// Validate validates a value using this schema
func (element Integer) Validate(value any) error {

	var err error

	intValue, ok := convert.Int64Ok(value, element.Default.Int64())

	if !ok {
		return ValidationError{Message: "field must be an integer"}
	}

	if element.Required {
		if intValue == 0 {
			return ValidationError{Message: "field is required"}
		}
	}

	if element.Minimum.IsPresent() {
		if intValue < element.Minimum.Int64() {
			err = derp.Append(err, ValidationError{Message: "minimum value is " + convert.String(element.Minimum)})
		}
	}

	if element.Maximum.IsPresent() {
		if intValue > element.Maximum.Int64() {
			err = derp.Append(err, ValidationError{Message: "maximum value is " + convert.String(element.Maximum)})
		}
	}

	if element.MultipleOf.IsPresent() {
		if (intValue % element.MultipleOf.Int64()) != 0 {
			err = derp.Append(err, ValidationError{Message: "must be a multiple of " + convert.String(element.MultipleOf)})
		}
	}

	if len(element.Enum) > 0 {
		if !compare.Contains(element.Enum, intValue) {
			err = derp.Append(err, ValidationError{Message: "must contain one of the specified values"})
		}
	}

	return err
}

// DefaultType returns the default type for this element
func (element Integer) DefaultType() reflect.Type {
	return reflect.TypeOf(int64(0))
}

// DefaultValue returns the default value for this element type
func (element Integer) DefaultValue() any {
	return element.Default.Int64()
}

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
		return derp.New(500, "schema.Integer.UnmarshalMap", "Data is not type 'integer'", data)
	}

	element.Default = convert.NullInt64(data["default"])
	element.Minimum = convert.NullInt64(data["minimum"])
	element.Maximum = convert.NullInt64(data["maximum"])
	element.MultipleOf = convert.NullInt64(data["multipleOf"])
	element.Required = convert.Bool(data["required"])
	element.Enum = convert.SliceOfInt(data["enum"])

	return err
}
