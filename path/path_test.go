package path

import (
	"testing"

	"github.com/benpate/rosetta/list"
	"github.com/stretchr/testify/require"
)

func TestProperties(t *testing.T) {

	d := getTestData()

	{
		value := Get(d, "name")
		require.Equal(t, "John Connor", value)
	}

	{
		value := Get(d, "email")
		require.Equal(t, "john@connor.mil", value)
	}

	{
		value := Get(d, "missing property")
		require.Nil(t, value)
	}
}

func TestSubProperties(t *testing.T) {

	d := getTestData()

	{
		value := Get(d, "relatives.mom")
		require.Equal(t, "Sarah Connor", value)
	}

	{
		value := Get(d, "relatives.dad")
		require.Equal(t, "Kyle Reese", value)
	}

	{
		value := Get(d, "relatives.sister")
		require.Nil(t, value)
	}
}

func TestArrays(t *testing.T) {

	d := getTestData()

	{
		value := Get(d, "enemies.0")
		require.Equal(t, "T-1000", value)
	}

	{
		value := Get(d, "enemies.1")
		require.Equal(t, "T-3000", value)
	}

	{
		value := Get(d, "enemies.2")
		require.Equal(t, "T-5000", value)
	}

	{
		value := Get(d, "enemies.-1")
		require.Nil(t, value)
	}

	{
		value := Get(d, "enemies.3")
		require.Nil(t, value)
	}

	{
		value := Get(d, "enemies.100000")
		require.Nil(t, value)
	}

	{
		value := Get(d, "enemies.fred")
		require.Nil(t, value)
	}
}

func TestError(t *testing.T) {

	{
		value := Get("unsupported data", "property")
		require.Nil(t, value)
	}

	{
		value := Get("string at the end of a path", "")
		require.Equal(t, "string at the end of a path", value)
	}
}

func TestGetter(t *testing.T) {

	d := getTestStruct()

	{
		value := Get(d, "name")
		require.Equal(t, "John Connor", value)
	}

	{
		value := Get(d, "email")
		require.Equal(t, "john@connor.mil", value)
	}

	{
		value := Get(d, "relatives.0.name")
		require.Equal(t, "Sarah Connor", value)
	}

	{
		value := Get(d, "relatives.1.relatives.1.name")
		require.Equal(t, "Sarah Connor", value)
	}

	{
		value := Get(d, "missing-property")
		require.Nil(t, value)
	}
}

func TestSet(t *testing.T) {

	d := getTestData()

	// Not implemented yet, so this should just error
	require.NotNil(t, Set(d, "anywhere.doesnt.matter", "1"))
}

func TestParseIntInRange(t *testing.T) {

	// Test valid index value
	{
		p := "7"
		index, error := Index(p, 10)
		require.Equal(t, 7, index)
		require.Nil(t, error)
	}

	// Test no maximum value
	{
		p := "7"
		index, error := Index(p, -1)
		require.Equal(t, 7, index)
		require.Nil(t, error)
	}

	// Test underflow minimum
	{
		p := "-1"
		index, error := Index(p, 5)
		require.Equal(t, 0, index)
		require.NotNil(t, error)
	}

	// Test overflow maximum
	{
		p := "7"
		index, error := Index(p, 5)
		require.Equal(t, 0, index)
		require.NotNil(t, error)
	}

	// Test non-integer
	{
		p := "non-integer"
		index, error := Index(p, 5)
		require.Equal(t, 0, index)
		require.NotNil(t, error)
	}

}

/********************************
 * SUPPORT FUNCS
 ********************************/

type testStruct struct {
	Name      string
	Email     string
	Relatives testStructArray
}

func (d testStruct) GetPath(path string) (interface{}, bool) {

	if path == "" {
		return d, true
	}

	head, tail := list.Split(path, ".")

	switch head {

	case "name":
		return d.Name, true

	case "email":
		return d.Email, true

	case "relatives":
		return d.Relatives.GetPath(tail)
	}

	return nil, false
}

type testStructArray []testStruct

func (d testStructArray) GetPath(path string) (interface{}, bool) {

	if path == "" {
		return d, true
	}

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(d))

	if err != nil {
		return nil, false
	}

	return d[index].GetPath(tail)
}

func getTestStruct() testStruct {

	return testStruct{
		Name:  "John Connor",
		Email: "john@connor.mil",
		Relatives: testStructArray{
			{
				Name:  "Sarah Connor",
				Email: "sarah@sky.net",
				Relatives: testStructArray{
					{Name: "John Connor"},
					{Name: "Kyle Reese"},
				},
			},
			{
				Name:  "Kyle Reese",
				Email: "kyle@resistance.mil",
				Relatives: testStructArray{
					{Name: "John Connor"},
					{Name: "Sarah Connor"},
				},
			},
		},
	}
}

func getTestData() map[string]interface{} {

	return map[string]interface{}{
		"name":  "John Connor",
		"email": "john@connor.mil",
		"relatives": map[string]interface{}{
			"mom": "Sarah Connor",
			"dad": "Kyle Reese",
		},
		"enemies": []interface{}{"T-1000", "T-3000", "T-5000"},
	}
}
