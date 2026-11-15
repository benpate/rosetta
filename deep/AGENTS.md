# deep

Agent notes for the `deep` package. See [README.md](README.md) for what the package is and how to use it.

## Unknown types are an error in both functions, on purpose

**`Copy` and `Remove` refuse a pointer, struct, function, or channel at any depth.** `Remove` strips `bto` and `bcc` before delivery, so skipping a container it cannot search would send a blind recipient out with no warning. Do not make either function pass an unknown value through. `time.Time` is the one struct both accept, as a simple value. A `*sliceof.Any`, which `schema` stores for new lists, is refused like any other pointer.

## The type lists in Copy and Remove must agree

**A type added to one switch belongs in the other.** They drifted once: `mapof.String` was copied but not searched. The reflection fallback now catches a type that is missing from a list, so a gap costs speed and not correctness, but only while both fallbacks handle the same kinds.

## Reflection covers the types rosetta cannot import

**`copyReflected` and `removeReflected` handle every map, slice, and array by kind.** That is how the MongoDB driver's `primitive.A`, `M`, `DateTime`, and `ObjectID` are supported without importing it. The tests use look-alike types in [testTypes_test.go](testTypes_test.go) for the same reason. `primitive.D` is a slice of structs, so it is refused.

## A copy keeps every container's type

**`copyMap` and `copySlice` take `~map` and `~[]` constraints so that they build the caller's own type.** With a plain `map[K]V` parameter, Go infers the underlying type and a `mapof.Any` comes back as a `map[string]any`. That then fails the type assertion in `Copy` for any typed container that holds one, such as a `sliceof.MapOfAny`.

## Copy is generic, and copyValue does the work

**`Copy[T]` asserts the result of `copyValue` back to the caller's type; every recursive call goes through one of the two.** `copyValue` is not named `copy`, because that would shadow the builtin. A `Copier` that returns a different type than its own is an error wherever the caller's type is concrete, and passes through where it is `any`.

## A nil item is stored as the zero value

**`Copy` and `copyReflectedItem` check for a nil copy before they assert its type.** A nil interface matches no type assertion, even `.(any)`, so without the check every JSON `null` inside a container fails the copy.

## Nil containers stay nil

**A nil map or slice copies to a nil of the same type, as `maps.Clone` and `slices.Clone` do.** The simple-type cases call those two functions directly, and `copyMap`, `copySlice`, and the reflected copies guard for nil to match them.

## Errors name the type, never the value

**Every error detail goes through `typeName`.** The values here are documents that hold blind recipients, and Emissary stores every reported error. On an error, `Copy` returns the zero value of the caller's type (nil for `any`) and never a partial copy.

## A self-referencing map or slice recurses until the stack overflows

**Neither function detects a cycle.** Refusing pointers does not prevent one, because maps and slices are references: `m["self"] = m` is a cycle with no pointer in it. A JSON or BSON decoder cannot build one, so only data assembled in Go can. A stack overflow is fatal and no `recover` catches it.

## Remove matches whole names only

**`Remove` has no prefix or predicate form yet.** Emissary's `assanitizer.Strip` removes by prefix, so `deep.Remove` can replace `assanitizer.StripKeys` but not `Strip`. Adding one is planned.

## deep imports mapof and sliceof

**Neither `mapof` nor `sliceof` may import `deep`,** which would be a cycle. `mapof.Template.Clone` keeps its own `deepCopy` for that reason.
