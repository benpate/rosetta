package translate

import "github.com/benpate/rosetta/schema"

// Runner in the interface for objects that implement Rules
type Runner interface {

	// Execute runs the rule on the input object, and writes the result to the output object
	Execute(inSchema schema.Schema, inObject any, outSchema schema.Schema, outObject any) error
}
