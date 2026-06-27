# first

The "coalesce" helper: `first.String(...)`, `first.Int(...)`, and `first.Int64(...)` each return the first non-zero/non-empty value from their variadic arguments, or the type's zero value if every argument is empty. Part of [rosetta](../README.md).

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/first)

Useful for fallback chains — `first.String(userValue, configValue, defaultValue)`. There's nothing subtle here; it's three short loops.
