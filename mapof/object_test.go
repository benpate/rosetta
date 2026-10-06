package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObjectConversion(t *testing.T) {
	var value map[string][]string = make(Object[[]string])
	require.NotNil(t, value)
}

func TestObjectZero(t *testing.T) {
	var value = NewObject[string]()
	require.Equal(t, "", value["this-is-zero"])
	require.True(t, value.IsZeroValue("this-is-zero"))

	value["exists"] = "true"
	require.False(t, value.IsZeroValue("exists"))
}
