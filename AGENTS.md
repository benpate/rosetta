# rosetta — Notes for AI Agents

Rosetta is a collection of data-manipulation packages (conversion, comparison, schema validation, typed maps/slices, template helpers) shared across the benpate/EmissarySocial ecosystem. See [README.md](README.md) for the package roster. The `null` package has its own notes in [null/AGENTS.md](null/AGENTS.md).

A recurring defect pattern in this repo: two implementations of one contract drifting apart (sibling groupers, sibling map types, two slice-growth routines). When you change one implementation of a shared contract, diff its siblings and keep them in agreement. And when adding fuzz targets, prefer differential tests against a stdlib oracle over "assert no panic", and bound any input that drives allocation.

## schema

- **`Schema.Set` validates before writing, and callers must expect errors.** It can return `derp.Validation` errors (unknown path, required-empty, pattern/enum violations), clamp integers to `Minimum`/`Maximum`, truncate strings to `MaxLength` on rune boundaries, and rewrite values through format functions. It is not a plain type-coercing setter; a Set that "obviously should work" can still fail or store a different value than was passed.
- **`Schema.Validate` returns `(bool, error)`.** The bool means "the value was coerced or rewritten during validation", not success. Element-level `Validate` methods are still error-only.
- **The per-type `Validate(any)` methods (Any, Boolean, Integer, Number, String) are intentional, not dead code.** They are outside the `Element` interface and have no in-repo callers, but they are deliberately retained. Do not delete them; keep their behavior and error messages consistent with the live `validate_*` free functions when changing validation.
- **Array-typed properties must implement `ArrayGetterSetter` via their POINTER.** `validate_Array` requires `GetIndex`/`SetIndex`/`Length`. Add a compile-time `var _ schema.ArrayGetterSetter = (*T)(nil)` assertion for every new slice type used in a schema. `addressableItem` in [schema/validate_array.go](schema/validate_array.go) exists so composite array items are validated against an addressable pointer copy (pointer-receiver accessors would otherwise be unreachable on by-value items); do not "simplify" it away.
- **Map-ness is declared, never inferred.** A map type opts in via the `MapTyper` interface (`IsMap() bool`); only then does `validate_Object` skip absent keys as "absent optional = empty". A struct missing a `GetStringOK`/`GetPointer` accessor still errors loudly, on purpose — that error is the signal that an accessor is missing.
- **All in-package inheritance goes through `inheritElement`, which allocates the child's `Properties` map before delegating.** `Object.Inherit` has a value receiver (the `Element` interface is satisfied by value types, so a pointer receiver would be a breaking change), which means a nil-map allocation inside `Object.Inherit` cannot stick. Calling `element.Inherit(parent)` directly on an Object with a nil `Properties` map silently does nothing.
- **The scalar `Inherit` no-ops report 0.0% coverage forever.** They have empty bodies, and `go tool cover -func` prints 0.0% for zero-statement functions. They are tested with `require.NotPanics`; don't chase them.
- **`validate_Object` ranges over a map, so multi-property validation failures are non-deterministic in order.** Run schema-heavy tests with `-count=3` when hunting flaky validation behavior.

## schema/format

- **Format functions rewrite values, and Set stores what they return.** A format is `func(string) (string, error)`; the returned string is what lands in the object, not the input.
- **The `webfinger` format cannot round-trip (open bug).** It requires input to start with `@` but stores the value with `@` stripped, so validating a stored value fails its own format. Do not put `Format: "webfinger"` fields in Set→Validate round-trip tests until the require/strip mismatch is resolved.
- **The `url` format is strict.** Absolute URLs only, scheme limited to http/https/mailto, max 2048 bytes. Unknown format names fail `Schema.ValidateFormats()`. Data written before this enforcement may fail validation on its next Set.
- **Validation regexps live at package scope.** Compiling a regexp inside the returned closure runs on every validation (a measured 233x slowdown). Keep patterns hoisted, like `hostnamePattern` in [schema/format/network.go](schema/format/network.go).
- **`Text`, `NoHTML`, and `Summary` inherit `html.RemoveTags`'s angle-bracket data loss.** See the html section; the known-bug pins live in [schema/format/Text_test.go](schema/format/Text_test.go).

## convert

- **`convert.String` formats floats with a fixed 2 decimal places — it is lossy.** `100` becomes `"100.00"`, `3.14159` becomes `"3.14"`. Never round-trip numeric data through `convert.String`; format explicitly when precision matters. `lenient.String` deliberately avoids this path to preserve exact source text.
- **The `*Ok` functions' second return means "the conversion was lossless", not "the value was naturally this type."** The non-Ok wrappers discard it.
- **`convert.Int` clamps to platform int width and clamps out-of-range floats.** Use `Int64` for wire data; a plain `Int` caps at 2^31 on 32-bit targets (including `GOARCH=wasm`).
- **There is no reflection fallback.** A named type (`type Foo string`, or any struct) that matches no concrete case or interface converts to `""`. This is why `mapof.Matchable[T].MapOfString()` loses every value for struct-typed `T` — pinned as a KNOWN LIMITATION in [mapof/matchable_test.go](mapof/matchable_test.go), not a bug to fix in mapof.

## list

- **Delimiters are effectively ASCII-only (open bug, BUG-115).** `PushHead`/`PushTail` join with `string(delimiter)` where delimiter is a byte, so Go rune-converts it: byte 200 is written as the two UTF-8 bytes `c3 88`, and `Index()` can never find it again. The constructors' `strings.Join(..., string(Delimiter))` calls share the flaw. Latent only because all five shipped list types (comma, dot, semicolon, slash, space) are ASCII. `FuzzList_Push` skips delimiters above 0x7f with a comment pointing here — delete that skip when the fix (`string([]byte{delimiter})`) lands.
- **`SplitTail`'s single-item behavior is load-bearing.** `"photo"` splits to `("photo", "")` — the whole value stays in the head, NOT the `RemoveLast`/`Last` decomposition `("", "photo")`. Real callers split `filename.ext` and depend on exactly this; its doc comment says so. Do not "fix" it.
- **`Last2`/`Last3`/`Last4` pad in the middle for lists shorter than N.** `Last3(["a","b"])` returns `("b", "", "a")` — they inherit `SplitTail`'s filename semantics, which makes them wrong for short lists. Zero known callers; the oddity is recorded in `TestLast_ShortListIsUnspecified`. Don't build on their short-list behavior.

## sliceof

- **A schema path segment can grow any sliceof slice without bound (open bug, BUG-116).** `SetIndex`/`GetPointer` parse an index out of the path with no maximum, then `growSlice` appends until it is addressable — `schema.Set(&obj, "tags.50000", "boom")` silently allocates 50,001 elements. Never pass attacker-controlled PATH keys (not just values) into `schema.Set`/`SetAll` on objects containing these types. The fuzz suite caps indexes at `growthLimit = 1000` for this reason; raise it once a bound exists.
- **Grouper `IsFooter` guards must be `index < 0`, matching `slice.Grouper`.** An `index <= 0` guard means the first group never closes. The four sliceof groupers and `slice.Grouper` implement one contract; keep them identical.
- **`FirstN` clamps negative n to 0.** The slice expression `x[:n]` panics otherwise; every sliceof type carries the clamp plus a regression test.

## mapof

- **Mutating methods take pointer receivers and call `makeNotNil()` first.** A value-receiver setter panics on a nil map (the zero value of every mapof type). `Slices.Add` follows the pattern; any new method that writes must too.
- **`Matchable[T].MapOfAny()` must copy entries directly — never hand `m` to `convert.MapOfAny`.** `Matchable[T]` is a `map[string]T`, which matches none of `MapOfAnyOk`'s concrete map cases (Go type switches match exact types, not underlying types) and falls through to its `MapOfAnyGetter` case, which calls straight back into the method: unbounded mutual recursion, fatal stack overflow, for every `T`. The in-code comment marks this; do not "simplify" it back.
- **`mapof.String.GetStringOK` returns `("", false)` for absent keys, on purpose.** Absent-key tolerance during schema validation is handled by `MapTyper`, not by loosening the getter contract. All mapof types implement `IsMap() bool` returning true, with compile-time assertions.

## funcmap

- **The trust-cast helpers each have a distinct trust contract; know it before piping data in.** `html` is a raw `template.HTML` cast — safety is per-call-site, so sanitize remote/federated content with `htmlMinimal` (bluemonday) first. `css` is a raw `template.CSS` cast reserved for owner/admin-authored stylesheets — never remote or query-string data. `attr` and `cssValue` are safe-by-construction: they reject any value containing context-breakout characters (returning empty) via `safeAttr`/`safeCSSValue`. `text` self-sanitizes through `html.FromText`; do not "harden" it further. `safeURL` guards navigation targets. When auditing templates, grep for `{{html`, `{{attr`, `{{css` and verify the piped value's trust.
- **`html/template` resolves function names at parse time.** A template calling a helper that a dependent's rosetta version lacks fails to PARSE, not just render. Adding a funcmap helper does nothing for downstream templates until they build against a rosetta that has it.

## html

- **`RemoveTags` mangles unescaped angle brackets in plain text (open bug, pinned in tests).** A bare `>` resets the write cursor and drops everything before it (`"6 > 5"` → `" 5"`); a doubled `<` duplicates preceding text (`"a<<br>b"` → `"aab"`), so output can exceed input length. It is NOT a sanitizer hole — no `<`/`>` survives, which `FuzzRemoveTags` proves. Reachable from user content via `schema/format.Text`/`NoHTML`/`Summary` and funcmap. The pins carry KNOWN BUG markers so a fix surfaces as a test failure; any fix must not let `<` through.
- **The fuzz file is [html/html_fuzz_test.go](html/html_fuzz_test.go), not `fuzz_test.go`.** A `fuzz*_test.go` glob misses it and you will write a duplicate.
- **When fuzzing sanitizer output, scope substring checks to tag spans only.** bluemonday escapes `<` in text to `&lt;`, so bare prose containing "onerror" is harmless; a naive `strings.Contains` assertion fails instantly on safe output.

## lenient

- **This package is the receiving half of Postel's law: scalars that accept whatever encoding a remote system actually sends.** Its strict counterpart is `null`, which never coerces across JSON types; the split is deliberate — don't widen `null` or tighten `lenient` to merge them.
- **`String` decodes with `json.Decoder.UseNumber()` to keep a number's exact source text** (`1.0` stays `"1.0"`). Never route it through `convert.String`, which would force two decimals.
- **`Int64` is `Int64`, not `Int`, on purpose, and its exact-integer fast path is load-bearing.** `convert.Int` clamps to platform int width (2^31 on 32-bit/wasm targets), and `json.Unmarshal` into `any` yields `float64`, silently rounding above 2^53 — the fast path preserves large integers. Precision beyond int64 is explicitly NOT a goal: larger values clamp, and integers nested inside arrays still round through `float64`.
- **A new tolerant type goes into `jsonTargets` in [lenient/fuzz_test.go](lenient/fuzz_test.go)** to inherit the universal properties (no panic, valid re-encode, fixed-point round trip). Gate cross-arch with `GOOS=linux GOARCH=386 go vet ./...` (plus `arm`, `js/wasm`).

## null

See [null/AGENTS.md](null/AGENTS.md) for this package's rules (IsZero vs IsNull, strict no-coercion JSON, the encoding/json marshal contract).

## delta

- **`Slice.SetValue` must copy the incoming and previous slices; keep the append-onto-fresh-slice pattern.** The Deleted list is compacted in place, which would write through to the caller's slice if aliased, and `slices.Clone(nil)` returns nil while `NewSlice`/`Reset` guarantee non-nil lists — the `append(make([]T, 0, ...), ...)` form satisfies both constraints at once.
- **`NewSlice` stores its variadic slice directly, so a later `SetIndex` can write through to caller memory.** Known residual aliasing; be aware of it when a caller retains the slice it passed in.
- **`Slice.SetIndex` shares sliceof's unbounded-growth defect (BUG-116).** It rejects negative indexes then grows to any requested length, and it implements `schema.ArraySetter`, so the same schema-path route reaches it. It cannot share `sliceof.growSlice` (different growth bookkeeping — appended slots are recorded as Added), so any bounding fix must cover both packages.

## channel

- **The package has open, structural defects — do not build new code on it without accepting them.** `Closed()` receives from the channel in its select, so it reports "closed" whenever a value is merely waiting and DISCARDS that value (its doc comment claims otherwise); Go offers no non-destructive closed test, so this cannot be fixed in general. Every producer (`Map`, `Pipe`, `Filter`, ...) sends without a cancellation select, so abandoned pipelines leak goroutines; even `PipeWithCancel` leaks, because its `output <- item` send is unguarded and cancellation never reaches upstream stages.

## Tooling

- **[revive.toml](revive.toml) sets `enableAllRules = true` with targeted overrides.** The high-signal revive categories in this repo are `suspicious assignment`, `identical branches`, `differs only by capitalization`, and `redefinition of the built-in`; most style findings on validators and test constants are accepted noise.
- **Check nilaway against the CI form before calling it red or green.** CI runs `nilaway -exclude-test-files=true` (see `.github/workflows/`); a plain local `nilaway ./...` adds findings from tests that deliberately call map setters on nil maps.
