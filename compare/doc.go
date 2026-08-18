// Package compare provides functions for comparing values of arbitrary types.
//
// The concrete helpers (Int, String, Float64, and their siblings) return the
// usual three-way result: negative, zero, or positive. Interface compares two
// values of unknown type by promoting them to a common base type, and the
// predicate helpers built on it — Equal, GreaterThan, Contains, BeginsWith,
// EndsWith, and friends — are what higher-level packages use to evaluate filter
// expressions. The Operator constants name those predicates as strings, so a
// query language can be parsed straight into a comparison.
//
// Custom types opt into all of this by implementing one of the small interfaces
// in interfaces.go: Booler, Inter, Floater, Stringer, Hexer, LengthGetter,
// Nuller, or ContainsInterfacer. IsZero consults them in a fixed order, which
// is why a null.Int reports as zero through its Nuller implementation rather
// than through its numeric value.
package compare
