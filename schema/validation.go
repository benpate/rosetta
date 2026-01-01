package schema

import (
	"github.com/benpate/derp"
)

// validate is a helper function that validates a sub-element of an object.
// This is used by Object and Array types to validate all of their child entries.
func validate(element Element, object any, name string) error {

	switch typed := element.(type) {

	case Any:
		return nil

	case Array:
		return validate_array(typed, object, name)

	case Boolean:
		return validate_boolean(typed, object, name)

	case Integer:
		return validate_integer(typed, object, name)

	case Number:
		return validate_number(typed, object, name)

	case Object:
		return validate_object(typed, object, name)

	case String:
		return validate_string(typed, object, name)

	default:
		return derp.Internal("schema.validate", "Unable to validate unknown type", element, object, name)
	}
}

// validate_array specifically validates Array sub-elements
func validate_array(element Array, object any, name string) error {

	if getter, ok := object.(PointerGetter); ok {
		if value, ok := getter.GetPointer(name); ok {
			return element.Validate(value)
		}
	}

	// Fall through means that we don't have a PointerGetter.  That's bad...
	return derp.Internal("schema.validate_array", "To validate this property, the Object must be a 'PointerGetter'", object, name)
}

// validate_boolean specifically validates Boolean sub-elements
func validate_boolean(element Boolean, object any, name string) error {

	if getter, ok := object.(BoolGetter); ok {
		if value, ok := getter.GetBoolOK(name); ok {
			return element.Validate(value)
		}
	}

	if getter, ok := object.(PointerGetter); ok {
		if value, ok := getter.GetPointer(name); ok {
			if typed, ok := value.(*bool); ok {
				return element.Validate(*typed)
			}
		}
	}

	if element.Required {
		return derp.Validation("schema.validate_boolean", "Required boolean property is missing", element, object, name)
	}

	return nil
}

// validate_integer specifically validates Integer sub-elements
func validate_integer(element Integer, object any, name string) error {

	if element.BitSize == 64 {
		return validate_int64(element, object, name)
	}

	return validate_int32(element, object, name)
}

// validate_number specifically validates int32 sub-elements
func validate_int32(element Integer, object any, name string) error {

	if getter, ok := object.(IntGetter); ok {
		if value, ok := getter.GetIntOK(name); ok {
			return element.Validate(value)
		}
	}

	if getter, ok := object.(PointerGetter); ok {
		if value, ok := getter.GetPointer(name); ok {
			if typed, ok := value.(*int); ok {
				return element.Validate(*typed)
			}
		}
	}

	if element.Required {
		return derp.Validation("schema.validate_int32", "Required int32 property is missing", element, object, name)
	}

	return nil
}

// validate_number specifically validates int64 sub-elements
func validate_int64(element Integer, object any, name string) error {
	if getter, ok := object.(Int64Getter); ok {
		if value, ok := getter.GetInt64OK(name); ok {
			return element.Validate(value)
		}
	}

	if getter, ok := object.(PointerGetter); ok {
		if value, ok := getter.GetPointer(name); ok {
			if typed, ok := value.(*int64); ok {
				return element.Validate(*typed)
			}
		}
	}

	if element.Required {
		return derp.Validation("schema.validate_int64", "Required int64 property is missing", element, object, name)
	}

	return nil
}

// validate_number specifically validates Number sub-elements
func validate_number(element Number, object any, name string) error {

	if getter, ok := object.(FloatGetter); ok {
		if value, ok := getter.GetFloatOK(name); ok {
			return element.Validate(value)
		}
	}

	if getter, ok := object.(PointerGetter); ok {
		if value, ok := getter.GetPointer(name); ok {
			if typed, ok := value.(*float64); ok {
				return element.Validate(*typed)
			}
		}
	}

	if element.Required {
		return derp.Validation("schema.validate_number", "Required number property is missing", element, object, name)
	}

	return nil
}

// validate_object specifically validates Object sub-elements
func validate_object(element Object, object any, name string) error {

	if getter, ok := object.(PointerGetter); ok {
		if value, ok := getter.GetPointer(name); ok {
			return element.Validate(value)
		}
	}

	return derp.Internal("schema.validate_object", "To validate this property, the Object must be a 'PointerGetter'", object, name)
}

// validate_string specifically validates String sub-elements
func validate_string(element String, object any, name string) error {
	if getter, ok := object.(StringGetter); ok {
		if value, ok := getter.GetStringOK(name); ok {
			return element.Validate(value)
		}
	}

	if getter, ok := object.(PointerGetter); ok {
		if value, ok := getter.GetPointer(name); ok {
			if typed, ok := value.(*string); ok {
				return element.Validate(*typed)
			}
		}
	}

	if element.Required {
		return derp.Validation("schema.validate_string", "Required string property is missing", element, object, name)
	}

	return nil
}
