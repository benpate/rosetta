package slice

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// AnyAppend adds a new (empty) item to the end of an array or slice of any type.
func AnyAppend(array any, newItem any) (any, int, error) {

	newItemType := reflect.TypeOf(newItem)

	// If the original value was just nil, then make a new slice using the type of the new item
	if array == nil {
		resultType := reflect.SliceOf(newItemType)
		result := reflect.MakeSlice(resultType, 1, 1)
		result.Index(0).Set(reflect.ValueOf(newItem))

		return result.Interface(), 0, nil
	}

	valueOfArray := convert.ReflectValue(array)
	arrayType := valueOfArray.Type()
	elementType := arrayType.Elem()

	if !newItemType.AssignableTo(elementType) {
		return nil, 0, derp.NewInternalError("slice.AppendToInterface", "New item is not assignable to the array element type", newItem)
	}

	switch valueOfArray.Kind() {

	// Create a new array that is one item larger, and copy the existing array into it
	case reflect.Array:

		// Create a new array that's one item longer
		oldLength := valueOfArray.Len()
		newLength := oldLength + 1
		resultType := reflect.ArrayOf(newLength, elementType)
		result := reflect.New(resultType).Elem()

		// Copy from the existing array into the new array
		for index := 0; index < oldLength; index++ {
			result.Index(index).Set(valueOfArray.Index(index))
		}

		// Add an "initialized" empty array to the end of the array
		result.Index(oldLength).Set(reflect.ValueOf(newItem))

		return result.Interface(), result.Len() - 1, nil

	// Append a new item to the end of the slice.
	case reflect.Slice:

		// Use the built-in "Append" method to add a new item to the slice
		result := reflect.Append(valueOfArray, reflect.ValueOf(newItem))
		return result.Interface(), result.Len() - 1, nil
	}

	// Otherwise, we don't know how to append to this type of value
	return nil, 0, derp.NewInternalError("slice.AppendInterface", "Value must be an array or slice", array)
}

// AnyRemove removes an item from an array or slice at the specified index.
func AnyRemove(value any, removeIndex int) (any, error) {

	valueOf := convert.ReflectValue(value)

	// Validate that we have the right kind of value (must be an array or slice)
	switch valueOf.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, derp.NewInternalError("slice.RemoveFromInterface", "Value must be an array or slice", value)
	}

	// Length for bounds checking...
	length := valueOf.Len()

	// Soft fail if index is out of bounds
	if (removeIndex < 0) || (removeIndex >= length) {
		return value, nil
	}

	// Create a new value (array or slice) that is one item shorter
	var result reflect.Value
	elementType := valueOf.Type().Elem()

	switch valueOf.Kind() {

	case reflect.Array:
		resultType := reflect.ArrayOf(length-1, elementType)
		result = reflect.Zero(resultType)

	case reflect.Slice:
		resultType := reflect.SliceOf(elementType)
		result = reflect.MakeSlice(resultType, length-1, length-1)

	default:
		return nil, derp.NewInternalError("slice.RemoveFromInterface", "Value must be an array or slice", value)
	}

	// Copy items from the old value into the new value, skipping the item at the specified index
	newIndex := 0
	for oldIndex := 0; oldIndex < length; oldIndex++ {

		if newIndex == removeIndex {
			continue
		}

		result.Index(newIndex).Set(valueOf.Index(oldIndex))
		newIndex += 1
	}

	// Success!!
	return result.Interface(), nil
}

// AnySort (will) sort an array or slice of any type.
func AnySort(value any) (any, error) {
	return value, nil
}
