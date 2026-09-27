package mapof

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/benpate/rosetta/loose"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

// templateData is a struct for templates to render against, so that a missing field is an
// error rather than "<no value>"
type templateData struct {
	ID string
}

// newTestTemplate parses a map holding one template beside plain values of every common type
func newTestTemplate(t *testing.T) LooseTemplate {

	t.Helper()

	return ParseLooseTemplate(map[string]any{
		"validator": "/v?id={{.ID}}",
		"plain":     "off",
		"rows":      6,
		"price":     1.5,
		"required":  true,
		"list":      []any{"a", "b"},
		"nested":    map[string]any{"key": "value"},
		"when":      "2026-09-26T00:00:00Z",
	})
}

// TestParseLooseTemplate requires that templates are compiled, and every other value, including
// every other string, stored as it is
func TestParseLooseTemplate(t *testing.T) {

	value := newTestTemplate(t)

	require.IsType(t, loose.Template{}, value["validator"])
	require.True(t, value["validator"].(loose.Template).IsTemplate())
	require.Equal(t, "off", value["plain"])
	require.Equal(t, 6, value["rows"])
	require.Equal(t, []any{"a", "b"}, value["list"])

	empty := ParseLooseTemplate(nil)
	require.NotNil(t, empty)
	require.True(t, empty.IsEmpty())
}

// TestParseLooseTemplate_Malformed requires that a malformed template is kept as a plain string
func TestParseLooseTemplate_Malformed(t *testing.T) {

	result := ParseLooseTemplate(map[string]any{"validator": "{{.ID | nosuchfunc}}"})
	require.Equal(t, "{{.ID | nosuchfunc}}", result["validator"])

	rendered, err := result.Evaluate("validator", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "{{.ID | nosuchfunc}}", rendered)
}

// TestLooseTemplate_Evaluate covers a compiled template, a Go-literal template, and plain values
func TestLooseTemplate_Evaluate(t *testing.T) {

	value := newTestTemplate(t)
	data := templateData{ID: "123"}

	result, err := value.Evaluate("validator", data)
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)

	result, err = value.Evaluate("plain", data)
	require.NoError(t, err)
	require.Equal(t, "off", result)

	result, err = value.Evaluate("rows", data)
	require.NoError(t, err)
	require.Equal(t, "6", result)

	result, err = value.Evaluate("missing", data)
	require.NoError(t, err)
	require.Equal(t, "", result)

	// A map literal written in Go is compiled when it is evaluated
	literal := LooseTemplate{"validator": "/v?id={{.ID}}", "plain": "off"}

	result, err = literal.Evaluate("validator", data)
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)
	require.Equal(t, "/v?id={{.ID}}", literal["validator"])

	result, err = literal.Evaluate("plain", data)
	require.NoError(t, err)
	require.Equal(t, "off", result)
}

// TestLooseTemplate_Evaluate_Errors requires that a failed render returns an error, parsed or
// literal, and that a malformed Go-literal template renders as itself
func TestLooseTemplate_Evaluate_Errors(t *testing.T) {

	parsed := ParseLooseTemplate(map[string]any{"broken": "{{.Missing}}"})

	result, err := parsed.Evaluate("broken", templateData{ID: "123"})
	require.Error(t, err)
	require.Equal(t, "", result)

	literal := LooseTemplate{"malformed": "{{.ID | nosuchfunc}}", "broken": "{{.Missing}}"}

	result, err = literal.Evaluate("malformed", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "{{.ID | nosuchfunc}}", result)

	_, err = literal.Evaluate("broken", templateData{ID: "123"})
	require.Error(t, err)
}

// TestLooseTemplate_Getters requires that every getter reads the original value, and never renders
func TestLooseTemplate_Getters(t *testing.T) {

	value := newTestTemplate(t)

	require.Equal(t, "/v?id={{.ID}}", value.GetString("validator"))
	require.Equal(t, "off", value.GetString("plain"))
	require.Equal(t, "/v?id={{.ID}}", value.GetAny("validator"))
	require.Equal(t, 6, value.GetAny("rows"))
	require.Equal(t, 6, value.GetInt("rows"))
	require.Equal(t, int64(6), value.GetInt64("rows"))
	require.Equal(t, 1.5, value.GetFloat("price"))
	require.True(t, value.GetBool("required"))
	require.Equal(t, time.Date(2026, 9, 26, 0, 0, 0, 0, time.UTC), value.GetTime("when").UTC())
	require.Equal(t, []string{"a", "b"}, value.GetSliceOfString("list"))
	require.Equal(t, []any{"a", "b"}, value.GetSliceOfAny("list"))
	require.Equal(t, []int{6}, value.GetSliceOfInt("rows"))
	require.Equal(t, []float64{1.5}, value.GetSliceOfFloat("price"))
	require.Equal(t, Any{"key": "value"}, value.GetMap("nested"))
	require.Equal(t, Any{"key": "value"}, value.GetMapOfAny("nested"))
	require.Equal(t, NewString(), value.GetMapOfString("nested"), "as mapof.Any does, a map of any is not a map of strings")
	require.Equal(t, String{"key": "value"}, LooseTemplate{"strings": map[string]string{"key": "value"}}.GetMapOfString("strings"))

	// The map-shaped getters match mapof.Any on the same original values
	original := Any(value.MapOfAny())
	for _, key := range []string{"nested", "list", "validator", "missing"} {
		require.Equal(t, original.GetSliceOfMap(key), value.GetSliceOfMap(key), key)
		require.Equal(t, original.GetSliceOfPlainMap(key), value.GetSliceOfPlainMap(key), key)
		require.Equal(t, original.GetMapOfAny(key), value.GetMapOfAny(key), key)
		require.Equal(t, original.GetMapOfString(key), value.GetMapOfString(key), key)
	}

	// A string that is a template does not coerce to a number
	require.Equal(t, 0, value.GetInt("validator"))

	// Every OK variant reports absence
	_, ok := value.GetAnyOK("missing")
	require.False(t, ok)
	_, ok = value.GetBoolOK("missing")
	require.False(t, ok)
	_, ok = value.GetFloatOK("missing")
	require.False(t, ok)
	_, ok = value.GetIntOK("missing")
	require.False(t, ok)
	_, ok = value.GetInt64OK("missing")
	require.False(t, ok)
	_, ok = value.GetStringOK("missing")
	require.False(t, ok)
	_, ok = value.GetTimeOK("missing")
	require.False(t, ok)

	pointer, ok := value.GetPointer("validator")
	require.True(t, ok)
	require.Equal(t, "/v?id={{.ID}}", pointer)
}

// TestLooseTemplate_MapManipulations covers the whole-map methods
func TestLooseTemplate_MapManipulations(t *testing.T) {

	value := newTestTemplate(t)

	require.Equal(t, 8, value.Length())
	require.Equal(t, []string{"list", "nested", "plain", "price", "required", "rows", "validator", "when"}, value.Keys())
	require.True(t, value.IsMap())
	require.True(t, value.NotEmpty())
	require.False(t, value.IsEmpty())
	require.True(t, value.IsZeroValue("missing"))
	require.False(t, value.IsZeroValue("validator"))

	// Equality compares original values, so templates compare by their source
	require.True(t, value.Equal(value.MapOfAny()))
	require.False(t, value.NotEqual(value.MapOfAny()))
	require.Equal(t, "/v?id={{.ID}}", value.MapOfAny()["validator"])
	require.Equal(t, "/v?id={{.ID}}", value.MapOfString()["validator"])
	require.Equal(t, "6", value.MapOfString()["rows"])

	var empty LooseTemplate
	require.Nil(t, empty.MapOfAny())
	require.True(t, NewLooseTemplate().IsEmpty())
}

// TestLooseTemplate_Setters requires that setters store values as mapof.Any does, and compile strings
func TestLooseTemplate_Setters(t *testing.T) {

	var value LooseTemplate

	require.True(t, value.SetString("validator", "/v?id={{.ID}}"))
	require.True(t, value["validator"].(loose.Template).IsTemplate())
	require.True(t, value.SetString("plain", "off"))
	require.Equal(t, "off", value["plain"])
	require.True(t, value.SetAny("any", "/a/{{.ID}}"))
	require.True(t, value["any"].(loose.Template).IsTemplate())
	require.True(t, value.SetAny("rows", 6))
	require.Equal(t, 6, value["rows"])
	require.True(t, value.SetBool("required", true))
	require.True(t, value.SetFloat("price", 1.5))
	require.True(t, value.SetInt("count", 3))
	require.True(t, value.SetInt64("big", 4))

	// A malformed template is stored as a plain string
	require.True(t, value.SetString("bad", "{{.ID | nosuchfunc}}"))
	require.Equal(t, "{{.ID | nosuchfunc}}", value["bad"])
	require.True(t, value.SetAny("bad", ""))

	// Zero values delete their keys, as they do in mapof.Any
	require.True(t, value.SetString("plain", ""))
	require.True(t, value.SetAny("rows", 0))
	require.True(t, value.SetFloat("price", 0))
	require.True(t, value.SetInt("count", 0))
	require.True(t, value.SetInt64("big", 0))
	require.Equal(t, []string{"any", "required", "validator"}, value.Keys())

	require.True(t, value.Remove("any"))
	require.Equal(t, []string{"required", "validator"}, value.Keys())

	var nilMap LooseTemplate
	require.True(t, nilMap.Remove("missing"))
	require.NotNil(t, nilMap)
}

// TestLooseTemplate_SetValue requires that a whole map is replaced and its strings compiled
func TestLooseTemplate_SetValue(t *testing.T) {

	var value LooseTemplate

	require.NoError(t, value.SetValue(Any{"validator": "/v?id={{.ID}}", "rows": 6}))
	require.True(t, value["validator"].(loose.Template).IsTemplate())
	require.Equal(t, 6, value["rows"])

	require.NoError(t, value.SetValue(newTestTemplate(t)))
	require.Equal(t, "/v?id={{.ID}}", value.GetString("validator"))

	require.Error(t, value.SetValue("not a map"))

	// mapof.Any accepts a LooseTemplate too, as its original values
	var plain Any
	require.NoError(t, plain.SetValue(newTestTemplate(t)))
	require.Equal(t, "/v?id={{.ID}}", plain["validator"])
	pointer := newTestTemplate(t)
	converted, ok := MapOfAny(&pointer)
	require.True(t, ok)
	require.Equal(t, "off", converted["plain"])
}

// TestLooseTemplate_Append requires that appending forms a slice of original values
func TestLooseTemplate_Append(t *testing.T) {

	value := newTestTemplate(t)

	value.Append("validator", "second")
	require.Equal(t, []any{"/v?id={{.ID}}", "second"}, value["validator"])

	value.Append("validator", "third")
	require.Equal(t, []any{"/v?id={{.ID}}", "second", "third"}, value["validator"])

	var empty LooseTemplate
	empty.Append("new", 1)
	require.Equal(t, []any{1}, empty["new"])
}

// TestLooseTemplate_SetObject covers setting a value directly and through a child map
func TestLooseTemplate_SetObject(t *testing.T) {

	element := schema.Object{Properties: schema.ElementMap{
		"validator": schema.String{},
		"child":     schema.Object{Properties: schema.ElementMap{"name": schema.String{}, "title": schema.String{}}},
	}}

	var value LooseTemplate

	require.NoError(t, value.SetObject(element, list.ByDot("validator"), "/v?id={{.ID}}"))
	require.True(t, value["validator"].(loose.Template).IsTemplate())

	require.NoError(t, value.SetObject(element, list.ByDot("child.name"), "Sarah"))
	require.Equal(t, Any{"name": "Sarah"}, value["child"])

	require.NoError(t, value.SetObject(element, list.ByDot("child.title"), "Captain"))
	require.Equal(t, Any{"name": "Sarah", "title": "Captain"}, value["child"], "an existing child map is reused")
	require.Error(t, value.SetObject(element, list.ByDot(""), "x"))
	require.Error(t, value.SetObject(element, list.ByDot("missing.name"), "x"))
	require.Error(t, value.SetObject(element, list.ByDot("child.nope"), "x"), "the child schema has no such property")
	require.NoError(t, value.SetObject(element, list.ByDot("validator"), "{{.ID | nosuchfunc}}"))
	require.Equal(t, "{{.ID | nosuchfunc}}", value["validator"])
}

// TestLooseTemplate_JSON requires that a LooseTemplate round-trips through JSON as its original values,
// and compiles its strings while decoding
func TestLooseTemplate_JSON(t *testing.T) {

	source := `{"validator":"/v?id={{.ID}}","plain":"off","rows":6,"required":true,"list":["a"]}`

	var value LooseTemplate
	require.NoError(t, json.Unmarshal([]byte(source), &value))
	require.True(t, value["validator"].(loose.Template).IsTemplate())

	result, err := value.Evaluate("validator", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)

	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, source, string(encoded))
}

// TestLooseTemplate_JSON_Invalid requires that malformed JSON fails, and a malformed template
// decodes as a plain string
func TestLooseTemplate_JSON_Invalid(t *testing.T) {

	var value LooseTemplate
	require.Error(t, value.UnmarshalJSON([]byte(`{"a":`)))
	require.Nil(t, value)

	require.NoError(t, json.Unmarshal([]byte(`{"a":"{{.ID | nosuchfunc}}"}`), &value))
	require.Equal(t, "{{.ID | nosuchfunc}}", value["a"])
}

// cloneTestItem is a struct stored inside a slice option, as form.LookupCode is
type cloneTestItem struct {
	Label string
}

// newCloneTestTemplate returns a map holding every kind of value a form option can hold
func newCloneTestTemplate(t *testing.T) LooseTemplate {

	t.Helper()

	return ParseLooseTemplate(map[string]any{
		"validator": "/v?id={{.ID}}",
		"rows":      6,
		"rules":     map[string]any{"width": 800, "types": []any{"webp"}, "empty": nil},
		"enum":      []cloneTestItem{{Label: "one"}, {Label: "two"}},
		"list":      []any{"a", map[string]any{"key": "value"}, nil},
		"strings":   String{"key": "value"},
		"nilMap":    map[string]any(nil),
		"nilSlice":  []string(nil),
	})
}

// TestLooseTemplate_Clone requires that a clone equals its original, and that changing it at any
// depth leaves the original unchanged
func TestLooseTemplate_Clone(t *testing.T) {

	original := newCloneTestTemplate(t)
	clone := original.Clone()
	require.Equal(t, newCloneTestTemplate(t).MapOfAny(), clone.MapOfAny())

	clone["rows"] = 7
	clone["added"] = true
	clone["rules"].(map[string]any)["width"] = 1
	clone["rules"].(map[string]any)["types"].([]any)[0] = "png"
	clone["enum"].([]cloneTestItem)[0].Label = "changed"
	clone["list"].([]any)[1].(map[string]any)["key"] = "changed"
	clone["strings"].(String)["key"] = "changed"

	require.Equal(t, newCloneTestTemplate(t).MapOfAny(), original.MapOfAny())
}

// TestLooseTemplate_Clone_Types requires that every value keeps its type, nil entries stay present,
// and nil containers stay nil
func TestLooseTemplate_Clone_Types(t *testing.T) {

	clone := newCloneTestTemplate(t).Clone()

	require.IsType(t, []cloneTestItem{}, clone["enum"])
	require.IsType(t, String{}, clone["strings"])

	value, exists := clone["rules"].(map[string]any)["empty"]
	require.True(t, exists)
	require.Nil(t, value)
	require.Nil(t, clone["list"].([]any)[2])

	require.Nil(t, clone["nilMap"].(map[string]any))
	require.Nil(t, clone["nilSlice"].([]string))
}

// TestLooseTemplate_Clone_SharesTemplates requires that a compiled template is shared, and still renders
func TestLooseTemplate_Clone_SharesTemplates(t *testing.T) {

	clone := newCloneTestTemplate(t).Clone()

	require.True(t, clone["validator"].(loose.Template).IsTemplate())

	result, err := clone.Evaluate("validator", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)
}

// TestLooseTemplate_Clone_Empty requires that nil and empty maps clone as they are
func TestLooseTemplate_Clone_Empty(t *testing.T) {

	var nilMap LooseTemplate
	require.Nil(t, nilMap.Clone())

	empty := NewLooseTemplate()
	clone := empty.Clone()
	require.NotNil(t, clone)
	clone["key"] = "value"
	require.Empty(t, empty)
}

// TestLooseTemplate_Clone_Pointer pins that a pointer is shared, because Clone copies only maps and slices
func TestLooseTemplate_Clone_Pointer(t *testing.T) {

	shared := &[]string{"original"}
	clone := LooseTemplate{"pointer": shared}.Clone()
	require.Same(t, shared, clone["pointer"])
}
