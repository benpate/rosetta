package schema

import (
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

// Array represents an array data type within a JSON-Schema.
type Array struct {
	Items     Element
	MinLength int
	MaxLength int
	Required  bool
}

/***********************************
 * Container Interface
 ***********************************/

func (element Array) GetProperty(name string) (Element, error) {

	index, err := strconv.Atoi(name)

	if err != nil {
		return nil, derp.Wrap(err, "schema.Array.GetProperty", "Invalid array index", name)
	}

	if index < 0 {
		return nil, derp.NewBadRequestError("schema.Array.GetProperty", "Array index must not be negative", name)
	}

	if index > element.MaxLength {
		return nil, derp.NewBadRequestError("schema.Array.GetProperty", "Array index must not be greater than the maximum", name, element.MaxLength)
	}

	return element.Items, nil
}

/***********************************
 * Element Interface
 ***********************************/

func (element Array) DefaultValue() any {
	// TODO: We can make a better default than this.
	return []any{}
}

// IsRequired returns TRUE if this element is a required field
func (element Array) IsRequired() bool {
	return element.Required
}

// Validate validates a value against this schema
func (element Array) Validate(object any) derp.MultiError {

	var err derp.MultiError

	lengthGetter, ok := object.(LengthGetter)

	if !ok {
		err.Append(derp.NewValidationError("Array must implement LengthGetter interface"))
		return err
	}

	// Check minimum/maximum lengths
	length := lengthGetter.Length()

	if element.Required && length == 0 {
		err.Append(derp.NewValidationError(" array value is required"))
		return err
	}

	if (element.MinLength > 0) && (length < element.MinLength) {
		err.Append(derp.NewValidationError(" minimum array length is " + convert.String(element.MinLength)))
		return err
	}

	if (element.MaxLength > 0) && (length > element.MaxLength) {
		err.Append(derp.NewValidationError(" maximum array length is " + convert.String(element.MaxLength)))
		return err
	}

	for index := 0; index < length; index = index + 1 {
		indexString := strconv.Itoa(index)
		err.Append(validate(element.Items, object, indexString))
	}

	return err
}

func (element Array) Clean(value any) derp.MultiError {
	// TODO: HIGH: Implement the Clean method for Array
	return nil
}

func (element Array) getElement(name string) (Element, bool) {

	if name == "" {
		return element, true
	}

	head, tail := list.Split(name, list.DelimiterDot)

	if _, ok := Index(head, element.MaxLength); ok {
		return element.Items.getElement(tail)
	}

	return nil, false
}

/***********************************
 * MARSHAL / UNMARSHAL METHODS
 ***********************************/

// MarshalMap populates object data into a map[string]any
func (element Array) MarshalMap() map[string]any {

	return map[string]any{
		"type":      TypeArray,
		"items":     element.Items.MarshalMap(),
		"minLength": element.MinLength,
		"maxLength": element.MaxLength,
	}
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Array) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "array" {
		return derp.NewInternalError("schema.Array.UnmarshalMap", "Data is not type 'array'", data)
	}

	element.Items, err = UnmarshalMap(data["items"])
	element.Required = convert.Bool(data["required"])
	element.MinLength = convert.Int(data["minLength"])
	element.MaxLength = convert.Int(data["maxLength"])

	return err
}
