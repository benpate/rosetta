# delta

Change-tracking containers. `delta.Slice[T]` wraps a `[]T` and records which elements were `Added` and `Deleted` relative to its starting state, so a caller can compute and apply a minimal diff. Part of [rosetta](../README.md).

## What matters here

- **Only `Values` is serialized; `Added` and `Deleted` are transient.** Both carry `json:"-" bson:"-"` tags — the diff bookkeeping exists only for the lifetime of the in-memory object and is intentionally not persisted. A round-trip through JSON/BSON yields a `Slice` with the current values but an empty change set, as if freshly constructed.
- **`SetValue(any)` is the schema-integration seam** (it satisfies the `ValueSetter` interface used by [schema](../schema/)). It accepts a generic value and coerces it into `[]T`; that is why the element type is constrained to `comparable` (membership checks drive the Added/Deleted tracking).
- **`SetValue` must never assume it is handed a `[]T`.** Multi-value form widgets (multiselect, check-button-group) post a `*sliceof.String`, because that is the only shape satisfying the `ArrayGetterSetter` validation that `schema.validate_Array` requires. `coerce.go` reshapes pointers, named slice types, and array-like values; a strict type assertion here silently discarded every multiselect write for six weeks.
- **Coercion is structural, not a value conversion.** Elements must already be `T` — a `Slice[int]` handed `[]string{"5"}` returns an error rather than guessing. The concrete `sliceof.*` types can be looser because they know their element type and delegate to `convert.SliceOfXxx`; a generic `Slice[T]` has no such helper.
- **A value that cannot be read leaves the Slice untouched.** Coercion completes before any state is written, so a failed `SetValue` cannot empty `Values` *and* fill `Deleted` — which would read downstream as "the user removed everything."
- **`Values` is copied, not aliased.** The incoming slice arrives straight from a `url.Values` entry that the caller still owns.
- **`UnmarshalJSON` parses untrusted input and is fuzzed** (`FuzzSliceUnmarshalJSON`). When changing the unmarshal path, keep that fuzzer green — it guards against panics on malformed JSON.
