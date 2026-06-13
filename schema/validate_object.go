package schema

import (
	"github.com/benpate/derp"
)

// validate_Object checks that the provided value meets the requirements of the Object schema element,
// and updates the value if necessary.
func validate_Object(element Object, value any) (any, bool, error) {

	const location = "schema.validate_Object"

	objectChanged := false

	// Validate each property IN THE ELEMENT (not the object)
	// This allows us to ignore properties that are not covered by the schema,
	// which facilitates partial updates and multiple, semi-overlapping schemas per object.
	for key, subElement := range element.Properties {

		// Get the property from the object
		propertyValue, err := getProperty(element, value, key)

		if err != nil {
			return nil, false, derp.Wrap(err, location, "Getting property")
		}

		// Validate the property value
		changedValue, itemChanged, err := validate(subElement, propertyValue)

		if err != nil {
			return nil, false, derp.Wrap(err, location, "Validating property")
		}

		// If changed, then set the new value in the object
		if itemChanged {

			if err := setProperty(element, value, key, changedValue); err != nil {
				return nil, false, derp.Wrap(err, location, "Setting property")
			}

			objectChanged = true
		}
	}

	return value, objectChanged, nil
}
