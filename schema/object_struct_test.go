package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStruct_Validate(t *testing.T) {

	type testStruct struct {
		Title    string `path:"title"`
		Content  string `path:"content"`
		Age      int    `path:"age"`
		Location struct {
			Latitude  float64 `path:"latitude"`
			Longitude float64 `path:"longitude"`
		} `path:"location"`
		Address struct {
			Address1 string `path:"address1"`
			Address2 string `path:"address2"`
			City     string `path:"city"`
			State    string `path:"state"`
			ZipCode  string `path:"zipCode"`
		} `path:"address"`
		Friends []string `path:"friends"`
	}

	schema := getTestSchema()

	// This should be valid
	{
		value := testStruct{
			Title:   "This Title is Valid",
			Content: "This Content is Valid",
			Age:     19,
		}

		require.Nil(t, schema.Validate(value))
	}

	// Missing required values
	{
		value := testStruct{
			Title:   "This Title is Valid",
			Content: "",
			Age:     21,
		}

		require.NotNil(t, schema.Validate(value))
		// t.Log(schema.Validate(value))
	}

	// Age does not meet minimum value
	{
		value := testStruct{
			Title:   "This Title is Valid",
			Content: "This is valid content",
			Age:     17,
		}

		require.NotNil(t, schema.Validate(value))
		// t.Log(schema.Validate(value))
	}
}

func TestStruct_Set(t *testing.T) {

	type testStruct struct {
		Title    string `path:"title"`
		Content  string `path:"content"`
		Age      int    `path:"age"`
		Location struct {
			Latitude  float64 `path:"latitude"`
			Longitude float64 `path:"longitude"`
		} `path:"location"`
		Address struct {
			Address1 string `path:"address1"`
			Address2 string `path:"address2"`
			City     string `path:"city"`
			State    string `path:"state"`
			ZipCode  string `path:"zipCode"`
		} `path:"address"`
		Friends []string `path:"friends"`
	}

	schema := getTestSchema()

	// This should be valid
	{
		value := testStruct{
			Title:   "This Title is Valid",
			Content: "This Content is Valid",
			Age:     19,
		}

		require.Nil(t, schema.Validate(value))

		// Set top-level fields
		require.Nil(t, schema.Set(&value, "content", "Mr. Fluffkins"))
		require.Equal(t, value.Content, "Mr. Fluffkins")

		// Set nested fields
		require.Nil(t, schema.Set(&value, "location.latitude", 123.456))
		require.Equal(t, 123.456, value.Location.Latitude)

		// Create/set array values
		require.Nil(t, schema.Set(&value, "friends.1", "Sara Mason"))
		require.Equal(t, "", value.Friends[0])
		require.Equal(t, "Sara Mason", value.Friends[1])

		// Verify other fields are unchanged.
		require.Equal(t, 19, value.Age)

		require.NotNil(t, schema.Set(&value, "invalid-field", "should break"))
	}
}
