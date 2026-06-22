package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalMap_AllTypes(t *testing.T) {
	cases := map[string]any{
		"any":     Any{},
		"array":   Array{},
		"boolean": Boolean{},
		"integer": Integer{},
		"number":  Number{},
		"object":  Object{},
		"string":  String{},
	}

	for typeName := range cases {
		data := map[string]any{"type": typeName}

		// "array" requires an "items" definition to unmarshal successfully
		if typeName == "array" {
			data["items"] = map[string]any{"type": "string"}
		}

		element, err := UnmarshalMap(data)
		require.NoError(t, err, typeName)
		require.NotNil(t, element, typeName)
	}
}

func TestUnmarshalMap_Any(t *testing.T) {
	// UnmarshalMap dispatches the "any" type to an Any element
	element, err := UnmarshalMap(map[string]any{"type": "any"})
	require.NoError(t, err)
	require.IsType(t, Any{}, element)
}

func TestUnmarshalMap_Nil(t *testing.T) {
	_, err := UnmarshalMap(nil)
	require.Error(t, err)
}

func TestUnmarshalMap_NotAMap(t *testing.T) {
	_, err := UnmarshalMap("not-a-map")
	require.Error(t, err)
}

func TestUnmarshalMap_UnknownType(t *testing.T) {
	_, err := UnmarshalMap(map[string]any{"type": "not-a-real-type"})
	require.Error(t, err)
}

func TestUnmarshalJSON_Element(t *testing.T) {
	element, err := UnmarshalJSON([]byte(`{"type":"string"}`))
	require.NoError(t, err)
	require.IsType(t, String{}, element)
}

func TestUnmarshalJSON_Element_InvalidJSON(t *testing.T) {
	_, err := UnmarshalJSON([]byte(`{ not valid json`))
	require.Error(t, err)
}

func TestUnmarshalJSON_Element_UnknownType(t *testing.T) {
	_, err := UnmarshalJSON([]byte(`{"type":"not-a-real-type"}`))
	require.Error(t, err)
}

func TestSchema_MarshalJSON(t *testing.T) {
	schema := New(String{})
	bytes, err := schema.MarshalJSON()

	require.NoError(t, err)
	require.JSONEq(t, `{"type":"string","required":false}`, string(bytes))
}

func TestSchema_MarshalJSON_Nil(t *testing.T) {
	bytes, err := Schema{}.MarshalJSON()

	require.NoError(t, err)
	require.Equal(t, "null", string(bytes))
}

func TestSchema_MarshalMap_WithMetadata(t *testing.T) {
	schema := Schema{ID: "my-id", Comment: "my-comment", Element: String{}}
	result := schema.MarshalMap()

	require.Equal(t, "my-id", result["$id"])
	require.Equal(t, "my-comment", result["$comment"])
}

func TestSchema_MarshalMap_Nil(t *testing.T) {
	require.Equal(t, map[string]any{}, Schema{}.MarshalMap())
}

func TestSchema_UnmarshalJSON(t *testing.T) {
	var schema Schema
	err := schema.UnmarshalJSON([]byte(`{"$id":"abc","type":"string"}`))

	require.NoError(t, err)
	require.Equal(t, "abc", schema.ID)
	require.IsType(t, String{}, schema.Element)
}

func TestSchema_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var schema Schema
	require.Error(t, schema.UnmarshalJSON([]byte(`{ not valid json`)))
}

func TestSchema_UnmarshalJSON_InvalidElement(t *testing.T) {
	var schema Schema
	require.Error(t, schema.UnmarshalJSON([]byte(`{"type":"not-a-real-type"}`)))
}
