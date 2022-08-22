package schema

import (
	"reflect"
	"testing"

	"github.com/benpate/rosetta/maps"
	"github.com/benpate/rosetta/null"
	"github.com/benpate/rosetta/slice"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestArrayCreate(t *testing.T) {

	value := map[string]any{}

	schema := getTestSchema()

	err := schema.Set(&value, "friends.1", "Sarah Connor")

	require.Nil(t, err)
	require.Equal(t, "", value["friends"].([]string)[0])
	require.Equal(t, "Sarah Connor", value["friends"].([]string)[1])
}

func TestArrayValidation(t *testing.T) {

	s := Array{
		Items: String{MaxLength: null.NewInt(10)},
	}

	{
		v := []string{"one", "two", "three", "valid"}
		require.Nil(t, s.Validate(v))
	}

	{
		v := []string{"one", "two", "three", "invalid because its way too long"}

		err := s.Validate(v)
		require.NotNil(t, err)
	}

	{
		err := s.Validate(17)
		require.NotNil(t, err)
	}
}

func TestArrayGet(t *testing.T) {
	s := &Array{Items: String{}}

	v := []string{"zero", "one", "two", "three"}

	{
		result, _, err := s.Get(reflect.ValueOf(v), "0")
		require.Nil(t, err)
		require.Equal(t, "zero", result)
	}

	{
		// Test that negaitve indexes are handled correctly
		result, _, err := s.Get(reflect.ValueOf(v), "-1")
		require.NotNil(t, err)
		require.Nil(t, result)
	}

	{
		// Test that negaitve indexes are handled correctly
		result, _, err := s.Get(reflect.ValueOf(v), "notanumber")
		require.NotNil(t, err)
		require.Nil(t, result)
	}
}

func TestComplexArrayOperations(t *testing.T) {

	schema := Schema{
		Element: Object{
			Properties: ElementMap{
				"name":  String{},
				"age":   Number{},
				"email": String{},
				"friends": Array{
					Items: Object{
						Properties: ElementMap{
							"name":  String{},
							"age":   Number{},
							"email": String{},
						},
					},
				},
			},
		},
	}

	data := make(maps.Map)

	schema.Set(&data, "name", "John Connor")
	schema.Set(&data, "age", 30)
	schema.Set(&data, "email", "john@connor.mil")

	friend := maps.Map{
		"name":  "Sarah Connor",
		"age":   30,
		"email": "sarah@sky.net",
	}

	friends, _, _ := slice.AnyAppend(data["friends"], friend)

	schema.Set(&data, "friends", friends)

	spew.Dump(data)
}
