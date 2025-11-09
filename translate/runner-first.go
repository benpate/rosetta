package translate

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// firstRunner executes its sub-steps, stopping with the first one that returns a non-zero value
type firstRunner struct {
	TargetPath string
	Rules      Pipeline
}

// First creates a new Rule that executes its sub-steps, stopping with the first one that returns a non-zero value
func First(targetPath string, rules []map[string]any) Rule {

	runner := firstRunner{}
	if err := runner.populate(targetPath, rules); err != nil {
		derp.Report(err)
	}

	return Rule{runner}
}

// Execute implements the Runner interface
func (runner firstRunner) Execute(sourceSchema schema.Schema, sourceValue any, targetSchema schema.Schema, targetValue any) error {

	// Loop through each rule in the array, stopping after the first rule that actually sets a value
	for _, rule := range runner.Rules {

		// Execute the rule. If there's an error, fail silently and continue to the next rule
		if err := rule.Execute(sourceSchema, sourceValue, targetSchema, targetValue); err != nil {
			continue
		}

		// If the target value is not zero, then stop processing
		if value, err := targetSchema.Get(targetValue, runner.TargetPath); err == nil {

			if !compare.IsZero(value) {
				break
			}
		}
	}

	return nil
}

/******************************************
 * Serialization Methods
 ******************************************/

func (runner firstRunner) MarshalMap() map[string]any {
	return map[string]any{
		"first": runner.TargetPath,
		"rules": runner.Rules.MarshalSliceOfMap(),
	}
}

func (runner *firstRunner) UnmarshalMap(data mapof.Any) error {

	return runner.populate(
		data.GetString("first"),
		data.GetSliceOfPlainMap("rules"),
	)
}

func (runner *firstRunner) populate(target string, rules []map[string]any) error {

	// Parse Rules
	pipeline, err := NewFromMap(rules...)

	if err != nil {
		return derp.Wrap(err, "rosetta.translate.firstRunner.populate", "Unable to create Pipeline", rules)
	}

	// Populate remaining fields
	runner.TargetPath = target
	runner.Rules = pipeline

	return nil
}
