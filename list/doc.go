// Package list treats a delimited string as a list, without splitting it into a
// slice.
//
// Head, Tail, Split, At, Last, PushHead, PushTail, and the rest read and
// rewrite a string in place around a single-byte delimiter, so walking a
// dotted path or a comma-separated field costs no allocations for the
// intermediate pieces. The named types — Dot, Comma, Slash, Semicolon, Space,
// Equal — bind a delimiter to a string so the same operations read as methods,
// and the List interface lets code accept any of them.
//
// The delimiter is a byte, not a rune. A multi-byte delimiter cannot be
// expressed here, and the functions are written for the ASCII punctuation that
// path and header syntaxes actually use.
package list
