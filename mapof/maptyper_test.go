package mapof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Every mapof type must report IsMap() == TRUE. schema's validate_Object relies on this
// to treat an absent key as a legitimately-empty optional property; a type that silently
// lost IsMap() would instead make absent keys fail validation.
//
// The compile-time interface assertions live in maptyper_assertions_test.go; these
// exercise the method bodies.
func TestIsMap_AllTypes(t *testing.T) {

	check := func(name string, isMap bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			require.True(t, isMap)
		})
	}

	// Populated maps
	check("Any", Any{"a": 1}.IsMap())
	check("Bool", Bool{"a": true}.IsMap())
	check("Float", Float{"a": 1.0}.IsMap())
	check("Int", Int{"a": 1}.IsMap())
	check("Int64", Int64{"a": 1}.IsMap())
	check("String", String{"a": "b"}.IsMap())
	check("Object", Object[string]{"a": "b"}.IsMap())
	check("Matchable", NewMatchable[testMatcher]().IsMap())

	// Empty maps
	check("Any/empty", Any{}.IsMap())
	check("Bool/empty", Bool{}.IsMap())
	check("Float/empty", Float{}.IsMap())
	check("Int/empty", Int{}.IsMap())
	check("Int64/empty", Int64{}.IsMap())
	check("String/empty", String{}.IsMap())
	check("Object/empty", Object[string]{}.IsMap())

	// NIL maps must answer just the same -- IsMap describes the TYPE, not the contents
	check("Any/nil", Any(nil).IsMap())
	check("Bool/nil", Bool(nil).IsMap())
	check("Float/nil", Float(nil).IsMap())
	check("Int/nil", Int(nil).IsMap())
	check("Int64/nil", Int64(nil).IsMap())
	check("String/nil", String(nil).IsMap())
	check("Object/nil", Object[string](nil).IsMap())
	check("Matchable/nil", Matchable[testMatcher](nil).IsMap())
}
