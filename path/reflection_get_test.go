package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testPerson struct {
	Name      string       `path:"name"`
	Email     string       `path:"email"`
	Relatives []testPerson `path:"relatives"`
	Age       int
}

func getTestPeople() []testPerson {

	return []testPerson{{
		Name:  "Michael Jackson",
		Email: "michael@jackson.com",
		Age:   24,
		Relatives: []testPerson{{
			Name:  "Tito",
			Age:   23,
			Email: "tito@jackson.com",
		}, {
			Name:  "Janet",
			Age:   22,
			Email: "janet@jackson.com",
		}},
	}, {
		Name:  "Andrew Jackson",
		Email: "andrew@jackson.com",
	}, {
		Name:  "Kendall Jackson",
		Email: "kendall@jackson.com",
	}}
}

func TestReflection(t *testing.T) {

	data := getTestPeople()

	// Test Struct First
	require.Equal(t, "Andrew Jackson", Get(data[1], "name"))

	// Test Array First
	require.Equal(t, "Michael Jackson", Get(data, "0.name"))
	require.Equal(t, "janet@jackson.com", Get(data, "0.relatives.1.email"))

	// Test Errors
	require.Nil(t, Get(data, "x.name"))     // Not integer index
	require.Nil(t, Get(data, "-1.name"))    // Negative index
	require.Nil(t, Get(data, "0.location")) // Invalid field name
	require.Nil(t, Get(data, "0.age"))      // Unexported field
	require.Nil(t, Get(data, "0.location")) // Invalid field name
}

func TestReflection_Pointer(t *testing.T) {

	data := getTestPeople()
	require.Equal(t, "Kendall Jackson", Get(&data, "2.name"))
}

func TestReflection_Arrays(t *testing.T) {

	value := [5]int{9, 8, 7, 6, 5}

	require.Equal(t, 9, Get(value, "0"))
	require.Equal(t, 8, Get(value, "1"))
	require.Equal(t, 7, Get(value, "2"))
	require.Equal(t, 6, Get(value, "3"))
	require.Equal(t, 5, Get(value, "4"))

	// Out of bounds errors
	require.Nil(t, Get(value, "-1"))
	require.Nil(t, Get(value, "5"))
	require.Nil(t, Get(value, "station"))
}
