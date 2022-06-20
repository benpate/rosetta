package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReflection_SetStruct(t *testing.T) {

	data := getTestPeople()

	{
		err := Set(&data, "0.name", "Jackson Browne")
		require.Nil(t, err)
		require.Equal(t, "Jackson Browne", data[0].Name)
	}

	{
		err := Set(data, "0.email", "jackson@browne.com")
		require.Nil(t, err)
		require.Equal(t, "jackson@browne.com", data[0].Email)
	}

	{
		err := Set(data, "0.relatives.0.name", "Neil Young")
		require.Nil(t, err)
		require.Equal(t, "Neil Young", data[0].Relatives[0].Name)
	}
}
