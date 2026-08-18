# first

The "coalesce" helper: `first.String(...)`, `first.Int(...)`, and `first.Int64(...)` each return the first non-zero/non-empty value from their variadic arguments, or the type's zero value if every argument is empty. Part of [rosetta](../README.md).

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/rosetta/first.svg)](https://pkg.go.dev/github.com/benpate/rosetta/first)

Useful for fallback chains — `first.String(userValue, configValue, defaultValue)`. There's nothing subtle here; it's three short loops.
