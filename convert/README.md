# convert

Best-effort coercion between arbitrary Go types. Scalar converters (`Bool`, `Int`, `Int32`, `Int64`, `Float`, `String`, `Time`, `Bytes`) plus container converters (`SliceOfString`, `SliceOfInt`, `MapOfAny`, …) and generic pointer helpers (`Pointer[T]`, `Element[T]`). Each scalar has plain, `Default`, and `…Ok` variants. Part of [rosetta](../README.md).

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/convert)

## What matters here

- **The `…Ok` suffix means "lossless?", NOT "succeeded?".** `StringOk`/`IntOk`/etc. return `(value, ok)` where `ok` is `true` only when the converted value round-trips back to the original input. A conversion that rounds, truncates, or falls back to the default returns a usable value with `ok=false`. Use the `Ok` form when you must distinguish an exact conversion from a lossy one.
- **`String` formats floats with a FIXED two decimal places — this is lossy.** `convert.String(3.14159)` yields `"3.14"`, and `StringOk` reports `ok=false` because it doesn't round-trip. This bites anywhere a float is rendered through `convert.String` (and transitively through `mapof`/`hannibal` string getters). It is intentional, not a bug; reach for `strconv.FormatFloat` directly when you need full precision.
- **The plain and `Default` forms swallow the lossless flag.** `String(v)` is `StringOk(v, "")` discarding the bool; `StringDefault(v, d)` substitutes `d` only when the value can't be converted at all. Neither tells you the result was rounded — only the `Ok` form does.
- **Converters recognize rosetta's getter interfaces** (`Booler`, `Inter`, `Floater`, `Stringer`, `Hexer`, plus `io.Reader` and `reflect.Value`). A type that implements one is converted via that method, so custom types flow through `convert` without a dedicated case. The type switch is order-sensitive — adding a case can change which interface a value matches first.
- **For typed pointer/deref helpers prefer `convert.Pointer[T]`/`Element[T]` over the reflection-based [pointer](../pointer/) package** when the type is known at compile time; `convert` owns the generic, type-safe versions.
