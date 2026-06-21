# sliceof

Generic, slice-backed types (`sliceof.Any`, `sliceof.Int`, `sliceof.Int64`, `sliceof.Float`, `sliceof.String`, `sliceof.Object[T]`, and the `MapOfAny`/`MapOfString` variants) that add typed accessors, grouping, and schema integration on top of a plain Go slice. The slice analogue of [mapof](../mapof/); part of [rosetta](../README.md).

## What matters here

- **`Shuffle` uses `math/rand`, not `crypto/rand` — intentionally.** Ordering randomization is not security-sensitive, so the cheaper PRNG is correct here. Do not "upgrade" it to `crypto/rand`.
- **These types implement the `schema` array interfaces** (`Length`, `GetIndex`, `SetIndex`, `GetPointer`, `SetValue`). That is the seam the [schema](../schema/) engine walks for array paths; index access there is bounds-checked by `schema.Index`, so out-of-range paths fail safely rather than panicking.
- **`sliceof.Object[T]` is generic over the element type**, while the scalar types (`Any`, `Int`, …) are concrete named slices. Reach for the concrete types when the element type is fixed; use `Object[T]` only when you genuinely need the type parameter.
- **Reads are bounds-safe; `SetIndex` grows the slice to fit.** `GetIndex` past the end returns the zero value with `ok=false` (via `slice.AtOK`); `SetIndex` at an out-of-range index *extends* the slice (`growSlice`) rather than failing, so setting index 100 on an empty slice yields a length-101 slice with zero-value gaps. Know this before using `SetIndex` with caller-supplied indices.
