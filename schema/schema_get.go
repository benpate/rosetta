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
		return object, nil
	}

	head, tail := path.Split()

	subElement, ok := element.GetElement(head)

	if !ok {
		return nil, derp.NewInternalError("schema.Schema.get", "Unknown property", head)
	}

	switch typed := subElement.(type) {

	case Array, Object:

		getter, ok := object.(ObjectGetter)
		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Object must be an ObjectGetter", object)
		}

		subObject, ok := getter.GetObject(head)
		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Unable to get object", head, object)
		}

		return schema.get(subObject, typed, tail)

	case Boolean:

		getter, ok := object.(BoolGetter)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Object must be a BoolGetter", object)
		}

		result, ok := getter.GetBoolOK(head)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Unable to get bool", head, object)
		}

		return result, nil

	case Integer:

		if typed.BitSize == 64 {

			getter, ok := object.(Int64Getter)

			if !ok {
				return nil, derp.NewInternalError("schema.Schema.get", "Object must be a Int64Getter", object)
			}

			result, ok := getter.GetInt64OK(head)

			if !ok {
				return nil, derp.NewInternalError("schema.Schema.get", "Unable to get int64", head, object)
			}
			return result, nil

		}

		getter, ok := object.(IntGetter)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Object must be a IntGetter", object)
		}

		result, ok := getter.GetIntOK(head)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Unable to get int", head, object)
		}

		return result, nil

	case Number:

		getter, ok := object.(FloatGetter)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Object must be a FloatGetter", object)
		}

		result, ok := getter.GetFloatOK(head)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Unable to get float", head, object)
		}

		return result, nil

	case String:

		getter, ok := object.(StringGetter)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Object must be a StringGetter", object)
		}

		result, ok := getter.GetStringOK(head)

		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Unable to get string", head, object)
		}

		return result, nil
	}

	return nil, derp.NewInternalError("schema.Schema.get", "Unable to get property", path, object)
}
