package translate

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/slice"
)

// Rule represents a single mapping rule
type Rule struct {
	Runner
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

		thenMap := slice.Map(data.GetSliceOfMap("then"), toPlainMap)
		elseMap := slice.Map(data.GetSliceOfMap("else"), toPlainMap)

		runner, err := newConditionRunner(condition, thenMap, elseMap)

		if err != nil {
			return derp.Wrap(err, location, "Error creating ConditionRunner", data)
		}

		rule.Runner = runner
		return nil
	}

	// ForEach Runner
	if sourcePath := data.GetString("forEach"); sourcePath != "" {

		targetPath := data.GetString("target")
		filter := data.GetString("filter")
		rulesMap := slice.Map(data.GetSliceOfMap("rules"), toPlainMap)

		runner, err := newForEachRunner(sourcePath, targetPath, filter, rulesMap)

		if err != nil {
			return derp.Wrap(err, location, "Error creating ForEachRunner", data)
		}

		rule.Runner = runner
		return nil
	}

	// From Runner
	if path, ok := data.GetStringOK("path"); ok {
		targetPath := data.GetString("target")
		rule.Runner = newPathRunner(path, targetPath)
		return nil
	}

	// Expression Runner
	if expression := data.GetString("expression"); expression != "" {

		targetPath := data.GetString("target")
		runner, err := newExpressionRunner(expression, targetPath)

		if err != nil {
			return derp.Wrap(err, location, "Error creating ExpressionRunner", expression)
		}

		rule.Runner = runner
		return nil
	}

	// Value Runner
	if value := data.GetAny("value"); value != nil {
		targetPath := data.GetString("target")
		rule.Runner = newValueRunner(value, targetPath)
		return nil
	}

	// Unrecognized Ruleset
	return derp.NewInternalError(location, "No valid runner found in rule", data)
}
