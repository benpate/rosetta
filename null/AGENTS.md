# null

Agent notes for the `null` package. See [README.md](README.md) for what the package is and how to use it.

## IsZero is not an alias for IsNull

`IsZero()` returns TRUE for a null value **and** for a present zero; `IsNull()`/`IsNil()` return TRUE only for null. Code that means "was this ever set" must use `IsNull()`. Do not "simplify" one into the other — they answer different questions, and the tests pin both.

## compare.IsZero answers differently than the types' own IsZero

[compare/zero.go](../compare/zero.go) detects the `Nuller` interface and calls `IsNull()`, so a present-but-zero `null.Int` is *not* zero to `compare` while it *is* zero to its own `IsZero()`. Expect the two to disagree; that is the design, not a bug.

## Object[T].IsZero follows Go's definition of zero, not "emptiness"

It reflects over `T` and calls `reflect.Value.IsZero()`, so an empty-but-non-nil slice or map reads as **not** zero. Only a nil slice/map, or an untyped nil when `T` is an interface, counts.

## Object[T] cannot represent a present nil pointer

A present `*Foo(nil)` marshals to `null`, which reads back as *unset*. The round trip silently converts "present, holding nil" into "absent". Don't build logic that depends on telling those apart through JSON.

## The scalars never coerce across JSON types

`String` rejects `123` and `true`; `Int`/`Int64`/`Float` reject `"123"`; `Bool` accepts only the literals `true`, `false`, and `null`. Widening any of these to "be helpful" breaks the package's contract — [lenient](../lenient/README.md) is the package that tolerates loose encodings, and the split between the two is deliberate.

## String marshals through encoding/json, not strconv.Quote

Load-bearing: it keeps escaping, invalid UTF-8, and HTML-sensitive runes identical to a plain string field. `string.go` carries a `RULE:` comment saying so. Swapping in `strconv.Quote` would change the wire format.

## Empty input means null for every type except Bool

`Int`, `Int64`, `Float`, `String`, and `Object[T]` treat `[]byte("")` the same as `null`. `Bool` has no such case and errors on it. If you add a type, match the majority; if you touch `Bool`, know the asymmetry is currently unpinned by any test.

## Float accepts NaN and Inf on a direct UnmarshalJSON call, and re-emits them

`Float.UnmarshalJSON` uses `strconv.ParseFloat`, which accepts `NaN`, `Inf`, and `+Infinity` — none of which are valid JSON. `MarshalJSON` then writes `f.String()` as raw bytes, so the value round-trips back out as the literal `NaN` and poisons the enclosing document. This is unreachable through `encoding/json` (the outer decoder rejects those tokens first), so it only bites code that calls `UnmarshalJSON` directly.

## Float always formats without an exponent

`Float.String()` uses `strconv.FormatFloat(value, 'f', -2, 64)` — shortest round-trippable digits, but `'f'` format, never `'e'`. A very small or very large magnitude renders as a long run of digits rather than exponent notation. `MarshalJSON` writes that string verbatim, so the format choice is the wire format.

## Three copies of the Nullable interface exist on purpose

`null.Nullable`, [schema.Nullable](../schema/interfaces.go), and [compare.Nuller](../compare/interfaces.go) all declare `IsNull() bool`. They are duplicated so those packages don't import `null`; don't consolidate them into one shared declaration.

## Keep the fuzz corpus green

[fuzz_test.go](fuzz_test.go) fuzzes all six `UnmarshalJSON` methods with a shared `assertRoundTrip` helper. Any new nullable type should get a target there and inherit the same no-panic and round-trip properties.
