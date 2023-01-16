package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

func (schema Schema) Set(object any, path string, value any) error {
	return schema.set(object, schema.Element, list.ByDot(path), value)
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

func (schema Schema) set(object any, element Element, path list.List, value any) error {

	if path.IsEmpty() {
		return derp.NewInternalError("schema.Schema.set", "Cannot set values on empty path")
	}

	head, tail := path.Split()

	subElement, ok := element.getProperty(head)

	if !ok {
		return derp.NewInternalError("schema.Schema.set", "Unknown property", head)
	}

	switch typed := subElement.(type) {

	case Array, Object:
		if getter, ok := object.(ObjectGetter); ok {
			if subObject, ok := getter.GetObjectOK(head); ok {
				return schema.set(subObject, typed, tail, value)
			}
		}

	case Boolean:
		if boolValue, ok := value.(bool); ok {
			if setter, ok := object.(BoolSetter); ok {
				if setter.SetBoolOK(head, boolValue) {
					return nil
				}
			}
		}

	case Integer:
		if typed.BitSize == 64 {
			if int64Value, ok := value.(int64); ok {
				if setter, ok := object.(Int64Setter); ok {
					if setter.SetInt64OK(head, int64Value) {
						return nil
					}
				}
			}
		}

		if intValue, ok := value.(int); ok {
			if setter, ok := object.(IntSetter); ok {
				if setter.SetIntOK(head, intValue) {
					return nil
				}
			}
		}

	case Number:
		if floatValue, ok := value.(float64); ok {
			if setter, ok := object.(FloatSetter); ok {
				if setter.SetFloatOK(head, floatValue) {
					return nil
				}
			}
		}

	case String:
		if stringValue, ok := value.(string); ok {
			if setter, ok := object.(StringSetter); ok {
				if setter.SetStringOK(head, stringValue) {
					return nil
				}
			}
		}
	}

	return derp.NewInternalError("schema.Schema.set", "Unable to set property", path, object)
}
