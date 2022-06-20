package schema

import (
	"testing"

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

func TestArrayValidation(t *testing.T) {

	s := &Array{
		Items: &String{MaxLength: null.NewInt(10)},
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
