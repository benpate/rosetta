package translate

import (
	"bytes"
	"text/template"

	"github.com/benpate/derp"
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
