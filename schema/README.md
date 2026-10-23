# schema

A fast, minimal subset of [JSON-Schema](http://json-schema.org). It unmarshals a schema from JSON, validates values against Array/Boolean/Integer/Number/Object/String element types (with custom string [formats](./format/)), and reads/writes data through a schema by dotted path (such as `child.name` or `tags.0`). This is the heart of [rosetta](../README.md) — the engine that lets dynamic data (form inputs, CMS records) be typed and validated at runtime. If you need a complete, rigorous JSON-Schema implementation, use another tool.

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/rosetta/schema.svg)](https://pkg.go.dev/github.com/benpate/rosetta/schema)

**Included:** unmarshal-from-JSON, the six element-type validators, custom `Format` rules, and a happy API for walking the schema tree and getting/setting data by path. **Left out:** `$ref` references and loading remote schemas by URI.

## How data access works

Schema does not reflect over your structs. Instead, the data structure must implement Getter/Setter interfaces so schema has a type-safe path to its values. [mapof](../mapof/) and [sliceof](../sliceof/) implement them, so a `mapof.Any` is the easiest carrier for dynamic data; structs can implement them by hand:

```go
type MyStruct struct {
	Name  string
	Email string
}

func (m MyStruct) GetStringOK(path string) (string, bool) {
	switch path {
	case "name":
		return m.Name, true
	case "email":
		return m.Email, true
	}
	return "", false
}

func (m *MyStruct) SetString(path string, value string) bool {
	switch path {
	case "name":
		m.Name = value
		return true
	case "email":
		m.Email = value
		return true
	}
	return false
}
```

Additional interfaces (`GetPointer`, `SetKey`, `IsMap`, `Keys`, array `Length`/`GetIndex`/`SetIndex`, …) let schema traverse nested objects and arrays.

## What matters here

- **`Schema.Set` VALIDATES and may COERCE the value before storing it — it is not a plain assignment.** Every `Set` runs the element's `validate` step, so a value can be clamped (min/max), truncated (maxLength), rewritten by a format, or rejected with an error. Callers that treat `Set` as a dumb setter can be surprised when the stored value differs from the input or when `Set` returns an error. (`SetAll`/`SetURLValues` additionally run `ValidateRequiredIf` after all values are set.)
- **Data structures must implement the Getter/Setter interfaces; schema never discovers fields by reflection.** A type missing the interface for a given path fails the get or set with an error rather than panicking. This is the central design constraint — the interfaces in `interfaces.go` are the contract every carrier type must satisfy.
- **String formats are registered with `UseFormat`, and the registry FREEZES on first read.** The first validation that looks up a format locks the registry; any `UseFormat` after that is rejected. Register all custom formats during `init`/startup, before validation runs. See [format/](./format/).
- **The per-element `Validate(any)` methods (`Boolean.Validate`, `Integer.Validate`, `Number.Validate`, `String.Validate`, `Any.Validate`) are intentionally retained.** They look like dead duplicates of the package-level `validate`/`Validate` flow, and they are outside the `Element` interface with no in-repo callers, but they are kept on purpose — do not prune them as unused.
- **Validation can rewrite values, and `Validate` reports every rewrite as an error.** It writes the rewritten value through a pointer or a map and returns a `derp.Validation` error, so its `changed` result is always `false`. `Normalize` applies the same rewrites without an error and returns the path of each one.
