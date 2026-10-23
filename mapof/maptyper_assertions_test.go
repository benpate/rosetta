package mapof_test

import (
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// Compile-time assertions that the mapof map types declare themselves maps via schema.MapTyper (not
// Matchable, whose test Matcher is internal, or Slices). validate_Object relies on this to treat an absent
// key as a legitimately-empty optional property; a type that lost IsMap() would fail absent keys instead.
var (
	_ schema.MapTyper = mapof.Any(nil)
	_ schema.MapTyper = mapof.Bool(nil)
	_ schema.MapTyper = mapof.Int(nil)
	_ schema.MapTyper = mapof.Int64(nil)
	_ schema.MapTyper = mapof.Float(nil)
	_ schema.MapTyper = mapof.String(nil)
	_ schema.MapTyper = mapof.Object[string](nil)
	_ schema.MapTyper = mapof.Template(nil)
)
