package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMap(t *testing.T) {

	value := map[string]any{
		"one":   1,
		"two":   "two",
		"three": []string{"a", "b", "c"},
	}

	s := New(Object{
		Properties: ElementMap{
			"one":   Integer{},
			"two":   String{},
			"three": Array{Items: String{}},
			"four":  String{},
		}},
	)

	{
		one, err := s.Get(value, "one")
		require.Nil(t, err)
		require.Equal(t, int(1), one)
	}

	{
		two, err := s.Get(value, "two")
		require.Nil(t, err)
		require.Equal(t, "two", two)
	}

	{
		three, err := s.Get(value, "three")
		require.Nil(t, err)
		require.Equal(t, []string{"a", "b", "c"}, three)
	}

	{
		four, err := s.Get(value, "four")
		require.Nil(t, err)
		require.Equal(t, "", four)
	}

	{
		five, err := s.Get(value, "five")
		require.NotNil(t, err)
		require.Equal(t, nil, five)
	}
}

func TestGetMap_Missing(t *testing.T) {

	value := map[string]any{
		"one":   1,
		"two":   "two",
		"three": []string{"a", "b", "c"},
	}

	s := New(Object{
		Properties: ElementMap{
			"one":   Integer{},
			"two":   String{},
			"three": Array{Items: String{}},
			"four":  String{},
		}},
	)

	{
		five, err := s.Get(value, "five")
		require.NotNil(t, err)
		require.Equal(t, nil, five)
	}
}
