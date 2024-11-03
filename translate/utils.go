package translate

import (
	"bytes"
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
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

// toPlainMap returns a mapof.Any as a map[string]any.
// It's defined here to make it easier to call in slice.Map operations
func toPlainMap(value mapof.Any) map[string]any {
	return value.MapOfAny()
}
