// Package iterator adapts cursor-style collections into slices, channels, and
// mapped results.
//
// An Iterator yields values by filling in a pointer — Next(any) bool — which is
// the shape a database cursor naturally has. Slice, Channel, and Map drain one,
// calling a caller-supplied constructor once per element; that constructor must
// return a fresh value every time, because a shared value would leave every
// result aliasing the same item.
//
// This is the legacy iteration interface, predating range-over-func. New code
// should prefer iter.Seq and the ranges package; reach for this one only to
// adapt a type that already implements Iterator.
package iterator
