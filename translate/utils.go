package translate

import (
	"bytes"
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
)

func executeTemplate(t *template.Template, data any) string {
	var buffer bytes.Buffer

	if err := t.Execute(&buffer, data); err != nil {
		derp.Report(err)
		return ""
	}

	return buffer.String()
}

func toPlainMap(value mapof.Any) map[string]any {
	return value.MapOfAny()
}

/* These functions are unused for now.  Hiding them so that the linter doesn't lose its mind.
func getArrayLength(value any) (int, bool) {

	switch typed := value.(type) {
	case []any:
		return len(typed), true
	case []mapof.Any:
		return len(typed), true
	case schema.LengthGetter:
		return typed.Length(), true
	default:
		return 0, false
	}
}

func getArrayIndex(value any, index int) (any, bool) {

	switch typed := value.(type) {
	case []any:
		return typed[index], true
	case []mapof.Any:
		return typed[index], true
	case schema.ArrayGetter:
		return typed.GetIndex(index)
	default:
		return nil, false
	}
}

func setArrayIndex(value any, index int, item any) bool {

	switch typed := value.(type) {
	case schema.ArraySetter:
		return typed.SetIndex(index, item)
	default:
		return false
	}
}
*/
