# delta

Change-tracking containers. `delta.Slice[T]` wraps a `[]T` and records which elements were `Added` and `Deleted` relative to its starting state, so a caller can compute and apply a minimal diff. Part of [rosetta](../README.md).

## What matters here

- **Only `Values` is serialized; `Added` and `Deleted` are transient.** Both carry `json:"-" bson:"-"` tags — the diff bookkeeping exists only for the lifetime of the in-memory object and is intentionally not persisted. A round-trip through JSON/BSON yields a `Slice` with the current values but an empty change set, as if freshly constructed.
- **`SetValue(any)` is the schema-integration seam** (it satisfies the `ValueSetter` interface used by [schema](../schema/)). It accepts a generic value and coerces it into `[]T`; that is why the element type is constrained to `comparable` (membership checks drive the Added/Deleted tracking).
- **`SetValue` must never assume it is handed a `[]T`.** Multi-value form widgets (multiselect, check-button-group) post a `*sliceof.String`, because that is the only shape satisfying the `ArrayGetterSetter` validation that `schema.validate_Array` requires. Conversion is delegated to [`convert.SliceOfOk[T]`](../convert/sliceOf.go), which unwraps pointers, named slice types, and `convert.ArrayGetter` values; a strict type assertion here silently discarded every multiselect write for six weeks.
- **Conversion is as loose as the `convert` package allows, including LOSSY.** `SetValue` discards convert's lossless flag, so `Slice[string]` handed `3.14159` stores `["3.14"]`. The flip side: a `Slice[int]` handed `"not a number"` stores `[0]` rather than erroring. Only a value that is not a collection at all — a struct, a map, a channel — is an error, checked up front with `convert.SliceOfAnyOk` so that a rejected write cannot empty `Values` and fill `Deleted`.
- **A value that cannot be read leaves the Slice untouched.** Coercion completes before any state is written, so a failed `SetValue` cannot empty `Values` *and* fill `Deleted` — which would read downstream as "the user removed everything."
- **`Values` is copied, not aliased.** The incoming slice arrives straight from a `url.Values` entry that the caller still owns.
- **`UnmarshalJSON` parses untrusted input and is fuzzed** (`FuzzSliceUnmarshalJSON`). When changing the unmarshal path, keep that fuzzer green — it guards against panics on malformed JSON.
