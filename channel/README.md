# channel

Generic combinators over Go channels: `Filter`, `Map`, `Limit`, `Reverse`, `Beep`, `Pipe`/`PipeWithCancel`, `Slice`/`FromSlice`, and `Closed`. Part of [rosetta](../README.md).

## What matters here

- **Each combinator spawns a goroutine and owns/closes its OUTPUT channel.** `Filter`, `Map`, etc. return a new channel and `defer close()` it from the goroutine they launch. Never close a channel returned by these functions yourself (double-close panics); the combinator closes it when its input drains. `Pipe` and `PipeWithCancel` close the output channel you pass in, and `Limit` closes the `done` channel you pass in (EmissarySocial BUG-289), so never close those yourself either.
- **The goroutine lives until the INPUT channel closes — so the caller must close the input, or the goroutine leaks.** A combinator chained off a channel that never closes is a goroutine leak. When you need to stop early without closing the source, use `PipeWithCancel`, which selects on a `done` channel while it waits for input. It does not watch `done` while it sends, so it still leaks when nothing reads its output.
- **Only `Filter` and `Map` buffer their output (size 1); every other channel this package creates is unbuffered.** The buffer decouples producer and consumer by one item but does not make the pipeline unbounded — a slow consumer still backpressures the goroutine.
- **`Closed` is a non-blocking best-effort check**, useful in tests and guards; do not use it to make correctness decisions in a concurrent producer, where the state can change immediately after the check. It receives from the channel, so on a channel that holds a value it consumes that value and reports `true`.
