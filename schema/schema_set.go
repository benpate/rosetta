package schema

import (
	"net/url"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

// Set sets the value at the specified path within the object according to this schema
func (schema Schema) Set(object any, path string, value any) error {

	const location = "schema.Schema.Set"

	// Find the element for this path
	element, ok := schema.GetElement(path)

	if !ok {
		return derp.BadRequest(location, "Invalid path", path)
	}

	// Validate the value (and update if necessary) against the schema rules for this element
	value, _, err := validate(element, value)

	if err != nil {
		return derp.Wrap(err, location, "Value is not valid for this schema", path, value)
	}

	// set the property value in the object
	return setProperty(schema.Element, object, path, value)
}

// SetAll iterates over Set to apply all of the values to the object one at a time, stopping
// at the first error it encounters.  If all values are addedd successfully, then SetAll
// also uses Validate() to confirm that the object is still correct.
func (schema Schema) SetAll(object any, values map[string]any) error {

	const location = "schema.Schema.SetAll"

	// Set each value in the schema
	for path, value := range values {

		// Try to set the value in the object.
		if err := schema.Set(object, path, value); err != nil {
			return derp.Wrap(err, location, "Setting value", path, value)
		}
	}

	// Validate the whole schema once all the values are set
	// TODO: Validate Required-If here...
	_, _, err := Validate(schema, object)
	if err != nil {
		return derp.Wrap(err, location, "Validation Error")
	}

	// Success!!
	return nil
}

// SetURLValues iterates over Set to apply all of the values to the object one at a time, stopping
// at the first error it encounters.  If all values are addedd successfully, then SetURLValues
// also uses Validate() to confirm that the object is still correct.
func (schema Schema) SetURLValues(object any, values url.Values) error {

	const location = "schema.Schema.SetURLValues"

	// Set each value in the schema
	for path, value := range values {

		// Errors are intentionally ignored here.
		// Unallowed data does not make it through the schema filter
		// nolint: errcheck
		_ = schema.Set(object, path, value)
	}

	// Validate the whole schema once all the values are set
	// TODO: Changes are not yet applied here
	_, _, err := Validate(schema, object)
	if err != nil {
		return derp.Wrap(err, location, "Validation Error")
	}

	// Success!!
	return nil
}

// setProperty sets the value at the specified path within the object according to the provided schema
func setProperty(element Element, object any, path string, value any) (err error) {

	const location = "schema.setProperty"

	defer func() {
		if recovered := recover(); recovered != nil {
			err = derp.Internal(location, "Panic while setting value", path, value, recovered)
		}
	}()

	// In rare cases, we may need to set an entire element in one call.
	// If the path is empty, then we're setting the entire element, which
	// must implement the "ValueSetter" interface
	if path == "" {
		if setter, ok := object.(ValueSetter); ok {
			if err := setter.SetValue(value); err != nil {
				return derp.Wrap(err, location, "Unable to set value", object, value)
			}
			return nil
		}
		return derp.Internal(location, "Cannot set values on empty path", object, element, path, value)
	}

	// Split the path into head and tail
	head, tail, _ := strings.Cut(path, ".")

	// FInd the property definition for the first path segment
	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.Internal(location, "Property does not exist in schema", head, path)
	}

	// Use the element type to find the correct property setter
	switch typed := subElement.(type) {

	case Any, Array, Object:
		return setProperty_Object(typed, object, path, head, tail, value)

	case Boolean:
		return setProperty_Boolean(object, path, value)

	case Integer:
		if typed.BitSize == 64 {
			return setProperty_Integer64(object, path, value)
		}

		return setProperty_Integer32(object, path, value)

	case Number:
		return setProperty_Number(object, path, value)

	case String:
		return setProperty_String(object, path, value)
	}

	return derp.Internal(location, "Unsupported element type", path, subElement, object)
}

// setProperty_Object sets a value in the object using either the ObjectSetter or PointerGetter interface.
func setProperty_Object(element Element, object any, path string, head string, tail string, value any) error {

	const location = "schema.setProperty_Object"

	// ObjectSetter interface is required for Maps
	if setter, ok := object.(ObjectSetter); ok {
		return setter.SetObject(element, list.ByDot(path), value)
	}

	// PointerGetter works for Structs, Slices, and Arrays
	if getter, ok := object.(PointerGetter); ok {
		if subPointer, ok := getter.GetPointer(head); ok {
			return setProperty(element, subPointer, tail, value)
		}
	}

	// Cannot set the value
	return derp.Internal(location, "Target Object must be an ObjectSetter or PointerGetter", path, object)
}

// setProperty_Boolean sets a boolean value in the object.
func setProperty_Boolean(object any, path string, value any) error {

	const location = "schema.setProperty_Boolean"

	// Convert the value to a bool
	boolValue := convert.Bool(value)

	// Try to set the value using the BoolSetter interface
	if setter, ok := object.(BoolSetter); ok {
		if setter.SetBool(path, boolValue) {
			return nil
		}
	}

	// Try to set the value using the PointerGetter interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(path); ok {
			if value, ok := pointer.(*bool); ok {
				*value = boolValue
				return nil
			}
		}
	}

	// Cannot set the value
	return derp.Internal(location, "Target Object must be a BoolSetter or PointerGetter", path, object)
}

// setProperty_Integer32 sets a 32-bit integer value in the object.
func setProperty_Integer32(object any, path string, value any) error {

	const location = "schema.setProperty_Integer32"

	// Convert the value to an int
	intValue := convert.Int(value)

	// Try to set the value using the IntSetter interface
	if setter, ok := object.(IntSetter); ok {
		if setter.SetInt(path, intValue) {
			return nil
		}
	}

	// Try to set the value using the PointerGetter interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(path); ok {
			if value, ok := pointer.(*int); ok {
				*value = intValue
				return nil
			}
		}
	}

	// Cannot set the value
	return derp.Internal(location, "Target Object must be an IntSetter or PointerGetter", path, object)
}

// setProperty_Integer64 sets a 64-bit integer value in the object.
func setProperty_Integer64(object any, path string, value any) error {

	const location = "schema.setProperty_Integer64"

	// Convert the value to an int64
	int64Value := convert.Int64(value)

	// Try to set the value using the Int64Setter interface
	if setter, ok := object.(Int64Setter); ok {
		if setter.SetInt64(path, int64Value) {
			return nil
		}
	}

	// Try to set the value using the PointerGetter interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(path); ok {
			if value, ok := pointer.(*int64); ok {
				*value = int64Value
				return nil
			}
		}
	}

	// Cannot set the value
	return derp.Internal(location, "Target Object must be an Int64Setter or PointerGetter", path, object)
}

func setProperty_Number(object any, path string, value any) error {

	const location = "schema.setProperty_Number"

	// Convert the value to a float
	floatValue := convert.Float(value)

	// Try to set the value using the NumberSetter interface
	if setter, ok := object.(FloatSetter); ok {
		if setter.SetFloat(path, floatValue) {
			return nil
		}
	}

	// Try to set the value using the PointerGetter interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(path); ok {
			if value, ok := pointer.(*float64); ok {
				*value = floatValue
				return nil
			}
		}
	}

	// Cannot set the value
	return derp.Internal(location, "Target Object must be a FloatSetter or PointerGetter", path, object)
}

func setProperty_String(object any, path string, value any) error {

	const location = "schema.setProperty_String"

	// Convert the value to a string
	stringValue := convert.String(value)

	// Try to set the value using the StringSetter interface
	if setter, ok := object.(StringSetter); ok {
		if setter.SetString(path, stringValue) {
			return nil
		}
	}

	// Try to set the value using the PointerGetter interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(path); ok {
			if value, ok := pointer.(*string); ok {
				*value = stringValue
				return nil
			}
		}
	}

	// Cannot set the value
	return derp.Internal(location, "Target Object must be a StringSetter or PointerGetter", path, object)
}
