package pointer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPointerTo(t *testing.T) {

	value := map[string]string{
		"key": "value",
		"foo": "bar",
	}

	result := To(value)

	value["other"] = "thing"

	require.Equal(t, &value, result)
	require.Equal(t, "thing", (*result.(*map[string]string))["other"])
	require.Equal(t, "thing", value["other"])
}
