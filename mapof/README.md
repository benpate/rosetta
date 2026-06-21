# mapof

Generic, map-backed types (`mapof.Any`, `mapof.Bool`, `mapof.Int`, `mapof.Int64`, `mapof.Float`, `mapof.String`) that add type-safe getters/setters and schema integration on top of a plain Go map. Part of [rosetta](../README.md).

These exist so a `map[string]any` can be passed through the [schema](../schema/) engine (and other rosetta tools) with typed accessors instead of repeated type assertions at every call site.

## What matters here

- **Setters use pointer receivers and lazily allocate via `makeNotNil()`.** A zero-value `var m mapof.Bool` is a nil map; calling `m.SetBool(...)` works because the setter takes `*Bool` and allocates the backing map on first write. This is the standard "guard the nil-map write" pattern — do not "simplify" setters to value receivers, and do not add a constructor requirement. Getters use value receivers (reading a nil map is safe).
- **nilaway reports false positives here.** It flags the `makeNotNil()` dereference as a potential nil panic because it cannot prove the addressable receiver is non-nil through the pointer (e.g. `var m Bool; m.SetBool(...)` takes `&m`). The `-race` tests cover these paths and pass; the findings are spurious — verify against the tests, not the linter.
- **`mapof.Any` implements the `schema` getter/setter interfaces** (`GetPointer`, `SetObject`, `SetValue`, the typed `GetXxxOK` family). This is the seam that lets schema walk and mutate a generic map; changing those method signatures breaks schema path access.
- **`GetXxx` returns the zero value; `GetXxxOK` returns `(value, ok)`.** Use the `OK` form when "absent" and "present-but-zero" must be distinguished — the plain getters cannot tell them apart.
