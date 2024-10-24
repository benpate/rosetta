package translate

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema"
)

// valueRunner applies a constant value to the output object
type valueRunner struct {
	Value  any
	Target string
}

// Value creates a new Rule that writes a constant value to the output object
func Value(value any, target string) Rule {
	return Rule{newValueRunner(value, target)}
}

// newValueRunner returns a fully initialized valueRunner object
func newValueRunner(value any, target string) valueRunner {
	return valueRunner{
		Value:  value,
		Target: target,
	}
}

// Execute implements the Runner interface
func (runner valueRunner) Execute(_ schema.Schema, _ any, targetSchema schema.Schema, targteObject any) error {

	if err := targetSchema.Set(targteObject, runner.Target, runner.Value); err != nil {
		return derp.Wrap(err, "rosetta.translate.valueRunner.Set", "Error setting value in target", runner.Target)
	}

	return nil
}
