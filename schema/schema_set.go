package schema

import (
	"net/url"

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
		// nolint: errcheck
		schema.Set(object, path, value)
	}

	// Validate the whole schema once all the values are set
	if err := schema.Validate(object); err != nil {
		return derp.Wrap(err, location, "Validation Error")
	}

	// Success!!
	return nil
}

// SetURLValues iterates over Set to apply all of the values to the object one at a time, stopping
// at the first error it encounters.  If all values are addedd successfully, then SetURLValues
// also uses Validate() to confirm that the object is still correct.
func (schema Schema) SetURLValues(object any, values url.Values) error {

	const location = "schema.Schema.SetAll"

	// Set each value in the schema
	for path, value := range values {

		// Errors are intentionally ignored here.
		// Unallowed data does not make it through the schema filter
		// nolint: errcheck
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
		return derp.NewInternalError(location, "Cannot set values on empty path", object, element, path, value)
	}

	// Otherwise, we're setting a sub-element within the object:

	head, tail := path.Split()
	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.NewInternalError(location, "Property does not exist in schema", head)
	}

	// Different interfaces are required for different types of objects
	switch typed := subElement.(type) {

	case Any, Array, Object:

		// ObjectSetter interface is required for Maps
		if setter, ok := object.(ObjectSetter); ok {
			return setter.SetObject(element, path, value)
		}

		// PointerGetter works for Structs, Slices, and Arrays
		if getter, ok := object.(PointerGetter); ok {
			if subPointer, ok := getter.GetPointer(head); ok {
				return SetElement(subPointer, typed, tail, value)
			}
		}
		return derp.NewInternalError(location, "To set an 'Any', 'Array', or 'Object' value, the target Object must be an ObjectSetter or PointerGetter", object)

	case Boolean:
		boolValue, _ := convert.BoolOk(value, false)
		if setter, ok := object.(BoolSetter); ok {
			if setter.SetBool(head, boolValue) {
				return nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*bool); ok {
					*value = boolValue
					return nil
				}
			}
		}

		return derp.NewInternalError(location, "To set a 'Boolean' value, the target Object must be a BoolSetter or PointerGetter", object)

	case Integer:
		if typed.BitSize == 64 {
			int64Value, _ := convert.Int64Ok(value, 0)
			if setter, ok := object.(Int64Setter); ok {
				if setter.SetInt64(head, int64Value) {
					return nil
				}
			}

			if getter, ok := object.(PointerGetter); ok {
				if pointer, ok := getter.GetPointer(head); ok {
					if value, ok := pointer.(*int64); ok {
						*value = int64Value
						return nil
					}
				}
			}

			return derp.NewInternalError(location, "To set a 64-bit 'Integer' value, the target Object must be an Int64Setter or PointerGetter", object)
		}

		intValue, _ := convert.IntOk(value, 0)
		if setter, ok := object.(IntSetter); ok {
			if setter.SetInt(head, intValue) {
				return nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*int); ok {
					*value = intValue
					return nil
				}
			}
		}

		return derp.NewInternalError(location, "To set an 'Integer' value, the target Object must be an IntSetter or PointerGetter", object)

	case Number:
		floatValue, _ := convert.FloatOk(value, 0)
		if setter, ok := object.(FloatSetter); ok {
			if setter.SetFloat(head, floatValue) {
				return nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*float64); ok {
					*value = floatValue
					return nil
				}
			}
		}

		return derp.NewInternalError(location, "To set a 'Number' value, the target Object must be a FloatSetter or PointerGetter", object)

	case String:
		stringValue, _ := convert.StringOk(value, "")
		if setter, ok := object.(StringSetter); ok {
			if setter.SetString(head, stringValue) {
				return nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*string); ok {
					*value = stringValue
					return nil
				}
			}
		}

		return derp.NewInternalError(location, "To set a 'String' value, the target Object must be a StringSetter or PointerGetter", object)
	}

	return derp.NewInternalError(location, "Unsupported element type", path, subElement, object)
}
