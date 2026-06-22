# slice

Generic free functions over plain `[]T` — membership (`Contains`, `ContainsAny`, `ContainsAll`), bounds-safe access (`At`, `AtOK`), set operations (`Difference`, `Equal`), transforms (`Map`, `Filter`, `Find`, `First`, `NonZero`), grouping (`Grouper`), and `Shuffle`. Part of [rosetta](../README.md).

These operate on raw slices and return new slices; contrast with [sliceof](../sliceof/), whose types attach methods to named slice types for schema integration.

## What matters here

- **`At` / `AtOK` are bounds-safe — they never panic.** A negative or out-of-range index returns the element type's zero value (`At`) or `(zero, false)` (`AtOK`). Use these instead of raw `slice[i]` when the index comes from untrusted or computed input.
- **`Shuffle` uses `math/rand`, intentionally** — ordering is not security-sensitive. Do not switch it to `crypto/rand`.
- **Transform functions return new slices and do not mutate the input** (`Filter`, `Map`, `NonZero`, `Difference`). The input slice's backing array is left alone, so callers can keep using it.
- **`Grouper` requires elements that expose a string field via the `stringOKGetter` interface.** It groups a slice by a named field for rendering; the element type constraint is interface-based, not `comparable`.
