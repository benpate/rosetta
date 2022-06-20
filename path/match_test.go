package path

import (
	"testing"

	"github.com/benpate/exp"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	d := getTestStruct()

	require.True(t, Match(d, exp.Equal("name", "John Connor")))
	require.True(t, Match(d, exp.Equal("relatives.0.name", "Sarah Connor")))
	require.True(t, Match(d, exp.Equal("relatives.0.relatives.1.name", "Kyle Reese")))
}
