# compare

Comparison helpers for values of different or unknown data types. `Interface(a, b)` coerces two arbitrary values to a common kind and returns `-1`/`0`/`1`; `WithOperator(a, op, b)` evaluates a named operator (`greater-than`, `contains`, `begins-with`, …) to a bool. Typed helpers (`Int`/`Int8`/…/`Int64`, `Float32`/`Float64`, `String`, `BeginsWith`, `Contains`, `EndsWith`) cover the known-type cases. Part of [rosetta](../README.md).

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/compare)

## What matters here

- **A `bool` on either side coerces the *other* operand to bool, and this runs before the numeric block.** `Interface(true, 1)` and `Interface(1, true)` both succeed symmetrically, and a bool is never treated as a number. Know this before comparing mixed bool/number data — the bool wins.
- **Numbers compare on a common numeric kind, not per-type-pair.** Any signed/unsigned/float combination is handled symmetrically through a single numeric path, so you don't get spurious "incompatible types" between, say, an `int` and a `float64`.
- **Incompatible types return an error, not a silent `0`.** `Interface` (and `WithOperator`) return a non-nil error when the two values can't be coerced to a common comparable kind — callers must check it rather than assuming `0` means "equal".
- **`BeginsWith`/`Contains`/`EndsWith` are substring operators handled outside the ordering path** (see `WithOperator`). They answer a containment question, not a `-1/0/1` ordering, so they're routed before the `Interface` coercion runs.

This package works around Go's pre-generics comparison gaps and will likely shrink as the standard library grows; lean on the typed helpers when the type is known and reserve `Interface` for genuinely dynamic data.
