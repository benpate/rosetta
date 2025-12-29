package translate

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// appendRunner applies a constant value to the output object
type appendRunner struct {
	Append any
	Target string
}

// Append creates a new Rule that writes a constant value to the output object
func Append(append any, target string) Rule {
	return Rule{newAppendRunner(append, target)}
}

// newAppendRunner returns a fully initialized appendRunner object
func newAppendRunner(append any, target string) appendRunner {
	return appendRunner{
		Append: append,
		Target: target,
	}
}

// Execute implements the Runner interface
func (runner appendRunner) Execute(_ schema.Schema, _ any, targetSchema schema.Schema, targteObject any) error {

	if err := targetSchema.Append(targteObject, runner.Target, runner.Append); err != nil {
		return derp.Wrap(err, "rosetta.translate.appendRunner.Set", "Unable to set value in target", runner.Target)
	}

	return nil
}

/******************************************
 * Serialization Methods
 ******************************************/

func (runner appendRunner) MarshalMap() map[string]any {
	return mapof.Any{
		"append": runner.Append,
		"target": runner.Target,
	}
}

func (runner *appendRunner) UnmarshalMap(data mapof.Any) error {

	runner.Append = upscale(data.GetAny("append"))
	runner.Target = data.GetString("target")

	return nil
}
