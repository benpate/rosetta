// Package channel provides generic combinators for working with Go channels:
// Filter, Map, Limit, Reverse, Beep, Pipe, and conversions to and from slices.
//
// Every combinator spawns a goroutine, returns a new output channel, and owns
// that channel's lifetime — it closes the output when its input drains. Callers
// must never close a channel returned from this package, and must close the
// input they pass in, or the goroutine outlives them. PipeWithCancel and Limit
// take a done channel for the cases where stopping early is the only way out.
//
// Output channels carry a buffer of one, which decouples producer from consumer
// by a single item without making the pipeline unbounded. A slow consumer still
// applies backpressure all the way up the chain.
package channel
