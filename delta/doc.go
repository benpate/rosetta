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
package delta
