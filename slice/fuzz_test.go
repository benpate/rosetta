package slice

import (
	"strings"
	"testing"
)

// FuzzIndexOperations drives every index-taking helper with arbitrary indexes against slices of
// arbitrary length.  Each one does its own bounds arithmetic, so an out-of-range or negative
// index must come back as a zero value rather than as a panic.
func FuzzIndexOperations(f *testing.F) {

	f.Add(0, 0)
	f.Add(0, 1)
	f.Add(1, 0)
	f.Add(3, -1)
	f.Add(3, 99)
	f.Add(0, -1)
	f.Add(1, -9223372036854775808)
	f.Add(1, 9223372036854775807)

	f.Fuzz(func(t *testing.T, length int, index int) {

		// Keep the slice small; the interesting variable is the INDEX
		if (length < 0) || (length > 64) {
			return
		}

		original := make([]string, length)

		for i := range original {
			original[i] = "item"
		}

		_ = At(original, index)

		value, ok := AtOK(original, index)

		// RULE: AtOK reporting success means the index really was in range
		if ok && ((index < 0) || (index >= length)) {
			t.Fatalf("AtOK reported success for out-of-range index %d (length %d)", index, length)
		}

		if !ok && (value != "") {
			t.Fatalf("AtOK failed for index %d but returned %q", index, value)
		}

		// RemoveAt must never grow the slice, and never panic
		if removed := RemoveAt(original, index); len(removed) > length {
			t.Fatalf("RemoveAt(%d) grew the slice from %d to %d", index, length, len(removed))
		}
	})
}

// FuzzSliceHelpers_NeverPanic runs the whole-slice helpers over arbitrary content, including
// the empty and nil cases that bounds arithmetic tends to miss.
func FuzzSliceHelpers_NeverPanic(f *testing.F) {

	f.Add("")
	f.Add("a")
	f.Add("a,b,c")
	f.Add(",")
	f.Add("\x00")
	f.Add(strings.Repeat("a,", 512))

	f.Fuzz(func(t *testing.T, text string) {

		for _, original := range [][]string{nil, {}, strings.Split(text, ",")} {

			_, _ = Split(original)
			_ = Reverse(original)
			_ = Shuffle(original)
			_ = Unique(original)
			_ = NonZero(original)
			_ = Difference(original, original)
			_ = Contains(original, text)
			_ = NotContains(original, text)
			_ = ContainsAny(original, text)
			_ = ContainsAll(original, text)
			_ = Equal(original, original)
			_ = NotEqual(original, original)
			_ = Filter(original, func(string) bool { return true })
			_ = First(original, func(string) bool { return false })
			_, _ = Find(original, func(string) bool { return false })
			_ = Map(original, func(value string) int { return len(value) })

			for range Range(original) {
				break
			}
		}
	})
}

// FuzzSlice_Laws asserts the structural laws the helpers claim.  These are the properties a
// caller relies on without checking: reversing twice restores the original, Unique preserves
// membership, and Difference removes exactly what it says.
func FuzzSlice_Laws(f *testing.F) {

	f.Add("")
	f.Add("a")
	f.Add("a,b,c")
	f.Add("a,a,b")
	f.Add(",,,")
	f.Add("\x00,\x00")

	f.Fuzz(func(t *testing.T, text string) {

		original := strings.Split(text, ",")

		// Reversing twice restores the original order
		if doubled := Reverse(Reverse(original)); !Equal(original, doubled) {
			t.Fatalf("Reverse is not an involution for %q: got %v", text, doubled)
		}

		// Unique keeps every distinct member, and adds nothing
		unique := Unique(original)

		if len(unique) > len(original) {
			t.Fatalf("Unique grew %v into %v", original, unique)
		}

		for _, value := range original {
			if !Contains(unique, value) {
				t.Fatalf("Unique(%v) dropped member %q", original, value)
			}
		}

		for _, value := range unique {
			if !Contains(original, value) {
				t.Fatalf("Unique(%v) invented member %q", original, value)
			}
		}

		// A slice differenced against itself keeps nothing
		if remainder := Difference(original, original); len(remainder) != 0 {
			t.Fatalf("Difference(x, x) left %v behind", remainder)
		}

		// Shuffle rearranges, but never adds or drops
		if shuffled := Shuffle(original); len(shuffled) != len(original) {
			t.Fatalf("Shuffle changed the length of %v to %d", original, len(shuffled))
		}
	})
}
