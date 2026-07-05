package schema

import (
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

// IsRegisteredFormat returns TRUE if a string format generator is registered under
// the given name. Like every registry read, the first call freezes the registry
// against further UseFormat calls, so it must not be called until all format
// registrations are complete.
func IsRegisteredFormat(name string) bool {
	_, exists := lookupFormat(name)
	return exists
}

// ValidateFormats confirms that every format name used by a String element in this
// schema resolves in the format registry. String validation silently skips names it
// does not recognize, so a typo'd format otherwise degrades without error; calling
// this from a test (or at startup, after all UseFormat registrations) makes that
// failure loud instead.
func (schema Schema) ValidateFormats() error {
	return ValidateElementFormats(schema.Element, "")
}

// ValidateElementFormats recursively checks the format names on an element and all
// of its children. The path argument locates errors within nested schemas; pass ""
// for the root element.
func ValidateElementFormats(element Element, path string) error {

	const location = "schema.ValidateElementFormats"

	switch typed := unwrapElement(element).(type) {

	case nil:
		return nil

	case String:
		for _, name := range formatNames(typed.Format) {
			if !IsRegisteredFormat(name) {
				return derp.Internal(location, unrecognizedFormatMessage(name, path), name, path)
			}
		}
		return nil

	case Array:
		return ValidateElementFormats(typed.Items, joinPath(path, "items"))

	case Object:

		for name, property := range typed.Properties {
			if err := ValidateElementFormats(property, joinPath(path, name)); err != nil {
				return err
			}
		}

		return ValidateElementFormats(typed.Wildcard, joinPath(path, "*"))
	}

	// All other element types (Any, Boolean, Integer, Number) carry no format names.
	return nil
}

// unwrapElement dereferences pointer elements (returning nil for nil pointers) so
// that ValidateElementFormats only has to handle value types.
func unwrapElement(element Element) Element {

	switch typed := element.(type) {

	case *String:
		if typed == nil {
			return nil
		}
		return *typed

	case *Array:
		if typed == nil {
			return nil
		}
		return *typed

	case *Object:
		if typed == nil {
			return nil
		}
		return *typed
	}

	return element
}

// formatNames extracts the format names from a String element's space-separated
// format definition, using the same parsing as formatFunctions. Empty tokens are
// dropped because the runtime skips them harmlessly (an all-empty definition just
// selects the no-html default).
func formatNames(definition string) []string {

	result := make([]string, 0)

	for _, arg := range strings.Split(definition, " ") {

		name, _ := list.Equal(arg).Split()

		if name == "" {
			continue
		}

		result = append(result, name)
	}

	return result
}

// unrecognizedFormatMessage names the offending format -- and its location within the
// schema, when it is not the root element -- so load-time error logs are actionable.
func unrecognizedFormatMessage(name string, path string) string {

	if path == "" {
		return "Unrecognized format name: " + name
	}

	return "Unrecognized format name: " + name + " (at " + path + ")"
}

// joinPath appends a child name to a dot-separated element path.
func joinPath(path string, name string) string {

	if path == "" {
		return name
	}

	return path + "." + name
}
