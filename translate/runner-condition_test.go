package translate

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalCondition_If(t *testing.T) {

	testJSON := `[
		{"if":"{{eq 42 .value}}", "then": [
			{"target":"name", "value":"forty two"},
			{"target":"age", "path":"value"}
		], "else": [
			{"target":"name", "value":"not forty two"},
			{"target":"age", "path":"value"}
		]}
	]`

	testJSONBytes := []byte(testJSON)

	pipeline := Pipeline{}

	if err := json.Unmarshal(testJSONBytes, &pipeline); err != nil {
		t.Error(err)
	}

	require.NotNil(t, pipeline[0].Runner.(conditionRunner).Condition)
	require.Equal(t, 2, len(pipeline[0].Runner.(conditionRunner).ThenRules))
	require.Equal(t, 2, len(pipeline[0].Runner.(conditionRunner).ElseRules))

	testData := mapof.Any{
		"label": "John Connor",
		"value": 42,
	}

	result := mapof.Any{}

	err := pipeline.Execute(schema.Wildcard(), testData, schema.Wildcard(), &result)
	require.Nil(t, err)

	require.Equal(t, 2, len(result))
	require.Equal(t, "forty two", result.GetString("name"))
	require.Equal(t, 42, result.GetInt("age"))
}

func TestUnmarshalCondition_Else(t *testing.T) {

	testJSON := `[
		{"if":"{{eq 42 .value}}", "then": [
			{"target":"name", "value":"forty two"}
		], "else": [
			{"target":"name", "value":"not forty two"}
		]}
	]`

	testJSONBytes := []byte(testJSON)

	pipeline := Pipeline{}

	if err := json.Unmarshal(testJSONBytes, &pipeline); err != nil {
		t.Error(err)
	}

	require.NotNil(t, pipeline[0].Runner.(conditionRunner).Condition)

	testData := mapof.Any{
		"label": "Sarah Connor",
		"value": 43,
	}

	err := pipeline.Execute(schema.Wildcard(), testData, schema.Wildcard(), &testData)
	require.Nil(t, err)

	require.Equal(t, 3, len(testData))
	require.Equal(t, "Sarah Connor", testData.GetString("label"))
	require.Equal(t, 43, testData.GetInt("value"))
	require.Equal(t, "not forty two", testData.GetString("name"))

}
