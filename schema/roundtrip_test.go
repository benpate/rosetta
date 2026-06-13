package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// roundTrip marshals an element to a map and unmarshals it back into a fresh
// element, returning the result so callers can assert that no fields were lost.
func roundTrip(t *testing.T, original Element) Element {
	t.Helper()

	result, err := UnmarshalMap(original.MarshalMap())
	require.NoError(t, err)

	return result
}

// TestRoundTrip_String confirms that every String field survives a
// MarshalMap -> UnmarshalMap round trip.
func TestRoundTrip_String(t *testing.T) {

	original := String{
		Default:    "hello",
		MinLength:  3,
		MaxLength:  20,
		MinValue:   "aaa",
		MaxValue:   "zzz",
		Enum:       []string{"hello", "world"},
		Format:     "no-html",
		Required:   true,
		RequiredIf: "other == true",
	}

	require.Equal(t, original, roundTrip(t, original))
}

// TestRoundTrip_Integer confirms that every Integer field survives a round trip.
func TestRoundTrip_Integer(t *testing.T) {

	original := Integer{
		Default:    null.NewInt64(5),
		Minimum:    null.NewInt64(1),
		Maximum:    null.NewInt64(100),
		MultipleOf: null.NewInt64(5),
		BitSize:    64,
		Enum:       []int{5, 10, 15},
		Required:   true,
		RequiredIf: "other == true",
	}

	require.Equal(t, original, roundTrip(t, original))
}

// TestRoundTrip_Number confirms that every Number field survives a round trip.
func TestRoundTrip_Number(t *testing.T) {

	original := Number{
		Default:    null.NewFloat(2.5),
		Minimum:    null.NewFloat(0.5),
		Maximum:    null.NewFloat(99.5),
		MultipleOf: null.NewFloat(0.5),
		BitSize:    64,
		Enum:       []float64{0.5, 1.5, 2.5},
		Required:   true,
		RequiredIf: "other == true",
	}

	require.Equal(t, original, roundTrip(t, original))
}

// TestRoundTrip_Boolean confirms that every Boolean field survives a round trip.
func TestRoundTrip_Boolean(t *testing.T) {

	original := Boolean{
		Default:    null.NewBool(true),
		Required:   true,
		RequiredIf: "other == true",
	}

	require.Equal(t, original, roundTrip(t, original))
}

// TestRoundTrip_Any confirms that every Any field survives a round trip.
func TestRoundTrip_Any(t *testing.T) {

	original := Any{
		Required:   true,
		RequiredIf: "other == true",
	}

	require.Equal(t, original, roundTrip(t, original))
}

// TestRoundTrip_Array confirms that every Array field survives a round trip.
func TestRoundTrip_Array(t *testing.T) {

	original := Array{
		Items:      String{MaxLength: 10, Enum: []string{"a", "b"}},
		MinLength:  1,
		MaxLength:  5,
		Required:   true,
		RequiredIf: "other == true",
	}

	require.Equal(t, original, roundTrip(t, original))
}

// TestRoundTrip_Object confirms that every Object field (including nested
// properties, wildcard, and required-if) survives a round trip.
//
// Note: enums are set explicitly on nested elements because a missing enum
// round-trips to an empty (rather than nil) slice, and the object's own
// Required is left false because Object.UnmarshalMap intentionally propagates
// a true Required down to child properties. Both are pre-existing behaviors
// unrelated to field preservation.
func TestRoundTrip_Object(t *testing.T) {

	original := Object{
		Properties: ElementMap{
			"name": String{MaxLength: 50, Enum: []string{"a", "b"}, Required: true},
			"age":  Integer{Minimum: null.NewInt64(0), BitSize: 32, Enum: []int{1, 2}},
		},
		Wildcard:   String{MaxLength: 100, Enum: []string{"x"}},
		Required:   false,
		RequiredIf: "other == true",
	}

	require.Equal(t, original, roundTrip(t, original))
}
