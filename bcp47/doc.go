// Package bcp47 turns any BCP 47 language tag into the one spelling that should be stored.
//
// Canonical rewrites "EN-us", "en_US", and "en-US" to the same "en-US", replaces deprecated
// codes with their current form ("iw" is "he"), and drops extensions and private-use subtags,
// which describe formatting preferences rather than a different language. It refuses values
// that name no language a person can write in, including "und", "mul", and "*". Parse returns
// the same result as a golang.org/x/text/language Tag, for callers that match or describe it.
//
// Every value is bounded before parsing, because language tags usually arrive from visitors
// and authors. The schema package's "language" string format is built on Canonical, so a
// schema and this package can never disagree about a non-empty value.
package bcp47
