package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

/***********************************
 * validate_String: Pattern branch (was uncovered)
 ***********************************/

// TestValidateString_Pattern_Match confirms a value that matches the pattern is accepted.
func TestValidateString_Pattern_Match(t *testing.T) {
	value, changed, err := validate_String(String{Pattern: "^[a-z]+$"}, "hello")
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, "hello", value)
}

// TestValidateString_Pattern_NoMatch confirms a value that does not match the pattern is rejected.
func TestValidateString_Pattern_NoMatch(t *testing.T) {
	_, _, err := validate_String(String{Pattern: "^[a-z]+$"}, "Hello123")
	require.Error(t, err)
}

// TestValidateString_Pattern_InvalidRegex confirms an unparseable pattern returns an error
// rather than panicking.
func TestValidateString_Pattern_InvalidRegex(t *testing.T) {
	_, _, err := validate_String(String{Pattern: "([a-z"}, "hello")
	require.Error(t, err)
}

// TestValidateString_Pattern_ViaSchema confirms the Pattern rule is reachable through the
// public Validate entry point too.
func TestValidateString_Pattern_ViaSchema(t *testing.T) {
	_, _, okErr := Validate(New(String{Pattern: "^[0-9]{3}$"}), "123")
	require.NoError(t, okErr)

	_, _, failErr := Validate(New(String{Pattern: "^[0-9]{3}$"}), "12")
	require.Error(t, failErr)
}

/***********************************
 * Schema.Validate (method) — was 0%
 ***********************************/

// TestSchemaValidate_Method_OK confirms the Schema.Validate method accepts a conforming value
// and reports no change.
func TestSchemaValidate_Method_OK(t *testing.T) {
	changed, err := New(Integer{Minimum: null.NewInt64(0), Maximum: null.NewInt64(10)}).Validate(5)
	require.NoError(t, err)
	require.False(t, changed)
}

// TestSchemaValidate_Method_Error confirms the Schema.Validate method rejects an out-of-range
// value (which Set would clamp) rather than silently rewriting it.
func TestSchemaValidate_Method_Error(t *testing.T) {
	_, err := New(Integer{Maximum: null.NewInt64(10)}).Validate(10000)
	require.Error(t, err)
}

// TestSchemaValidate_Method_RequiredIf confirms the Schema.Validate method also enforces
// required-if conditions on the whole object.
func TestSchemaValidate_Method_RequiredIf(t *testing.T) {

	schema := New(testStructA_Schema())

	// name triggers the condition ("name is Aethelflad") but optional is empty -> error.
	object := testStructA{Name: "Aethelflad", Optional: ""}

	_, err := schema.Validate(&object)
	require.Error(t, err)
}

/***********************************
 * validate_Object: write-back branch (was 79%)
 ***********************************/

// TestValidateObject_WritesBackRewrittenValue confirms that when a property is rewritten during
// validation, validate_Object stores the rewritten value back into the object and reports the
// object as changed. This exercises validate_Object directly (the top-level Validate wrapper
// turns "changed" into an error, so it would not reach this branch).
func TestValidateObject_WritesBackRewrittenValue(t *testing.T) {

	element := Object{
		Properties: ElementMap{
			"name": String{MaxLength: 4},
		},
	}

	// "Aethelflad" exceeds MaxLength 4, so validation truncates it and writes the result back.
	object := testStructA{Name: "Aethelflad"}

	value, changed, err := validate(element, &object)
	require.NoError(t, err)
	require.True(t, changed)

	// The rewritten (truncated) value must have been stored back into the object.
	stored, err := New(element).Get(value, "name")
	require.NoError(t, err)
	require.Equal(t, "Aeth", stored)
}

/***********************************
 * Any.ValidateRequiredIf (was ~21%)
 ***********************************/

// anyRequiredObject is a backing object for exercising Any.ValidateRequiredIf: a trigger field
// drives the condition and a value field is read back via the getter interfaces.
type anyRequiredObject struct {
	Trigger string
	Value   string
}

func (o anyRequiredObject) GetStringOK(name string) (string, bool) {
	switch name {
	case "trigger":
		return o.Trigger, true
	case "value":
		return o.Value, true
	}
	return "", false
}

func (o *anyRequiredObject) GetPointer(name string) (any, bool) {
	switch name {
	case "trigger":
		return &o.Trigger, true
	case "value":
		// Return the string value itself (not a pointer) so an empty value reads as a zero
		// value, which is how Any.ValidateRequiredIf detects a missing required field.
		return o.Value, true
	}
	return nil, false
}

// TestAny_ValidateRequiredIf covers the Any branches: condition false (skip), condition true with
// a missing value (error), and condition true with a present value (success).
func TestAny_ValidateRequiredIf(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"trigger": String{},
			"value":   Any{RequiredIf: "trigger is go"},
		},
	})

	// Condition false -> no error.
	require.NoError(t, schema.ValidateRequiredIf(&anyRequiredObject{Trigger: "stop"}))

	// Condition true, value missing -> error.
	require.Error(t, schema.ValidateRequiredIf(&anyRequiredObject{Trigger: "go", Value: ""}))

	// Condition true, value present -> no error.
	require.NoError(t, schema.ValidateRequiredIf(&anyRequiredObject{Trigger: "go", Value: "here"}))
}
