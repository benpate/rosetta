package schema

import (
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

// Array represents an array data type within a JSON-Schema.
type Array struct {
	Items      Element `json:"items"`
	MinLength  int     `json:"minLength"`
	MaxLength  int     `json:"maxLength"`
	Required   bool    `json:"required"`
	RequiredIf string  `json:"required-if"`
}

/******************************************
 * Container Interface
 ******************************************/

// GetProperty returns the property with the specified name
func (element Array) GetProperty(name string) (Element, error) {

	const location = "schema.Array.GetProperty"

	index, err := strconv.Atoi(name)

	if err != nil {
		return nil, derp.Wrap(err, location, "Invalid array index", name)
	}

	if index < 0 {
		return nil, derp.BadRequestError(location, "Array index must not be negative", name)
	}

	if index > element.MaxLength {
		return nil, derp.BadRequestError(location, "Array index must not be greater than the maximum", name, element.MaxLength)
	}

	return element.Items, nil
}

/******************************************
 * Element Interface
 ******************************************/

// DefaultValue implements the Element interface
// It returns the default value for this element type
func (element Array) DefaultValue() any {
	// TODO: We can make a better default than this.
	return []any{}
}

// IsRequired implements the Element interface
// It returns TRUE if this element is a required field
func (element Array) IsRequired() bool {
	return element.Required
}

// Validate implements the Element interface
// It validates a value against this schema
func (element Array) Validate(object any) error {

	length, isLengthGetter := getLength(object)

	if !isLengthGetter {
		return derp.InternalError("schema.Array.Validate", "Array must implement LengthGetter interface")
	}

	// Check minimum/maximum lengths
	if element.Required && length == 0 {
		return derp.ValidationError(" array value is required")
	}

	if (element.MinLength > 0) && (length < element.MinLength) {
		return derp.ValidationError(" minimum array length is " + convert.String(element.MinLength))
	}

	if (element.MaxLength > 0) && (length > element.MaxLength) {
		return derp.ValidationError(" maximum array length is " + convert.String(element.MaxLength))
	}

	for index := 0; index < length; index = index + 1 {
		indexString := strconv.Itoa(index)

		if err := validate(element.Items, object, indexString); err != nil {
			return derp.Wrap(err, "schema.Array.Validate", "Error Validating object at index", index)
		}
	}

	return nil
}

// ValidateRequiredIf implements the Element interface
// It returns an error if the conditional expression is true but the value is empty
func (element Array) ValidateRequiredIf(schema Schema, path list.List, globalValue any) error {

	const location = "schema.Array.ValidateRequiredIf"

	if element.RequiredIf != "" {

		localValue, err := schema.get(globalValue, element, path)

		if err != nil {
			return derp.Wrap(err, location, "Error getting value for path", path)
		}

		length, ok := getLength(localValue)

		if !ok {
			return derp.ValidationError("Array must implement LengthGetter interface")
		}

		if length == 0 {
			isRequired, err := schema.Match(globalValue, exp.Parse(element.RequiredIf))

			if err != nil {
				return derp.Wrap(err, location, "Error evaluating condition", element.RequiredIf)
			}

			if isRequired {
				return derp.ValidationError("field: " + path.String() + " is required based on condition: " + element.RequiredIf)
			}
		}

		for index := range length {
			subPath := path.PushTail(strconv.Itoa(index))

			if element.Items == nil {
				return derp.InternalError(location, "Array items cannot be nil", path)
			}

			if err := element.Items.ValidateRequiredIf(schema, subPath, globalValue); err != nil {
				return derp.Wrap(err, "schema.Array.ValidateRequiredIf", "Error Validating object at index", index)
			}
		}
	}

	return nil
}

// GetElement implements the Element interface
// It returns the element at the specified path
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

// Inherit implements the Element interface
// It is a no-op for Array elements
func (element Array) Inherit(_ Element) {
	// Do nothing
}

// AllProperties implements the Element interface
// It returns a map of all properties for this element
func (element Array) AllProperties() ElementMap {
	return ElementMap{
		"": element,
	}
}

/******************************************
 * Array-Specific Methods
 ******************************************/

// GetLength returns the length of the array value (if the object implements ArrayGetter)
func (element Array) GetLength(value any) (int, bool) {
	return getLength(value)
}

// GetIndex returns the value at a specific index in the array (if the object implements ArrayGetter)
func (element Array) GetIndex(value any, index int) (any, bool) {
	return getIndex(value, index)
}

// SetIndex sets the value at a specific index in the array (if the object implements ArraySetter)
func (element Array) SetIndex(value any, index int, item any) bool {

	if setter, ok := value.(ArraySetter); ok {
		return setter.SetIndex(index, item)
	}

	return false
}

// Append adds a new item to the end of the array (if the object implements ArraySetter)
func (element Array) Append(value ArraySetter, item any) error {

	const location = "schema.Array.Append"

	// Try to set the value at the end of the array
	if success := value.SetIndex(value.Length(), item); !success {
		return derp.InternalError(location, "Unable to set value at end of array", value)
	}

	// Success
	return nil
}

/******************************************
 * Marshal / Unmarshal Methods
 ******************************************/

// MarshalMap populates object data into a map[string]any
func (element Array) MarshalMap() map[string]any {

	return map[string]any{
		"type":        TypeArray,
		"items":       element.Items.MarshalMap(),
		"minLength":   element.MinLength,
		"maxLength":   element.MaxLength,
		"required":    element.Required,
		"required-if": element.RequiredIf,
	}
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Array) UnmarshalMap(data map[string]any) error {

	var err error

	// RULE: `type` must be "array"
	if convert.String(data["type"]) != "array" {
		return derp.InternalError("schema.Array.UnmarshalMap", "Data is not type 'array'", data)
	}

	// Try to retrieve the array items from the data map
	items, err := UnmarshalMap(data["items"])

	if err != nil {
		return derp.Wrap(err, "schema.Array.UnmarshalMap", "Unable to unmarshal 'items'", data["items"])
	}

	if items == nil {
		return derp.InternalError("schema.Array.UnmarshalMap", "'items' cannot be nil", data)
	}

	// Populate the element
	element.Items = items
	element.MinLength = convert.Int(data["minLength"])
	element.MaxLength = convert.Int(data["maxLength"])
	element.Required = convert.Bool(data["required"])
	element.RequiredIf = convert.String(data["required-if"])

	return err
}
