# delta

Change-tracking containers. `delta.Slice[T]` wraps a `[]T` and records which elements were `Added` and `Deleted` relative to its starting state, so a caller can compute and apply a minimal diff. Part of [rosetta](../README.md).

## What matters here

- **Only `Values` is serialized; `Added` and `Deleted` are transient.** Both carry `json:"-" bson:"-"` tags — the diff bookkeeping exists only for the lifetime of the in-memory object and is intentionally not persisted. A round-trip through JSON/BSON yields a `Slice` with the current values but an empty change set, as if freshly constructed.
- **`SetValue(any)` is the schema-integration seam** (it satisfies the `ValueSetter` interface used by [schema](../schema/)). It accepts a generic value and coerces it into `[]T`; that is why the element type is constrained to `comparable` (membership checks drive the Added/Deleted tracking).
- **`UnmarshalJSON` parses untrusted input and is fuzzed** (`FuzzSliceUnmarshalJSON`). When changing the unmarshal path, keep that fuzzer green — it guards against panics on malformed JSON.
