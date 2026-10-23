# tests

Cross-package integration tests that exercise [schema](../schema/) against the [mapof](../mapof/) and [sliceof](../sliceof/) carrier types together. Part of [rosetta](../README.md).

## What matters here

- **This directory does not break an import cycle, and the tests are not special.** `schema` imports `mapof` and `sliceof` so it can create their containers (so neither may import `schema`), and schema's own `_test` files import both too. The round-trip tests here import all three and verify that the implementations satisfy schema's path get/set contract; they could live in `schema` just as well.
- **`testInline` asserts a lossless round-trip: `Set` then `Get` returns the same value.** It's the shared helper for "this value survives a schema write/read unchanged." Because `schema.Set` validates and may coerce, a test that fails here often means the value was rewritten by the schema, not that get/set is broken — check the element's validation rules first.
