package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

func (schema Schema) Set(object any, path string, value any) error {
	return SetElement(object, schema.Element, list.ByDot(path), value)
}

// SetAll iterates over Set to apply all of the values to the object one at a time, stopping
// at the first error it encounters.  If all values are addedd successfully, then SetAll
// also uses Validate() to confirm that the object is still correct.
func (schema Schema) SetAll(object any, values map[string]any) error {

	const location = "schema.Schema.SetAll"

	// Set each value in the schema
	for path, value := range values {

		// Errors are intentionally ignored here.
		// Unallowed data does not make it through the schema filter
		schema.Set(object, path, value)
	}

	// Validate the whole schema once all the values are set
	if err := schema.Validate(object); err != nil {
		return derp.Wrap(err, location, "Validation Error")
	}

	// Success!!
	return nil
}

func SetElement(object any, element Element, path list.List, value any) error {

	// In rare cases, we may need to set an entire element in one call.
	// If the path is empty, then we're setting the entire element, which
	// must implement the "ValueSetter" interface
	if path.IsEmpty() {
		if setter, ok := object.(ValueSetter); ok {
			if err := setter.SetValue(value); err != nil {
				return derp.Wrap(err, "schema.SetElement", "Error setting value", object, value)
			}
			return nil
		}
		return derp.NewInternalError("schema.SetElement", "Cannot set values on empty path")
	}

	// Otherwise, we're setting a sub-element within the object:

	head, tail := path.Split()
	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.NewInternalError("schema.SetElement", "Unknown property", head, element)
	}

	// Different interfaces are required for different types of objects
	switch typed := subElement.(type) {

	case Array, Object:

		// ObjectSetter interface is required for Maps
		if setter, ok := object.(ObjectSetter); ok {
			return setter.SetObject(element, path, value)
		}

		// ObjectGetter works for Structs, Slices, and Arrays
		if getter, ok := object.(ObjectGetter); ok {
			if subObject, ok := getter.GetObjectOK(head); ok {
				return SetElement(subObject, typed, tail, value)
			}
		}

	case Boolean:
		boolValue, _ := convert.BoolOk(value, false)
		if setter, ok := object.(BoolSetter); ok {
			if setter.SetBool(head, boolValue) {
				return nil
			}
		}

	case Integer:
		if typed.BitSize == 64 {
			int64Value, _ := convert.Int64Ok(value, 0)
			if setter, ok := object.(Int64Setter); ok {
				if setter.SetInt64(head, int64Value) {
					return nil
				}
			}
		}

		intValue, _ := convert.IntOk(value, 0)
		if setter, ok := object.(IntSetter); ok {
			if setter.SetInt(head, intValue) {
				return nil
			}
		}

	case Number:
		floatValue, _ := convert.FloatOk(value, 0)
		if setter, ok := object.(FloatSetter); ok {
			if setter.SetFloat(head, floatValue) {
				return nil
			}
		}

	case String:
		stringValue, _ := convert.StringOk(value, "")
		if setter, ok := object.(StringSetter); ok {
			if setter.SetString(head, stringValue) {
				return nil
			}
		}
	}

	return derp.NewInternalError("schema.SetElement", "Unable to set property", path, object)
}
