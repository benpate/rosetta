// Package ranges provides combinators for iter.Seq iterators: Filter, Map,
// Limit, Join, Unique, Empty, and Values, plus Slice to collect one.
//
// Everything here is lazy. Each combinator returns a new iter.Seq that pulls
// from its source only as the consumer ranges over it, so chaining them
// composes a single pass rather than building intermediate slices, and stopping
// early stops the whole chain.
//
// This is the modern counterpart to the iterator package, which adapts the
// older cursor-style interface. New code should start here.
package ranges
