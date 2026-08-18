// Package lenient provides scalar types that are forgiving about how a remote
// system encodes them. They exist for the receiving half of Postel's law: be
// conservative in what you send, be liberal in what you accept.
//
// A plain `int64` field fails an entire document the first time a peer sends
// "480" instead of 480. Substituting lenient.Int64 absorbs that, and every
// other numeric spelling seen in the wild, without weakening the strict
// validation you apply to documents you author. lenient.String does the same
// for text fields that peers sometimes send as bare numbers or booleans.
//
// The numeric type is Int64, not Int, on purpose: a sender's JSON has no idea
// what GOARCH you built for, so the accepted range must not narrow on a 32-bit
// target. Precision beyond int64 is not a goal.
//
// Tolerance stops at structure. Both types reject input that is not a JSON
// scalar at all, because coercing an object or array would invent a value the
// sender never wrote. Everything short of that quietly degrades to the zero
// value rather than erroring, so one sloppy field never costs you the
// document.
package lenient
