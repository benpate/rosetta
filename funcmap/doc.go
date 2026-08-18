// Package funcmap provides a registry of helper functions for Go templates.
//
// All returns a map[string]any suitable for passing to Template.Funcs, covering
// comparison, arrays, dates, currency, HTML, math, logic, and string helpers.
// It rebuilds that map on every call, so call it once when templates are
// parsed rather than once per request.
//
// These helpers report their errors rather than returning them. A template
// function that returns a non-nil error aborts the whole render, so a helper
// that fails logs through derp and yields an empty or otherwise safe value
// instead — one bad field should not cost the reader the entire page.
package funcmap
