# schema/format

String-format validators for the [schema](../) package. Each exported function (`Email`, `Color`, `ObjectID`, `IPv4`, `IPv6`, `Hostname`, `URI`, `Date`, `DateTime`, `Time`, `ISO8601`, `HTML`, `NoHTML`, `In`, `NotIn`, `MatchRegex`, `Token`, `Username`, `WebFinger`, `UnsafeAny`) is a `Generator` — it takes a configuration `arg` and returns a `StringFormat` closure that validates (and may rewrite) a string. See the parent [schema README](../README.md).

## What matters here

- **A format is registered via `schema.UseFormat`, which is only safe during `init`/startup.** The registry freezes the first time validation reads from it; registering a format afterward is rejected (reported, then ignored). Register all custom formats before any schema validation runs.
- **Empty string is always allowed.** Every validator returns early on `""` without error — "required" is enforced by the schema element (`String.Required`), not by the format. A format only constrains the *shape* of a non-empty value.
- **A `StringFormat` returns `(value, error)` and may REWRITE the value.** Validators are part of a coercion pipeline (see `validate_String_Formats` in the parent package), so returning a changed string is normal, not a bug — the caller tracks the change.
- **Many `Generator`s ignore their `arg`** (e.g. `Email`, `IPv4`); the parameter exists to satisfy the uniform `Generator` signature. `In`/`NotIn`/`MatchRegex` are the ones that actually consume `arg` (the allowed set / pattern).
- **The validators here are the trust boundary for untrusted strings.** `NoHTML`/`HTML` gate markup (and lean on [../../html](../../html/) + bluemonday); `URI`/`Hostname`/`IP*` constrain network identifiers. Tighten here, not at call sites.
