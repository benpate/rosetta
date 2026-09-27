package loose

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// templateData is a struct for templates to render against, so that a missing field is an
// error rather than "<no value>"
type templateData struct {
	ID string
}

// TestNewTemplate_Template requires that a templated string compiles, renders, and keeps its source
func TestNewTemplate_Template(t *testing.T) {

	value := NewTemplate("/v?id={{.ID}}")
	require.True(t, value.IsTemplate())
	require.Equal(t, "/v?id={{.ID}}", value.String())

	result, err := value.Execute(templateData{ID: "123"})
	require.NoError(t, err)
	require.Equal(t, "/v?id=123", result)

	// A map renders too
	result, err = value.Execute(map[string]string{"ID": "456"})
	require.NoError(t, err)
	require.Equal(t, "/v?id=456", result)
}

// TestNewTemplate_Plain requires that a string that is not a template, including a malformed
// one, is kept and executed as it is
func TestNewTemplate_Plain(t *testing.T) {

	for _, source := range []string{
		"off",
		"",
		"{{.ID",
		".ID}}",
		"}} {{",
		"{{.ID | nosuchfunc}}",
		"/v?id={{.ID",
		"{{end}}",
	} {
		value := NewTemplate(source)
		require.False(t, value.IsTemplate(), source)
		require.Equal(t, source, value.String())

		result, err := value.Execute(templateData{ID: "123"})
		require.NoError(t, err, source)
		require.Equal(t, source, result)
	}
}

// TestIsTemplate requires a "{{" followed later by a "}}"
func TestIsTemplate(t *testing.T) {
	require.True(t, looksLikeTemplate("{{.ID}}"))
	require.True(t, looksLikeTemplate("/a/{{.ID}}/b"))
	require.True(t, looksLikeTemplate("}} {{.ID}}"))
	require.False(t, looksLikeTemplate(""))
	require.False(t, looksLikeTemplate("/a/b"))
	require.False(t, looksLikeTemplate("{{.ID"))
	require.False(t, looksLikeTemplate(".ID}}"))
	require.False(t, looksLikeTemplate("}} {{"))
	require.False(t, looksLikeTemplate("{{}"))
}

// TestTemplate_Execute_Error requires that a template naming a missing field fails
func TestTemplate_Execute_Error(t *testing.T) {

	value := NewTemplate("/v?id={{.Missing}}")
	require.True(t, value.IsTemplate())

	result, err := value.Execute(templateData{ID: "123"})
	require.Error(t, err)
	require.Equal(t, "", result)
}

// TestTemplate_Zero requires that the zero Template executes as an empty string
func TestTemplate_Zero(t *testing.T) {

	var value Template
	require.False(t, value.IsTemplate())
	require.Equal(t, "", value.String())

	result, err := value.Execute(nil)
	require.NoError(t, err)
	require.Equal(t, "", result)
}

// TestTemplate_Concurrent executes one Template from many goroutines, which the race detector
// reports if execution writes shared state
func TestTemplate_Concurrent(t *testing.T) {

	value := NewTemplate("{{.ID}}")
	results := make(chan bool, 20)

	var wait sync.WaitGroup
	for range 20 {
		wait.Go(func() {
			result, err := value.Execute(templateData{ID: "x"})
			results <- (err == nil) && (result == "x")
		})
	}

	wait.Wait()
	close(results)

	for ok := range results {
		require.True(t, ok)
	}
}

// TestTemplate_JSON requires that a Template decodes any scalar as a string, and encodes as that
// string
func TestTemplate_JSON(t *testing.T) {

	for _, test := range []struct {
		source  string
		encoded string
	}{
		{`"/v?id={{.ID}}"`, `"/v?id={{.ID}}"`},
		{`"off"`, `"off"`},
		{`"{{.ID | nosuchfunc}}"`, `"{{.ID | nosuchfunc}}"`},
		{`6`, `"6"`},
		{`1.0`, `"1.0"`},
		{`true`, `"true"`},
		{`null`, `""`},
	} {
		var value Template
		require.NoError(t, json.Unmarshal([]byte(test.source), &value), test.source)

		encoded, err := json.Marshal(value)
		require.NoError(t, err)
		require.Equal(t, test.encoded, string(encoded), test.source)
	}

	var value Template
	require.NoError(t, json.Unmarshal([]byte(`"/v?id={{.ID}}"`), &value))
	require.True(t, value.IsTemplate())
}

// TestTemplate_JSON_Invalid requires that malformed JSON, objects, and arrays fail to decode,
// leaving the target unchanged
func TestTemplate_JSON_Invalid(t *testing.T) {

	for _, source := range []string{`{"a":`, `{"a":1}`, `["a"]`} {
		value := NewTemplate("unchanged")
		require.Error(t, value.UnmarshalJSON([]byte(source)), source)
		require.Equal(t, "unchanged", value.String())
	}
}
