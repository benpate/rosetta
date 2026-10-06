package translate

import (
	"bytes"
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// executeTemplate is just some syntax sugar to execute a template and return the result as a string
func executeTemplate(t *template.Template, data any) string {

	const location = "rosetta.translate.executeTemplate"

	// A nil template means construction failed earlier (e.g. a malformed condition);
	// guard here so a reported-but-unrecoverable parse error can't panic at execution time.
	if t == nil {
		derp.Report(derp.Internal(location, "Template must not be nil"))
		return ""
	}

	var buffer bytes.Buffer

	if err := t.Execute(&buffer, data); err != nil {
		derp.Report(derp.Wrap(err, location, "Executing template", data))
		return ""
	}

	return buffer.String()
}

// upscale boosts standarg Golang types to their Rosetta equivalents
func upscale(value any) any {

	switch typed := value.(type) {

	case []any:
		return sliceof.Any(typed)

	case map[string]any:
		return mapof.Any(typed)

	default:
		return typed
	}
}

// copyValue returns a deep copy of a rule's constant value, so that no output object ever
// shares a map or a list with the rule, which every caller reuses
func copyValue(value any) any {

	switch typed := value.(type) {

	case mapof.Any:
		return mapof.Any(copyMap(typed))

	case map[string]any:
		return copyMap(typed)

	case sliceof.Any:
		return sliceof.Any(copySlice(typed))

	case []any:
		return copySlice(typed)
	}

	// Scalars are copied by value already
	return value
}

// copyMap returns a deep copy of a map
func copyMap(value map[string]any) map[string]any {

	result := make(map[string]any, len(value))

	for key, item := range value {
		result[key] = copyValue(item)
	}

	return result
}

// copySlice returns a deep copy of a slice
func copySlice(value []any) []any {

	result := make([]any, len(value))

	for index, item := range value {
		result[index] = copyValue(item)
	}

	return result
}
