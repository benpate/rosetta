package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// validate_Array checks that the provided value meets the requirements of the schema element, and updates the value if necessary.
func validate_Array[T any](element Array, value T) (T, rewriteList, error) {

	const location = "schema.validate_Array"

	// Use the ArrayGetterSetter interface to retrieve and update values in this array.
	getterSetter, isGetterSetter := any(value).(ArrayGetterSetter)

	if !isGetterSetter {
		return value, nil, derp.Internal(location, "Value must implement ArrayGetterSetter interface")
	}

	// Validate array length
	if err := validate_Array_length(getterSetter, element); err != nil {
		return value, nil, err
	}

	// Track every value that is changed during validation
	rewrites := make(rewriteList, 0)

	// Validate each item in the array
	for index := 0; index < getterSetter.Length(); index = index + 1 {

		// Get the current value in the Array
		indexValue, isValid := getterSetter.GetIndex(index)

		if !isValid {
			return value, nil, derp.Internal(location, "Getting value at index", index)
		}

		// RULE: Composite items (Object/Array) reach their schema accessors through
		// pointer-receiver interfaces (PointerGetter, etc.). GetIndex hands back a copy
		// by value, whose address is not the slice element's, so wrap it in an addressable
		// pointer before validating. Scalars validate fine by value and are left alone.
		itemValue, restore := addressableItem(element.Items, indexValue)

		// Validate the value using the schema's "Items" definition
		validatedValue, itemRewrites, err := validate(element.Items, itemValue)

		if err != nil {
			return value, nil, derp.Wrap(err, location, "Validating object at index", index)
		}

		// If the value has been changed, then update the value in the array, and record
		// the rewrites (prefixed with this item's index) for the caller.
		if len(itemRewrites) > 0 {
			getterSetter.SetIndex(index, restore(validatedValue))
			rewrites = append(rewrites, itemRewrites.prefix(convert.String(index))...)
		}
	}

	// Return results to caller
	return value, rewrites, nil
}

// addressableItem prepares an array element for nested validation. For composite item
// schemas (Object/Array/Any) whose value is a non-pointer, it returns an addressable
// pointer to a copy so that pointer-receiver accessors are reachable, along with a
// "restore" function that converts the validated result back into the element's original
// shape for SetIndex. Scalars (and values that are already pointers) pass through unchanged.
func addressableItem(items Element, indexValue any) (any, func(any) any) {

	identity := func(v any) any { return v }

	// Only composite item schemas need pointer-based access.
	switch items.(type) {
	case Object, Array, Any:
		// proceed
	default:
		return indexValue, identity
	}

	reflectValue := reflect.ValueOf(indexValue)

	// A nil value, or one that is already a pointer, needs no wrapping.
	if !reflectValue.IsValid() || reflectValue.Kind() == reflect.Pointer {
		return indexValue, identity
	}

	// Build an addressable pointer to a copy of the value.
	pointer := addressableCopy(reflectValue)

	// The restore function dereferences the (possibly-mutated) copy so SetIndex receives
	// the same concrete type GetIndex returned.
	restore := func(validated any) any {
		return restoreType(validated, reflectValue.Type())
	}

	return pointer, restore
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

// addressableCopy returns a pointer to a copy of the value.  A plain Go map has no getters behind
// its pointer, so it is copied into a mapof.Any instead, exactly as Get reads it.
func addressableCopy(original reflect.Value) any {

	pointer := reflect.New(original.Type())
	pointer.Elem().Set(original)

	if original.Kind() != reflect.Map {
		return pointer.Interface()
	}

	if _, isGetter := pointer.Interface().(PointerGetter); isGetter {
		return pointer.Interface()
	}

	// pointerToCopy returns a mapof.Any for a plain map, and the map itself when none can hold it
	if readable, ok := pointerToCopy(original.Interface()).(PointerGetter); ok {
		return readable
	}

	return pointer.Interface()
}

// restoreType dereferences a validated pointer, and returns its value as the original type.  Any
// other value is returned unchanged.
func restoreType(validated any, originalType reflect.Type) any {

	validatedPointer := reflect.ValueOf(validated)

	if (validatedPointer.Kind() != reflect.Pointer) || validatedPointer.IsNil() {
		return validated
	}

	// A plain map read through a mapof.Any is stored back as its own type
	result := validatedPointer.Elem()

	if (result.Type() != originalType) && result.Type().ConvertibleTo(originalType) {
		result = result.Convert(originalType)
	}

	return result.Interface()
}
