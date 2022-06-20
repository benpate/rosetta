package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteMapOfString(t *testing.T) {

	value := map[string]string{
		"first":  "1",
		"second": "2",
		"third":  "3",
	}

	Delete(value, "second")

	require.Equal(t, value["first"], "1")
	require.Equal(t, value["third"], "3")
	require.Empty(t, value["second"])
}

func TestDeleteMapOfInterface(t *testing.T) {

	value := map[string]interface{}{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	Delete(value, "second")

	require.Equal(t, value["first"], 1)
	require.Equal(t, value["third"], 3)
	require.Nil(t, value["second"])
}
