package translate

import (
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// expressionRunner is a Runner that executes a template expression
type expressionRunner struct {
	Expression    *template.Template
	ExpressionRaw string
	Target        string
}

// Expression creates a new Rule that executes a template expression
func Expression(expression string, target string) Rule {

	runner := expressionRunner{}
	err := runner.populate(expression, target)
	derp.Report(err)

	return Rule{runner}
}

// Execute implements the Runner interface
func (runner expressionRunner) Execute(_ schema.Schema, sourceValue any, targetSchema schema.Schema, targetValue any) error {

	value := executeTemplate(runner.Expression, sourceValue)

	if err := targetSchema.Set(targetValue, runner.Target, value); err != nil {
		return derp.Wrap(err, "rosetta.translate.expressionRunner.Set", "Error setting value in target", runner.Target)
	}

	return nil
}

/******************************************
 * Serialization Methods
 ******************************************/

func (runner expressionRunner) MarshalMap() map[string]any {
	return map[string]any{
		"expression": runner.ExpressionRaw,
		"target":     runner.Target,
	}
}

// newExpressionRunner returns a fully initialized expressionRunner
func (runner *expressionRunner) UnmarshalMap(data mapof.Any) error {

	return runner.populate(
		data.GetString("expression"),
		data.GetString("target"),
	)
}

func (runner *expressionRunner) populate(expression string, target string) error {

	const location = "rosetta.translate.expressionRunner.populate"

	// Parse the template
	expressionTemplate, err := template.New("").Parse(expression)

	if err != nil {
		return derp.Wrap(err, location, "Error parsing `expression` template", expression)
	}

	// Populate values
	runner.Expression = expressionTemplate
	runner.ExpressionRaw = expression
	runner.Target = target

	// Great Success!
	return nil
}
