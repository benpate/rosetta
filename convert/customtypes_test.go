package convert

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

type customA []any

type customM map[string]any

func TestCustomTypes_SliceOfMap(t *testing.T) {

	value := customA{
		customM{
			"first":  1,
			"second": 2,
			"third":  3,
		},
		customM{
			"first":  4,
			"second": 5,
			"third":  6,
		},
		customM{
			"first":  7,
			"second": 8,
			"third":  9,
		},
	}

	result := SliceOfMap(value)

	require.Equal(t, 3, len(result))
	spew.Dump(result)
}

func TestCustomTypes_MapOfAny(t *testing.T) {

	value := customM{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	result := MapOfAny(value)
	require.Equal(t, 3, len(value))
	spew.Dump(result)
}

func TestCustomTypes_MapOfString(t *testing.T) {

	value := customM{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	result := MapOfString(value)
	require.Equal(t, 3, len(value))
	spew.Dump(result)
}

func TestCustomTypes_MapOfInt(t *testing.T) {

	value := customM{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	result := MapOfInt(value)
	require.Equal(t, 3, len(value))
	spew.Dump(result)
}

func TestCustomTypes_MapOfInt32(t *testing.T) {

	value := customM{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	result := MapOfInt32(value)
	require.Equal(t, 3, len(value))
	spew.Dump(result)
}
