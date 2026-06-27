# schema

A fast, minimal subset of [JSON-Schema](http://json-schema.org). It unmarshals a schema from JSON, validates values against Array/Boolean/Integer/Number/Object/String element types (with custom string [formats](./format/)), and reads/writes data through a schema by [JSON-Pointer](https://tools.ietf.org/html/rfc6901) path. This is the heart of [rosetta](../README.md) — the engine that lets dynamic data (form inputs, CMS records) be typed and validated at runtime. If you need a complete, rigorous JSON-Schema implementation, use another tool.

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/schema)

**Included:** unmarshal-from-JSON, the six element-type validators, custom `Format` rules, and a happy API for walking the schema tree and getting/setting data by path. **Left out:** `$ref` references and loading remote schemas by URI.

## How data access works

Schema does not reflect over your structs. Instead, the data structure must implement Getter/Setter interfaces so schema has a type-safe path to its values. [mapof](../mapof/) and [sliceof](../sliceof/) implement these fully, so a `mapof.Any` is the easiest carrier for dynamic data; structs can implement them by hand:

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

func (m *MyStruct) SetStringOK(path string, value string) bool {
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

Additional interfaces (`GetPointer`, `GetObject`, array `GetIndex`/`SetIndex`, …) let schema traverse nested objects and arrays.

## What matters here

- **`Schema.Set` VALIDATES and may COERCE the value before storing it — it is not a plain assignment.** Every `Set` runs the element's `validate` step, so a value can be clamped (min/max), truncated (maxLength), rewritten by a format, or rejected with an error. Callers that treat `Set` as a dumb setter can be surprised when the stored value differs from the input or when `Set` returns an error. (`SetAll`/`SetURLValues` additionally run `ValidateRequiredIf` after all values are set.)
- **Data structures must implement the Getter/Setter interfaces; schema never reflects.** A type missing the interface for a given path silently fails the get/set (`ok=false`) rather than panicking. This is the central design constraint — the interfaces in `interfaces.go` are the contract every carrier type must satisfy.
- **String formats are registered with `UseFormat`, and the registry FREEZES on first read.** The first validation that looks up a format locks the registry; any `UseFormat` after that is rejected. Register all custom formats during `init`/startup, before validation runs. See [format/](./format/).
- **The per-element `Validate(any)` methods (`Boolean.Validate`, `Integer.Validate`, `Number.Validate`, `String.Validate`, `Any.Validate`) are intentionally retained.** They look like dead duplicates of the package-level `validate`/`Validate` flow but are part of the element interface and kept on purpose — do not prune them as unused.
- **Validation can rewrite values, so the validate path returns `(value, changed, error)`.** Treat a returned value as authoritative even when there's no error; the `changed` flag tells you whether coercion occurred.
