# null

Nullable scalar primitives for Go: `null.Bool`, `null.Int`, `null.Int64`, and `null.Float`. Each is a struct carrying a value plus a `present` flag, so it distinguishes "absent/null" from "present-but-zero". All satisfy the `Nullable` interface (`IsNull() bool`) and implement JSON marshal/unmarshal. Part of [rosetta](../README.md).

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/null)

```go
// A zero-value null.Bool is null and ready to use — no constructor needed
var b null.Bool

b.Set(true)   // now present, value true
b.Bool()      // read the value
b.IsNull()    // false
b.Unset()     // back to null

b := null.NewBool(true) // or construct a non-null value directly
```

## What matters here

- **The zero value is a valid "null" — no constructor required.** `var x null.Int` is immediately usable and reads as null (`IsNull() == true`). This is the whole point of the package: it separates "never set" from "set to zero", which a plain `int`/`bool` cannot.
- **`Unset()` resets BOTH the value and the present flag.** It zeroes the stored value and clears `present`, so an unset value reads back as the type's zero — don't rely on a stale value surviving an `Unset`.
- **`MarshalJSON` emits the bare value or `null`; there is no `omitempty`-style elision.** A null value marshals to the literal `null`, and a present zero marshals to `0`/`false`/`0`. To omit a null field from JSON output entirely, the *containing* struct must handle that — this type always renders something.
- **`UnmarshalJSON` is fuzzed** (`fuzz_test.go`). Keep the fuzzer green when touching the parse path; it guards against panics on malformed numeric/boolean JSON.
