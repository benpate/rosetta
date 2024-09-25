package translate

import (
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema"
)

// expressionRunner is a Runner that executes a template expression
type expressionRunner struct {
	Expression *template.Template
	Target     string
}

// Expression creates a new Rule that executes a template expression
func Expression(expression string, target string) Rule {

	r, err := newExpressionRunner(expression, target)
	derp.Report(err)

	return Rule{r}
}

// newExpressionRunner returns a fully initialized expressionRunner
func newExpressionRunner(expression string, target string) (expressionRunner, error) {

	t, err := template.New("expression").Parse(expression)

	if err != nil {
		return expressionRunner{}, derp.Wrap(err, "mapper.NewexpressionRunner", "Error parsing template", expression)
	}

	return expressionRunner{
		Expression: t,
		Target:     target,
	}, nil
}

// Execute implements the Runner interface
func (runner expressionRunner) Execute(_ schema.Schema, sourceValue any, targetSchema schema.Schema, targetValue any) error {

	value := executeTemplate(runner.Expression, sourceValue)

	if err := targetSchema.Set(targetValue, runner.Target, value); err != nil {
		return derp.Wrap(err, "mapper.expressionRunner.Set", "Error setting value in target", runner.Target)
	}

	return nil
}
