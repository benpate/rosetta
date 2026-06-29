package mapof_test

import (
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// Compile-time assertions that every mapof type declares itself a map via schema.MapTyper.
// validate_Object relies on this to treat an absent key as a legitimately-empty optional property;
// a map type that silently lost IsMap() would instead make absent keys fail validation.
var (
	_ schema.MapTyper = mapof.Any(nil)
	_ schema.MapTyper = mapof.Bool(nil)
	_ schema.MapTyper = mapof.Int(nil)
	_ schema.MapTyper = mapof.Int64(nil)
	_ schema.MapTyper = mapof.Float(nil)
	_ schema.MapTyper = mapof.String(nil)
	_ schema.MapTyper = mapof.Object[string](nil)
)
