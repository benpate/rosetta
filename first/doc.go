// Package first returns the first non-empty value from a list of candidates.
//
// It exists for the fallback chains that show up all over configuration and
// display code — use the value the caller passed, otherwise the one on the
// record, otherwise a hard-coded default — written as a single expression
// instead of a ladder of if statements. Emptiness means the zero value for the
// type, and if every candidate is empty, that zero value is what comes back.
package first
