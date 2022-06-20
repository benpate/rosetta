package schema

import (
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

// ValidationErrorCode represents HTTP Status Code: 422 "Unproccessable Entity"
const ValidationErrorCode = 422

// ValidationError represents an input validation error, and includes fields necessary to
// report problems back to the end user.
type ValidationError struct {
	Path    string `json:"path"`    // Identifies the PATH (or variable name) that has invalid input
	Message string `json:"message"` // Human-readable message that explains the problem with the input value.
}

// Invalid returns a fully populated ValidationError to the caller
func Invalid(message string) ValidationError {
	return ValidationError{
		Path:    "",
		Message: message,
	}
}

// Error returns a string representation of this ValidationError, and implements
// the builtin errors.error interface.
func (v ValidationError) Error() string {
	return strings.Trim(v.Path+" "+v.Message, " ")
}

// ErrorCode returns CodeValidationError for this ValidationError
// It implements the ErrorCodeGetter interface.
func (v ValidationError) ErrorCode() int {
	return ValidationErrorCode
}

// addPath pushes the path string into the provided error (if the types match)
func addPath(path string, err error) error {

	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case ValidationError:
		e.Path = list.PushHead(e.Path, path, ".")
		return e

	case derp.MultiError:

		for index := range e {
			e[index] = addPath(path, e[index])
		}

		return e
	}

	return err
}
