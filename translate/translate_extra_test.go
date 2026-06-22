package translate

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestPipeline_EmptyNotEmpty(t *testing.T) {
	require.True(t, New().IsEmpty())
	require.False(t, New().NotEmpty())

	p := New(Value("x", "target"))
	require.False(t, p.IsEmpty())
	require.True(t, p.NotEmpty())
}

func TestPipeline_NewFromMap(t *testing.T) {

	pipeline, err := NewFromMap(
		map[string]any{"path": "a", "target": "b"},
		map[string]any{"value": "v", "target": "c"},
	)

	require.NoError(t, err)
	require.Equal(t, 2, len(pipeline))
}

func TestPipeline_NewFromMap_Error(t *testing.T) {
	// A map with no recognized runner key is an error
	_, err := NewFromMap(map[string]any{"bogus": "data"})
	require.Error(t, err)
}

func TestNewSliceOfPipelines(t *testing.T) {

	result, err := NewSliceOfPipelines([][]map[string]any{
		{{"path": "a", "target": "b"}},
		{{"value": "v", "target": "c"}, {"path": "x", "target": "y"}},
	})

	require.NoError(t, err)
	require.Equal(t, 2, len(result))
	require.Equal(t, 1, len(result[0]))
	require.Equal(t, 2, len(result[1]))
}

func TestNewSliceOfPipelines_Error(t *testing.T) {
	_, err := NewSliceOfPipelines([][]map[string]any{
		{{"bogus": "data"}},
	})
	require.Error(t, err)
}

// TestPipeline_RoundTrip builds a pipeline using each constructor, marshals it to
// JSON, and reads it back, confirming the runner types survive the round trip.
func TestPipeline_RoundTrip(t *testing.T) {

	original := New(
		Path("email", "email"),
		Value("person", "type"),
		Expression("{{.firstName}} {{.lastName}}", "fullName"),
		Append("tag", "tags"),
		Condition(`{{eq "M" .gender}}`,
			[]Rule{Expression("{{.firstName}} is Male", "comment")},
			[]Rule{Expression("{{.firstName}} is not Male", "comment")},
		),
		First("winner", []map[string]any{
			{"path": "nickname", "target": "winner"},
			{"path": "firstName", "target": "winner"},
		}),
		ForEach("items", "results", "", []map[string]any{
			{"path": "name", "target": "name"},
		}),
	)

	data, err := json.Marshal(original)
	require.NoError(t, err)

	restored, err := NewFromJSON(string(data))
	require.NoError(t, err)
	require.Equal(t, len(original), len(restored))

	// Confirm each runner deserialized to the expected concrete type
	require.IsType(t, pathRunner{}, restored[0].Runner)
	require.IsType(t, valueRunner{}, restored[1].Runner)
	require.IsType(t, expressionRunner{}, restored[2].Runner)
	require.IsType(t, appendRunner{}, restored[3].Runner)
	require.IsType(t, conditionRunner{}, restored[4].Runner)
	require.IsType(t, firstRunner{}, restored[5].Runner)
	require.IsType(t, forEachRunner{}, restored[6].Runner)
}

func TestPipeline_MarshalSliceOfMap(t *testing.T) {

	pipeline := New(Value("person", "type"))

	result := pipeline.MarshalSliceOfMap()
	require.Equal(t, 1, len(result))
	require.Equal(t, "person", result[0]["value"])
	require.Equal(t, "type", result[0]["target"])
}

func TestRule_MarshalJSON_Nil(t *testing.T) {

	// A rule with no runner marshals to null and an empty map
	rule := Rule{}

	data, err := rule.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "null", string(data))

	require.Equal(t, map[string]any{}, rule.MarshalMap())
}

func TestRule_UnmarshalJSON_Error(t *testing.T) {
	var rule Rule
	require.Error(t, json.Unmarshal([]byte("not json"), &rule))
}

func TestFirstRunner_Execute(t *testing.T) {

	targetSchema := schema.New(schema.Object{
		Properties: schema.ElementMap{
			"winner": schema.String{},
		},
	})

	rule := First("winner", []map[string]any{
		{"path": "nickname", "target": "winner"},
		{"path": "firstName", "target": "winner"},
	})

	pipeline := New(rule)

	// nickname is empty, so the second rule (firstName) wins
	source := mapof.Any{"firstName": "John", "nickname": ""}
	target := mapof.Any{}

	require.NoError(t, pipeline.Execute(schema.Wildcard(), source, targetSchema, &target))
	require.Equal(t, "John", target.GetString("winner"))
}

func TestValueRunner_Execute(t *testing.T) {

	targetSchema := schema.New(schema.Object{
		Properties: schema.ElementMap{
			"type": schema.String{},
		},
	})

	pipeline := New(Value("person", "type"))

	target := mapof.Any{}
	require.NoError(t, pipeline.Execute(schema.Wildcard(), mapof.Any{}, targetSchema, &target))
	require.Equal(t, "person", target.GetString("type"))
}
