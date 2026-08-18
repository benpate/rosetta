// Package html converts between HTML and plain text, sanitizes untrusted
// markup, and extracts summaries from HTML documents.
//
// FromText escapes plain text into lightly-formatted HTML. ToText and
// ToSearchText strip markup back down to text, and Summary produces a short
// plain-text excerpt truncated on a rune boundary. RemoveTags, RemoveAnchors,
// RemoveSpecialCharacters, and CollapseWhitespace handle the narrower cleanups,
// and IsHTML reports whether a string appears to contain markup at all.
//
// Minimal is the trust boundary. It is the only function here that sanitizes:
// a bluemonday policy plus an explicit element allow-list, and the one thing to
// run untrusted markup through before rendering it. Everything else in this
// package cleans or escapes without making any safety promise about tags that
// were already present.
package html
