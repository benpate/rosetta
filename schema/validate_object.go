package schema

import (
	"github.com/benpate/derp"
)

// validate_Object checks that the provided value meets the requirements of the Object schema element,
// and updates the value if necessary.
func validate_Object(element Object, value any) (any, rewriteList, error) {

	const location = "schema.validate_Object"

	rewrites := make(rewriteList, 0)
	allowMissingKeys := false

	// Maps are allowed to have missing keys, but non-maps are not.
	// "Map-ness" must be declared, not inferred.
	if mapTyper, ok := value.(MapTyper); ok {
		allowMissingKeys = mapTyper.IsMap()
	}

	// Validate each property IN THE SCHEMA ELEMENT (not the object)
	// This allows us to ignore properties that are not covered by the schema,
	// which facilitates partial updates and multiple, semi-overlapping schemas per object.
	for key, subElement := range element.Properties {

		// Get the property from the object
		propertyValue, err := getProperty(element, value, key)

		if err != nil {

			// If this is not a map, then this is a legitimate error to return to the caller
			if !allowMissingKeys {
				return nil, nil, derp.Wrap(err, location, "Getting property", key)
			}

			// For maps, a missing property may not be an error (but required values are still required)
			// An absent REQUIRED key fails validation with a clear "required" message,
			// the same outcome as a present-but-empty required value.
			if subElement.IsRequired() {
				return nil, nil, derp.Validation("Required property is missing", key)
			}

			// Otherwise, this property is not required, and an empty map value is fine.
			continue
		}

		// Validate the property value
		changedValue, childRewrites, err := validate(subElement, propertyValue)

		if err != nil {
			return nil, nil, derp.Wrap(err, location, "Validating property", key)
		}

		// If changed, then set the new value in the object, and record
		// the rewrites (prefixed with this property's name) for the caller.
		if len(childRewrites) > 0 {

			if err := SetProperty(element, value, key, changedValue); err != nil {
				return nil, nil, derp.Wrap(err, location, "Setting property", key)
			}

			rewrites = append(rewrites, childRewrites.prefix(key)...)
		}
	}

	return value, rewrites, nil
}
