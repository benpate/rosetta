# ranges

Composable combinators over Go 1.23 `iter.Seq[T]` range-over-func iterators: `Values`, `Empty`, `Filter`/`FilterPointer`, `Map`, `Join`, `Limit`, `Unique`, and `Slice`. Part of [rosetta](../README.md).

This is the **modern** iteration toolkit (built on `iter.Seq`); for the older cursor-style `Next(any) bool` interface, see [iterator](../iterator/).

## What matters here

- **Every combinator propagates early termination via the `yield` return value.** Each wraps the source sequence and returns from its own `yield` loop when the downstream `yield` returns `false`. This is what makes `for x := range ranges.Limit(10, seq)` stop the *source* after 10 items instead of draining it — `break` in the consumer must reach the producer. The tests assert this (`*_EarlyTermination`); a combinator that ignores `yield`'s bool would silently break it.
- **`iter.Seq` values are lazy and re-runnable.** Ranging the same sequence twice re-executes it from the start; combinators build a pipeline that does no work until ranged. Don't assume single-shot semantics.
- **`FilterPointer` passes `*T` to its predicate** (to avoid copying large elements), while `Filter` passes `T` by value. Pick `FilterPointer` for big structs; the predicate must not retain the pointer past the call.
- **`Slice` is the terminal that materializes a sequence into `[]T`** — use it when you need a concrete slice; otherwise keep things lazy.
