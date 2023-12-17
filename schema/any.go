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
	if element.RequiredIf != "" {
		if schema.Match(globalValue, exp.Parse(element.RequiredIf)) {
			if localValue, err := schema.Get(globalValue, path.String()); err != nil {
				return derp.Wrap(err, "schema.Any.ValidateRequiredIf", "Error getting value for path", path)
			} else if compare.IsZero(localValue) {
				return derp.NewValidationError("field: " + path.String() + " is required based on condition: " + element.RequiredIf)
			}
		}
	}
	return nil
}

// Clean updates a value to match the schema.  The value must be a pointer.
func (element Any) Clean(value any) error {
	return nil
}

func (element Any) GetElement(name string) (Element, bool) {
	return element, true
}

func (element Any) Inherit(parent Element) {
	// Do nothing
}

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
		return derp.NewInternalError("schema.String.UnmarshalMap", "Data is not type 'string'", data)
	}

	element.Required = convert.Bool(data["required"])

	return err
}
