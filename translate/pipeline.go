package translate

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema"
)

// Pipeline represents a slice of Rule objects
type Pipeline []Rule

// New returns a new Pipeline object, populated with the provided rules
func New(rules ...Rule) Pipeline {
	return Pipeline(rules)
}

// NewFromMap parses a slice of mapof.Any objects into a Pipeline
func NewFromMap(rules ...map[string]any) (Pipeline, error) {

	result := make(Pipeline, len(rules))

	for index, data := range rules {

		if err := result[index].UnmarshalMap(data); err != nil {
			return result, derp.Wrap(err, "rosetta.translate.NewPipeline", "Error unmarshalling rule", data)
		}
	}

	return result, nil
}

// NewSliceOfPipelines parses a slice of maps into a slice of Pipelines.
// If any of the maps DOES NOT represent a valid Pipeline, then the function will return an error.
func NewSliceOfPipelines(slices [][]map[string]any) ([]Pipeline, error) {

	const location = "rosetta.translate.NewSliceOfPipelines"

	result := make([]Pipeline, len(slices))

	for index, item := range slices {

		pipeline, err := NewFromMap(item...)

		if err != nil {
			return result, derp.Wrap(err, location, "Unable to create metadata pipeline", item)
		}

		result[index] = pipeline
	}

	return result, nil
}

// NewFromJSON reads a JSON string and returns a Pipeline object
func NewFromJSON(jsonString string) (Pipeline, error) {
	result := make([]Rule, 0)

	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return result, derp.Wrap(err, "rosetta.translate.Parse", "Error parsing JSON", jsonString)
	}

	return result, nil
}

// Execute runs each of the rules in the Pipeline
func (pipeline Pipeline) Execute(inSchema schema.Schema, inObject any, outSchema schema.Schema, outObject any) error {

	for _, rule := range pipeline {
		if err := rule.Execute(inSchema, inObject, outSchema, outObject); err != nil {
			return derp.Wrap(err, "rosetta.translate.Pipeline.Execute", "Error executing rule")
		}
	}

	return nil
}

func (pipeline Pipeline) IsEmpty() bool {
	return len(pipeline) == 0
}

func (pipeline Pipeline) NotEmpty() bool {
	return len(pipeline) > 0
}

/******************************************
 * Serialization Methods
 ******************************************/

func (pipeline Pipeline) MarshalJSON() ([]byte, error) {
	return json.Marshal(pipeline.MarshalSliceOfMap())
}

func (pipeline Pipeline) MarshalSliceOfMap() []map[string]any {

	result := make([]map[string]any, len(pipeline))

	for index, rule := range pipeline {
		result[index] = rule.MarshalMap()
	}

	return result
}
