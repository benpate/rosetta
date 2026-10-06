package mapof

import (
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/benpate/rosetta/loose"
	"github.com/stretchr/testify/require"
)

// templateData is a struct for templates to render against, so that a missing field is an
// error rather than "<no value>"
type templateData struct {
	ID string
}

// newTestTemplate parses a map holding one template beside plain values of every common type
func newTestTemplate(t *testing.T) Template {

	t.Helper()

	return ParseTemplate(map[string]any{
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

// TestParseTemplate requires that templates are compiled, and every other value, including
// every other string, stored as it is
func TestParseTemplate(t *testing.T) {

	value := newTestTemplate(t)

	require.IsType(t, loose.Template{}, value["validator"])
	require.True(t, value["validator"].(loose.Template).IsTemplate())
	require.Equal(t, "off", value["plain"])
	require.Equal(t, 6, value["rows"])
	require.Equal(t, []any{"a", "b"}, value["list"])

	empty := ParseTemplate(nil)
	require.NotNil(t, empty)
	require.True(t, empty.IsEmpty())
}

// TestParseTemplate_Malformed requires that a malformed template is kept as a plain string
func TestParseTemplate_Malformed(t *testing.T) {

	result := ParseTemplate(map[string]any{"validator": "{{.ID | nosuchfunc}}"})
	require.Equal(t, "{{.ID | nosuchfunc}}", result["validator"])

	rendered, err := result.Evaluate("validator", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "{{.ID | nosuchfunc}}", rendered)
}

// TestTemplate_Evaluate covers a compiled template, a Go-literal template, and plain values
func TestTemplate_Evaluate(t *testing.T) {

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
	literal := Template{"validator": "/v?id={{.ID}}", "plain": "off"}

	result, err = literal.Evaluate("validator", data)
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)
	require.Equal(t, "/v?id={{.ID}}", literal["validator"])

	result, err = literal.Evaluate("plain", data)
	require.NoError(t, err)
	require.Equal(t, "off", result)
}

// TestTemplate_Evaluate_Errors requires that a failed render returns an error, parsed or
// literal, and that a malformed Go-literal template renders as itself
func TestTemplate_Evaluate_Errors(t *testing.T) {

	parsed := ParseTemplate(map[string]any{"broken": "{{.Missing}}"})

	result, err := parsed.Evaluate("broken", templateData{ID: "123"})
	require.Error(t, err)
	require.Equal(t, "", result)

	literal := Template{"malformed": "{{.ID | nosuchfunc}}", "broken": "{{.Missing}}"}

	result, err = literal.Evaluate("malformed", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "{{.ID | nosuchfunc}}", result)

	_, err = literal.Evaluate("broken", templateData{ID: "123"})
	require.Error(t, err)
}

// TestTemplate_Getters_NilData requires that every getter reads the source value when it has
// no data to render against
func TestTemplate_Getters_NilData(t *testing.T) {

	value := newTestTemplate(t)

	require.Equal(t, "/v?id={{.ID}}", value.GetString("validator", nil))
	require.Equal(t, "off", value.GetString("plain", nil))
	require.Equal(t, "/v?id={{.ID}}", value.GetAny("validator", nil))
	require.Equal(t, 6, value.GetAny("rows", nil))
	require.Equal(t, 6, value.GetInt("rows", nil))
	require.Equal(t, int64(6), value.GetInt64("rows", nil))
	require.Equal(t, 1.5, value.GetFloat("price", nil))
	require.True(t, value.GetBool("required", nil))
	require.Equal(t, time.Date(2026, 9, 26, 0, 0, 0, 0, time.UTC), value.GetTime("when", nil).UTC())
	require.Equal(t, []string{"a", "b"}, value.GetSliceOfString("list", nil))
	require.Equal(t, []any{"a", "b"}, value.GetSliceOfAny("list", nil))
	require.Equal(t, []int{6}, value.GetSliceOfInt("rows", nil))
	require.Equal(t, []float64{1.5}, value.GetSliceOfFloat("price", nil))
	require.Equal(t, Any{"key": "value"}, value.GetMap("nested", nil))
	require.Equal(t, Any{"key": "value"}, value.GetMapOfAny("nested", nil))
	require.Equal(t, NewString(), value.GetMapOfString("nested", nil), "as mapof.Any does, a map of any is not a map of strings")
	require.Equal(t, String{"key": "value"}, Template{"strings": map[string]string{"key": "value"}}.GetMapOfString("strings", nil))

	// The map-shaped getters match mapof.Any on the same original values
	original := Any(value.MapOfAny())
	for _, key := range []string{"nested", "list", "validator", "missing"} {
		require.Equal(t, original.GetSliceOfMap(key), value.GetSliceOfMap(key, nil), key)
		require.Equal(t, original.GetSliceOfPlainMap(key), value.GetSliceOfPlainMap(key, nil), key)
		require.Equal(t, original.GetMapOfAny(key), value.GetMapOfAny(key, nil), key)
		require.Equal(t, original.GetMapOfString(key), value.GetMapOfString(key, nil), key)
	}

	// The source of a template does not coerce to a number
	require.Equal(t, 0, value.GetInt("validator", nil))

	// Every OK variant reports absence
	_, ok := value.GetAnyOK("missing", nil)
	require.False(t, ok)
	_, ok = value.GetBoolOK("missing", nil)
	require.False(t, ok)
	_, ok = value.GetFloatOK("missing", nil)
	require.False(t, ok)
	_, ok = value.GetIntOK("missing", nil)
	require.False(t, ok)
	_, ok = value.GetInt64OK("missing", nil)
	require.False(t, ok)
	_, ok = value.GetStringOK("missing", nil)
	require.False(t, ok)
	_, ok = value.GetTimeOK("missing", nil)
	require.False(t, ok)

	pointer, ok := value.GetPointer("validator")
	require.True(t, ok)
	require.Equal(t, "/v?id={{.ID}}", pointer)
}

// renderData is a struct for getters to render templates against, holding a value of every type
type renderData struct {
	ID    string
	On    bool
	Rows  int
	Price float64
	When  string
}

// TestTemplate_Getters_Render requires that every getter renders a template against the data,
// and converts the result to its own type
func TestTemplate_Getters_Render(t *testing.T) {

	value := ParseTemplate(map[string]any{
		"url":     "/v?id={{.ID}}",
		"on":      "{{if .On}}true{{else}}false{{end}}",
		"rows":    "{{.Rows}}",
		"price":   "{{.Price}}",
		"when":    "{{.When}}",
		"padded":  "  {{.Rows}}\n",
		"plain":   "off",
		"number":  6,
		"list":    []any{"{{.ID}}"},
		"missing": "{{.Nope}}",
	})

	data := renderData{ID: "123", On: true, Rows: 6, Price: 1.5, When: "2026-09-26T00:00:00Z"}

	require.Equal(t, "/v?id=123", value.GetString("url", data))
	require.Equal(t, "/v?id=123", value.GetAny("url", data))
	require.True(t, value.GetBool("on", data))
	require.False(t, value.GetBool("on", renderData{}))
	require.Equal(t, 6, value.GetInt("rows", data))
	require.Equal(t, int64(6), value.GetInt64("rows", data))
	require.Equal(t, 1.5, value.GetFloat("price", data))
	require.Equal(t, time.Date(2026, 9, 26, 0, 0, 0, 0, time.UTC), value.GetTime("when", data).UTC())
	require.Equal(t, []string{"/v?id=123"}, value.GetSliceOfString("url", data))
	require.Equal(t, []int{6}, value.GetSliceOfInt("rows", data))
	require.Equal(t, []float64{1.5}, value.GetSliceOfFloat("price", data))
	require.Equal(t, []any{"6"}, value.GetSliceOfAny("rows", data))

	// One template converts to every type it can, because the stored value never changes
	require.Equal(t, "6", value.GetString("rows", data))
	require.Equal(t, 6, value.GetInt("rows", data))
	require.Equal(t, 6.0, value.GetFloat("rows", data))
	require.Equal(t, "{{.Rows}}", value.GetString("rows", nil))

	// Scalar getters trim whitespace from a rendered template; GetString keeps it
	require.Equal(t, 6, value.GetInt("padded", data))
	require.Equal(t, "  6\n", value.GetString("padded", data))

	// Plain strings and other values are returned as they are
	require.Equal(t, "off", value.GetString("plain", data))
	require.Equal(t, 6, value.GetAny("number", data))
	require.Equal(t, 6, value.GetInt("number", data))

	// Templates are top-level only: strings inside a list are not rendered
	require.Equal(t, []string{"{{.ID}}"}, value.GetSliceOfString("list", data))

	// A failed render is swallowed: the value is present, but renders as ""
	result, ok := value.GetStringOK("missing", data)
	require.True(t, ok)
	require.Equal(t, "", result)
	_, ok = value.GetBoolOK("missing", data)
	require.False(t, ok)
	_, ok = value.GetIntOK("missing", data)
	require.False(t, ok)
}

// TestTemplate_Getters_Literal requires that getters render a template in a Go map literal,
// which is never parsed, without changing the map
func TestTemplate_Getters_Literal(t *testing.T) {

	literal := Template{"url": "/literal/{{.ID}}", "rows": "{{.Rows}}"}
	data := renderData{ID: "123", Rows: 6}

	require.Equal(t, "/literal/123", literal.GetString("url", data))
	require.Equal(t, 6, literal.GetInt("rows", data))
	require.Equal(t, "/literal/{{.ID}}", literal.GetString("url", nil))
	require.Equal(t, Template{"url": "/literal/{{.ID}}", "rows": "{{.Rows}}"}, literal)
}

// TestTemplate_MapManipulations covers the whole-map methods
func TestTemplate_MapManipulations(t *testing.T) {

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

	var empty Template
	require.Nil(t, empty.MapOfAny())
	require.True(t, NewTemplate().IsEmpty())
}

// TestTemplate_Setters requires that setters store values as mapof.Any does, and compile strings
func TestTemplate_Setters(t *testing.T) {

	var value Template

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

	var nilMap Template
	require.True(t, nilMap.Remove("missing"))
	require.NotNil(t, nilMap)
}

// TestTemplate_SetValue requires that a whole map is replaced and its strings compiled
func TestTemplate_SetValue(t *testing.T) {

	var value Template

	require.NoError(t, value.SetValue(Any{"validator": "/v?id={{.ID}}", "rows": 6}))
	require.True(t, value["validator"].(loose.Template).IsTemplate())
	require.Equal(t, 6, value["rows"])

	require.NoError(t, value.SetValue(newTestTemplate(t)))
	require.Equal(t, "/v?id={{.ID}}", value.GetString("validator", nil))

	require.Error(t, value.SetValue("not a map"))

	// mapof.Any accepts a Template too, as its original values
	var plain Any
	require.NoError(t, plain.SetValue(newTestTemplate(t)))
	require.Equal(t, "/v?id={{.ID}}", plain["validator"])
	pointer := newTestTemplate(t)
	converted, ok := MapOfAny(&pointer)
	require.True(t, ok)
	require.Equal(t, "off", converted["plain"])
}

// TestTemplate_Append requires that appending forms a slice of original values
func TestTemplate_Append(t *testing.T) {

	value := newTestTemplate(t)

	value.Append("validator", "second")
	require.Equal(t, []any{"/v?id={{.ID}}", "second"}, value["validator"])

	value.Append("validator", "third")
	require.Equal(t, []any{"/v?id={{.ID}}", "second", "third"}, value["validator"])

	var empty Template
	empty.Append("new", 1)
	require.Equal(t, []any{1}, empty["new"])
}

// TestTemplate_SetKey requires that SetKey compiles a string that holds a template, and stores
// every other value as given, including a zero value
func TestTemplate_SetKey(t *testing.T) {

	var value Template

	require.NoError(t, value.SetKey("validator", "/v?id={{.ID}}"))
	require.True(t, value["validator"].(loose.Template).IsTemplate())

	require.NoError(t, value.SetKey("child", Any{"name": "Sarah"}))
	require.Equal(t, Any{"name": "Sarah"}, value["child"])

	require.NoError(t, value.SetKey("empty", ""))
	require.Equal(t, "", value["empty"])

	require.NoError(t, value.SetKey("validator", "{{.ID | nosuchfunc}}"))
	require.Equal(t, "{{.ID | nosuchfunc}}", value["validator"])
}

// TestTemplate_JSON requires that a Template round-trips through JSON as its original values,
// and compiles its strings while decoding
func TestTemplate_JSON(t *testing.T) {

	source := `{"validator":"/v?id={{.ID}}","plain":"off","rows":6,"required":true,"list":["a"]}`

	var value Template
	require.NoError(t, json.Unmarshal([]byte(source), &value))
	require.True(t, value["validator"].(loose.Template).IsTemplate())

	result, err := value.Evaluate("validator", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)

	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, source, string(encoded))
}

// TestTemplate_JSON_Invalid requires that malformed JSON fails, and a malformed template
// decodes as a plain string
func TestTemplate_JSON_Invalid(t *testing.T) {

	var value Template
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
func newCloneTestTemplate(t *testing.T) Template {

	t.Helper()

	return ParseTemplate(map[string]any{
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

// TestTemplate_Clone requires that a clone equals its original, and that changing it at any
// depth leaves the original unchanged
func TestTemplate_Clone(t *testing.T) {

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

// TestTemplate_Clone_Types requires that every value keeps its type, nil entries stay present,
// and nil containers stay nil
func TestTemplate_Clone_Types(t *testing.T) {

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

// TestTemplate_Clone_SharesTemplates requires that a compiled template is shared, and still renders
func TestTemplate_Clone_SharesTemplates(t *testing.T) {

	clone := newCloneTestTemplate(t).Clone()

	require.True(t, clone["validator"].(loose.Template).IsTemplate())

	result, err := clone.Evaluate("validator", templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)
}

// TestTemplate_Clone_Empty requires that nil and empty maps clone as they are
func TestTemplate_Clone_Empty(t *testing.T) {

	var nilMap Template
	require.Nil(t, nilMap.Clone())

	empty := NewTemplate()
	clone := empty.Clone()
	require.NotNil(t, clone)
	clone["key"] = "value"
	require.Empty(t, empty)
}

// TestTemplate_Clone_Pointer pins that a pointer is shared, because Clone copies only maps and slices
func TestTemplate_Clone_Pointer(t *testing.T) {

	shared := &[]string{"original"}
	clone := Template{"pointer": shared}.Clone()
	require.Same(t, shared, clone["pointer"])
}

// TestTemplate_Evaluate_Repeated requires that evaluating a Go literal repeatedly renders the
// same result, and leaves the map itself unchanged
func TestTemplate_Evaluate_Repeated(t *testing.T) {

	literal := Template{"url": "/cache/evaluate/{{.ID}}"}

	for range 3 {
		result, err := literal.Evaluate("url", templateData{ID: "123"})
		require.NoError(t, err)
		require.Equal(t, "/cache/evaluate/123", result)
	}

	require.Equal(t, Template{"url": "/cache/evaluate/{{.ID}}"}, literal)
}

// TestTemplate_Evaluate_Concurrent evaluates one shared Go-literal map from many goroutines,
// beside readers of the same map, which the race detector reports if Evaluate writes to it
func TestTemplate_Evaluate_Concurrent(t *testing.T) {

	shared := Template{
		"url":   "/cache/concurrent/{{.ID}}",
		"class": "wide",
	}

	results := make(chan bool, 150)

	var wait sync.WaitGroup
	for index := range 50 {
		id := strconv.Itoa(index)

		wait.Go(func() {
			result, err := shared.Evaluate("url", templateData{ID: id})
			results <- (err == nil) && (result == "/cache/concurrent/"+id)
		})

		wait.Go(func() {
			results <- shared.GetString("class", nil) == "wide"
		})

		wait.Go(func() {
			results <- shared.Clone().GetString("url", nil) == "/cache/concurrent/{{.ID}}"
		})
	}

	wait.Wait()
	close(results)

	for ok := range results {
		require.True(t, ok)
	}
}
