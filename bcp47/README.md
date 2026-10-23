# bcp47

Turns any [BCP 47](https://www.rfc-editor.org/info/bcp47) language tag into the one spelling that should be stored, so that `EN-us`, `en_US`, and `en-US` are never three different languages.

| Function | What it does |
| -------- | ------------ |
| `Canonical(value)` | Returns the canonical spelling (`EN-us` → `en-US`, `iw` → `he`), dropping extensions and private-use subtags. Returns `false` for anything that names no language, including `und`, `mul`, and `*`. |
| `Parse(value)` | Returns the same result as a `golang.org/x/text/language` `Tag`, for callers that match or describe it. |

Values longer than 64 bytes are refused before parsing. Variants stay, because `ca-valencia` is a different written form of Catalan; extensions such as `-u-ca-gregory` go, because they describe formatting preferences.

The schema package's `language` string format is built on `Canonical`, so `Set` stores `EN-us` as `en-US`, keeps an empty value, and refuses any other value that is not a language.
