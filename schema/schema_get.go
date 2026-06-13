package schema

import (
	"strings"

	"github.com/benpate/derp"
)

// Get retrieves a generic value from the object.  If the object is nil,
// Get still tries to return a default value if provided by the schema
func (schema Schema) Get(object any, path string) (any, error) {
	return getPropertyRecursive(schema.Element, object, path)
}

func getPropertyRecursive(element Element, object any, path string) (any, error) {

	const location = "schema.getPropertyRecursive"

	// If the path is empty, then try to get the value of the entire object.
	if path == "" {

		if getter, ok := object.(ValueGetter); ok {
			return getter.GetValue(), nil
		}

		return object, nil
	}

	// Split the path into head and tail
	head, tail, _ := strings.Cut(path, ".")

	// Get the property value from the object
	result, err := getProperty(element, object, head)

	if err != nil {
		return nil, derp.Wrap(err, location, "Getting property", object, head)
	}

	// If this is the end of the path, then return the result
	if tail == "" {
		return result, nil
	}

	// Otherwise, get the property definition for the intermediate result we have right now.
	subElement, ok := element.GetElement(head)

	if !ok {
		return nil, derp.BadRequest(location, "Invalid property", head)
	}

	// Continue digging until we exhaust the path.
	return getPropertyRecursive(subElement, result, tail)
}

// getProperty retrieves a generic value from the object.
func getProperty(element Element, object any, name string) (any, error) {

	const location = "schema.getProperty"

	// Find the property definition for the first path segment
	subElement, ok := element.GetElement(name)

	if !ok {
		return nil, derp.BadRequest(location, "Invalid property", name)
	}

	// Use the element type to find the correct property getter
	switch typed := subElement.(type) {

	case Any, Array, Object:
		return getProperty_PointerOnly(object, name)

	case Boolean:
		return getProperty_Boolean(object, name)

	case Integer:

		if typed.BitSize == 64 {
			return getProperty_Integer64(object, name)
		}

		return getProperty_Integer32(object, name)

	case Number:
		return getProperty_Number(object, name)

	case String:
		return getProperty_String(object, name)
	}

	// You suck.
	return nil, derp.Internal(location, "Unable to get property", name, object)
}

// getProperty_PointerOnly retrieves a value from the object using only the PointerGetter interface.
func getProperty_PointerOnly(object any, name string) (any, error) {

	const location = "schema.getProperty_PointerOnly"

	// Try to make a PointerGetter
	getter, ok := object.(PointerGetter)
	if !ok {
		return nil, derp.Internal(location, "Object must be a PointerGetter", object)
	}

	// Use the PointerGetter to retrieve the value for this property
	result, ok := getter.GetPointer(name)

	if !ok {
		return nil, derp.Internal(location, "Getting pointer to property", name, object)
	}

	return result, nil
}

func getProperty_Boolean(object any, name string) (any, error) {

	const location = "schema.getProperty_Boolean"

	// If this is a BoolGetter, then try to retrieve the value using that interface
	if getter, ok := object.(BoolGetter); ok {
		if result, ok := getter.GetBoolOK(name); ok {
			return result, nil
		}
	}

	// If this is a PointerGetter, then try to retrieve the value using that interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(name); ok {
			if value, ok := pointer.(*bool); ok {
				return *value, nil
			}
		}
	}

	// Failure.
	return nil, derp.Internal(location, "Object must be a BoolGetter or a PointerGetter", object)
}

// getProperty_Integer32 retrieves an int value from the object.
func getProperty_Integer32(object any, name string) (any, error) {

	const location = "schema.getProperty_Integer32"

	// If this is an IntGetter, then try to retrieve the value using that interface
	if getter, ok := object.(IntGetter); ok {
		if result, ok := getter.GetIntOK(name); ok {
			return result, nil
		}
	}

	// If this is a PointerGetter, then try to retrieve the value using that interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(name); ok {
			if value, ok := pointer.(*int); ok {
				return *value, nil
			}
		}
	}

	// Failure.
	return nil, derp.Internal(location, "Object must be an IntGetter or PointerGetter", object)
}

// getProperty_Integer64 retrieves an int64 value from the object.
func getProperty_Integer64(object any, name string) (any, error) {

	const location = "schema.getProperty_Integer64"

	// If this is an Int64Getter, then try to retrieve the value using that interface
	if getter, ok := object.(Int64Getter); ok {
		if result, ok := getter.GetInt64OK(name); ok {
			return result, nil
		}
	}

	// If this is a PointerGetter, then try to retrieve the value using that interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(name); ok {
			if value, ok := pointer.(*int64); ok {
				return *value, nil
			}
		}
	}

	// Failure.
	return nil, derp.Internal(location, "Object must be an Int64Getter or PointerGetter", object)
}

// getProperty_Number retrieves a number value from the object.
func getProperty_Number(object any, name string) (any, error) {

	const location = "schema.getProperty_Number"

	// If this is a FloatGetter, then try to retrieve the value using that interface
	if getter, ok := object.(FloatGetter); ok {
		if result, ok := getter.GetFloatOK(name); ok {
			return result, nil
		}
	}

	// If this is a PointerGetter, then try to retrieve the value using that interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(name); ok {
			if value, ok := pointer.(*float64); ok {
				return *value, nil
			}
		}
	}

	// Failure.
	return nil, derp.Internal(location, "Object must be a FloatGetter or PointerGetter", object)
}

// getProperty_String retrieves a string value from the object.
func getProperty_String(object any, name string) (any, error) {

	const location = "schema.getProperty_String"

	// If this is a StringGetter, then try to retrieve the value using that interface
	if getter, ok := object.(StringGetter); ok {
		if result, ok := getter.GetStringOK(name); ok {
			return result, nil
		}
	}

	// If this is a PointerGetter, then try to retrieve the value using that interface
	if getter, ok := object.(PointerGetter); ok {
		if pointer, ok := getter.GetPointer(name); ok {
			if value, ok := pointer.(*string); ok {
				return *value, nil
			}
		}
	}

	// Failure.
	return nil, derp.Internal(location, "Object must be a StringGetter or PointerGetter", object)
}
