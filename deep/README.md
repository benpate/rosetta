# deep

Functions that reach every level of nested data: `deep.Copy` makes a copy that shares nothing with the original, and `deep.Remove` deletes named properties wherever they appear. Both are written for the maps and slices that a JSON or database decoder returns. Part of [rosetta](../README.md).

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/rosetta/deep.svg)](https://pkg.go.dev/github.com/benpate/rosetta/deep)

A shallow copy such as `maps.Clone` still shares every nested map and slice, so a change to the copy's `object` changes the original's too. `deep.Copy` shares none of them, which makes it safe to change a document that something else still reads.

```go
// Copy the activity, so that the stored original keeps its blind recipients
outbound, err := deep.Copy(activity)

if err != nil {
	return err
}

// Remove `bto` and `bcc` from the activity, its object, and everything inside them
if err := deep.Remove(outbound, "bto", "bcc"); err != nil {
	return err
}
```

### What can be copied and searched

Booleans, numbers, strings, `time.Time`, and every map, slice, and array of them, to any depth. That includes the named types in [mapof](../mapof/README.md) and [sliceof](../sliceof/README.md), and types this package has never heard of, such as the MongoDB driver's `primitive.M` and `primitive.A`. A copy keeps the type of every container, and a nil map or slice stays nil.

`deep.Copy` returns the type it was given, so a `map[string]any` comes back as a `map[string]any` with no type assertion. An untyped `nil` has no type to return, so a literal is written as `deep.Copy[any](nil)`.

### What is refused

A pointer, struct, function, or channel is an error in both functions, wherever it is nested. A copy that still shared a pointer would not be a deep copy, and a property left behind inside a struct would go unnoticed. A type can opt in by implementing `deep.Copier` or `deep.Remover`.

`deep.Remove` matches whole, case-sensitive property names. It changes the value it is given, so copy first when the original must survive.
