// Package list treats a delimited string as a list, without splitting it into a
// slice.
//
// Head, Tail, Split, At, Last, and the rest slice a string around a single-byte
// delimiter, so walking a dotted path or a comma-separated field with these
// generic functions costs no allocations for the intermediate pieces. PushHead
// and PushTail build a new string. The named types — Dot, Comma, Slash,
// Semicolon, Space, Equal — bind a delimiter to a string so the same operations
// read as methods, and the List interface lets code accept any of them. Those
// methods return a List, which allocates for each non-empty list they return.
//
// The delimiter is a byte, not a rune. A multi-byte delimiter cannot be
// expressed here, and the functions are written for the ASCII punctuation that
// path and header syntaxes actually use.
package list
