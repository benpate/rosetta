package translate

import (
	"text/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/funcmap"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
)

// forEachRunner retrieves a value from the input object, and writes it to the output object
type forEachRunner struct {
	SourcePath string
	TargetPath string
	Filter     *template.Template
	FilterRaw  string
	Pipeline   Pipeline
}

// ForEach creates a new Rule that copies a value from one location to another
func ForEach(sourcePath string, targetPath string, filter string, rulesMap []map[string]any) Rule {

	runner := forEachRunner{}

	if err := runner.populate(sourcePath, targetPath, filter, rulesMap); err != nil {
		derp.Report(err)
	}

	return Rule{runner}
}

// Execute implements the Runner interface
func (runner forEachRunner) Execute(sourceSchema schema.Schema, sourceValue any, targetSchema schema.Schema, targetValue any) error {

	const location = "rosetta.translate.forEachRunner.Execute"

	// Get the element of each source item from the sourceSchema
	sourceItemElement, err := itemElement(sourceSchema, runner.SourcePath)

	if err != nil {
		return derp.Wrap(err, location, "Reading source element", runner.SourcePath)
	}

	// RULE: A source that is missing or nil has no items, so the following rules still run.
	// This matches the path rule, which writes "" for a source it cannot read.
	sourceArray, err := sourceSchema.Get(sourceValue, runner.SourcePath)

	if (err != nil) || (sourceArray == nil) {
		return nil
	}

	sourceGetter, ok := sourceArray.(schema.KeysGetter)

	if !ok {
		return derp.Internal(location, "Source value must implement schema.KeysGetter", sourceValue)
	}

	// Get the list of keys from the source array.  If none, then exit
	sourceKeys := sourceGetter.Keys()
	if len(sourceKeys) == 0 {
		return nil
	}

	// Get the target element from the targetSchema
	targetElement, ok := targetSchema.GetArrayElement(runner.TargetPath)

	if !ok {
		return derp.Internal(location, "Target element must exist in targetSchema", runner.TargetPath)
	}

	// Create the new schemas for the source/target array items.  Item rules read each
	// source item as {key, value}, so the source item's own element describes "value".
	sourceItemSchema := schema.New(schema.Object{Properties: schema.ElementMap{
		"key":   schema.String{},
		"value": sourceItemElement,
	}})
	targetItemSchema := schema.New(targetElement.Items)

	targetSlice := sliceof.NewObject[mapof.Any]()

	// Loop through each element in the array
	for _, key := range sourceKeys {

		// Get the value of the source array at the current index
		sourcePath := list.ByDot(runner.SourcePath).PushTail(key).String()
		sourceItemValue, _ := sourceSchema.Get(sourceValue, sourcePath)

		// If the filter exists, and returns false, then skip this record
		if runner.Filter != nil {
			if !convert.Bool(executeTemplate(runner.Filter, sourceItemValue)) {
				continue
			}
		}

		// Create a new item to add to the end of the target array
		sourceMap := mapof.Any{
			"key":   key,
			"value": sourceItemValue,
		}

		targetMap := mapof.NewAny()

		// Apply the pipeline to from the source to the target
		if err := runner.Pipeline.Execute(sourceItemSchema, sourceMap, targetItemSchema, &targetMap); err != nil {
			return derp.Wrap(err, location, "Error executing pipeline", runner.Pipeline)
		}

		// Add the new item to the target array
		targetSlice.Append(targetMap)
	}

	// Write the updated targetValue back to the targetValue. We pass a pointer
	// so that Set's validation can treat the slice as an ArrayGetterSetter.
	if len(targetSlice) > 0 {
		if err := targetSchema.Set(targetValue, runner.TargetPath, &targetSlice); err != nil {
			return derp.Wrap(err, location, "Unable to set value in target", runner.TargetPath)
		}
	}

	return nil
}

// itemElement returns the element of each item in the source at path: an Array's items, an
// Object's wildcard (or Any, for named properties), or Any beneath an Any
func itemElement(sourceSchema schema.Schema, path string) (schema.Element, error) {

	const location = "rosetta.translate.itemElement"

	element, ok := sourceSchema.GetElement(path)

	if !ok {
		return nil, derp.Internal(location, "Source element must exist in sourceSchema", path)
	}

	switch typed := element.(type) {

	case schema.Array:
		return typed.Items, nil

	case schema.Object:
		if typed.Wildcard != nil {
			return typed.Wildcard, nil
		}
		return schema.Any{}, nil

	case schema.Any:
		return schema.Any{}, nil
	}

	return nil, derp.Internal(location, "Source element must be a list or a map", path)
}

/******************************************
 * Serialization Methods
 ******************************************/

// MarshalMap converts this `forEach` rule back into its map representation.
func (runner forEachRunner) MarshalMap() map[string]any {
	return map[string]any{
		"forEach": runner.SourcePath,
		"target":  runner.TargetPath,
		"filter":  runner.FilterRaw,
		"rules":   runner.Pipeline.MarshalSliceOfMap(),
	}
}

// UnmarshalMap populates this `forEach` rule from its map representation.
func (runner *forEachRunner) UnmarshalMap(data mapof.Any) error {

	return runner.populate(
		data.GetString("forEach"),
		data.GetString("target"),
		data.GetString("filter"),
		data.GetSliceOfPlainMap("rules"),
	)
}

// populate assigns the parsed fields of this `forEach` rule, and is the one place that
// both UnmarshalMap and the package constructors go through.
func (runner *forEachRunner) populate(source string, target string, filter string, rules []map[string]any) error {

	// Populate Filter
	runner.Filter = nil
	if runner.FilterRaw != "" {
		filterTemplate, err := template.New("").Funcs(funcmap.All()).Parse(runner.FilterRaw)

		if err != nil {
			return derp.Wrap(err, "rosetta.translate.forEachRunner.populate", "Error parsing `filter` template", runner.FilterRaw)
		}
		runner.Filter = filterTemplate
	} else {
		runner.Filter = nil
	}

	// Parse Rules
	pipeline, err := NewFromMap(rules...)

	if err != nil {
		return derp.Wrap(err, "rosetta.translate.forEachRunner.populate", "Unable to create Pipeline", rules)
	}

	// Populate remaining fields
	runner.SourcePath = source
	runner.TargetPath = target
	runner.FilterRaw = filter
	runner.Pipeline = pipeline

	return nil
}
