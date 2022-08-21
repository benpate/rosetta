package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/maps"
	"github.com/davecgh/go-spew/spew"
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
func (element Object) Get(object any, path string) (any, Element, error) {
	return element.GetReflect(convert.ReflectValue(object), path)
}

// Find locates a child of this element
func (element Object) GetReflect(object reflect.Value, path string) (any, Element, error) {

	if path == "" {
		return convert.Interface(object), element, nil
	}

	// Find the property in the schema
	head, tail := list.Dot(path).Split()
	property, ok := element.Properties[head]

	if !ok {
		return nil, element, derp.NewInternalError("schema.Object.Get", "Property does not exist in schema", path)
	}

	switch object.Kind() {
	case reflect.Pointer:
		return element.GetReflect(object.Elem(), path)

	case reflect.Interface:
		return element.GetReflect(object.Elem(), path)

	case reflect.Map:
		valueOf := object.MapIndex(reflect.ValueOf(head))
		return property.GetReflect(valueOf, tail.String())

	case reflect.Struct:
		valueOf, err := findFieldByTag(object, head)

		if err != nil {
			return nil, property, derp.NewInternalError("schema.Object.Get", "Property does not exist in object", object, path)
		}

		return property.GetReflect(valueOf, tail.String())
	}

	return property.GetReflect(reflect.ValueOf(nil), tail.String())
}

// Set validates/formats a value using this schema
func (element Object) Set(object any, path string, value any) error {

	// If we've been passed a NIL, then cast it as a map[string]any
	if object == nil {
		object = make(maps.Map)
	}

	// Shortcut if the object is a PathSetter.  Just call the SetPath function and we're good.
	if setter, ok := object.(PathSetter); ok {
		return setter.SetPath(path, value)
	}

	return element.SetReflect(convert.ReflectValue(object), path, value)
}

// Set validates/formats a value using this schema
func (element Object) SetReflect(object reflect.Value, path string, value any) error {

	spew.Dump("<<<<<<<<<<<<<< object.SetReflect >>>>>>>>>>>", path, object.Interface(), value)

	// Otherwise, use reflection to push the value into the object
	switch object.Kind() {
	case reflect.Pointer:
		spew.Dump("DEREFERENCE POINTER")
		return element.SetReflect(object.Elem(), path, value)

	case reflect.Interface:
		spew.Dump("DEREFERENCE INTERFACE")
		return element.SetReflect(object.Elem(), path, value)

	case reflect.Map:
		result := element.setMap(object, path, value)
		spew.Dump("MAP result...", object.Interface())
		return result

	case reflect.Struct:
		result := element.setStruct(object, path, value)
		spew.Dump("STRUCT result...", object.Interface())
		return result

	case reflect.Invalid:
		newMap := make(maps.Map)
		result := element.setMap(reflect.ValueOf(newMap), path, value)
		spew.Dump("INVALID result...", object.Interface())
		return result
	}

	return derp.NewInternalError("schema.Object.Set", "object must be a struct or a map", object.Kind().String(), object.Interface(), path, value)
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
