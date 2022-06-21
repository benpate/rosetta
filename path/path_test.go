package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {

	head, tail := Split("head.tail.tail.tail")

	require.Equal(t, "head", head)
	require.Equal(t, "tail.tail.tail", tail)
}

func TestIndex(t *testing.T) {

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

	head, tail := Split(path)

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

	head, tail := Split(path)
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
