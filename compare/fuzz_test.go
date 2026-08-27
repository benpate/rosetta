package compare

import (
	"strings"
	"testing"
)

// fuzzValues returns the awkward operands worth throwing at every comparison.  They are built
// per call so that a target which mutates one cannot affect another.
func fuzzValues(text string) []any {

	var nilPointer *int
	var nilMap map[string]any

	return []any{
		nil,
		text,
		len(text),
		float64(len(text)),
		text == "true",
		nilPointer,
		nilMap,
		[]string{text},
		[]any{},
		map[string]any{text: text},
		make(chan int),
		func() {},
		struct{ Name string }{Name: text},
		&struct{ Name string }{Name: text},
	}
}

// FuzzWithOperator drives the operator dispatcher with an arbitrary operator string and every
// awkward operand pairing.  This is the engine behind `show-if` and `required-if` expressions,
// whose operators come out of template files rather than out of Go, so an unrecognized or
// malformed operator must produce an error -- never a panic that takes down form rendering.
func FuzzWithOperator(f *testing.F) {

	f.Add("=", "value")
	f.Add("==", "value")
	f.Add("!=", "value")
	f.Add(">", "1")
	f.Add("<", "1")
	f.Add(">=", "1")
	f.Add("<=", "1")
	f.Add("", "")
	f.Add("bogus", "value")
	f.Add("\x00", "value")
	f.Add(strings.Repeat("=", 512), "value")
	f.Add("contains", "value")
	f.Add("beginsWith", "value")

	f.Fuzz(func(t *testing.T, operator string, text string) {

		values := fuzzValues(text)

		for _, value1 := range values {
			for _, value2 := range values {

				// An unknown operator or an incomparable pair must ERROR, never panic
				_, _ = WithOperator(value1, operator, value2)
			}
		}
	})
}

// FuzzComparisons_NeverPanic runs every exported comparison over every awkward operand pairing.
// These are called from expression evaluation and sorting, where the operands are whatever the
// surrounding data happened to hold.
func FuzzComparisons_NeverPanic(f *testing.F) {

	f.Add("")
	f.Add("value")
	f.Add("0")
	f.Add("true")
	f.Add("\x00")
	f.Add("\xff\xfe")
	f.Add(strings.Repeat("a", 4096))

	f.Fuzz(func(t *testing.T, text string) {

		values := fuzzValues(text)

		for _, value1 := range values {

			_ = IsNil(value1)
			_ = NotNil(value1)
			_ = IsZero(value1)
			_ = NotZero(value1)

			for _, value2 := range values {
				_ = Equal(value1, value2)
				_ = LessThan(value1, value2)
				_ = GreaterThan(value1, value2)
				_ = Contains(value1, value2)
				_ = NotContains(value1, value2)
				_ = BeginsWith(value1, value2)
				_ = EndsWith(value1, value2)
				_, _ = Interface(value1, value2)
			}
		}
	})
}

// FuzzCompare_Antisymmetry asserts the ordering law that sorting depends on: if value1 sorts
// before value2, then value2 must sort after value1.  A comparison that reports "less than" in
// both directions makes sort.Slice produce a different order depending on input order, which
// surfaces as a list that reshuffles itself between page loads.
func FuzzCompare_Antisymmetry(f *testing.F) {

	f.Add("a", "b")
	f.Add("", "")
	f.Add("1", "2")
	f.Add("10", "9")
	f.Add("\x00", "\x00")
	f.Add("a", "a")

	f.Fuzz(func(t *testing.T, text1 string, text2 string) {

		// Only comparable operand types can be held to the ordering law
		for _, pair := range [][2]any{
			{text1, text2},
			{len(text1), len(text2)},
			{float64(len(text1)), float64(len(text2))},
		} {
			forward, err1 := Interface(pair[0], pair[1])
			reverse, err2 := Interface(pair[1], pair[0])

			if (err1 != nil) || (err2 != nil) {
				continue
			}

			// RULE: Reversing the operands must reverse the sign of the result
			if forward != -reverse {
				t.Fatalf("Interface(%v, %v) = %d but Interface(%v, %v) = %d",
					pair[0], pair[1], forward, pair[1], pair[0], reverse)
			}
		}
	})
}
