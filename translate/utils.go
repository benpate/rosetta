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
	var buffer bytes.Buffer

	if err := t.Execute(&buffer, data); err != nil {
		derp.Report(derp.Wrap(err, "rosetta.translate.executeTemplate", "Error executing template", data))
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
