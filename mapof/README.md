# mapof

Generic, map-backed types (`mapof.Any`, `mapof.Bool`, `mapof.Int`, `mapof.Int64`, `mapof.Float`, `mapof.String`) that add type-safe getters/setters and schema integration on top of a plain Go map. Part of [rosetta](../README.md).

These exist so a `map[string]any` can be passed through the [schema](../schema/) engine (and other rosetta tools) with typed accessors instead of repeated type assertions at every call site.

## What matters here

- **Setters use pointer receivers and lazily allocate via `makeNotNil()`.** A zero-value `var m mapof.Bool` is a nil map; calling `m.SetBool(...)` works because the setter takes `*Bool` and allocates the backing map on first write. This is the standard "guard the nil-map write" pattern — do not "simplify" setters to value receivers, and do not add a constructor requirement. Getters use value receivers (reading a nil map is safe).
- **The nilaway CI step runs with `-exclude-test-files`.** `.github/workflows/nilaway.yml` invokes `nilaway -exclude-test-files=true ./...` directly, because the `qbaware/nilaway-action` hardcodes `nilaway $PACKAGE_TO_SCAN` and cannot pass the flag. The flag scopes analysis to production code; nil flows that exist only in test fixtures (e.g. a setter called on a zero-value map variable) are out of scope. The gate must exit with zero output — if it reports anything, fix the code, do not relax the gate.
- **`mapof.Any`, `mapof.Template`, and `mapof.Object[T]` implement the `schema` getter/setter interfaces** (`GetPointer`, `SetKey`, `SetValue`, the typed `GetXxxOK` family) without importing schema. Schema walks every nested path itself and stores one key at a time through `SetKey`, so it can create `mapof.Any` and `sliceof.Any` children; importing schema from here would make that an import cycle. `GetPointer` on these maps returns the child value itself, which schema reads before it writes into a child. Changing those method signatures breaks schema path access.
- **`GetXxx` returns the zero value; `GetXxxOK` returns `(value, ok)`.** Use the `OK` form when "absent" and "present-but-zero" must be distinguished — the plain getters cannot tell them apart.
