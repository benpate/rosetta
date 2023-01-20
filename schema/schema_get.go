package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

// Get retrieves a generic value from the object.  If the object is nil,
// Get still tries to return a default value if provided by the schema
func (schema Schema) Get(object any, path string) (any, error) {
	return schema.get(object, schema.Element, list.ByDot(path))
}

func (schema Schema) get(object any, element Element, path list.List) (any, error) {

	if path.IsEmpty() {
		return nil, derp.NewInternalError("schema.Schema.get", "Cannot get values on empty path")
	}

	head, tail := path.Split()

	subElement, ok := element.GetElement(head)

	if !ok {
		return nil, derp.NewInternalError("schema.Schema.get", "Unknown property", head)
	}

	switch typed := subElement.(type) {

	case Array, Object:
		if getter, ok := object.(ObjectGetter); ok {
			if subObject, ok := getter.GetObjectOK(head); ok {
				return schema.get(subObject, typed, tail)
			}
		}

	case Boolean:
		if getter, ok := object.(BoolGetter); ok {
			if result, ok := getter.GetBoolOK(head); ok {
				return result, nil
			}
		}

	case Integer:
		if typed.BitSize == 64 {
			if getter, ok := object.(Int64Getter); ok {
				if result, ok := getter.GetInt64OK(head); ok {
					return result, nil
				}
			}
		}
		if getter, ok := object.(IntGetter); ok {
			if result, ok := getter.GetIntOK(head); ok {
				return result, nil
			}
		}

	case Number:
		if getter, ok := object.(FloatGetter); ok {
			if result, ok := getter.GetFloatOK(head); ok {
				return result, nil
			}
		}

	case String:
		if getter, ok := object.(StringGetter); ok {
			if result, ok := getter.GetStringOK(head); ok {
				return result, nil
			}
		}
	}

	return nil, derp.NewInternalError("schema.Schema.get", "Unable to get property", path, object)
}
