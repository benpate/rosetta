// Package pointer returns a pointer to a value whose type is not known at
// compile time.
//
// To is reflection-based and untyped — any in, any out — which is what makes it
// usable when building values dynamically for the schema engine. Inputs that
// are already pointers or interfaces pass through unchanged, so To is
// idempotent and never produces a pointer to a pointer.
//
// When the type is known at compile time, prefer convert.Pointer, which is
// generic, type-safe, and skips the reflection entirely. The convert package
// owns the typed pointer helpers; this package exists only for the untyped
// case.
package pointer
