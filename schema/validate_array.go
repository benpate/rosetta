package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// validate_Array checks that the provided value meets the requirements of the schema element.
func validate_Array[T any](element Array, value T) (T, bool, error) {

	const location = "schema.validate_Array"

	// Use the ArrayGetterSetter interface to retrieve and update values in this array.
	getterSetter, isGetterSetter := any(value).(ArrayGetterSetter)

	if !isGetterSetter {
		return value, false, derp.Internal(location, "Value must implement ArrayGetterSetter interface")
	}

	// Validate array length
	if err := validate_Array_length(getterSetter, element); err != nil {
		return value, false, err
	}

	// Track whether any values have been changed during validation
	changed := false

	// Validate each item in the array
	for index := 0; index < getterSetter.Length(); index = index + 1 {

		// Get the current value in the Array
		indexValue, isValid := getterSetter.GetIndex(index)

		if !isValid {
			return value, false, derp.Internal(location, "Getting value at index", index)
		}

		// Validate the value using the schema's "Items" definition
		indexValue, indexChanged, err := validate(element.Items, indexValue)

		if err != nil {
			return value, false, derp.Wrap(err, location, "Validating object at index", index)
		}

		// If the value has been changed, then update the value in the array
		if indexChanged {
			getterSetter.SetIndex(index, indexValue)
			changed = true
		}
	}

	// Return results to caller
	return value, changed, nil
}

// validate_Array_length checks that the length of the array meets the requirements of the schema element.
func validate_Array_length(getter ArrayGetter, element Array) error {

	length := getter.Length()

	// RULE: If the array is required, then it must contain at least one item.
	if element.Required && length == 0 {
		return derp.Validation("Array value is required")
	}

	// RULE: If the array has a minimum length, then it must contain at least that many items.
	if (element.MinLength > 0) && (length < element.MinLength) {
		return derp.Validation("Minimum array length is " + convert.String(element.MinLength))
	}

	// RULE: If the array has a maximum length, then it must contain no more than that many items.
	if (element.MaxLength > 0) && (length > element.MaxLength) {
		return derp.Validation("Maximum array length is " + convert.String(element.MaxLength))
	}

	// Success!
	return nil
}
