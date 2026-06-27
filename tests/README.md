# tests

Cross-package integration tests that exercise [schema](../schema/) against the [mapof](../mapof/) and [sliceof](../sliceof/) carrier types together. Part of [rosetta](../README.md).

## What matters here

- **This directory exists to break an import cycle, not because the tests are special.** `schema` defines the Getter/Setter interfaces, and `mapof`/`sliceof` implement them — but `schema` cannot import `mapof`/`sliceof` in its own `_test` files without a cycle. Putting the round-trip tests in a separate `tests` package lets them import all three and verify that the implementations actually satisfy schema's path get/set contract.
- **`testInline` asserts a lossless round-trip: `Set` then `Get` returns the same value.** It's the shared helper for "this value survives a schema write/read unchanged." Because `schema.Set` validates and may coerce, a test that fails here often means the value was rewritten by the schema, not that get/set is broken — check the element's validation rules first.
