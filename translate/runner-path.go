package translate

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	deepcopy "github.com/tiendc/go-deepcopy"
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

	const location = "rosetta.translate.pathRunner.Execute"

	var duplicateValue any = ""

	// Get the source value from the source object
	sourceValue, err := sourceSchema.Get(sourceObject, runner.Path)

	if err == nil {

		// Make a deep copy of the source (to prevent pointer shenanigans)
		if err := deepcopy.Copy(&duplicateValue, sourceValue); err != nil {
			return derp.Wrap(err, location, "Unable to deep copy value", runner.Path, sourceValue)
		}

		// Dereference pointers (if applicable)
		duplicateValue = convert.Element(duplicateValue)
	}

	// Set the value in the target object
	if err := targetSchema.Set(targetObject, runner.Target, duplicateValue); err != nil {
		return derp.Wrap(err, location, "Unable to set target value", runner.Target, sourceValue)
	}

	// Success.
	return nil
}

/******************************************
 * Serialization Methods
 ******************************************/

func (runner pathRunner) MarshalMap() map[string]any {
	return map[string]any{
		"path":   runner.Path,
		"target": runner.Target,
	}
}

func (runner *pathRunner) UnmarshalMap(data mapof.Any) error {

	runner.Path = data.GetString("path")
	runner.Target = data.GetString("target")

	return nil
}
