// Package schema defines, validates, and manipulates JSON-Schema-like
// structures in Go.
//
// A Schema wraps a single Element — Object, Array, String, Integer, Number,
// Boolean, or Any — each describing one node of a document along with its
// validation rules. Schemas are built in Go or unmarshaled from JSON, and are
// used for three jobs: validating a value, normalizing it into the shape the
// schema describes, and reading or writing individual values by path.
//
// Path access is interface-driven rather than reflective. A target object opts
// in by implementing the narrow Getter and Setter interfaces it can support —
// StringGetter, IntSetter, ArrayGetter, and the rest — and the schema calls
// only what the type advertises. The mapof and sliceof packages implement these
// already, which makes them the default carriers when there is no Go struct to
// bind to.
//
// Setting a value validates it, so a write can coerce, clamp, or truncate what
// it was given rather than storing it verbatim. String elements delegate that
// work to a named format from the format subpackage, and UseFormat registers
// additional ones.
package schema
