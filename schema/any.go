package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

// Any represents a schema that accepts any data type
type Any struct {
	Required   bool
	RequiredIf string
}

// DefaultValue returns the default value for this element
func (element Any) DefaultValue() any {
	return nil
}

// IsRequired returns true if this a value is required for this element
func (element Any) IsRequired() bool {
	return element.Required
}

// Validate validates the provided value
func (element Any) Validate(_ any) error {
	return nil
}

// ValidateRequiredIf returns an error if the conditional expression is true but the value is empty
func (element Any) ValidateRequiredIf(schema Schema, path list.List, value any) error {

	const location = "schema.Any.ValidateRequiredIf"

	// If RequiredIf is not set, then exit
	if element.RequiredIf == "" {
		return nil
	}

	// Evaluate the RequiredIf expression
	isRequired, err := schema.Match(value, exp.Parse(element.RequiredIf))

	if err != nil {
		return derp.Wrap(err, location, "Error evaluating condition", element.RequiredIf)
	}

	// If the expression did not evaluat to TRUE, then exit
	if !isRequired {
		return nil
	}

	// Find the propertyValue in the global object
	propertyValue, err := schema.Get(value, path.String())

	if err != nil {
		return derp.Wrap(err, location, "Error getting value for path", path)
	}

	// The value is required, but missing, so.. error.
	if compare.IsZero(propertyValue) {
		return derp.Validation("field: " + path.String() + " is required based on condition: " + element.RequiredIf)
	}

	// The value is required, but present, so.. success.
	return nil
}

// GetElement implements the Element interface
// It returns the element at the specified path
func (element Any) GetElement(_ string) (Element, bool) {
	return element, true
}

// Inherit implements the Element interface
// It is a no-op for Any elements
func (element Any) Inherit(_ Element) {
	// Do nothing
}

// AllProperties implements the Element interface
// It returns a map of all properties for this element
func (element Any) AllProperties() ElementMap {
	return ElementMap{
		"": element,
	}
}

/***********************************
 * Marshal / Unmarshal Methods
 ***********************************/

// MarshalMap populates the object data into a map[string]any
func (element Any) MarshalMap() map[string]any {

	return map[string]any{
		"type":        TypeAny,
		"required":    element.Required,
		"required-if": element.RequiredIf,
	}
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Any) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "any" {
		return derp.Internal("schema.String.UnmarshalMap", "Data is not type 'string'", data)
	}

	element.Required = convert.Bool(data["required"])

	return err
}
