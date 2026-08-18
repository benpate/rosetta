# null

Nullable types for Go: the scalars `null.Bool`, `null.String`, `null.Int`, `null.Int64`, and `null.Float`, plus the generic `null.Object[T]` for everything else. Each is a struct carrying a value plus a `present` flag, so it distinguishes "absent/null" from "present-but-zero". All satisfy the `Nullable` interface (`IsNull() bool`, aliased as `IsNil()`) and implement JSON marshal/unmarshal. Part of [rosetta](../README.md).

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/rosetta/null.svg)](https://pkg.go.dev/github.com/benpate/rosetta/null)

```go
// A zero-value null.Bool is null and ready to use — no constructor needed
var b null.Bool

b.Set(true)   // now present, value true
b.Bool()      // read the value
b.IsNull()    // false
b.Unset()     // back to null

b := null.NewBool(true) // or construct a non-null value directly

// null.Object[T] wraps any type with the same API
var o null.Object[time.Time]

o.Set(time.Now()) // now present
o.Object()        // read the value (T)
o.IsNull()        // false

o := null.NewObject(myStruct) // or construct a non-null value directly
```

### Null, zero, and present

A plain `int` or `bool` cannot tell you whether it was ever assigned. That is the entire reason this package exists: the zero value of every type here *is* a valid null, immediately usable, and `IsNull()` reports it as such until someone calls `Set()`.

Three predicates answer three different questions. `IsNull()` (and its alias `IsNil()`) asks *was this ever set*. `IsPresent()` is its inverse. `IsZero()` asks *is this empty* — it is TRUE for a null value **and** for a present zero, which makes it the right test for "nothing meaningful here" and the wrong test for "unassigned". `Interface()` rounds this out for code that works in `any`: it returns the value when present and an untyped `nil` when not.

`Unset()` returns a value to null, clearing the stored value along with the flag.

### JSON

These types are strict on the way in and plain on the way out.

Marshaling emits either the bare value or the literal `null`. There is no `omitempty`-style elision — a null field always renders as `null`, and a present zero always renders as `0`/`false`/`""`. Dropping the field entirely is the containing struct's decision, not this package's.

Unmarshaling never coerces across JSON types. A bare `123` is an error for `null.String`, not the text `"123"`; a quoted `"123"` is an error for `null.Int`. This is the deliberate opposite of the sibling [lenient](../lenient/README.md) package, which exists for the receiving half of Postel's law — reach for `lenient` when a remote peer's encoding is out of your control, and for `null` when the schema is yours to enforce.

`null.Object[T]` is the exception in one respect: it hands `T` to `encoding/json` rather than parsing it itself, so `T` marshals and unmarshals exactly as it would anywhere else, and only the `null` literal is special-cased.

### Fuzzing

All six `UnmarshalJSON` implementations are fuzzed in [fuzz_test.go](fuzz_test.go). Each target asserts that decoding arbitrary bytes never panics, and that any input the type *accepts* marshals back to bytes that decode to an equal value. Keep them green when touching a parse path.
