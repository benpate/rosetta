package translate

import (
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/schema"
)

// conditionRunner evaluates a condition, and executes a list of rules if the condition is true
type conditionRunner struct {
	Condition *template.Template
	ThenRules Pipeline
	ElseRules Pipeline
}

// Condition creates a new Rule that executes a condition, and then runs a set of rules based on the result
func Condition(condition string, thenRules []Rule, elseRules []Rule) Rule {

	conditionTemplate, err := template.New("").Parse(condition)
	derp.Report(err)

	return Rule{
		conditionRunner{
			Condition: conditionTemplate,
			ThenRules: thenRules,
			ElseRules: elseRules,
		},
	}
}

// newConditionRunner returns a fully initialized conditionRunner
func newConditionRunner(condition string, thenMap []map[string]any, elseMap []map[string]any) (conditionRunner, error) {

	const location = "rosetta.translate.newConditionRunner"

	// Parse the "if" template
	conditionTemplate, err := template.New("").Parse(condition)

	if err != nil {
		return conditionRunner{}, derp.Wrap(err, location, "Error parsing template", condition)
	}

	// Parse the "then" rules
	thenRules, err := NewFromMap(thenMap...)

	if err != nil {
		return conditionRunner{}, derp.Wrap(err, location, "Error creating Pipeline", thenMap)
	}

	// Parse the "else" rules
	elseRules, err := NewFromMap(elseMap...)

	if err != nil {
		return conditionRunner{}, derp.Wrap(err, location, "Error creating Pipeline", elseMap)
	}

	return conditionRunner{
		Condition: conditionTemplate,
		ThenRules: thenRules,
		ElseRules: elseRules,
	}, nil
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
