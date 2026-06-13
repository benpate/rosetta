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
		return derp.Internal(location, "Schema element must be an array", path)
	}

	// Get Original Value
	originalValue, err := schema.Get(object, path)

	if err != nil {
		return derp.Wrap(err, location, "Getting value", path)
	}

	// This will only work for ArraySetter objects
	setter, isSetter := pointer.To(originalValue).(ArraySetter)

	if !isSetter {
		return derp.Internal(location, "Value must implement ArraySetter interface", setter, isSetter, path)
	}

	// Append to the array
	if err := element.Append(setter, value); err != nil {
		return derp.Wrap(err, location, "Appending value", path)
	}

	// Write the grown array back into the object. We pass the pointer (setter)
	// rather than a dereferenced value so that Set's validation can treat it as
	// an ArrayGetterSetter; convert dereferences the pointer when persisting.
	if err := schema.Set(object, path, setter); err != nil {
		return derp.Wrap(err, location, "Re-populating value", path)
	}

	return nil
}
