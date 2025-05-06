package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

// Object represents an object data type within a JSON-Schema.
type Object struct {
	Properties ElementMap `json:"properties"`
	Wildcard   Element    `json:"wildcard"`
	Required   bool       `json:"required"`
	RequiredIF string     `json:"required-if"`
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
	result := make(map[string]any)

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
			return derp.Wrap(err, "schema.Object.Validate", "Error validating property", name)
		}
	}

	return nil
}

// ValidateRequiredIf returns an error if the conditional expression is true but the value is empty
func (element Object) ValidateRequiredIf(schema Schema, path list.List, globalValue any) error {

	for name, subElement := range element.Properties {
		subPath := path.PushTail(name)
		if err := subElement.ValidateRequiredIf(schema, subPath, globalValue); err != nil {
			return derp.Wrap(err, "schema.Object.ValidateRequiredIf", "Error validating property", subPath.String())
		}
	}

	return nil
}

func (element Object) Inherit(parent Element) {

	if element.Properties == nil {
		element.Properties = make(ElementMap)
	}

	// Inherit each property from the parent
	if parentObject, ok := parent.(Object); ok {
		for propertyName, parentProperty := range parentObject.Properties {
			if property, ok := element.Properties[propertyName]; ok {
				property.Inherit(parentProperty)
				element.Properties[propertyName] = property
			} else {
				element.Properties[propertyName] = parentProperty
			}
		}
	}
}

// AllProperties returns a flat slice of all properties in this element
// (in this case, it returns all properties of this object)
func (element Object) AllProperties() ElementMap {

	result := make(ElementMap, len(element.Properties))

	for name, property := range element.Properties {

		if object, ok := property.(Object); ok {

			for subName, subProperty := range object.AllProperties() {
				result[name+"."+subName] = subProperty
			}

		} else {
			result[name] = property
		}
	}

	return result
}

/***********************************
 * MARSHAL / UNMARSHAL METHODS
 ***********************************/

// MarshalMap populates object data into a map[string]any
func (element Object) MarshalMap() map[string]any {

	properties := make(map[string]any, len(element.Properties))

	for key, element := range element.Properties {
		properties[key] = element.MarshalMap()
	}

	result := map[string]any{
		"type":       TypeObject,
		"properties": properties,
	}

	if element.Wildcard != nil {
		result["wildcard"] = element.Wildcard.MarshalMap()
	}

	if element.Required {
		result["required"] = true
	}

	if element.RequiredIF != "" {
		result["required-if"] = element.RequiredIF
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Object) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "object" {
		return derp.InternalError("schema.Object.UnmarshalMap", "Data is not type 'object'", data)
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

	return err
}
