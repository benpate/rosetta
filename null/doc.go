// Package null provides nullable wrappers around Go primitive types: Bool,
// String, Int, Int64, and Float, plus the generic Object[T] for everything
// else.
//
// Each type carries a value alongside a flag recording whether that value was
// ever set, which is the one thing a plain int or bool cannot tell you. The
// zero value of every type here is a valid null and needs no constructor.
// IsNull (aliased as IsNil) asks whether a value was ever set; IsZero asks
// whether it is empty, and answers TRUE for a null value and for a present zero
// alike.
//
// JSON handling is strict in and plain out. Marshaling emits either the bare
// value or the literal null, with no omitempty-style elision — dropping the
// field entirely is the containing struct's decision. Unmarshaling never
// coerces across JSON types: a bare number is an error for String, and a quoted
// number is an error for Int. When a remote peer's encoding is not yours to
// control, the lenient package is the tolerant counterpart to this one.
package null
