package schema

import (
	"sync/atomic"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema/format"
)

// formats is the registry of named string-formatting functions. It is
// populated during init() and (by contract) must not be modified once
// validation has begun reading from it -- see UseFormat.
var formats = make(map[string]format.Generator)

// formatsFrozen latches to TRUE the first time the registry is read. Because
// UseFormat is only intended to be called during startup, this lets us detect
// (and reject) any registration that arrives after validation has started,
// without paying for a lock on the hot read path.
var formatsFrozen atomic.Bool

// UseFormat adds a custom FormatFunc function to this library.  Used to register custom validators.
//
// UseFormat is only safe to call during program startup (e.g. from init), before
// any schema validation reads the registry. A call that arrives after the registry
// has been read is ignored and reported, because mutating the map concurrently with
// validation reads would be a data race.
func UseFormat(name string, fn format.Generator) {

	if fn == nil {
		return
	}

	// The registry is read-only once validation has begun. A late registration
	// is dropped (rather than racing the readers) and reported for the developer.
	if formatsFrozen.Load() {
		derp.Report(derp.Internal("schema.UseFormat", "format registered after validation began; ignoring (UseFormat must be called during init)", name))
		return
	}

	formats[name] = fn
}

// lookupFormat returns the registered generator for the given name. The first
// lookup freezes the registry against further writes via UseFormat.
func lookupFormat(name string) (format.Generator, bool) {

	if !formatsFrozen.Load() {
		formatsFrozen.Store(true)
	}

	fn, ok := formats[name]
	return fn, ok
}

func init() {

	// Calendar
	UseFormat("iso8601", format.ISO8601)
	UseFormat("date", format.Date)
	UseFormat("dateTime", format.DateTime)
	UseFormat("time", format.Time)

	// Databases
	UseFormat("objectId", format.ObjectID)

	// Email
	UseFormat("email", format.Email)

	// HTML
	UseFormat("html", format.HTML)
	UseFormat("no-html", format.NoHTML)
	UseFormat("text", format.Text)

	// Network
	UseFormat("ipv4", format.IPv4)
	UseFormat("ipv6", format.IPv6)
	UseFormat("hostname", format.Hostname)
	UseFormat("uri", format.URI)

	// Passwords
	UseFormat("lower", format.HasLowercase)
	UseFormat("upper", format.HasUppercase)
	UseFormat("number", format.HasNumbers)
	// UseFormat("symbol", format.HasSymbols)
	// UseFormat("entropy", format.HasEntropy)

	// Regex
	UseFormat("regex", format.MatchRegex)

	// Sets
	UseFormat("in", format.In)
	UseFormat("notin", format.NotIn)

	// Text
	UseFormat("color", format.Color)
	UseFormat("token", format.Token)
	UseFormat("unsafe-any", format.UnsafeAny)
	UseFormat("username", format.Username)
	UseFormat("webfinger", format.WebFinger)
}
