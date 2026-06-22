package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Integer.ValidateRequiredIf -------------------------------------------------

func TestRequiredIf_Integer_ConditionFalse(t *testing.T) {
	// The condition is false, so the (empty) value is not required
	schema := New(Object{Properties: ElementMap{
		"name": String{},
		"age":  Integer{BitSize: 32, RequiredIf: "name is Bob"},
	}})

	require.NoError(t, schema.ValidateRequiredIf(testStructB{Name: "Alice", Age: 0}))
}

func TestRequiredIf_Integer_RequiredButMissing(t *testing.T) {
	// The condition is true and the value is zero, so validation fails
	schema := New(Object{Properties: ElementMap{
		"name": String{},
		"age":  Integer{BitSize: 32, RequiredIf: "name is Bob"},
	}})

	require.Error(t, schema.ValidateRequiredIf(testStructB{Name: "Bob", Age: 0}))
}

func TestRequiredIf_Integer_RequiredAndPresent(t *testing.T) {
	// The condition is true and the value is present, so validation passes
	schema := New(Object{Properties: ElementMap{
		"name": String{},
		"age":  Integer{BitSize: 32, RequiredIf: "name is Bob"},
	}})

	require.NoError(t, schema.ValidateRequiredIf(testStructB{Name: "Bob", Age: 42}))
}

// Boolean.ValidateRequiredIf -------------------------------------------------

func TestRequiredIf_Boolean_RequiredButMissing(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":   String{},
		"active": Boolean{RequiredIf: "name is Bob"},
	}})

	require.Error(t, schema.ValidateRequiredIf(testStructA{Name: "Bob", Active: false}))
}

func TestRequiredIf_Boolean_RequiredAndPresent(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":   String{},
		"active": Boolean{RequiredIf: "name is Bob"},
	}})

	require.NoError(t, schema.ValidateRequiredIf(testStructA{Name: "Bob", Active: true}))
}

func TestRequiredIf_Boolean_ConditionFalse(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":   String{},
		"active": Boolean{RequiredIf: "name is Bob"},
	}})

	require.NoError(t, schema.ValidateRequiredIf(testStructA{Name: "Alice", Active: false}))
}

// Number.ValidateRequiredIf --------------------------------------------------

func TestRequiredIf_Number_RequiredButMissing(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":     String{},
		"latitude": Number{BitSize: 64, RequiredIf: "name is Bob"},
	}})

	require.Error(t, schema.ValidateRequiredIf(testStructA{Name: "Bob", Latitude: 0}))
}

func TestRequiredIf_Number_RequiredAndPresent(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":     String{},
		"latitude": Number{BitSize: 64, RequiredIf: "name is Bob"},
	}})

	require.NoError(t, schema.ValidateRequiredIf(testStructA{Name: "Bob", Latitude: 45.5}))
}

func TestRequiredIf_Number_ConditionFalse(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":     String{},
		"latitude": Number{BitSize: 64, RequiredIf: "name is Bob"},
	}})

	require.NoError(t, schema.ValidateRequiredIf(testStructA{Name: "Alice", Latitude: 0}))
}

// String.ValidateRequiredIf --------------------------------------------------

func TestRequiredIf_String_RequiredButMissing(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":     String{},
		"optional": String{RequiredIf: "name is Bob"},
	}})

	require.Error(t, schema.ValidateRequiredIf(testStructA{Name: "Bob", Optional: ""}))
}

// Array.ValidateRequiredIf ---------------------------------------------------

func TestRequiredIf_Array_NoCondition(t *testing.T) {
	// With no condition, ValidateRequiredIf returns immediately
	schema := New(Object{Properties: ElementMap{
		"array": Array{Items: String{}},
	}})

	value := testStructA{Array: testArrayA{}}
	require.NoError(t, schema.ValidateRequiredIf(&value))
}

// NOTE: Array.ValidateRequiredIf resolves the field via
// getPropertyRecursive(element, ...) using the *array's own* element rather than
// the schema root, so a nested array with a real `required-if` condition fails
// to resolve its path ("Invalid property") before reaching its length/required
// logic. This test pins that current behavior; the deeper branches of
// Array.ValidateRequiredIf are not reachable for a nested array until that is
// fixed. See the summary notes for details.
func TestRequiredIf_Array_NestedConditionCannotResolvePath(t *testing.T) {
	schema := New(Object{Properties: ElementMap{
		"name":  String{},
		"array": Array{Items: String{}, RequiredIf: "name is Bob"},
	}})

	value := testStructA{Name: "Bob", Array: testArrayA{}}
	require.Error(t, schema.ValidateRequiredIf(&value))
}
