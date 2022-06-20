package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasUppercase_None(t *testing.T) {

	{
		value, err := HasUppercase("")("There are no uppercase")
		require.Nil(t, err)
		require.Equal(t, "There are no uppercase", value)
	}

	{
		value, err := HasUppercase("")("there are No uppercase")
		require.Nil(t, err)
		require.Equal(t, "there are No uppercase", value)
	}

	{
		value, err := HasUppercase("")("there are no uppercasE")
		require.Nil(t, err)
		require.Equal(t, "there are no uppercasE", value)
	}

	{
		value, err := HasUppercase("")("there are no uppercase")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}
}

func TestHasUppercase_One(t *testing.T) {

	{
		value, err := HasUppercase("1")("there are no uppercase")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasUppercase("1")("There is one uppercase")
		require.Nil(t, err)
		require.Equal(t, "There is one uppercase", value)
	}

	{
		value, err := HasUppercase("1")("there Is one uppercase")
		require.Nil(t, err)
		require.Equal(t, "there Is one uppercase", value)
	}

	{
		value, err := HasUppercase("1")("there is one uppercasE")
		require.Nil(t, err)
		require.Equal(t, "there is one uppercasE", value)
	}
}

func TestHasUppercase_Four(t *testing.T) {

	{
		value, err := HasUppercase("4")("there are no uppercase")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasUppercase("4")("There are no uppercase")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasUppercase("4")("There Are no uppercase")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasUppercase("4")("There aRe fouR uppercasE")
		require.Nil(t, err)
		require.Equal(t, "There aRe fouR uppercasE", value)
	}
}

func TestHasLowercase_None(t *testing.T) {

	{
		value, err := HasLowercase("")("THERE ARE NO LOWERCASE")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasLowercase("")("tHERE IS ONE LOWERCASE")
		require.Nil(t, err)
		require.Equal(t, "tHERE IS ONE LOWERCASE", value)
	}

	{
		value, err := HasLowercase("")("THERE IS ONE LoWERCASE")
		require.Nil(t, err)
		require.Equal(t, "THERE IS ONE LoWERCASE", value)
	}

	{
		value, err := HasLowercase("")("THERE ARe TWO LOWERcASE")
		require.Nil(t, err)
		require.Equal(t, "THERE ARe TWO LOWERcASE", value)
	}
}

func TestHasLowercase_One(t *testing.T) {

	{
		value, err := HasLowercase("1")("THERE ARE NO LOWERCASE")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasLowercase("1")("tHERE IS ONE LOWERCASE")
		require.Nil(t, err)
		require.Equal(t, "tHERE IS ONE LOWERCASE", value)
	}

	{
		value, err := HasLowercase("1")("THERE IS ONE LoWERCASE")
		require.Nil(t, err)
		require.Equal(t, "THERE IS ONE LoWERCASE", value)
	}

	{
		value, err := HasLowercase("1")("THERE ARe TWO LOWERcASE")
		require.Nil(t, err)
		require.Equal(t, "THERE ARe TWO LOWERcASE", value)
	}
}

func TestHasLowercase_Four(t *testing.T) {

	{
		value, err := HasLowercase("4")("THERE ARE NO LOWERCASE")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasLowercase("4")("THeRE IS ONE LOWERCASE")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasLowercase("4")("THErE ArE ThREE LOWERCASE")
		require.NotNil(t, err)
		require.Equal(t, "", value)
	}

	{
		value, err := HasLowercase("4")("ThERE ARE fOuR LoWERCASE")
		require.Nil(t, err)
		require.Equal(t, "ThERE ARE fOuR LoWERCASE", value)
	}
	{
		value, err := HasLowercase("4")("THeRe aRE SIx LOwERCASE")
		require.Nil(t, err)
		require.Equal(t, "THeRe aRE SIx LOwERCASE", value)
	}
}
