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

	const location = "schema.SetElement"

	// In rare cases, we may need to set an entire element in one call.
	// If the path is empty, then we're setting the entire element, which
	// must implement the "ValueSetter" interface
	if path.IsEmpty() {
		if setter, ok := object.(ValueSetter); ok {
			if err := setter.SetValue(value); err != nil {
				return derp.Wrap(err, location, "Error setting value", object, value)
			}
			return nil
		}
		return derp.NewInternalError(location, "Cannot set values on empty path")
	}

	// Otherwise, we're setting a sub-element within the object:

	head, tail := path.Split()
	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.NewInternalError(location, "Unknown property", head, element)
	}

	// Different interfaces are required for different types of objects
	switch typed := subElement.(type) {

	case Any:

		// ObjectSetter interface is required for Maps
		if setter, ok := object.(ObjectSetter); ok {
			return setter.SetObject(element, path, value)
		}
		return derp.NewInternalError(location, "To set an 'Any' value, the target Object must be an ObjectSetter", object)

	case Array, Object:

		// ObjectSetter interface is required for Maps
		if setter, ok := object.(ObjectSetter); ok {
			return setter.SetObject(element, path, value)
		}

		// ObjectGetter works for Structs, Slices, and Arrays
		if getter, ok := object.(ObjectGetter); ok {
			if subObject, ok := getter.GetObject(head); ok {
				return SetElement(subObject, typed, tail, value)
			}
		}

		return derp.NewInternalError(location, "To set an 'Array' or 'Object' value, the target Object must be an ObjectSetter or ObjectGetter", object)

	case Boolean:
		if setter, ok := object.(BoolSetter); ok {
			boolValue, _ := convert.BoolOk(value, false)
			if setter.SetBool(head, boolValue) {
				return nil
			} else {
				return derp.NewInternalError(location, "Unable to set boolean value", path, subElement, object)
			}
		}

		return derp.NewInternalError(location, "To set a 'Boolean' value, the target Object must be a BoolSetter", object)

	case Integer:
		if typed.BitSize == 64 {
			if setter, ok := object.(Int64Setter); ok {
				int64Value, _ := convert.Int64Ok(value, 0)
				if setter.SetInt64(head, int64Value) {
					return nil
				} else {
					return derp.NewInternalError(location, "Unable to set int64 value", path, subElement, object)
				}
			}
			return derp.NewInternalError(location, "To set a 64-bit 'Integer' value, the target Object must be an Int64Setter", object)
		}

		intValue, _ := convert.IntOk(value, 0)
		if setter, ok := object.(IntSetter); ok {
			if setter.SetInt(head, intValue) {
				return nil
			} else {
				return derp.NewInternalError(location, "Unable to set int value", path, subElement, object)
			}
		}
		return derp.NewInternalError(location, "To set an 'Integer' value, the target Object must be an IntSetter", object)

	case Number:
		floatValue, _ := convert.FloatOk(value, 0)
		if setter, ok := object.(FloatSetter); ok {
			if setter.SetFloat(head, floatValue) {
				return nil
			} else {
				return derp.NewInternalError(location, "Unable to set float value", path, subElement, object)
			}
		}
		return derp.NewInternalError(location, "To set a 'Number' value, the target Object must be a FloatSetter", object)

	case String:
		stringValue, _ := convert.StringOk(value, "")
		if setter, ok := object.(StringSetter); ok {
			if setter.SetString(head, stringValue) {
				return nil
			} else {
				return derp.NewInternalError(location, "Unable to set string value", path, subElement, object)
			}
		}

		return derp.NewInternalError(location, "To set a 'String' value, the target Object must be a StringSetter", object)
	}

	return derp.NewInternalError(location, "Unsupported element type", path, subElement, object)
}
