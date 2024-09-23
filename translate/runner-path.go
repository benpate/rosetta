package translate

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema"
)

// pathRunner retrieves a value from the input object, and writes it to the output object
type pathRunner struct {
	Path   string
	Target string
}

// Path creates a new Rule that copies a value from one location to another
func Path(from string, target string) Rule {
	return Rule{newPathRunner(from, target)}
}

// newPathRunner returns a fully initialized pathRunner object
func newPathRunner(path string, target string) pathRunner {
	return pathRunner{
		Path:   path,
		Target: target,
	}
}

// Execute implements the Runner interface
func (runner pathRunner) Execute(sourceSchema schema.Schema, sourceObject any, targetSchema schema.Schema, targetObject any) error {

	const location = "mapper.pathRunner.Execute"

	value, err := sourceSchema.Get(sourceObject, runner.Path)

	if err != nil {
		return derp.Wrap(err, location, "Error getting value from source", runner.Path)
	}

	if err := targetSchema.Set(targetObject, runner.Target, value); err != nil {
		return derp.Wrap(err, location, "Error setting value in target", runner.Target)
	}

	return nil
}
