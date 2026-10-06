package translate

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

func TestPipelineUnmarshal(t *testing.T) {

	rulesJSON := []byte(`[
		{"path":"name1", "target":"name1"},
		{"value":"application/json", "target":"mimeType"},
		{"expression":"{{.firstName}} {{.lastName}}", "target":"fullName"}
	]`)

	rules := Pipeline{}

	if err := json.Unmarshal(rulesJSON, &rules); err != nil {
		t.Error(err)
	}

	require.Equal(t, 3, len(rules))

	require.Equal(t, "name1", rules[0].Runner.(pathRunner).Path)
	require.Equal(t, "name1", rules[0].Runner.(pathRunner).Target)

	require.Equal(t, "application/json", rules[1].Runner.(valueRunner).Value)
	require.Equal(t, "mimeType", rules[1].Runner.(valueRunner).Target)

	require.NotNil(t, rules[2].Runner.(expressionRunner).Expression)
	require.Equal(t, "{{.firstName}} {{.lastName}}", rules[2].Runner.(expressionRunner).ExpressionRaw)
	require.Equal(t, "fullName", rules[2].Runner.(expressionRunner).Target)
}

func TestExecuteRules(t *testing.T) {

	// TARGET CONFIGURATION
	targetSchema := schema.New(schema.Object{
		Properties: schema.ElementMap{
			"fullName": schema.String{},
			"email":    schema.String{Format: "email"},
			"type":     schema.String{},
			"comment":  schema.String{},
		},
	})

	// MAPPING RULES
	rules, err := NewFromJSON(`[
		{"expression": "{{.firstName}} {{.lastName}}", "target": "fullName"},
		{"path": "email", "target": "email"},
		{"value": "person", "target": "type"},
		{"if": "{{eq \"M\" .gender}}", "then": [
			{"target": "comment", "expression": "{{.firstName}} is Male"}
		], "else": [
			{"target": "comment", "expression": "{{.firstName}} is not Male"}
		]}
	]`)

	require.Nil(t, err)

	// TEST JOHN
	{
		sourceValue := mapof.Any{
			"firstName": "John",
			"lastName":  "Connor",
			"email":     "john@connor.mil",
			"gender":    "M",
		}

		targetValue := mapof.Any{}

		err = rules.Execute(schema.Wildcard(), sourceValue, targetSchema, &targetValue)
		require.Nil(t, err)

		require.Equal(t, "John Connor", targetValue.GetString("fullName"))
		require.Equal(t, "john@connor.mil", targetValue.GetString("email"))
		require.Equal(t, "person", targetValue.GetString("type"))
		require.Equal(t, "John is Male", targetValue.GetString("comment"))
	}

	// TEST SARAH
	{
		sourceValue := mapof.Any{
			"firstName": "Sarah",
			"lastName":  "Connor",
			"email":     "sarah@sky.net",
			"gender":    "F",
		}

		targetValue := mapof.Any{}

		err = rules.Execute(schema.Wildcard(), sourceValue, targetSchema, &targetValue)
		require.Nil(t, err)

		require.Equal(t, "Sarah Connor", targetValue.GetString("fullName"))
		require.Equal(t, "sarah@sky.net", targetValue.GetString("email"))
		require.Equal(t, "person", targetValue.GetString("type"))
		require.Equal(t, "Sarah is not Male", targetValue.GetString("comment"))
	}
}

func ExampleNew() {

	// Define rules for the translation Pipeline
	rules := []Rule{
		Expression("{{.firstName}} {{.lastName}}", "fullName"),
		Path("email", "email"),
		Value("person", "type"),
		Condition(`{{eq "M" .gender}}`, []Rule{
			Expression("{{.firstName}} is Male", "comment"),
		}, []Rule{
			Expression("{{.firstName}} is not Malr", "comment"),
		}),
	}

	// Add all of the rules to a new Pipeline
	translator := New(rules...)
	fmt.Println(translator)
}

func ExampleNewFromJSON() {

	// Import JSON from external source
	rulesJSON := `[` + // nolint:scopeguard
		`{"expression": "{{.firstName}} {{.lastName}}", "target": "fullName"},
		{"path": "email", "target": "email"},
		{"value": "person", "target": "type"},
		{"if": "{{eq \"M\" .gender}}", "then": [
			{"target": "comment", "expression": "{{.firstName}} is Male"}
		], "else": [
			{"target": "comment", "expression": "{{.firstName}} is not Male"}
		]}
	]`

	// Unmarshal JSON directly into a Pipeline
	rules, _ := NewFromMap()
	if json.Unmarshal([]byte(rulesJSON), &rules) != nil {
		fmt.Println("Error parsing JSON")
	}

	// Success!
	fmt.Println(rules)
}

func ExampleNewFromMap() {

	// Define rules as a mapof.Any
	rules := []map[string]any{
		{"expression": "{{.firstName}} {{.lastName}}", "target": "fullName"},
		{"path": "email", "target": "email"},
		{"value": "person", "target": "type"},
		{"if": `{{eq "M" .gender}}`, "then": []mapof.Any{
			{"target": "comment", "expression": "{{.firstName}} is Male"},
		}, "else": []mapof.Any{
			{"target": "comment", "expression": "{{.firstName}} is not Male"},
		}},
	}

	// Create a new Pipeline from the rules
	if translator, err := NewFromMap(rules...); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(translator)
	}
}

func ExamplePipeline_Execute() {

	// SOURCE DATA
	sourceValue := mapof.Any{
		"firstName": "John",
		"lastName":  "Connor",
		"email":     "john@connor.mil",
		"gender":    "M",
	}

	// TARGET CONFIGURATION
	targetSchema := schema.New(schema.Object{
		Properties: schema.ElementMap{
			"fullName": schema.String{},
			"email":    schema.String{Format: "email"},
			"type":     schema.String{},
			"comment":  schema.String{},
		},
	})
	targetValue := mapof.Any{}

	// CREATE MAPPING RULES
	rules, err := NewFromMap(
		mapof.Any{"target": "fullName", "expression": "{{.firstName}} {{.lastName}}"},
		mapof.Any{"target": "email", "path": "email"},
		mapof.Any{"target": "type", "value": "person"},
		mapof.Any{"if": "{{eq \"M\" .gender}}", "then": []mapof.Any{
			{"target": "comment", "expression": "{{.firstName}} is Male"},
		}, "else": []mapof.Any{
			{"target": "comment", "expression": "{{.firstName}} is not Male"},
		}},
	)
	derp.Report(err)

	// MAP DATA FROM SOURCE TO TARGET
	err = rules.Execute(schema.Wildcard(), sourceValue, targetSchema, &targetValue)
	derp.Report(err)

	// OUTPUT RESULTS
	fmt.Println(targetValue.GetString("fullName"))
	fmt.Println(targetValue.GetString("email"))
	fmt.Println(targetValue.GetString("type"))
	fmt.Println(targetValue.GetString("comment"))

	// Output:
	// John Connor
	// john@connor.mil
	// person
	// John is Male
}

// TestExecuteRules_IndexedTarget pins Bandwagon's album "artists" rules: the indexed target
// writes into the first item of the list the first two rules built
func TestExecuteRules_IndexedTarget(t *testing.T) {

	// BUG-234: these rules used to replace the list with {"0": {...}}

	source := mapof.Any{"attributedTo": mapof.Any{"profileUrl": "https://example.com/@a", "name": "Lime Bar"}}

	t.Run("the first two rules build a list", func(t *testing.T) {
		rules, err := NewFromJSON(`[
			{"target": "artists", "value": []},
			{"target": "artists", "append": {"type": "Artist"}}
		]`)
		require.NoError(t, err)

		target := mapof.Any{}
		require.NoError(t, rules.Execute(activityStreamSchema(), source, activityStreamSchema(), &target))
		require.Equal(t, mapof.Any{"artists": &sliceof.Any{mapof.Any{"type": "Artist"}}}, target)
	})

	t.Run("the indexed rules write into its first item", func(t *testing.T) {
		rules, err := NewFromJSON(`[
			{"target": "artists", "value": []},
			{"target": "artists", "append": {"type": "Artist"}},
			{"target": "artists.0.id", "path": "attributedTo.profileUrl"},
			{"target": "artists.0.name", "path": "attributedTo.name"}
		]`)
		require.NoError(t, err)

		target := mapof.Any{}
		require.NoError(t, rules.Execute(activityStreamSchema(), source, activityStreamSchema(), &target))
		require.Equal(t, mapof.Any{"artists": &sliceof.Any{mapof.Any{"type": "Artist", "id": "https://example.com/@a", "name": "Lime Bar"}}}, target)
	})
}

// TestExecuteRules_StopsAtFirstError pins that a failed rule ends the pipeline: earlier rules
// keep their writes, and later rules never run
func TestExecuteRules_StopsAtFirstError(t *testing.T) {

	rules := New(
		Value("b", "before"),
		ForEach("missing", "items", "", []map[string]any{}),
		Value("a", "after"),
	)

	target := mapof.Any{}
	err := rules.Execute(activityStreamSchema(), mapof.Any{"missing": "not a list"}, activityStreamSchema(), &target)
	require.Error(t, err)
	require.Equal(t, mapof.Any{"before": "b"}, target)
}
