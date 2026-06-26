package schema

import (
	"math"
	"testing"

	"github.com/benpate/rosetta/null"
)

// constrainedObject is an addressable getter/setter backing store with one field
// of each constrained scalar type, used to prove the Set -> Validate invariant.
type constrainedObject struct {
	Str string
	Num float64
	Int int64
}

func (o *constrainedObject) GetPointer(name string) (any, bool) {
	switch name {
	case "str":
		return &o.Str, true
	case "num":
		return &o.Num, true
	case "int":
		return &o.Int, true
	}
	return nil, false
}

func (o *constrainedObject) GetStringOK(name string) (string, bool) {
	if name == "str" {
		return o.Str, true
	}
	return "", false
}

func (o *constrainedObject) GetFloatOK(name string) (float64, bool) {
	if name == "num" {
		return o.Num, true
	}
	return 0, false
}

func (o *constrainedObject) GetInt64OK(name string) (int64, bool) {
	if name == "int" {
		return o.Int, true
	}
	return 0, false
}

func (o *constrainedObject) SetString(name, value string) bool {
	if name == "str" {
		o.Str = value
		return true
	}
	return false
}

func (o *constrainedObject) SetFloat(name string, value float64) bool {
	if name == "num" {
		o.Num = value
		return true
	}
	return false
}

func (o *constrainedObject) SetInt64(name string, value int64) bool {
	if name == "int" {
		o.Int = value
		return true
	}
	return false
}

// constrainedSchema exercises the rewriting rules (clamp + format + length) that
// Set applies but Validate rejects.
func constrainedSchema() Schema {
	return New(Object{
		Properties: ElementMap{
			"str": String{MinLength: 1, MaxLength: 8},
			"num": Number{Minimum: null.NewFloat(-100), Maximum: null.NewFloat(100), BitSize: 64},
			"int": Integer{Minimum: null.NewInt64(-100), Maximum: null.NewInt64(100), BitSize: 64},
		},
	})
}

// assertStoredValueIsValid is the core invariant: after a successful Set, the value
// stored in the object MUST pass Validate against the same element. Set is allowed to
// rewrite (clamp/truncate/format); whatever it stores must already be conformant.
func assertStoredValueIsValid(t *testing.T, schema Schema, object *constrainedObject, path string) {
	t.Helper()

	stored, err := schema.Get(object, path)
	if err != nil {
		t.Fatalf("Get(%q) after Set failed: %v", path, err)
	}

	element, ok := schema.GetElement(path)
	if !ok {
		t.Fatalf("GetElement(%q) failed", path)
	}

	// The value Set stored must validate with no error and no further rewriting.
	if _, changed, err := Validate(New(element), stored); err != nil {
		t.Fatalf("Set(%q) stored %#v, but Validate rejected it: %v", path, stored, err)
	} else if changed {
		t.Fatalf("Set(%q) stored %#v, but Validate still reports it as needing changes", path, stored)
	}
}

// FuzzSetThenValidate_String confirms that any string Set successfully stores a value
// that subsequently passes Validate (Set may truncate/format; the result must be valid).
func FuzzSetThenValidate_String(f *testing.F) {

	f.Add("hello")
	f.Add("")
	f.Add("way too long to fit in eight")
	f.Add("<b>html tags</b>")
	f.Add("   whitespace   ")
	f.Add("日本語のテキストはとても長いです")

	f.Fuzz(func(t *testing.T, value string) {
		schema := constrainedSchema()
		object := &constrainedObject{}

		// A rejected Set is fine; we only assert the invariant when Set succeeds.
		if err := schema.Set(object, "str", value); err != nil {
			return
		}

		assertStoredValueIsValid(t, schema, object, "str")
	})
}

// FuzzSetThenValidate_Number confirms the invariant for numeric values, which Set clamps
// to the configured Minimum/Maximum.
func FuzzSetThenValidate_Number(f *testing.F) {

	f.Add(0.0)
	f.Add(50.0)
	f.Add(100.0)
	f.Add(1000.0)
	f.Add(-1000.0)
	f.Add(99.999)
	f.Add(math.NaN())
	f.Add(math.Inf(1))
	f.Add(math.Inf(-1))

	f.Fuzz(func(t *testing.T, value float64) {
		schema := constrainedSchema()
		object := &constrainedObject{}

		if err := schema.Set(object, "num", value); err != nil {
			return
		}

		assertStoredValueIsValid(t, schema, object, "num")
	})
}

// FuzzValidateNumber_Finite drives validate_Number directly across schemas with and without
// bounds, asserting the core finiteness guarantee: whenever validation succeeds, the returned
// value is a finite number (never NaN or ±Inf). A non-finite input must either be clamped to a
// finite bound or rejected with an error.
func FuzzValidateNumber_Finite(f *testing.F) {

	f.Add(0.0)
	f.Add(50.0)
	f.Add(math.NaN())
	f.Add(math.Inf(1))
	f.Add(math.Inf(-1))
	f.Add(math.MaxFloat64)

	f.Fuzz(func(t *testing.T, value float64) {

		// Exercise the bounded and unbounded shapes, since clamping is what lets a
		// non-finite value through as a (finite) rewrite rather than an error.
		elements := []Number{
			{},
			{Minimum: null.NewFloat(-100), Maximum: null.NewFloat(100)},
			{Minimum: null.NewFloat(-100)},
			{Maximum: null.NewFloat(100)},
		}

		for _, element := range elements {
			result, _, err := validate_Number(element, value)

			// A rejected value is always acceptable; we only constrain the success case.
			if err != nil {
				continue
			}

			if math.IsNaN(result) || math.IsInf(result, 0) {
				t.Fatalf("validate_Number(%#v, %v) succeeded with non-finite result %v", element, value, result)
			}
		}
	})
}

// FuzzSetThenValidate_Integer confirms the invariant for integer values, which Set clamps
// to the configured Minimum/Maximum.
func FuzzSetThenValidate_Integer(f *testing.F) {

	f.Add(int64(0))
	f.Add(int64(50))
	f.Add(int64(100))
	f.Add(int64(100000))
	f.Add(int64(-100000))

	f.Fuzz(func(t *testing.T, value int64) {
		schema := constrainedSchema()
		object := &constrainedObject{}

		if err := schema.Set(object, "int", value); err != nil {
			return
		}

		assertStoredValueIsValid(t, schema, object, "int")
	})
}

// FuzzSetAllThenValidate confirms the invariant through SetAll across every field at once:
// after SetAll succeeds, each stored value must independently pass Validate.
func FuzzSetAllThenValidate(f *testing.F) {

	f.Add("ok", 5.0, int64(5))
	f.Add("way too long for the field", 9999.0, int64(-9999))
	f.Add("<i>x</i>", -50.5, int64(100))
	f.Add("", 0.0, int64(0))

	f.Fuzz(func(t *testing.T, str string, num float64, i int64) {
		schema := constrainedSchema()
		object := &constrainedObject{}

		values := map[string]any{
			"str": str,
			"num": num,
			"int": i,
		}

		// SetAll stops at the first error; a rejected value is an acceptable outcome.
		if err := schema.SetAll(object, values); err != nil {
			return
		}

		// Every field SetAll wrote must now pass Validate.
		assertStoredValueIsValid(t, schema, object, "str")
		assertStoredValueIsValid(t, schema, object, "num")
		assertStoredValueIsValid(t, schema, object, "int")
	})
}
