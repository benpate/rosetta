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

	s := New(ElementMap{
		"one":   Integer{},
		"two":   String{},
		"three": Array{Items: String{}},
		"four":  Any{},
	})

	{
		one, _, err := s.Get(value, "one")
		require.Nil(t, err)
		require.Equal(t, int64(1), one)
	}

	{
		two, _, err := s.Get(value, "two")
		require.Nil(t, err)
		require.Equal(t, "two", two)
	}

	{
		three, _, err := s.Get(value, "three")
		require.Nil(t, err)
		require.Equal(t, []string{"a", "b", "c"}, three)
	}

	{
		four, _, err := s.Get(value, "four")
		require.Nil(t, err)
		require.Equal(t, nil, four)
	}

	{
		five, _, err := s.Get(value, "five")
		require.NotNil(t, err)
		require.Equal(t, nil, five)
	}
}
