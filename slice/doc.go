// Package slice provides generic utility functions for working with slices:
// membership tests (Contains, ContainsAll, ContainsAny), transformations
// (Map, Filter, Unique, Reverse, Shuffle, Difference, NonZero), and safe
// element access (At, AtOK, Find, First, Split, RemoveAt).
//
// At and its siblings are the reason most of this exists: they read an
// out-of-range index as the zero value rather than panicking, which lets
// calling code skip the bounds check that would otherwise wrap every access.
//
// These are functions over ordinary []T values. The sliceof package provides
// named slice types that carry methods and integrate with schema.
package slice
