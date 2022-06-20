package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {

	schema := getTestSchema()

	object := schema.Element.(*Object)

	assert.NotNil(t, object)
}

func TestSetComplexObject(t *testing.T) {

	schema := getTestSchema()

	data := map[string]any{}

	{
		err := schema.Set(&data, "title", "John Doe")
		require.Nil(t, err)
		require.Equal(t, "John Doe", data["title"])
	}

	{
		err := schema.Set(&data, "address.address1", "1600 Pennsylvania Avenue")
		require.Nil(t, err)
		require.Equal(t, "1600 Pennsylvania Avenue", data["address"].(map[string]any)["address1"])
	}
}

func TestValidateComplexObject(t *testing.T) {

	schema := Schema{
		Element: Object{
			Properties: map[string]Element{
				"name": String{},
				"age":  Integer{},
				"contact": Object{
					Properties: map[string]Element{
						"email":  String{},
						"phone":  String{},
						"domain": String{},
					},
				},
			},
		},
	}

	data := map[string]any{
		"name":     "John Connor",
		"age":      42,
		"badValue": "very bad value",
		"contact": map[string]any{
			"email":    "john@connor.com",
			"phone":    "123-456-7890",
			"domain":   "connor.com",
			"badValue": "We don't talk about bruno.",
		},
	}

	{
		err := schema.Validate(data)
		require.Nil(t, err)
	}

	{
		require.Equal(t, "John Connor", data["name"])
		require.Equal(t, int(42), data["age"])
		// require.Nil(t, data["badValue"])
		// require.Nil(t, data["neverPutThisValueIn"])

		innerMap := data["contact"].(map[string]any)
		require.Equal(t, "john@connor.com", innerMap["email"])
		require.Equal(t, "123-456-7890", innerMap["phone"])
		require.Equal(t, "connor.com", innerMap["domain"])
		// require.Nil(t, innerMap["badValue"])
	}
}
