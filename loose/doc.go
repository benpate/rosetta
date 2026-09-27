// Package loose provides types that are forgiving about the values they
// accept. They exist for the receiving half of Postel's law: be conservative in
// what you send, be liberal in what you accept.
//
// A plain `int64` field fails an entire document the first time a peer sends
// "480" instead of 480. Int64 absorbs that, and every other numeric spelling
// seen in the wild. String does the same for text fields that peers sometimes
// send as bare numbers or booleans. Int64 is not Int on purpose: the accepted
// range must not narrow on a 32-bit target. Both reject input that is not a
// JSON scalar, and quietly degrade anything else to the zero value.
//
// Template is a string that compiles as a text/template when it holds a "{{"
// followed later by a "}}". A string that does not compile is kept as plain
// text rather than rejected. Execute renders a template against a value, or
// returns plain text as it is. It encodes to JSON as its source string, and is
// safe to execute from many goroutines at once. CachedTemplate compiles each
// source string once, for strings that are rendered repeatedly.
package loose
