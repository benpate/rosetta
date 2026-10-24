package translate

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/derp/plugins"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

// plantedSecret is a value that no error may carry, whatever rule fails
const plantedSecret = "hunter2-planted-secret"

// secretRecord is a source object that holds a secret beside the fields the rules read
type secretRecord struct {
	Label    string
	Password string
}

// errorRecorder is a derp Reporter that keeps every error reported to it
type errorRecorder struct {
	errors []error
}

// Report implements the derp.Reporter interface
func (recorder *errorRecorder) Report(err error) {
	recorder.errors = append(recorder.errors, err)
}

// TestForEach_ErrorOmitsSource verifies that a forEach over a value that is not a list
// names the path, and never attaches the object it was reading
func TestForEach_ErrorOmitsSource(t *testing.T) {

	sourceSchema := schema.New(schema.Object{Properties: schema.ElementMap{
		"label":    schema.Any{},
		"password": schema.String{},
	}})
	source := mapof.Any{"label": "a single value", "password": plantedSecret}

	targetSchema := schema.New(schema.Object{Properties: schema.ElementMap{
		"items": schema.Array{Items: schema.Object{Properties: schema.ElementMap{"value": schema.String{}}}},
	}})
	target := mapof.NewAny()

	rules := New(ForEach("label", "items", "", []map[string]any{{"path": "value", "target": "value"}}))
	err := rules.Execute(sourceSchema, source, targetSchema, &target)

	require.Error(t, err)
	requireSecretAbsent(t, err)

	// The error still says which path failed, and what it held
	encoded, encodeErr := json.Marshal(err)
	require.NoError(t, encodeErr)
	require.Contains(t, string(encoded), `"details":["label","string"]`)
}

// TestPath_ErrorOmitsValue verifies that a path rule whose target rejects the value
// does not attach the value it read
func TestPath_ErrorOmitsValue(t *testing.T) {

	sourceSchema := schema.New(schema.Object{Properties: schema.ElementMap{
		"passwordReset": schema.Any{},
	}})
	source := mapof.Any{"passwordReset": mapof.Any{"authcode": plantedSecret}}

	targetSchema := schema.New(schema.Object{Properties: schema.ElementMap{
		"count": schema.Integer{},
	}})
	target := mapof.NewAny()

	err := New(Path("passwordReset", "count")).Execute(sourceSchema, source, targetSchema, &target)

	require.Error(t, err)
	requireSecretAbsent(t, err)
}

// TestExecuteTemplate_ReportOmitsData verifies that a template that fails to execute
// reports the error without attaching the object it was rendering
func TestExecuteTemplate_ReportOmitsData(t *testing.T) {

	recorder := &errorRecorder{}
	derp.SetPlugins(recorder)
	defer derp.SetPlugins(plugins.JSON{})

	sourceSchema := schema.New(schema.Object{Properties: schema.ElementMap{
		"Label":    schema.String{},
		"Password": schema.String{},
	}})
	source := &secretRecord{Label: "Profile", Password: plantedSecret}

	targetSchema := schema.New(schema.Object{Properties: schema.ElementMap{
		"name": schema.String{},
	}})
	target := mapof.NewAny()

	rules := New(
		Expression("{{.Missing}}", "name"),
		Condition("{{.Missing}}", nil, nil),
	)

	require.NoError(t, rules.Execute(sourceSchema, source, targetSchema, &target))
	require.Len(t, recorder.errors, 2)

	for _, err := range recorder.errors {
		requireSecretAbsent(t, err)
	}
}

// requireSecretAbsent fails the test if the planted secret appears in the error's log line,
// or in its encoding, which holds the message and details of every layer
func requireSecretAbsent(t *testing.T, err error) {

	t.Helper()

	encoded, encodeErr := json.Marshal(err)
	require.NoError(t, encodeErr)
	require.NotContains(t, string(encoded), plantedSecret)
	require.NotContains(t, fmt.Sprintf("%+v", err), plantedSecret)
}
