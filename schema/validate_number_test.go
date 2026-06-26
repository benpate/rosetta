package schema

import (
	"math"
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// testFloatGetter tests getting/setting float values
type testFloatGetter struct {
	value float64
}

func (t testFloatGetter) GetFloatOK(_ string) (float64, bool) {
	return t.value, true
}

func (t *testFloatGetter) SetFloat(_ string, value float64) bool {
	t.value = value
	return true
}

func TestFloatGetter(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Number{},
		},
	})

	var getter testFloatGetter

	value, err := schema.Get(&getter, "value")

	require.NoError(t, err)
	require.Equal(t, float64(0), value)

	require.NoError(t, schema.Set(&getter, "value", 12345678))
	require.Equal(t, float64(12345678), getter.value)
}

// testFloatPointer tests getting/setting bool values via a pointer
type testFloatPointer struct {
	value float64
}

// GetPointer gets a pointer to the float value
func (test *testFloatPointer) GetPointer(_ string) (any, bool) {
	return &test.value, true
}

// TestFloatPointer tests getting/setting float values via a pointer
func TestFloatPointer(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Number{BitSize: 64, Minimum: null.NewFloat(100), Maximum: null.NewFloat(1000000)}},
	},
	)

	var getter testFloatPointer

	// Validate correct value
	require.NoError(t, schema.Set(&getter, "value", 123456))
	require.Equal(t, float64(123456), getter.value)

	_, _, err := Validate(schema, &getter)
	require.NoError(t, err)

	// Value below the minimum is clamped to the minimum (does not fail)
	require.NoError(t, schema.Set(&getter, "value", float64(1)))
	require.Equal(t, float64(100), getter.value)

	_, _, err = Validate(schema, &getter)
	require.NoError(t, err)

	// Value above the maximum is clamped to the maximum (does not fail)
	require.NoError(t, schema.Set(&getter, "value", float64(1000000000)))
	require.Equal(t, float64(1000000), getter.value)

	_, _, err = Validate(schema, &getter)
	require.NoError(t, err)
}

// TestValidateNumber_Infinity pins the handling of ±Infinity: when a matching bound is
// present the value is clamped to it (and reported as changed); when no matching bound
// exists the value is rejected. +Inf clamps to Maximum, -Inf clamps to Minimum.
func TestValidateNumber_Infinity(t *testing.T) {

	// clamps asserts that value is rewritten to expected, with changed=true and no error.
	clamps := func(element Number, value float64, expected float64) {
		t.Helper()
		result, changed, err := validate_Number(element, value)
		require.NoError(t, err)
		require.True(t, changed, "value %v should be clamped", value)
		require.Equal(t, expected, result)
	}

	// rejects asserts that value is rejected with a validation error and not rewritten.
	rejects := func(element Number, value float64) {
		t.Helper()
		_, changed, err := validate_Number(element, value)
		require.Error(t, err, "value %v should be rejected", value)
		require.False(t, changed)
	}

	// +Inf clamps to the Maximum when one is present.
	clamps(Number{Maximum: null.NewFloat(100)}, math.Inf(1), 100)
	clamps(Number{Minimum: null.NewFloat(10), Maximum: null.NewFloat(100)}, math.Inf(1), 100)

	// -Inf clamps to the Minimum when one is present.
	clamps(Number{Minimum: null.NewFloat(10)}, math.Inf(-1), 10)
	clamps(Number{Minimum: null.NewFloat(10), Maximum: null.NewFloat(100)}, math.Inf(-1), 10)

	// +Inf is rejected when there is no Maximum to clamp to (a Minimum does not help).
	rejects(Number{}, math.Inf(1))
	rejects(Number{Minimum: null.NewFloat(10)}, math.Inf(1))

	// -Inf is rejected when there is no Minimum to clamp to (a Maximum does not help).
	rejects(Number{}, math.Inf(-1))
	rejects(Number{Maximum: null.NewFloat(100)}, math.Inf(-1))
}

// TestValidateNumber_NaN pins that NaN is always rejected and never clamped, regardless of
// which bounds are present (there is no sensible finite value to clamp NaN to).
func TestValidateNumber_NaN(t *testing.T) {

	rejectsNaN := func(element Number) {
		t.Helper()
		_, changed, err := validate_Number(element, math.NaN())
		require.Error(t, err)
		require.False(t, changed)
	}

	rejectsNaN(Number{})
	rejectsNaN(Number{Minimum: null.NewFloat(0)})
	rejectsNaN(Number{Maximum: null.NewFloat(100)})
	rejectsNaN(Number{Minimum: null.NewFloat(0), Maximum: null.NewFloat(100)})
	rejectsNaN(Number{MultipleOf: null.NewFloat(2)})
}
