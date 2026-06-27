# translate

A JSON-defined mapping engine that copies and transforms data from a source object into a target object. A `Pipeline` is an ordered list of `Rule`s; `pipeline.Execute(inSchema, inObject, outSchema, outObject)` runs them in order. Both objects are accessed through [schema](../schema/) Getter/Setter interfaces, so the easiest carriers are [mapof.Any](../mapof/) values, which implement those interfaces already. Part of [rosetta](../README.md).

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/translate)

## Rule types

Each rule is a JSON object whose keys select the rule kind:

- **`value`** — write a static value to the target; ignores the source. `{"value":"FIXED", "target":"target.path"}`
- **`path`** — copy a value from a source path to a target path, unchanged. `{"path":"source.path", "target":"target.path"}`
- **`expression`** — run a Go template against the source and write the result. `{"expression":"{{ … }}", "target":"target.path"}`
- **`append`** — append a value to a slice/collection at the target path. `{"append":"VALUE", "target":"target.path"}`
- **`if`** — evaluate a Go template; run `then` rules when it returns `"true"`, otherwise `else`. `{"if":"{{ … }}", "then":[…], "else":[…]}`
- **`forEach`** — loop a source map/array, running `rules` for each item under the target path (optionally filtered). `{"forEach":"source.path", "target":"target.path", "filter":"{{ … }}", "rules":[…]}`
- **`first`** — run a list of rules, stopping after the first to set a non-zero value at the target. `{"first":"target.path", "rules":[…]}`

## What matters here

- **The rule kind is chosen by which key is present, not by an explicit `type` field.** `Rule.UnmarshalMap` checks for `append`, `value`, `path`, `expression`, `if`, `forEach`, `first` in turn. A rule map carrying two of these keys is ambiguous — the first matched dispatch wins. Author one rule kind per object.
- **Source and target must each implement the schema Getter/Setter interfaces.** Plain structs work only if they implement those interfaces (see [schema](../schema/)); a `mapof.Any` is the path of least resistance because it implements them out of the box.
- **`expression`, `if`, and `forEach`'s `filter` are Go templates evaluated against the SOURCE object.** They can read anything in the source but write only via their rule's target — a template that "returns true" for `if` must emit the literal string `"true"`.
- **Writes go through `schema.Set`/`Append`, which validate and may coerce.** A translated value can be clamped, truncated, or rewritten by the target schema's rules (see the [schema](../schema/) README) — the target schema is the final authority on stored values, not the rule.
- **`NewFromJSON` parses untrusted JSON and is fuzzed** (`FuzzNewFromJSON`). Keep that fuzzer green when changing rule unmarshalling.
