package schema

import (
	"fmt"
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/require"
)

func TestValidationError(t *testing.T) {

	e := Invalid("name is required")

	require.Equal(t, "name is required", e.Message)
	require.Equal(t, "name is required", e.Error())
	require.Equal(t, ValidationErrorCode, e.ErrorCode())
}

func ExampleValidationError() {

	// Derp includes a custom error type for data validation, that tracks
	// the name (or path) of the invalid field and the reason that it is invalid

	err := Invalid("Field is required, or is too short, or is something else we don't like.")

	// ValidationErrors work anywhere that a standard error works
	fmt.Println(err.Error())

	// Derp can also calculates the HTTP error code for ValidationErrors, which is 422 "Unprocessable Entity".
	fmt.Println(derp.ErrorCode(err))

	// Output: Field is required, or is too short, or is something else we don't like.
	// 422
}
