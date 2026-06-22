# channel

Generic combinators over Go channels: `Filter`, `Map`, `Limit`, `Beep`, `Pipe`/`PipeWithCancel`, and `Closed`. Part of [rosetta](../README.md).

## What matters here

- **Each combinator spawns a goroutine and owns/closes its OUTPUT channel.** `Filter`, `Map`, etc. return a new channel and `defer close()` it from the goroutine they launch. Never close a channel returned by these functions yourself (double-close panics); the combinator closes it when its input drains.
- **The goroutine lives until the INPUT channel closes — so the caller must close the input, or the goroutine leaks.** A combinator chained off a channel that never closes is a goroutine leak. When you need to stop early without closing the source, use `PipeWithCancel`, which selects on a `done` channel and returns when signaled.
- **Output channels are buffered (size 1).** This decouples producer and consumer by one item but does not make the pipeline unbounded — a slow consumer still backpressures the goroutine.
- **`Closed` is a non-blocking best-effort check**, useful in tests and guards; do not use it to make correctness decisions in a concurrent producer, where the state can change immediately after the check.
