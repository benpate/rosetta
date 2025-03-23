package translate

import (
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/funcmap"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// conditionRunner evaluates a condition, and executes a list of rules if the condition is true
type conditionRunner struct {
	Condition    *template.Template
	ConditionRaw string
	ThenRules    Pipeline
	ElseRules    Pipeline
}

// Condition creates a new Rule that executes a condition, and then runs a set of rules based on the result
func Condition(condition string, thenRules []Rule, elseRules []Rule) Rule {

	conditionTemplate, err := template.New("").Funcs(funcmap.All()).Parse(condition)
	derp.Report(err)

	return Rule{
		conditionRunner{
			Condition:    conditionTemplate,
			ConditionRaw: condition,
			ThenRules:    thenRules,
			ElseRules:    elseRules,
		},
	}
}

// Execute implements the Runner interface
func (runner conditionRunner) Execute(sourceSchema schema.Schema, sourceValue any, targetSchema schema.Schema, targetValue any) error {

	const location = "rosetta.translate.conditionRunner.Execute"

	condition := executeTemplate(runner.Condition, sourceValue)

	if convert.Bool(condition) {

		// If the condition is true, then execute the rules
		if err := runner.ThenRules.Execute(sourceSchema, sourceValue, targetSchema, targetValue); err != nil {
			return derp.Wrap(err, location, "Error executing rules")
		}

		return nil
	}

	// If the condition is true, then execute the rules
	if err := runner.ElseRules.Execute(sourceSchema, sourceValue, targetSchema, targetValue); err != nil {
		return derp.Wrap(err, location, "Error executing rules")
	}

	return nil
}

/******************************************
 * Serialization Methods
 ******************************************/

func (runner conditionRunner) MarshalMap() map[string]any {
	return map[string]any{
		"if":   runner.ConditionRaw,
		"then": runner.ThenRules.MarshalSliceOfMap(),
		"else": runner.ElseRules.MarshalSliceOfMap(),
	}
}

func (runner *conditionRunner) UnmarshalMap(data mapof.Any) error {

	return runner.populate(
		data.GetString("if"),
		data.GetSliceOfPlainMap("then"),
		data.GetSliceOfPlainMap("else"),
	)
}

func (runner *conditionRunner) populate(condition string, thenRules []map[string]any, elseRules []map[string]any) error {

	const location = "rosetta.translate.conditionRunner.UnmarshalMap"

	// Parse Condition
	conditionTemplate, err := template.New("").Parse(condition)

	if err != nil {
		return derp.Wrap(err, location, "Error parsing `if` condition template", condition)
	}

	// Parse Then Rules
	thenPipeline, err := NewFromMap(thenRules...)

	if err != nil {
		return derp.Wrap(err, location, "Error creating `then` rules", thenRules)
	}

	// Parse Else Rules
	elsePipeline, err := NewFromMap(elseRules...)

	if err != nil {
		return derp.Wrap(err, location, "Error creating `else` rules", elseRules)
	}

	// Apply values
	runner.Condition = conditionTemplate
	runner.ConditionRaw = condition
	runner.ThenRules = thenPipeline
	runner.ElseRules = elsePipeline

	return nil
}
