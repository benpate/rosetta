# pointer

A single helper, `pointer.To(value any) any`, that returns a pointer to its argument (and returns the value unchanged if it is already a pointer or interface). Part of [rosetta](../README.md).

## What matters here

- **`To` is reflection-based and untyped (`any` in, `any` out).** It exists for cases where the type isn't known at compile time (e.g. building values dynamically for the [schema](../schema/) engine). When the type *is* known, prefer the generic `convert.Pointer[T](v) *T`, which is type-safe and avoids reflection — `convert` owns the generic pointer/dereference helpers (`Pointer`, `Element`), not this package.
- **Already-pointer and interface inputs pass through unchanged**, so `To` is idempotent — calling it on a `*T` returns the same `*T`, not a `**T`.
