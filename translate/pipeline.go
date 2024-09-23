package translate

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// Pipeline represents a slice of Rule objects
type Pipeline []Rule

// New returns a new Pipeline object, populated with the provided rules
func New(rules ...Rule) Pipeline {
	return Pipeline(rules)
}

// NewFromMap parses a slice of mapof.Any objects into a Pipeline
func NewFromMap(rules ...mapof.Any) (Pipeline, error) {

	result := make(Pipeline, len(rules))

	for index, data := range rules {

		if err := result[index].UnmarshalMap(data); err != nil {
			return result, derp.Wrap(err, "mapper.NewPipeline", "Error unmarshalling rule", data)
		}
	}

	return result, nil
}

// NewFromJSON reads a JSON string and returns a Pipeline object
func NewFromJSON(jsonString string) (Pipeline, error) {
	result := make([]Rule, 0)

	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return result, derp.Wrap(err, "mapper.Parse", "Error parsing JSON", jsonString)
	}

	return result, nil
}

// Execute runs each of the rules in the Pipeline
func (pipeline Pipeline) Execute(inSchema schema.Schema, inObject any, outSchema schema.Schema, outObject any) error {

	for _, rule := range pipeline {
		if err := rule.Execute(inSchema, inObject, outSchema, outObject); err != nil {
			return derp.Wrap(err, "mapper.Pipeline.Execute", "Error executing rule")
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
