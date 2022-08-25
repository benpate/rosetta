package schema

import (
	"testing"

	"github.com/benpate/rosetta/maps"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {

	schema := getTestSchema()

	object := schema.Element.(Object)

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
		require.Equal(t, "1600 Pennsylvania Avenue", data["address"].(maps.Map)["address1"])
	}
}

func TestComplexObjectSet(t *testing.T) {

	var err error

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
				"friends": Array{Items: Object{
					Properties: map[string]Element{
						"name": String{},
						"age":  Integer{},
						"friends": Array{Items: Object{
							Properties: map[string]Element{
								"name": String{},
								"age":  Integer{},
							},
						}},
					},
				}},
			},
		},
	}

	var data interface{}

	err = schema.Set(&data, "contact.email", "john@doemain.com")
	require.Nil(t, err)

	require.Equal(t, maps.Map{
		"name": "",
		"age":  int64(0),
		"contact": maps.Map{
			"domain": "",
			"email":  "john@doemain.com",
			"phone":  "",
		},
		"friends": []maps.Map{},
	}, data)
}

func TestComplexObjectSet2(t *testing.T) {

	var err error

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
				"friends": Array{Items: Object{
					Properties: map[string]Element{
						"name": String{},
						"age":  Integer{},
						"friends": Array{Items: Object{
							Properties: map[string]Element{
								"name": String{},
								"age":  Integer{},
							},
						}},
					},
				}},
			},
		},
	}

	data := maps.Map{
		"contact": maps.Map{
			"domain": "doemain.com",
		},
	}

	err = schema.Set(&data, "name", "John Doe")
	require.Nil(t, err)

	err = schema.Set(&data, "age", 42)
	require.Nil(t, err)

	err = schema.Set(&data, "contact.email", "john@doe.com")
	require.Nil(t, err)

	err = schema.Set(&data, "contact.phone", "123-456-7890")
	require.Nil(t, err)

	err = schema.Set(&data, "contact.domain", "doe.com")
	require.Nil(t, err)

	err = schema.Set(&data, "friends.2.friends.2.name", "John Doe")
	require.Nil(t, err)

	err = schema.Set(&data, "contact.badValue", "We don't talk about bruno.")
	require.NotNil(t, err)

	require.Equal(t, maps.Map{
		"name": "John Doe",
		"age":  int64(42),
		"contact": maps.Map{
			"email":  "john@doe.com",
			"phone":  "123-456-7890",
			"domain": "doe.com",
		},
		"friends": []maps.Map{
			nil,
			nil,
			{
				"name": "",
				"age":  int64(0),
				"friends": []maps.Map{
					nil,
					nil,
					{
						"name": "John Doe",
						"age":  int64(0),
					},
				},
			},
		},
	}, data)
}

func TestObjectRemove(t *testing.T) {

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

	data := maps.Map{
		"name": "John Doe",
		"age":  42,
		"contact": map[string]any{
			"email": "john@doe.com",
		},
	}

	// Remove non-existing property
	err := schema.Remove(&data, "contact.bad-property")
	require.NotNil(t, err)
	require.Equal(t, maps.Map{
		"name": "John Doe",
		"age":  42,
		"contact": map[string]any{
			"email": "john@doe.com",
		},
	}, data)

	// Remove non-existing key
	err = schema.Remove(&data, "contact.phone")
	require.Nil(t, err)
	require.Equal(t, maps.Map{
		"name": "John Doe",
		"age":  42,
		"contact": map[string]any{
			"email": "john@doe.com",
		},
	}, data)

	// Remove nested key
	err = schema.Remove(&data, "contact.email")
	require.Nil(t, err)
	require.Equal(t, maps.Map{
		"name":    "John Doe",
		"age":     42,
		"contact": map[string]any{},
	}, data)

	// Remove top-level key
	err = schema.Remove(&data, "contact")
	require.Nil(t, err)
	require.Equal(t, maps.Map{
		"name": "John Doe",
		"age":  42,
	}, data)

	// Remove top-level key
	err = schema.Remove(&data, "age")
	require.Nil(t, err)
	require.Equal(t, maps.Map{
		"name": "John Doe",
	}, data)

	// Remove top-level key
	err = schema.Remove(&data, "name")
	require.Nil(t, err)
	require.Equal(t, maps.Map{}, data)

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

func TestGroupID(t *testing.T) {

	spew.Config.DisableMethods = true

	value := maps.Map{
		"rule": "private",
		"groupIds": []string{
			"5fafafafafafafafafafaf",
			"5fbfbfbfbfbfbfbfbfbfbf",
			"5fcfcfcfcfcfcfcfcfcfcf",
		},
	}

	s := Schema{
		Element: Object{
			Properties: map[string]Element{
				"rule":     String{Default: "anonymous"},
				"groupIds": Array{Items: String{Format: "objectId"}},
			},
		},
	}

	{
		result, element, err := s.Get(value, "rule")
		require.Nil(t, err)
		require.NotNil(t, element)
		require.Equal(t, "private", result)
	}

	{
		result, element, err := s.Get(value, "groupIds")
		require.Nil(t, err)
		require.NotNil(t, element)
		require.Equal(t, []string{"5fafafafafafafafafafaf", "5fbfbfbfbfbfbfbfbfbfbf", "5fcfcfcfcfcfcfcfcfcfcf"}, result)
	}
}
