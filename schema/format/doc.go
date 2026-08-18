// Package format provides the named string formats used by the schema package.
//
// Each function here is a Generator: it takes the argument written in the
// schema (the part after the format name) and returns a StringFormat closure
// that validates one value. That closure returns the value it accepts, so a
// format may normalize as well as validate — trimming, case-folding, or
// rewriting into a canonical form on the way through.
//
// The set covers the usual identifier and network formats (email, hostname,
// IPv4, IPv6, URI, URL, ObjectID, token, username, WebFinger), dates and times,
// constraints on the text itself (In, NotIn, MatchRegex, HasUppercase,
// HasLowercase, HasNumbers), and the content formats that decide how much
// markup survives (Text, NoHTML, HTML, Markdown, UnsafeAny).
//
// Those content formats are a trust boundary: the format named on a string
// element is what decides whether a value is escaped, sanitized, or passed
// through untouched. Choose UnsafeAny only for values that are already trusted.
package format
