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

	case Any, Array, Object:

		getter, ok := object.(PointerGetter)
		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Object must be a PointerGetter", object)
		}

		subObject, ok := getter.GetPointer(head)
		if !ok {
			return nil, derp.NewInternalError("schema.Schema.get", "Unable to get object", head, object)
		}

		return schema.get(subObject, typed, tail)

	case Boolean:

		if getter, ok := object.(BoolGetter); ok {
			if result, ok := getter.GetBoolOK(head); ok {
				return result, nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*bool); ok {
					return *value, nil
				}
			}
		}

		return nil, derp.NewInternalError("schema.Schema.get", "Unable to get bool from BoolGetter or a PointerGettter", object)

	case Integer:

		if typed.BitSize == 64 {

			if getter, ok := object.(Int64Getter); ok {
				if result, ok := getter.GetInt64OK(head); ok {
					return result, nil
				}
			}

			if getter, ok := object.(PointerGetter); ok {
				if pointer, ok := getter.GetPointer(head); ok {
					if value, ok := pointer.(*int64); ok {
						return *value, nil
					}
				}
			}

			return nil, derp.NewInternalError("schema.Schema.get", "Unable to get int64 from Int64Getter or PointerGetter", object)
		}

		if getter, ok := object.(IntGetter); ok {
			if result, ok := getter.GetIntOK(head); ok {
				return result, nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*int); ok {
					return *value, nil
				}
			}
		}

		return nil, derp.NewInternalError("schema.Schema.get", "Unable to get int from IntGetter or PointerGetter", object)

	case Number:

		if getter, ok := object.(FloatGetter); ok {
			if result, ok := getter.GetFloatOK(head); ok {
				return result, nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*float64); ok {
					return *value, nil
				}
			}
		}

		return nil, derp.NewInternalError("schema.Schema.get", "Unable to get float from FloatGetter or PointerGetter", object)

	case String:

		if getter, ok := object.(StringGetter); ok {
			if result, ok := getter.GetStringOK(head); ok {
				return result, nil
			}
		}

		if getter, ok := object.(PointerGetter); ok {
			if pointer, ok := getter.GetPointer(head); ok {
				if value, ok := pointer.(*string); ok {
					return *value, nil
				}
			}
		}

		return nil, derp.NewInternalError("schema.Schema.get", "Unable to get string from StringGetter or PointerGetter", object)
	}

	return nil, derp.NewInternalError("schema.Schema.get", "Unable to get property", path, object)
}
