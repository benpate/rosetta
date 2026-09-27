# loose

Types that are forgiving about the values they accept: `loose.Int64`, `loose.String`, and `loose.Template`. They exist for the receiving half of Postel's law — be conservative in what you send, be liberal in what you accept. Part of [rosetta](../README.md).

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/rosetta/loose.svg)](https://pkg.go.dev/github.com/benpate/rosetta/loose)

A plain `int64` field fails an *entire document* the first time a peer sends `"480"` instead of `480`. Swap in `loose.Int64` and that document parses, along with every other numeric spelling in the wild — quoted integers, floats, nulls, empty strings, `"100%"`. `loose.String` does the same for text fields that peers sometimes send as bare numbers or booleans, which is how oEmbed's `"version": 1.0` arrives from SoundCloud.

```go
type Response struct {
	Version loose.String `json:"version"` // arrives as 1.0, a JSON *number*
	Width   loose.Int64  `json:"width"`   // arrives as "480", a JSON *string*
	Height  loose.Int64  `json:"height"`  // arrives as null
}
```

Tolerance never leaks into what you publish: `Int64` and `String` marshal to the plain, correct JSON encoding of their underlying value.

### Templates

`loose.Template` is a string that `NewTemplate` compiles as a `text/template` when it holds a `{{` followed later by a `}}`. A string that does not compile is kept as plain text, not rejected, so a typo renders literally instead of failing. `Execute` renders a template against a value, or returns plain text as it is. It decodes any JSON scalar the way `loose.String` does, and marshals to its source string. `loose.CachedTemplate` compiles each source string once for the life of the process, up to a fixed limit, for template strings that are rendered repeatedly; `NewTemplate` never caches. It is compiled once and never changed, so it is safe to execute from many goroutines. [mapof.Template](../mapof/template.go) is a map whose strings are stored this way.

### What survives, and what doesn't

`Int64` defers to [convert.Int64](../convert/README.md), so floats truncate toward zero, out-of-range values clamp to the `int64` bounds, and anything unparseable quietly becomes zero — "not provided" rather than an error. Top-level integers are read from their **source text**, so IDs and timestamps above 2^53 stay exact instead of rounding through a `float64`.

It is `Int64` rather than `Int` deliberately. `convert.Int` clamps to the *platform* int width, so on a 32-bit target — `GOARCH=wasm` is one — a plain `Int` would silently cap at 2^31. A sender's JSON has no idea what you built for. Precision beyond int64 is explicitly not a goal: larger values clamp, and an integer nested inside an array rounds.

`String` keeps a number's source text too, which is the whole point: `1.0` must stay `"1.0"` and never become `"1"` or `"1.00"`.

Tolerance stops at structure. `String` rejects objects and arrays outright, because coercing one would invent a value the sender never wrote. Both types error only on input that isn't valid JSON at all.

### Fuzzing

The invariants are pinned by property-based fuzz targets in [fuzz_test.go](fuzz_test.go): decode never panics, a successful decode always re-encodes to valid JSON, values are a fixed point under repeated round trips, exact integers stay exact, and both types behave as struct fields without desynchronizing the decoder for their neighbors. The corpus in `testdata/fuzz` includes every crasher these targets have found — the round-trip property is what caught the `float64` precision loss.

Add a new loose type to `jsonTargets` in that file and it inherits the whole universal safety net.
