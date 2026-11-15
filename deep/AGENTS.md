# deep

Agent notes for the `deep` package. See [README.md](README.md) for what the package is and how to use it.

## Unknown types are an error in both functions, on purpose

**`Copy` and `Remove` refuse a pointer, function, channel, or a struct that is not simple, at any depth.** `Remove` strips `bto` and `bcc` before delivery, so skipping a container it cannot search would send a blind recipient out with no warning. Do not make either function pass an unknown value through. The error is `derp.Internal`, because a value of the wrong type is the program's mistake and not the client's.

## Simple types are copied by value and never searched

**`isSimpleType` accepts a boolean, number, or string, and any array or struct made only of those.** Such a value holds no map, slice, pointer, or interface, so nothing in it can be shared or hide a property. This is how the MongoDB driver's `ObjectID`, `Decimal128`, `Timestamp`, and `Regex` pass. A type that implements `Copier` or `Remover` is never simple, so its own method is still called inside a typed slice. `time.Time` holds a pointer, so it has its own case in both switches.

## A *sliceof.Any is the one pointer both functions follow

**It has its own case in both switches because `schema` stores every new list as one.** `Copy` returns a new pointer to a copied slice. Every other pointer is refused.

## Map keys are copied as they are, and never checked

**`copyReflectedMap` deep-copies a map's values and reuses its keys, and `Remove` never looks inside a key.** A map keyed by a pointer, or by a struct that holds one, therefore copies without an error and the copy shares that pointer, which is the one gap in the rule above. No JSON or BSON decoder builds such a key. Closing it means refusing a key type that `isSimpleType` rejects, while still allowing the `any` keys of a `map[any]any`.

## The type lists in Copy and Remove must agree

**A type added to one switch belongs in the other.** They drifted once: `mapof.String` was copied but not searched. The reflection fallback now catches a type that is missing from a list, so a gap costs speed and not correctness, but only while both fallbacks handle the same kinds.

## Reflection covers the types rosetta cannot import

**`copyReflected` and `removeReflected` handle every map, slice, and array by kind.** That is how the MongoDB driver's `primitive.A`, `M`, `DateTime`, and `ObjectID` are supported without importing it. The tests use look-alike types in [testTypes_test.go](testTypes_test.go) for the same reason. `primitive.D` is refused: its items are structs that hold an interface, and `Remove` could not shorten a slice it was handed by value.

## A copy keeps every container's type

**`copyMap` and `copySlice` take `~map` and `~[]` constraints so that they build the caller's own type.** With a plain `map[K]V` parameter, Go infers the underlying type and a `mapof.Any` comes back as a `map[string]any`. That then fails the type assertion in `copyAs` for any typed container that holds one, such as a `sliceof.MapOfAny`.

## Copy is generic, and copyValue does the work

**`copyAs[T]` asserts the result of `copyValue` back to the caller's type; every recursive call goes through one of the two, and `Copy` only starts the walk.** `copyValue` is not named `copy`, because that would shadow the builtin. A `Copier` that returns a different type than its own is an error wherever the caller's type is concrete, and passes through where it is `any`.

## A nil item is stored as the zero value

**`copyAs` and `copyReflectedItem` check for a nil copy before they assert its type.** A nil interface matches no type assertion, even `.(any)`, so without the check every JSON `null` inside a container fails the copy.

## Nil containers stay nil

**A nil map or slice copies to a nil of the same type, as `maps.Clone` and `slices.Clone` do.** The simple-type cases call those two functions directly, and `copyMap`, `copySlice`, and the reflected copies guard for nil to match them.

## Errors name the type, never the value

**Every error detail goes through `typeName`.** The values here are documents that hold blind recipients, and Emissary stores every reported error. On an error, `Copy` returns the zero value of the caller's type (nil for `any`) and never a partial copy.

## The depth limit is the only cycle guard

**`copyValue` and `removeValue` compare their depth to `maxDepth` first, and every container passes `depth+1` to what it holds.** Refusing pointers does not prevent a cycle, because maps and slices are references: `m["self"] = m` has no pointer in it. Without the limit that value overflows the stack, which is fatal and which no `recover` catches. A helper that recurses without adding to the depth brings the crash back.

## Copier and Remover carry the depth

**Both methods receive the depth of the value they are called on, and must pass `depth+1` to `CopyWithDepth` or `RemoveWithDepth` for everything they hold.** A `Remover` passes the `Matcher` it was given along unchanged. An implementation that calls `deep.Copy` or `deep.Remove` on its children restarts the count at zero, and a cycle through it overflows the stack. `CopyWithDepth` and `RemoveWithDepth` are exported for this one purpose, and return their errors unwrapped.

## Errors are wrapped once, at Copy and Remove

**The recursive helpers return an error as it is, including one from a `Copier` or `Remover`, and only `Copy` and `Remove` call `derp.Wrap`.** Wrapping at every level would nest the error as deeply as the value, and the depth error would arrive 10,000 layers deep.

## A failed Remove leaves the value half-stripped

**`Remove` changes the value in place and stops at the first error.** Which names are already gone depends on map order. A caller must treat a value that failed as unusable, and never send it.

## Every removal runs on a Matcher

**`Remove` is `RemoveFunc` with `matchNames`, and the walk asks the `Matcher` about every string key.** A map of simple values is ranged over key by key, which costs more than deleting each name directly: 825 ns against 27 ns for a no-match `Remove` on a 100-key `map[string]string`. Give `Named` a fast path only with a benchmark that shows a caller needs it.

## Benchmark baseline

**[benchmark_test.go](benchmark_test.go) holds the baseline, measured 2026-11-14 on darwin/arm64 with go1.27.4.** Allocations are the number to compare, since they do not vary between runs.

| Benchmark | ns/op | allocs/op |
|---|---|---|
| `Copy_Listed` | 650 | 12 |
| `Copy_Reflected` | 1,700 | 45 |
| `Remove_Listed` | 255 | 2 |
| `Remove_Reflected` | 763 | 24 |
| `Remove_LargeStringMap` | 860 | 2 |

The listed types are why the two switches are long: reflection costs more than twice the time and allocations.

## The simple-type clause returns the value it was given

**`copyValue` lists every simple type in one clause and returns `value`, never a typed copy of it.** A clause per type that returned its own typed variable boxed each value into a new `any`, which cost 10 more allocations on the sample document (22 against 12).

## The long switches stay flagged by revive

**`copyValue` and `removeValue` exceed revive's length and cyclomatic limits, and that is accepted.** A clause that calls a generic helper needs its own typed variable, so those clauses cannot be merged, and moving them into helpers only moves the finding. Clauses that all do the same thing are merged into one.

## A Matcher sees the property name and nothing else

**It is not told the value, the path, or the depth.** A rule that needs one of those is a new function type, not a change to `Matcher`, because `Remover` implementations depend on its signature.

## Matchers refuse the inputs that would match everything or nothing by accident

**`HasPrefix` drops an empty prefix, and `RemoveFunc` and `RemoveWithDepth` return an error for a nil `Matcher`.** An empty prefix begins every name, so it would delete every property. A nil `Matcher` that removed nothing would look like a successful strip.

## Functions are ordered callers first

**Each file lists its entry point, then that function's callees in the order they are first called, level by level.** A new helper goes where that order puts it, not at the end of the file. `copierType` and `removerType` sit directly above `isSimpleType`, which is the only function that reads them.

## deep imports mapof and sliceof

**Neither `mapof` nor `sliceof` may import `deep`,** which would be a cycle. `mapof.Template.Clone` keeps its own `deepCopy` for that reason.
