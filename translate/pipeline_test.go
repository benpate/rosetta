package translate

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
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
	rulesJSON := `[
		{"expression": "{{.firstName}} {{.lastName}}", "target": "fullName"},
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
	rules := []mapof.Any{
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
		mapof.Any{"expression": "{{.firstName}} {{.lastName}}", "target": "fullName"},
		mapof.Any{"path": "email", "target": "email"},
		mapof.Any{"value": "person", "target": "type"},
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
