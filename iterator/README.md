# iterator

Adapters that drain a cursor-style `Iterator` (`Next(any) bool` + `Count() int`) into a slice, channel, or mapped result. Part of [rosetta](../README.md).

This is the **legacy** iteration interface, predating Go 1.23 range-over-func. For new code using `iter.Seq`, see [ranges](../ranges/); reach for this package only when adapting a type that already implements the `Iterator` interface (e.g. a database cursor).

## What matters here

- **`Next` populates a value through a pointer, so the `constructor func() T` must return a FRESH value each call.** `Slice`, `Channel`, and `Map` call the constructor once per element and pass `&value` to `Next`. If the constructor returns a shared/aliased value, every slice entry ends up pointing at the same mutated item. Allocate a new zero value each time.
- **`Channel` spawns a goroutine that closes its output and runs until the iterator is exhausted.** If the consumer stops reading early, the goroutine blocks forever (leak). Use `ChannelWithCancel` and signal its `cancel` channel to stop early.
- **`Slice` pre-allocates capacity from `iterator.Count()`.** An `Iterator` whose `Count()` lies (returns fewer than it yields) just causes a re-grow — correctness is fine, but `Count()` should be accurate for the allocation to help.
