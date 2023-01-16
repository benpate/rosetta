package schema

import (
	"github.com/benpate/derp"
)

// validate is a helper function that validates a sub-element of an object.
// This is used by Object and Array types to validate all of their child entries.
func validate(element Element, object any, name string) derp.MultiError {

	switch typed := element.(type) {

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
		return derp.MultiError{
			derp.NewInternalError("schema.validate", "Unable to validate unknown type", element, object, name),
		}
	}
}

// validate_array specifically validates Array sub-elements
func validate_array(element Array, object any, name string) derp.MultiError {

	if getter, ok := object.(ObjectGetter); ok {
		if value, ok := getter.GetObjectOK(name); ok {
			return element.Validate(value)
		}
	}

	return derp.MultiError{
		derp.NewInternalError("schema.validate_array", "Unable to validate array property", object, name),
	}
}

// validate_boolean specifically validates Boolean sub-elements
func validate_boolean(element Boolean, object any, name string) derp.MultiError {

	if getter, ok := object.(BoolGetter); ok {
		if value, ok := getter.GetBoolOK(name); ok {
			return element.Validate(value)
		}
	}

	return derp.MultiError{
		derp.NewInternalError("schema.validate_boolean", "Unable to validate bool property", object, name),
	}
}

// validate_integer specifically validates Integer sub-elements
func validate_integer(element Integer, object any, name string) derp.MultiError {

	if element.BitSize == 64 {
		return validate_int64(element, object, name)
	}

	return validate_int32(element, object, name)
}

// validate_number specifically validates int32 sub-elements
func validate_int32(element Integer, object any, name string) derp.MultiError {
	if getter, ok := object.(IntGetter); ok {
		if value, ok := getter.GetIntOK(name); ok {
			return element.Validate(value)
		}
	}

	return derp.MultiError{
		derp.NewInternalError("schema.validate_int32", "Unable to validate integer(32) property", object, name),
	}
}

// validate_number specifically validates int64 sub-elements
func validate_int64(element Integer, object any, name string) derp.MultiError {
	if getter, ok := object.(Int64Getter); ok {
		if value, ok := getter.GetInt64OK(name); ok {
			return element.Validate(value)
		}
	}

	return derp.MultiError{
		derp.NewInternalError("schema.validate_int64", "Unable to validate integer(64) property", object, name),
	}
}

// validate_number specifically validates Number sub-elements
func validate_number(element Number, object any, name string) derp.MultiError {
	if getter, ok := object.(FloatGetter); ok {
		if value, ok := getter.GetFloatOK(name); ok {
			return element.Validate(value)
		}
	}

	return derp.MultiError{
		derp.NewInternalError("schema.validate_number", "Unable to validate number property", object, name),
	}
}

// validate_object specifically validates Object sub-elements
func validate_object(element Object, object any, name string) derp.MultiError {

	if getter, ok := object.(ObjectGetter); ok {
		if value, ok := getter.GetObjectOK(name); ok {
			return element.Validate(value)
		}
	}

	return derp.MultiError{
		derp.NewInternalError("schema.validate_object", "Unable to validate object property", object, name),
	}
}

// validate_string specifically validates String sub-elements
func validate_string(element String, object any, name string) derp.MultiError {
	if getter, ok := object.(StringGetter); ok {
		if value, ok := getter.GetStringOK(name); ok {
			return element.Validate(value)
		}
	}

	return derp.MultiError{
		derp.NewInternalError("schema.validate_string", "Unable to validate string property", object, name),
	}
}
