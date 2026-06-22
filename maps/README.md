# maps

Generic free functions over plain Go maps: `Equal`/`NotEqual` (value comparison), `Keys`/`Values` (extract slices), and `KeysSorted`. Part of [rosetta](../README.md).

This is the map analogue of [slice](../slice/); for typed map *values* with accessors, see [mapof](../mapof/).

## What matters here

- **`Keys` and `Values` return results in map-iteration order, which Go randomizes — they do NOT sort.** If you need deterministic output (for comparison, display, or golden tests), use `KeysSorted`. Reaching for `Keys` and expecting stable order is a silent flakiness bug.
- **`Equal` compares values with `==`,** so the value type is constrained to `comparable`. It is not a deep/reflect comparison — maps of slices or maps of maps are out of scope.
