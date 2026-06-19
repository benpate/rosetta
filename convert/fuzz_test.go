package convert

import "testing"

// FuzzStringParsers feeds arbitrary strings to the string-based parsers to confirm
// that none of them panic, regardless of input.
func FuzzStringParsers(f *testing.F) {

	f.Add("")
	f.Add("123")
	f.Add("-456")
	f.Add("3.14159")
	f.Add("true")
	f.Add("false")
	f.Add("ff")
	f.Add("2026-03-04T13:02:00Z")
	f.Add("not a number")
	f.Add("99999999999999999999999999999")

	f.Fuzz(func(t *testing.T, value string) {
		// None of these conversions should ever panic.
		_ = Int(value)
		_ = Int32(value)
		_ = Int64(value)
		_ = Float(value)
		_ = Bool(value)
		_ = Bytes(value)
		_ = String(value)
		_ = Time(value)
		_ = SliceOfString(value)
		_ = SliceOfInt(value)
		_ = SliceOfFloat(value)
		_ = SliceOfAny(value)
	})
}
