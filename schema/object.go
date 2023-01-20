package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/maps"
)

// Object represents an object data type within a JSON-Schema.
type Object struct {
	Properties    ElementMap
	Wildcard      Element
	RequiredProps []string
	Required      bool
}

/***********************************
 * Container Interface
 ***********************************/

func (element Object) GetElement(name string) (Element, bool) {

	if name == "" {
		return element, true
	}

	head, tail := list.Split(name, list.DelimiterDot)

	if child, ok := element.Properties[head]; ok {
		return child.GetElement(tail)
	}

	if element.Wildcard != nil {
		return element.Wildcard.GetElement(tail)
	}

	return nil, false
}

/***********************************
 * Element Interface
 ***********************************/

// DefaultValue returns the default value for this element type.  In a special case for objects,
// which can be represented as both Go structs and maps, this returns a map[string]any that has been
// populated with any known default keys.
func (element Object) DefaultValue() any {
	result := maps.Map{}

	for key, element := range element.Properties {
		result[key] = element.DefaultValue()
	}

	return result
}

// IsRequired returns TRUE if this element is a required field
func (element Object) IsRequired() bool {
	return element.Required
}

// Validate validates a value against this schema
func (element Object) Validate(object any) error {

	for name, subElement := range element.Properties {
		if err := validate(subElement, object, name); err != nil {
			return err
		}
	}

	return nil
}

func (element Object) Clean(value any) error {
	// TODO: HIGH: Implement the Clean method for the Object type
	return nil
}

/***********************************
 * MARSHAL / UNMARSHAL METHODS
 ***********************************/

// MarshalMap populates object data into a map[string]any
func (element Object) MarshalMap() map[string]any {

	properties := make(maps.Map, len(element.Properties))

	for key, element := range element.Properties {
		properties[key] = element.MarshalMap()
	}

	result := maps.Map{
		"type":       TypeObject,
		"properties": properties,
		"required":   element.RequiredProps,
	}

	if element.Wildcard != nil {
		result["wildcard"] = element.Wildcard.MarshalMap()
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Object) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "object" {
		return derp.NewInternalError("schema.Object.UnmarshalMap", "Data is not type 'object'", data)
	}

	// Handle "simple" required as a boolean
	if required, ok := data["required"].(bool); ok {
		element.Required = required
	}

	// Handle property map
	if properties, ok := data["properties"].(map[string]any); ok {

		element.Properties = make(map[string]Element, len(properties))

		for key, value := range properties {

			if propertyMap, ok := value.(map[string]any); ok {

				if _, ok := propertyMap["required"]; !ok && element.Required {
					propertyMap["required"] = true
				}

				if propertyObject, err := UnmarshalMap(propertyMap); err == nil {

					element.Properties[key] = propertyObject
				}
			}
		}
	}

	// Handle other "wildcard" properties
	if wildcard, ok := data["wildcard"].(map[string]any); ok {
		if wildcardObject, err := UnmarshalMap(wildcard); err == nil {
			element.Wildcard = wildcardObject
		}
	}

	// Handle "standards" required as an array of strings.
	if required, ok := data["required"].([]any); ok {

		element.RequiredProps = convert.SliceOfString(required)

		for _, name := range element.RequiredProps {

			if property, ok := element.Properties[name]; ok {

				switch p := property.(type) {
				case Array:
					p.Required = true
				case Boolean:
					p.Required = true
				case Integer:
					p.Required = true
				case Number:
					p.Required = true
				case Object:
					p.Required = true
				case String:
					p.Required = true
				}
			}
		}
	}

	return err
}
