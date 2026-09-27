package loose

import (
	"sync"
	"sync/atomic"
)

// templateCache holds the templates that CachedTemplate compiles, keyed by source, so that each one
// is compiled once for the life of the process.
var templateCache sync.Map

// templateCacheSize counts the entries in templateCache, which sync.Map cannot report.
var templateCacheSize atomic.Int64

// CachedTemplate returns the Template for a source string, compiling it only if the cache does not
// already hold it. Use it for template strings that are rendered repeatedly.
func CachedTemplate(source string) Template {

	// Use the cached template, if there is one
	if cached, ok := templateCache.Load(source); ok {
		if compiled, ok := cached.(Template); ok {
			return compiled
		}
	}

	compiled := NewTemplate(source)

	// Plain text costs nothing to wrap, and would crowd real templates out of the cache
	if !compiled.IsTemplate() {
		return compiled
	}

	return storeTemplate(source, compiled)
}

// storeTemplate adds a compiled template to the cache if there is room, and returns the
// template the cache holds for its source, or the compiled template if the cache is full.
func storeTemplate(source string, compiled Template) Template {

	// RULE: The cache never grows past its limit, so arbitrary strings cannot exhaust memory
	if templateCacheSize.Add(1) > templateCacheLimit {
		templateCacheSize.Add(-1)
		return compiled
	}

	// Another goroutine may have compiled the same source first; keep whichever was stored
	existing, loaded := templateCache.LoadOrStore(source, compiled)

	if !loaded {
		return compiled
	}

	templateCacheSize.Add(-1)

	if stored, ok := existing.(Template); ok {
		return stored
	}

	// Unreachable: only this file writes to the cache, and it writes only Templates
	return compiled
}
