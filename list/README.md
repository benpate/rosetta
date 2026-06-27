# list

Treats a *delimited string* as a list, without ever allocating a `[]string`. Generic free functions (`Head`, `Tail`, `First`, `Last`, `Split`, `SplitTail`, `At`, `PushHead`, `PushTail`, `RemoveLast`, `Index`, `LastIndex`, `IsEmpty`) operate on any `~string | []byte` value plus a delimiter byte. Named wrapper types (`Comma`, `Dot`, `Slash`, `Semicolon`, `Space`) bind a delimiter so you can call the same operations as methods. Part of [rosetta](../README.md).

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/list)

## What matters here

- **This is string slicing, not slice building — there is no `[]string` in sight.** `Head`/`Tail` walk the string to the next delimiter and return substrings; iterating a list is a `Head`/`Tail` loop, not a `for range` over a slice. This avoids allocation for the common "peel off the first path segment" case (e.g. routing on a `/`-delimited path). If you actually want a `[]string`, use `strings.Split` instead.
- **The delimiter is a single `byte`, not a string.** Multi-character delimiters aren't supported by design; the named types (`Comma` = `,`, `Dot` = `.`, `Slash` = `/`, …) each pin one byte via a `Delimiter…` constant.
- **`Tail`/`SplitTail`/`PushHead`/`PushTail` return the list's own type `T`, while `Head`/`First`/`Last`/`At` return a plain `string`.** The "remaining list" stays in the list type so you can keep chaining; an extracted single item drops to `string`. Mind the return types when composing them.
- **The named types are immutable value types.** `PushHead`/`PushTail`/`RemoveLast` return a *new* list rather than mutating the receiver — assign the result.
