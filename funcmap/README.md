# funcmap

A registry of helper functions for Go `text/template` / `html/template`. `All()` returns a `map[string]any` (a `template.FuncMap`) covering compare, arrays, date, currency, HTML, math, logic, and string helpers. Part of [rosetta](../README.md).

Wire it in with `template.New("…").Funcs(funcmap.All())`.

## What matters here

- **Template functions report-and-swallow their errors; they do not return them.** Go template functions that return an `error` abort rendering, so helpers like `json`, `jsonIndent`, and `markdown` instead call `derp.Report(...)` on failure and return an empty/safe string. This is deliberate — a single bad value should not blow up a whole page render. Don't "fix" these to return errors.
- **`All()` rebuilds the map on every call.** It is meant to be called once at template-parse time, not per-request in a hot path. Each call allocates a fresh map and re-registers every category.
- **HTML-producing helpers depend on the [html](../html/) package's escaping/sanitizing.** Output safety (e.g. `markdown`, the HTML helpers) is only as strong as that package — changes to escaping there affect what these template funcs emit.
