// Package compare provides functions for comparing values of arbitrary types.
//
// The concrete helpers (Int, String, Float64, and their siblings) return the
// usual three-way result: negative, zero, or positive. Interface compares two
// values of unknown type by promoting them to a common base type. Equal,
// LessThan, and GreaterThan are built on it, while Contains, BeginsWith, and
// EndsWith test containment without it. Higher-level packages evaluate filter
// expressions with these predicates, and the Operator constants name them as
// strings, so a query language can be parsed straight into a comparison.
//
// Custom types opt into all of this by implementing one of the small interfaces
// in interfaces.go: Booler, Inter, Floater, Stringer, Hexer, LengthGetter,
// Nuller, or ContainsInterfacer. IsZero consults them in a fixed order, which
// is why a null.Int reports as zero through its Nuller implementation rather
// than through its numeric value.
package compare
