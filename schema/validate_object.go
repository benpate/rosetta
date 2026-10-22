package schema

import (
	"slices"

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

		// Validate the property value.  A nested map held by value is wrapped in a pointer first,
		// because map setters such as mapof.Any's take pointer receivers.
		childValue, restore := addressableItem(subElement, propertyValue)
		changedValue, childRewrites, err := validate(subElement, childValue)

		if err != nil {
			return nil, nil, derp.Wrap(err, location, "Validating property", key)
		}

		// If changed, then set the new value in the object, and record
		// the rewrites (prefixed with this property's name) for the caller.
		if len(childRewrites) > 0 {

			if err := SetProperty(element, value, key, restore(changedValue)); err != nil {
				return nil, nil, derp.Wrap(err, location, "Setting property", key)
			}

			rewrites = append(rewrites, childRewrites.prefix(key)...)
		}
	}

	// Validate every key that only the Wildcard describes
	wildcardRewrites, err := validate_ObjectWildcard(element, value)

	if err != nil {
		return nil, nil, derp.Wrap(err, location, "Validating wildcard properties")
	}

	rewrites = append(rewrites, wildcardRewrites...)

	return value, rewrites, nil
}

// validate_ObjectWildcard validates each key of a declared map that no named property describes,
// using the Object's Wildcard element.
func validate_ObjectWildcard(element Object, value any) (rewriteList, error) {

	const location = "schema.validate_ObjectWildcard"

	rewrites := make(rewriteList, 0)

	for _, key := range wildcardKeys(element, value) {

		// Read the key itself, never as a path: a wildcard key may contain dots, or be empty
		propertyValue, err := getPropertyAs(element.Wildcard, value, key)

		// A value the Wildcard cannot read is skipped, as an unreadable named property of a map is,
		// unless the Wildcard is required
		if err != nil {

			if element.Wildcard.IsRequired() {
				return nil, derp.Validation("Required property is missing", key)
			}

			continue
		}

		// Validate the property value against the Wildcard.  A nested map held by value is wrapped
		// in a pointer first, because map setters such as mapof.Any's take pointer receivers.
		childValue, restore := addressableItem(element.Wildcard, propertyValue)
		changedValue, childRewrites, err := validate(element.Wildcard, childValue)

		if err != nil {
			return nil, derp.Wrap(err, location, "Validating property", key)
		}

		if len(childRewrites) == 0 {
			continue
		}

		// Write the changed value back to the same key, and record the rewrites under it
		err = setPropertyKey(element, element.Wildcard, value, key, restore(changedValue))

		if err != nil {
			return nil, derp.Wrap(err, location, "Setting property", key)
		}

		rewrites = append(rewrites, childRewrites.prefix(key)...)
	}

	return rewrites, nil
}

// wildcardKeys returns the keys of a declared map that no named property describes, sorted so
// that rewrites are reported in a stable order.  It returns nothing when there is no Wildcard.
func wildcardKeys(element Object, value any) []string {

	// RULE: Wildcard element is required
	if isNil(element.Wildcard) {
		return nil
	}

	// RULE: Only a value that declares itself a map has keys beyond its named properties
	mapTyper, ok := value.(MapTyper)

	if !ok {
		return nil
	}

	if !mapTyper.IsMap() {
		return nil
	}

	// RULE: The map must be able to list its keys
	keysGetter, ok := value.(KeysGetter)

	if !ok {
		return nil
	}

	// Named properties were already validated by their own elements.  Keys may return the map's
	// own slice, so filter a copy.
	keys := slices.DeleteFunc(slices.Clone(keysGetter.Keys()), func(key string) bool {
		_, isNamed := element.Properties[key]
		return isNamed
	})

	slices.Sort(keys)
	return keys
}
