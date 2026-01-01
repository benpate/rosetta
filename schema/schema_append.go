package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/pointer"
)

// Append appends a value to the array at the specified path within the object according to this schema
func (schema Schema) Append(object any, path string, value any) error {

	const location = "schema.Schema.Append"

	element, exists := schema.GetArrayElement(path)

	if !exists {
		return derp.Internal(location, "Element must be an ArrayError finding schema element", path)
	}

	// Get Original Value
	originalValue, err := schema.Get(object, path)

	if err != nil {
		return derp.Wrap(err, location, "Error getting value", path)
	}

	// This will only work for ArraySetter objects
	setter, isSetter := pointer.To(originalValue).(ArraySetter)

	if !isSetter {
		return derp.Internal(location, "Value must implement ArraySetter interface", setter, isSetter, path)
	}

	// Append to the array
	if err := element.Append(setter, value); err != nil {
		return derp.Wrap(err, location, "Error appending value", path)
	}

	// Set the value back into the object
	finalValue := indirect(setter)
	if err := schema.Set(object, path, finalValue); err != nil {
		return derp.Wrap(err, location, "Error re-populating value", path)
	}

	return nil
}
