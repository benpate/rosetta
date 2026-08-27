// Package delta provides collection types that track their own changes.
//
// A delta.Slice[T] wraps an ordinary []T and records which elements have been
// added and which have been deleted relative to the values it started with, so
// a caller can compute a minimal diff and apply only what actually moved. The
// element type is constrained to comparable because membership checks are what
// drive that bookkeeping.
//
// Only the values are serialized. The Added and Deleted sets are marked
// json:"-" and bson:"-" on purpose: they describe one in-memory editing
// session, not durable state, so a value that round-trips through storage comes
// back with its current contents and an empty change set.
//
// SetValue is the seam that the schema package writes through, and it converts
// whatever it is handed rather than asserting a []T: form widgets deliver array
// values as a *sliceof.String, and several other shapes are possible besides.
// The conversion is delegated to convert.SliceOfOk, so it is as forgiving as
// that package allows -- lossy renderings included -- and only a value that is
// not a collection at all is reported, without disturbing the slice.
package delta
