package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

type Any struct {
	Required   bool
	RequiredIf string
}

// Default returns the default value for this element
func (element Any) DefaultValue() any {
	return nil
}

// IsRequired returns true if this a value is required for this element
func (element Any) IsRequired() bool {
	return element.Required
}

// Validate validates the provided value
func (element Any) Validate(value any) error {
	return nil
}

// ValidateRequiredIf returns an error if the conditional expression is true but the value is empty
func (element Any) ValidateRequiredIf(schema Schema, path list.List, globalValue any) error {

	const location = "schema.Any.ValidateRequiredIf"

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

func (element Any) GetElement(name string) (Element, bool) {
	return element, true
}

func (element Any) Inherit(parent Element) {
	// Do nothing
}

// AllProperties returns a map of all properties for this element
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
		return derp.InternalError("schema.String.UnmarshalMap", "Data is not type 'string'", data)
	}

	element.Required = convert.Bool(data["required"])

	return err
}
