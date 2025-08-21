package schema

import (
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/schema/format"
)

// String represents a string data type within a JSON-Schema.
type String struct {
	Default    string   `json:"default"`
	MinLength  int      `json:"minLength"`
	MaxLength  int      `json:"maxLength"`
	Enum       []string `json:"enum"`
	MinValue   string   `json:"minValue"`
	MaxValue   string   `json:"maxValue"`
	Pattern    string   `json:"pattern"`
	Format     string   `json:"format"`
	Required   bool     `json:"required"`
	RequiredIf string   `json:"required-if"`
}

/***********************************
 * Element Interface
 ***********************************/

// DefaultValue returns the default value for this element type
func (element String) DefaultValue() any {
	return element.Default
}

// IsRequired returns TRUE if this element is a required field
func (element String) IsRequired() bool {
	return element.Required
}

// Validate compares a generic data value using this Schema
func (element String) Validate(value any) error {

	stringValue, ok := value.(string)

	if !ok {
		return derp.ValidationError(" must be a string")
	}

	// Verify required fields (after format functions are applied)
	if element.Required {
		if stringValue == "" {
			return derp.ValidationError(" string field is required")
		}
	}

	// Validate minimum value
	if element.MinValue != "" {
		if stringValue < element.MinValue {
			return derp.ValidationError(" minimum string value is " + element.MinValue)
		}
	}

	// Validate maximum value
	if element.MaxValue != "" {
		if stringValue > element.MaxValue {
			return derp.ValidationError(" maximum string value is " + element.MaxValue)
		}
	}

	// Validate minimum length
	if element.MinLength > 0 {
		if len(stringValue) < element.MinLength {
			return derp.ValidationError(" minimum string length is " + convert.String(element.MinLength))
		}
	}

	// Validate maximum length
	if element.MaxLength > 0 {
		if len(stringValue) > element.MaxLength {
			return derp.ValidationError(" maximum string length is " + convert.String(element.MaxLength))
		}
	}

	// Validate enumerated values
	if len(element.Enum) > 0 {
		if (stringValue != "") && (!compare.Contains(element.Enum, stringValue)) {
			return derp.ValidationError(" string must match one of the required values", stringValue, element.Enum)
		}
	}

	// Validate against all formatting functions
	for _, format := range element.formatFunctions() {
		if _, err := format(stringValue); err != nil {
			return err
		}
	}

	return nil
}

// ValidateRequiredIf returns an error if the conditional expression is true but the value is empty
func (element String) ValidateRequiredIf(schema Schema, path list.List, globalValue any) error {

	const location = "schema.String.ValidateRequiredIf"

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

func (element String) GetElement(name string) (Element, bool) {

	if name == "" {
		return element, true
	}
	return nil, false
}

func (element String) Inherit(parent Element) {
	// Do nothing
}

// AllProperties returns a map of all properties for this element
func (element String) AllProperties() ElementMap {
	return ElementMap{
		"": element,
	}
}

/***********************************
 * Enumerator Interface
 ***********************************/

// Enumerate implements the "Enumerator" interface
func (element String) Enumerate() []string {
	return element.Enum
}

/***********************************
 * Marshal / Unmarshal Methods
 ***********************************/

// MarshalMap populates object data into a map[string]any
func (element String) MarshalMap() map[string]any {

	result := map[string]any{
		"type":     TypeString,
		"required": element.Required,
	}

	if element.Default != "" {
		result["default"] = element.Default
	}

	if element.MinLength > 0 {
		result["minLength"] = element.MinLength
	}

	if element.MaxLength > 0 {
		result["maxLength"] = element.MaxLength
	}

	if element.Pattern != "" {
		result["pattern"] = element.Pattern
	}

	if element.Format != "" {
		result["format"] = element.Format
	}

	if len(element.Enum) > 0 {
		result["enum"] = element.Enum
	}

	if element.RequiredIf != "" {
		result["required-if"] = element.RequiredIf
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *String) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "string" {
		return derp.InternalError("schema.String.UnmarshalMap", "Data is not type 'string'", data)
	}

	element.Default = convert.String(data["default"])
	element.MinLength = convert.Int(data["minLength"])
	element.MaxLength = convert.Int(data["maxLength"])
	element.Pattern = convert.String(data["pattern"])
	element.Format = convert.String(data["format"])
	element.Enum = convert.SliceOfString(data["enum"])
	element.Required = convert.Bool(data["required"])
	element.RequiredIf = convert.String(data["required-if"])

	return err
}

/***********************************
 * Helper Methods
 ***********************************/

// formatFunctions parses all of the formatting functions
// used by this string value.
func (element String) formatFunctions() []format.StringFormat {

	result := make([]format.StringFormat, 0)

	// Split multiple formats into individual function calls
	params := strings.Split(element.Format, " ")

	for _, arg := range params {

		name, value := list.Equal(arg).Split()

		if formatFunction, ok := formats[name]; ok {
			result = append(result, formatFunction(value.String()))
		}
	}

	// If there are no valid formats defined, then default to
	// no-html, which strictly removes all HTML tags from the value.
	if len(result) == 0 {
		result = []format.StringFormat{formats["no-html"]("")}
	}

	return result
}
