package schema

import (
	"reflect"

	"github.com/benpate/derp"
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
	Required   bool
}

/***********************************
 * ELEMENT META-DATA
 ***********************************/

// Type returns the data type of this Schema
func (element Integer) Type() reflect.Type {
	return element.intSize(0).Type()
}

// DefaultValue returns the default value for this element type
func (element Integer) DefaultValue() any {

	if element.Default.IsPresent() {
		return element.intSize(element.Default.Int64()).Interface()
	}

	return element.Default.Interface()
}

// IsRequired returns TRUE if this element is a required field
func (element Integer) IsRequired() bool {
	return element.Required
}

/***********************************
 * PRIMARY INTERFACE METHODS
 ***********************************/

func (element Integer) Get(object reflect.Value, path list.List) (reflect.Value, error) {

	if !path.IsEmpty() {
		return reflect.ValueOf(nil), derp.NewInternalError("schema.Integer.Find", "Can't find sub-properties on an 'integer' type", path)
	}

	if intValue, ok := convert.Int64Ok(object, 0); ok {
		return element.intSize(intValue), nil
	}

	if element.Default.IsPresent() {
		defaultValue := convert.Int64Default(object, element.Default.Int64())
		return element.intSize(defaultValue), nil
	}

	return reflect.ValueOf(nil), nil
}

// GetElement returns a sub-element of this schema
func (element Integer) GetElement(path list.List) (Element, error) {

	if path.IsEmpty() {
		return element, nil
	}

	return nil, derp.NewInternalError("schema.Integer.GetElement", "Can't find sub-properties on an 'integer' type", path)
}

// Set formats a value and applies it to the provided object/path
func (element Integer) Set(object reflect.Value, path list.List, value any) (reflect.Value, error) {

	// RULE: Cannot set sub-properties of an Integer
	if !path.IsEmpty() {
		return reflect.ValueOf(nil), derp.NewInternalError("schema.Integer.Set", "Can't set sub-properties on an integer", path, value)
	}

	// Convert and return the new value
	intValue, ok := convert.Int64Ok(value, element.Default.Int64())

	if !ok {
		return reflect.ValueOf(nil), derp.NewBadRequestError("schema.Integer.Set", "Value must be convertable to an integer", value)
	}

	return element.intSize(intValue), nil
}

// Remove removes a value from the provided object/path.  In the case of integers, this is a no-op.
func (element Integer) Remove(_ reflect.Value, _ list.List) (reflect.Value, error) {
	return reflect.ValueOf(nil), derp.NewInternalError("schema.Integer.Remove", "Can't remove properties from a integer.  This should never happen.")
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

func (element Integer) Clean(value any) error {
	// TODO: HIGH: Implement the "Clean" method for Integer
	return nil
}

func (element Integer) intSize(value int64) reflect.Value {

	switch element.BitSize {
	case 8:
		return reflect.ValueOf(int8(value))
	case 16:
		return reflect.ValueOf(int16(value))
	case 32:
		return reflect.ValueOf(int32(value))
	case 64:
		return reflect.ValueOf(int64(value))
	default:
		return reflect.ValueOf(int(value))
	}
}

/***********************************
 * ENUMERATOR INTERFACE
 ***********************************/

// Enumerate implements the "Enumerator" interface
func (element Integer) Enumerate() []string {
	return convert.SliceOfString(element.Enum)
}

/***********************************
 * MARSHAL / UNMARSHAL METHODS
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
