// Package sliceof provides named slice types that carry helper methods: Any,
// String, Int, Float, MapOfAny, MapOfString, and the generic Object[T].
//
// Each is a plain []T under its methods, so it ranges, indexes, and serializes
// like the slice it wraps. What it adds is the Getter and Setter
// implementations that let the schema package address elements by path, along
// with the GroupBy helpers that report where a sorted slice changes value in a
// named field, for rendering group headers.
//
// Together with mapof, these are the carrier types a JSON-Schema-shaped
// document is read into when there is no Go struct to bind to. For free
// functions over ordinary slices, see the slice package.
package sliceof
