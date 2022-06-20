package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

// Object represents an object data type within a JSON-Schema.
type Object struct {
	Properties    ElementMap
	RequiredProps []string
	Required      bool
}

// Type returns the data type of this Element
func (element Object) Type() reflect.Type {
	return reflect.TypeOf(map[string]any{})
}

// IsRequired returns TRUE if this element is a required field
func (element Object) IsRequired() bool {
	return element.Required
}

// Find locates a child of this element
func (element Object) Get(object reflect.Value, path string) (any, Element, error) {

	if path == "" {
		return convert.Interface(object), element, nil
	}

	// Find the property in the schema
	head, tail := list.Split(path, ".")
	property, ok := element.Properties[head]

	if !ok {
		return nil, element, derp.NewInternalError("schema.Object.Get", "Property does not exist in schema", path)
	}

	switch object.Kind() {
	case reflect.Pointer:
		return element.Get(object.Elem(), path)

	case reflect.Map:
		valueOf := object.MapIndex(reflect.ValueOf(head))
		return property.Get(valueOf, tail)

	case reflect.Struct:
		valueOf, err := findFieldByTag(object, head)

		if err != nil {
			return nil, property, derp.NewInternalError("schema.Object.Get", "Property does not exist in object", object, path)
		}

		return property.Get(valueOf, tail)
	}

	return property.Get(reflect.ValueOf(nil), tail)
}

// Set validates/formats a value using this schema
func (element Object) Set(object reflect.Value, path string, value any) error {

	// Otherwise, use reflection to push the value into the object
	switch object.Kind() {
	case reflect.Pointer:
		return element.Set(object.Elem(), path, value)

	case reflect.Map:
		return element.setMap(object, path, value)

	case reflect.Struct:
		return element.setStruct(object, path, value)

	case reflect.Invalid:
		newMap := map[string]any{}
		return element.setMap(reflect.ValueOf(newMap), path, value)
	}

	return derp.NewInternalError("schema.Object.Set", "object must be a struct or a map", object.Kind(), path, value)
}

// Validate validates a value against this schema
func (element Object) Validate(value any) error {

	const location = "schema.Object.Validate"
	var errorReport error

	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {

	case reflect.Pointer:
		return element.Validate(convert.Interface(valueOf.Elem()))

	case reflect.Map:

		for key, child := range element.Properties {
			keyValue := reflect.ValueOf(key)
			mapValue := valueOf.MapIndex(keyValue)

			errorReport = derp.Append(errorReport, addPath(key, child.Validate(convert.Interface(mapValue))))
		}

	case reflect.Struct:

		for name, child := range element.Properties {
			if field, err := findFieldByTag(valueOf, name); err != nil {
				return derp.Wrap(err, location, "field not found", name)
			} else {
				errorReport = derp.Append(errorReport, addPath(name, child.Validate(field.Interface())))
			}
		}

	default:
		return Invalid("Element must be a struct or a map")
	}

	return errorReport
}

// MarshalMap populates object data into a map[string]any
func (element Object) MarshalMap() map[string]any {

	properties := make(map[string]any, len(element.Properties))

	for key, element := range element.Properties {
		properties[key] = element.MarshalMap()
	}

	return map[string]any{
		"type":       TypeObject,
		"properties": properties,
		"required":   element.RequiredProps,
	}
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Object) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "object" {
		return derp.New(500, "schema.Object.UnmarshalMap", "Data is not type 'object'", data)
	}

	// Handle "simple" required as a boolean
	if required, ok := data["required"].(bool); ok {
		element.Required = required
	}

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

	// Handle "standards" required as an array of strings.
	if required, ok := data["required"].([]any); ok {

		element.RequiredProps = convert.SliceOfString(required)

		for _, name := range element.RequiredProps {

			if property, ok := element.Properties[name]; ok {

				switch p := property.(type) {
				case *Any:
					p.Required = true
				case *Array:
					p.Required = true
				case *Boolean:
					p.Required = true
				case *Integer:
					p.Required = true
				case *Number:
					p.Required = true
				case *Object:
					p.Required = true
				case *String:
					p.Required = true
				}
			}
		}
	}

	return err
}
