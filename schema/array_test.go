package schema

import (
	"reflect"
	"testing"

	"github.com/benpate/rosetta/maps"
	"github.com/benpate/rosetta/null"
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

func TestArrayGetElement(t *testing.T) {

	{
		s := New(Array{Items: String{}})
		element, err := s.GetElement("")
		require.Nil(t, err)
		require.Equal(t, s.Element, element)
	}

	{
		s := New(Array{Items: String{}})
		element, err := s.GetElement("0")
		require.Nil(t, err)
		require.Equal(t, s.Element.(Array).Items, element)
	}
}

func TestArrayValidation(t *testing.T) {

	s := New(Array{
		Items: String{MaxLength: 10},
	})

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
	s := New(Array{Items: String{}})

	v := []string{"zero", "one", "two", "three"}

	{
		result, err := s.Get(reflect.ValueOf(v), "0")
		require.Nil(t, err)
		require.Equal(t, "zero", result)
	}

	{
		// Test that negaitve indexes are handled correctly
		result, err := s.Get(reflect.ValueOf(v), "-1")
		require.NotNil(t, err)
		require.Equal(t, nil, result)
	}

	{
		// Test that negaitve indexes are handled correctly
		result, err := s.Get(reflect.ValueOf(v), "notanumber")
		require.NotNil(t, err)
		require.Equal(t, nil, result)
	}
}

func TestArrayNil1(t *testing.T) {

	s := Schema{
		Element: Array{
			Items: Integer{Minimum: null.NewInt64(0), Maximum: null.NewInt64(10)},
		},
	}

	value := make([]int, 0)

	element, err := s.GetElement("0")
	require.Equal(t, element, s.Element.(Array).Items)
	require.Nil(t, err)

	result, err := s.Get(value, "0")

	require.Equal(t, nil, result)
	require.NotNil(t, err)
}

func TestArrayNil2(t *testing.T) {

	s := Schema{
		Element: Array{
			Items: Object{
				Properties: ElementMap{
					"firstName": String{Default: "John"},
					"lastName":  String{Default: "Connor", Enum: []string{"Connor", "Jackson", "Jones", "Smith"}},
				},
			},
		},
	}

	value := make([]maps.Map, 0)

	element, err := s.GetElement("0.firstName")
	require.Equal(t, element, s.Element.(Array).Items.(Object).Properties["firstName"])
	require.Nil(t, err)

	result, err := s.Get(value, "0.firstName")
	require.Equal(t, nil, result)
	require.NotNil(t, err)
}

func TestComplexArrayOperations(t *testing.T) {

	schema := Schema{
		Element: Object{
			Properties: ElementMap{
				"name":  String{},
				"age":   Integer{},
				"email": String{},
				"notes": String{},
				"friends": Array{
					Items: Object{
						Properties: ElementMap{
							"name":  String{},
							"age":   Integer{},
							"email": String{},
						},
					},
				},
			},
		},
	}

	data := make(maps.Map)

	schema.Set(&data, "name", "John Connor")
	schema.Set(&data, "email", "john@connor.mil")
	schema.Set(&data, "age", 30)

	schema.Set(&data, "friends.0", maps.Map{
		"name":  "Sarah Connor",
		"age":   30,
		"email": "sarah@sky.net",
	})

	require.Equal(t, "John Connor", data["name"])
	require.Equal(t, "john@connor.mil", data["email"])
	require.Equal(t, int(30), data["age"])
	require.Equal(t, 1, len(data["friends"].([]maps.Map)))
	require.Equal(t, "Sarah Connor", data["friends"].([]maps.Map)[0]["name"])
	require.Equal(t, "sarah@sky.net", data["friends"].([]maps.Map)[0]["email"])
	require.Equal(t, 30, data["friends"].([]maps.Map)[0]["age"])

	notes, err := schema.Get(data, "notes")
	require.Equal(t, nil, data["notes"])
	require.Equal(t, "", notes)
	require.Nil(t, err)
}

func TestArrayRemoveSimple(t *testing.T) {
	s := New(Array{Items: String{}})

	// Remove first item
	{
		v := []string{"zero", "one", "two", "three"}

		err := s.Remove(&v, "0")
		require.Nil(t, err)
		require.Equal(t, 3, len(v))
		require.Equal(t, []string{"one", "two", "three"}, v)
	}

	// Remove last item
	{
		v := []string{"zero", "one", "two", "three"}

		err := s.Remove(&v, "3")
		require.Nil(t, err)
		require.Equal(t, 3, len(v))
		require.Equal(t, []string{"zero", "one", "two"}, v)
	}

	// Remove middle item
	{
		v := []string{"zero", "one", "two", "three"}

		err := s.Remove(&v, "1")
		require.Nil(t, err)
		require.Equal(t, 3, len(v))
		require.Equal(t, []string{"zero", "two", "three"}, v)
	}
}

func TestArrayRemoveComplex(t *testing.T) {

	type testPerson struct {
		Name     string       `path:"name"`
		Email    string       `path:"email"`
		Children []testPerson `path:"children"`
	}

	data := []testPerson{}

	schema := New(Array{Items: Object{Properties: ElementMap{
		"name":  String{},
		"email": String{},
		"children": Array{Items: Object{Properties: ElementMap{
			"name":  String{},
			"email": String{},
		}}},
	}}})

	// Alfred
	require.Nil(t, schema.Set(&data, "0.name", "Alfred"))
	require.Nil(t, schema.Set(&data, "0.email", "aelfred@wessex.gov"))

	// Aethelflad
	require.Nil(t, schema.Set(&data, "0.children.0.name", "Aethelflad"))
	require.Nil(t, schema.Set(&data, "0.children.0.email", "aethelflad@mercia.gov"))

	// Edward
	require.Nil(t, schema.Set(&data, "0.children.1.name", "Edward"))
	require.Nil(t, schema.Set(&data, "0.children.1.email", "slick-eddie@wessex.gov"))

	// Aelfgiffu
	require.Nil(t, schema.Set(&data, "1.name", "Aelfgiffu"))
	require.Nil(t, schema.Set(&data, "1.email", "aelfgiffu@wessex.gov"))

	// Verify inserts
	require.Equal(t, []testPerson{{
		Name:  "Alfred",
		Email: "aelfred@wessex.gov",
		Children: []testPerson{
			{
				Name:  "Aethelflad",
				Email: "aethelflad@mercia.gov",
			},
			{
				Name:  "Edward",
				Email: "slick-eddie@wessex.gov",
			},
		}},
		{
			Name:  "Aelfgiffu",
			Email: "aelfgiffu@wessex.gov",
		},
	}, data)

	require.Nil(t, schema.Remove(&data, "0.children.0"))

	// Verify deletes
	require.Equal(t, []testPerson{{
		Name:  "Alfred",
		Email: "aelfred@wessex.gov",
		Children: []testPerson{
			{
				Name:  "Edward",
				Email: "slick-eddie@wessex.gov",
			},
		}},
		{
			Name:  "Aelfgiffu",
			Email: "aelfgiffu@wessex.gov",
		},
	}, data)

	require.Nil(t, schema.Remove(&data, "0.children.0"))

	// Verify deletes
	require.Equal(t, []testPerson{
		{
			Name:     "Alfred",
			Email:    "aelfred@wessex.gov",
			Children: []testPerson{},
		},
		{
			Name:  "Aelfgiffu",
			Email: "aelfgiffu@wessex.gov",
		},
	}, data)
}
