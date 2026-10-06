package translate

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
	deepcopy "github.com/tiendc/go-deepcopy"
)

func TestDeepCopy_Map(t *testing.T) {

	source := map[string]any{
		"name": "John",
		"age":  30,
		"address": map[string]any{
			"street": "123 Main St",
			"city":   "Anytown",
		},
	}

	var target any

	if err := deepcopy.Copy(&target, source); err != nil {
		t.Errorf("Error during deepcopy: %v", err)
	}

	t.Log(source, target)
}

func TestDeepCopy_Struct(t *testing.T) {

	source := struct {
		Name string
		Age  int
	}{
		Name: "John",
		Age:  30,
	}

	var target any

	if err := deepcopy.Copy(&target, source); err != nil {
		t.Errorf("Error during deepcopy: %v", err)
	}

	t.Log(source, target)
}

// pathSource is a struct whose fields are reached through GetPointer, as Emissary's Stream is
type pathSource struct {
	Label string
	Rank  int
	Data  mapof.Any
}

// GetPointer implements the schema PointerGetter interface
func (source *pathSource) GetPointer(name string) (any, bool) {

	switch name {
	case "label":
		return &source.Label, true
	case "rank":
		return &source.Rank, true
	case "data":
		return &source.Data, true
	}

	return nil, false
}

// TestPathRunner_StructSourceGuard pins the JSON that path rules produce when they read a struct
// through a wildcard schema, the way Emissary builds a stream's ActivityPub document. Get returns
// a struct's fields by pointer today; a change to that must not change this output.
func TestPathRunner_StructSourceGuard(t *testing.T) {

	rules, err := NewFromJSON(`[
		{"target": "name", "path": "label"},
		{"target": "track.position", "path": "rank"},
		{"target": "content", "path": "data.lyrics"},
		{"target": "missing", "path": "data.nothing"}
	]`)
	require.NoError(t, err)

	source := pathSource{Label: "Headlights", Rank: 3, Data: mapof.Any{"lyrics": "la la"}}
	target := mapof.Any{}

	ruleSchema := activityStreamSchema()
	require.NoError(t, rules.Execute(ruleSchema, &source, ruleSchema, &target))

	encoded, err := json.Marshal(target)
	require.NoError(t, err)
	expected := `{"name":"Headlights","track":{"position":3},"content":"la la","missing":""}`
	require.JSONEq(t, expected, string(encoded))
}
