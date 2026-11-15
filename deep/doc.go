// Package deep copies and searches nested data: the maps and slices that a JSON
// or database decoder returns, at every level of nesting.
//
// Copy returns a copy that shares no map or slice with the original, so that a
// caller can change one without changing the other. Every container in the copy
// keeps its own type, and the copy is returned as the type that Copy was given.
//
// Remove deletes the named properties from every map within a value, in place.
// To leave the original untouched, Copy it first and Remove from the copy.
//
// Both functions list the common types for speed, and reach every other map,
// slice, and array by reflection. A value that neither can walk, such as a
// pointer or a struct, is an error, unless it implements Copier or Remover.
package deep
