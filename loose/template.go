package loose

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/benpate/derp"
)

// Template is a string that may be a text/template, compiled when it is created
type Template struct {
	source   string
	template *template.Template
}

// NewTemplate wraps a string, compiling it as a template if it holds a "{{" followed later by
// a "}}". A malformed template is kept as a plain string.
func NewTemplate(source string) Template {

	result := Template{source: source}

	if looksLikeTemplate(source) {
		if compiled, err := template.New("").Parse(source); err == nil {
			result.template = compiled
		}
	}

	return result
}

// IsTemplate returns TRUE if the string was compiled as a template
func (t Template) IsTemplate() bool {
	return t.template != nil
}

// Execute renders a template against the provided data, or returns a plain string as it is
func (t Template) Execute(data any) (string, error) {

	const location = "loose.Template.Execute"

	if t.template == nil {
		return t.source, nil
	}

	var result strings.Builder

	if err := t.template.Execute(&result, data); err != nil {
		return "", derp.Wrap(err, location, "Unable to execute template", t.source)
	}

	return result.String(), nil
}

// String returns the original string, without rendering it
func (t Template) String() string {
	return t.source
}

// MarshalJSON encodes the original string
func (t Template) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.source)
}

// UnmarshalJSON decodes any scalar JSON value as a String, then compiles it if it holds a template
func (t *Template) UnmarshalJSON(data []byte) error {

	const location = "loose.Template.UnmarshalJSON"

	var source String

	if err := source.UnmarshalJSON(data); err != nil {
		return derp.Wrap(err, location, "Invalid template value")
	}

	*t = NewTemplate(string(source))

	// Look ma, no hands.
	return nil
}

// looksLikeTemplate returns TRUE if a string holds a "{{" followed later by a "}}"
func looksLikeTemplate(value string) bool {

	_, after, found := strings.Cut(value, "{{")

	if !found {
		return false
	}

	return strings.Contains(after, "}}")
}
