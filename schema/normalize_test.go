package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNormalize_RewritesInPlace confirms that Normalize applies rewrites to the
// object (instead of rejecting it, as Validate does) and reports the rewritten paths.
func TestNormalize_RewritesInPlace(t *testing.T) {

	s := New(pointerGetterStruct_Schema())

	// "TooLongName" exceeds MaxLength:8 and must be truncated by normalization.
	value := pointerGetterStruct{Name: "TooLongName", URL: "https://example.com/x"}

	rewrites, err := s.Normalize(&value)
	require.Nil(t, err)
	require.Equal(t, []string{"name"}, rewrites)
	require.Equal(t, "TooLongN", value.Name)
}

// TestNormalize_CleanValue confirms that a conforming value passes through
// unchanged with no rewrites reported.
func TestNormalize_CleanValue(t *testing.T) {

	s := New(pointerGetterStruct_Schema())

	value := pointerGetterStruct{Name: "Alice", URL: "https://example.com/alice"}

	rewrites, err := s.Normalize(&value)
	require.Nil(t, err)
	require.Empty(t, rewrites)
	require.Equal(t, "Alice", value.Name)
}

// TestNormalize_HardViolation confirms that values that cannot be made to
// conform still return an error.
func TestNormalize_HardViolation(t *testing.T) {

	s := New(Object{
		Properties: ElementMap{
			"name": String{Enum: []string{"red", "green", "blue"}},
		},
	})

	value := pointerGetterStruct{Name: "purple"}

	_, err := s.Normalize(&value)
	require.NotNil(t, err)
}

// TestNormalize_ArrayPaths confirms that rewrites inside arrays are reported
// with their full dot-path.
func TestNormalize_ArrayPaths(t *testing.T) {

	s := New(Array{Items: pointerGetterStruct_Schema()})

	value := pointerGetterSlice{
		{Name: "Alice", URL: "https://example.com/alice"},
		{Name: "TooLongName", URL: "https://example.com/x"},
	}

	rewrites, err := s.Normalize(&value)
	require.Nil(t, err)
	require.Equal(t, []string{"1.name"}, rewrites)
	require.Equal(t, "TooLongN", value[1].Name)
}
