// Package convert converts values from one type into another, accepting generic
// input and producing a specific Go type with sensible defaults.
//
// Most conversions come in three shapes. The bare form (Int, String, Bool)
// always returns a value, falling back to the zero value when the input cannot
// be converted. The Default form substitutes a caller-supplied fallback
// instead. The Ok form returns that same value alongside a boolean reporting
// whether the conversion actually succeeded, for callers that need to tell
// "absent" from "converted to zero".
//
// Conversions reach through pointers, interfaces, and the small accessor
// interfaces shared with the compare package, so a type that can describe
// itself as an int or a string is converted through that description rather
// than by reflection over its fields.
package convert
