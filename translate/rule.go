package translate

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
)

// Rule represents a single mapping rule
type Rule struct {
	Runner
}

func (rule *Rule) MarshalJSON() ([]byte, error) {

	if rule.Runner == nil {
		return []byte("null"), nil
	}

	return json.Marshal(rule.Runner)
}

func (rule *Rule) MarshalMap() map[string]any {

	if rule.Runner == nil {
		return make(map[string]any)
	}

	return rule.Runner.MarshalMap()
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (rule *Rule) UnmarshalJSON(data []byte) error {

	const location = "rosetta.translate.Rule.UnmarshalJSON"

	temp := mapof.NewAny()

	if err := json.Unmarshal(data, &temp); err != nil {
		return derp.Wrap(err, location, "Error unmarshalling JSON", string(data))
	}

	return rule.UnmarshalMap(temp)
}

// UnmarshalMap populates this object from a mapof.Any
func (rule *Rule) UnmarshalMap(data mapof.Any) error {

	const location = "rosetta.translate.Rule.UnmarshalMap"

	// Condition Runner
	if condition := data.GetString("if"); condition != "" {

		runner := conditionRunner{}

		if err := runner.UnmarshalMap(data); err != nil {
			return derp.Wrap(err, location, "Error unmarshalling ConditionRunner", data)
		}

		rule.Runner = runner
		return nil
	}

	// ForEach Runner
	if sourcePath := data.GetString("forEach"); sourcePath != "" {

		runner := forEachRunner{}

		if err := runner.UnmarshalMap(data); err != nil {
			return derp.Wrap(err, location, "Error unmarshalling ForEachRunner", data)
		}

		rule.Runner = runner
		return nil
	}

	// From Runner
	if path := data.GetString("path"); path != "" {

		runner := pathRunner{}

		if err := runner.UnmarshalMap(data); err != nil {
			return derp.Wrap(err, location, "Error unmarshalling ForEachRunner", data)
		}

		rule.Runner = runner
		return nil
	}

	// Expression Runner
	if expression := data.GetString("expression"); expression != "" {

		runner := expressionRunner{}

		if err := runner.UnmarshalMap(data); err != nil {
			return derp.Wrap(err, location, "Error unmarshalling ExpressionRunner", data)
		}

		rule.Runner = runner
		return nil
	}

	// Value Runner
	if value := data.GetAny("value"); value != nil {

		runner := valueRunner{}

		if err := runner.UnmarshalMap(data); err != nil {
			return derp.Wrap(err, location, "Error unmarshalling ValueRunner", data)
		}

		rule.Runner = runner
		return nil
	}

	// Unrecognized Ruleset
	return derp.NewInternalError(location, "No valid runner found in rule", data)
}
