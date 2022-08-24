package schema

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/maps"
	"github.com/stretchr/testify/require"
)

func TestSchema(t *testing.T) {

	s := getTestSchema()

	require.Equal(t, s.ID, "https://pate.org/example/article")
	require.Equal(t, s.Comment, "I had to copy this one to make it work right.")
}

func TestEmptySchema(t *testing.T) {

	var s Schema

	object := make(map[string]string)

	// This should throw an error because the schema has no root-level element
	require.NotNil(t, s.Validate(object))

	// Try marshalling the empty schema to JSON
	schemaJSON, err := json.Marshal(s)

	require.Nil(t, err)
	require.Equal(t, []byte("null"), schemaJSON)

}

func TestSet(t *testing.T) {

	s := getTestSchema()

	object := maps.Map{}

	// Test setting values using the schema
	err := s.Set(&object, "title", "This is the title")
	require.Nil(t, err)
	require.Equal(t, "This is the title", object["title"])

	err = s.Set(&object, "content", "This is the content")
	require.Nil(t, err)
	require.Equal(t, "This is the content", object["content"])

	err = s.Set(&object, "age", 21)
	require.Nil(t, err)
	require.Equal(t, int64(21), object["age"])

	// Test values that will not get set (and should return an error)
	err = s.Set(&object, "this-path-doesn't-exist", "so it won't get set")
	require.NotNil(t, err)

	err = s.Set(&object, "age", "this is not an integer")
	require.NotNil(t, err)

	// Verify object contents now.
	// require.Equal(t, 3, len(object))
	require.Equal(t, "This is the title", object["title"])
	require.Equal(t, "This is the content", object["content"])
	require.Equal(t, int64(21), object["age"])
}

func getTestSchema() Schema {

	var result Schema

	data := []byte(`{
		"$id": "https://pate.org/example/article",
		"$comment" : "I had to copy this one to make it work right.",
		"title": "Article",
		"type": "object",
		"properties": {
			"title": {
				"type": "string",
				"required": true
			},
			"content": {
				"type": "string",
				"required": true
			},
			"age": {
				"description": "Age in years",
				"type": "integer",
				"minimum": 18
			},
			"location": {
				"type": "object",
				"properties":{
					"latitude" : {"type":"number"},
					"longitude": {"type":"number"}
				}
			},
			"friends": {
			  "type" : "array",
			  "items" : { "type" : "string"}
			},
			"address": {
				"type": "object",
				"properties": {
					"address1": {"type": "string", "$id":"addr1"},
					"address2": {"type": "string", "$id":"addr2"},
					"city": {"type": "string", "$id":"city"},
					"state": {"type": "string", "$id":"state"},
					"zipCode": {"type": "string", "$id":"zip"}
				}
			}
		},
		"required": ["title", "content"]
	  }`)

	if err := json.Unmarshal(data, &result); err != nil {
		panic(err.Error())
	}

	return result
}
