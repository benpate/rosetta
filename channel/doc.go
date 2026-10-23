// Package channel provides generic combinators for working with Go channels:
// Filter, Map, Limit, Reverse, Beep, Pipe, and conversions to and from slices.
//
// Every combinator spawns a goroutine that closes its output when it finishes,
// usually when its input drains. Most return a new output channel; Pipe and
// PipeWithCancel write to, and close, the one they are given. Callers must never
// close a channel this package closes, and must close the input they pass in, or
// the goroutine outlives them. PipeWithCancel takes a done channel to stop early,
// while Limit closes the done channel it is given when it stops.
//
// Filter and Map buffer one item, which decouples producer from consumer by a
// single item; every other channel it creates is unbuffered. A slow consumer still applies
// backpressure all the way up the chain.
package channel
