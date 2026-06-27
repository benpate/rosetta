# html

Utilities for manipulating HTML content: `FromText` (plain text → lightly-formatted, escaped HTML), `ToText`/`ToSearchText` (strip markup back to text), `Minimal` (sanitize to a small allow-list of tags), `Summary` (a short plain-text excerpt), `RemoveTags`/`RemoveAnchors`/`RemoveSpecialCharacters`, `CollapseWhitespace`, and `IsHTML`. Part of [rosetta](../README.md).

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/html)

## What matters here

- **`FromText` escapes `&` FIRST, before the other entity replacements.** Order is load-bearing: escaping `<`/`>`/`"`/`'` introduces `&…;` entities, so if `&` were escaped after them it would double-escape those entities into `&amp;lt;`. Don't reorder the `ReplaceAll` chain in `fromText.go`.
- **`Minimal` is the only sanitizer, and it's built on a bluemonday `UGCPolicy` plus an explicit element allow-list** (`br p b i u img div pre code ol ul li`). It's the trust boundary for untrusted markup — widen the allow-list with care, and treat anything *not* passed through `Minimal` as unsanitized. `FromText` only escapes; it does not sanitize existing tags.
- **`Summary` truncates by RUNE, not byte (200 runes + ellipsis).** It strips tags and collapses whitespace first, so the cut never lands mid-multibyte-character. The rune-boundary behavior is fuzz-tested — keep `FuzzSummary` green when changing the truncation.
- **The text-stripping functions are fuzzed against malformed input** (`FuzzFromText`, `FuzzRemoveTags`, `FuzzCollapseWhitespace`, `FuzzHTMLDecoders`). When changing any parsing/stripping path, run the fuzzers — they guard against panics and decode mismatches on garbage HTML.
