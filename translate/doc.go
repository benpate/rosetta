// Package translate maps data from one object into another using rules that
// can be defined in JSON.
//
// A Pipeline is an ordered list of Rules, executed in order against a source
// and a target. Rules cover the mapping vocabulary: copy a value from one path
// to another, write a static value, evaluate a Go template expression, append
// to a collection, take the first rule that produces a result, branch on a
// condition, or iterate a source collection and apply nested rules to each
// element.
//
// Both objects are read and written through the schema package's Getter and
// Setter interfaces, so a pipeline never needs to know the concrete Go types
// involved. That makes mapof.Any the easiest carrier on either end, since it
// implements those interfaces already.
//
// Because a pipeline is just data, the mapping between two systems can be
// stored, shipped, and edited without recompiling anything.
package translate
