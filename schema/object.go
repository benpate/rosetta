package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/maps"
)

// Object represents an object data type within a JSON-Schema.
type Object struct {
	Properties    ElementMap
	Unlisted      Element
	RequiredProps []string
	Required      bool
}

/***********************************
 * ELEMENT META-DATA
 ***********************************/

// Type returns the data type of this Element
func (element Object) Type() reflect.Type {
	return reflect.TypeOf(maps.Map{})
}

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

/***********************************
 * PRIMARY INTERFACE METHODS
 ***********************************/

// Find locates a child of this element
func (element Object) Get(object reflect.Value, path list.List) (reflect.Value, error) {

	switch object.Kind() {

	// If value is nil, then substitute the default value
	case reflect.Invalid:
		return element.Get(reflect.ValueOf(element.DefaultValue()), path)

	// Dereference pointers
	case reflect.Pointer:
		return element.Get(object.Elem(), path)

	// Dereference interfaces
	case reflect.Interface:
		return element.Get(object.Elem(), path)

	// look up values in Maps
	case reflect.Map:
		return element.getFromMap(object, path)

	// look up values in Structs
	case reflect.Struct:
		return element.getFromStruct(object, path)
	}

	return reflect.ValueOf(nil), derp.NewInternalError("schema.Object.Get", "object must be a struct or a map", object.Kind().String(), object.Interface(), path)
}

// GetElement returns a sub-element of this schema
func (element Object) GetElement(path list.List) (Element, error) {

	if path.IsEmpty() {
		return element, nil
	}

	if property, ok := element.Properties[path.Head()]; ok {
		return property.GetElement(path.Tail())
	}

	return nil, derp.NewInternalError("schema.Integer.GetElement", "Property does not exist in this object", path)
}

// Set validates/formats a value using this schema
func (element Object) Set(object reflect.Value, path list.List, value any) (reflect.Value, error) {

	switch object.Kind() {

	// Dereference pointers
	case reflect.Pointer:
		return element.Set(object.Elem(), path, value)

	// Dereference interfaces
	case reflect.Interface:
		return element.Set(object.Elem(), path, value)

	// Set Maps
	case reflect.Map:
		return element.setToMap(object, path, value)

	// Set Structs
	case reflect.Struct:
		return element.setToStruct(object, path, value)

	// If the value is nil, then create a default value and set the new property in it.
	case reflect.Invalid:
		return element.Set(reflect.ValueOf(element.DefaultValue()), path, value)

	default:
		return reflect.ValueOf(nil), derp.Report(derp.NewInternalError("schema.Object.Set", "object must be a struct or a map", object.Kind().String(), object.Interface(), path, value))
	}

}

// Remove removes a child of the target object.  In the case of Maps, the map key is removed.
// In the case of Structs, the field is set to its default value.
func (element Object) Remove(object reflect.Value, path list.List) (reflect.Value, error) {

	switch object.Kind() {

	// If the value is nil, then create a default value and set the new property in it.
	case reflect.Invalid:
		return element.Remove(reflect.ValueOf(element.DefaultValue()), path)

	// Dereference pointers
	case reflect.Pointer:
		return element.Remove(object.Elem(), path)

	// Dereference interfaces
	case reflect.Interface:
		return element.Remove(object.Elem(), path)

	// Set Maps
	case reflect.Map:
		return element.removeFromMap(object, path)

	// Set Structs
	case reflect.Struct:
		return element.removeFromStruct(object, path)

	default:
		return reflect.ValueOf(nil), derp.Report(derp.NewInternalError("schema.Object.Set", "object must be a struct or a map", object.Kind().String(), object.Interface(), path))
	}
}

// Validate validates a value against this schema
func (element Object) Validate(value any) error {

	const location = "schema.Object.Validate"
	var errorReport error

	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {

	case reflect.Pointer:
		return element.Validate(convert.Interface(valueOf.Elem()))

	case reflect.Interface:
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
				errorReport = derp.Append(errorReport, derp.Wrap(err, location, "field not found", name))
			} else {
				errorReport = derp.Append(errorReport, addPath(name, child.Validate(field.Interface())))
			}
		}

	default:
		return Invalid("Element must be a struct or a map")
	}

	return errorReport
}

func (element Object) Clean(value any) error {
	// TODO: HIGH: Implement this
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

	return maps.Map{
		"type":       TypeObject,
		"properties": properties,
		"required":   element.RequiredProps,
	}
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
