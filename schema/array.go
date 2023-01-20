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
func (element Array) Validate(object any) error {

	lengthGetter, ok := object.(LengthGetter)

	if !ok {
		return derp.NewValidationError("Array must implement LengthGetter interface")
	}

	// Check minimum/maximum lengths
	length := lengthGetter.Length()

	if element.Required && length == 0 {
		return derp.NewValidationError(" array value is required")
	}

	if (element.MinLength > 0) && (length < element.MinLength) {
		return derp.NewValidationError(" minimum array length is " + convert.String(element.MinLength))
	}

	if (element.MaxLength > 0) && (length > element.MaxLength) {
		return derp.NewValidationError(" maximum array length is " + convert.String(element.MaxLength))
	}

	for index := 0; index < length; index = index + 1 {
		indexString := strconv.Itoa(index)

		if err := validate(element.Items, object, indexString); err != nil {
			return derp.Wrap(err, "schema.Array.Validate", "Error Validating object at index", index)
		}
	}

	return nil
}

func (element Array) Clean(value any) error {
	// TODO: HIGH: Implement the Clean method for Array
	return nil
}

func (element Array) GetElement(name string) (Element, bool) {

	if name == "" {
		return element, true
	}

	head, tail := list.Split(name, list.DelimiterDot)

	var ok bool

	if element.MaxLength > 0 {
		_, ok = Index(head, element.MaxLength)

	} else {
		_, ok = Index(head)
	}

	if ok {
		return element.Items.GetElement(tail)
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
