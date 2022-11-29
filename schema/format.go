package schema

import "github.com/benpate/rosetta/schema/format"

var formats map[string]format.Function

// UseFormat adds a custom FormatFunc function to this library.  Used to register custom validators
func UseFormat(name string, fn format.Function) {
	formats[name] = fn
}

func init() {

	formats = make(map[string]format.Function)

	// Calendar
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
	UseFormat("token", format.Token)
	UseFormat("unsafe-any", format.UnsafeAny)
}
