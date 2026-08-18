package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecond(t *testing.T) {

	check := func(name string, list List, expected string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			require.Equal(t, expected, Second(list))
		})
	}

	check("three items", ByDot("a", "b", "c"), "b")
	check("exactly two items", ByDot("a", "b"), "b")
	check("only one item, so there is no second", ByDot("a"), "")
	check("empty list", ByDot(), "")
	check("empty second item", Dot("a..c"), "")
}

func TestThird(t *testing.T) {

	check := func(name string, list List, expected string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			require.Equal(t, expected, Third(list))
		})
	}

	check("four items", ByDot("a", "b", "c", "d"), "c")
	check("exactly three items", ByDot("a", "b", "c"), "c")
	check("only two items, so there is no third", ByDot("a", "b"), "")
	check("only one item", ByDot("a"), "")
	check("empty list", ByDot(), "")
}

func TestFirst2(t *testing.T) {

	check := func(name string, list List, want1 string, want2 string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			got1, got2 := First2(list)
			require.Equal(t, want1, got1)
			require.Equal(t, want2, got2)
		})
	}

	check("three items", ByDot("a", "b", "c"), "a", "b")
	check("exactly two items", ByDot("a", "b"), "a", "b")
	check("only one item", ByDot("a"), "a", "")
	check("empty list", ByDot(), "", "")
}

func TestFirst3(t *testing.T) {

	check := func(name string, list List, want1 string, want2 string, want3 string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			got1, got2, got3 := First3(list)
			require.Equal(t, want1, got1)
			require.Equal(t, want2, got2)
			require.Equal(t, want3, got3)
		})
	}

	check("four items", ByDot("a", "b", "c", "d"), "a", "b", "c")
	check("exactly three items", ByDot("a", "b", "c"), "a", "b", "c")
	check("only two items", ByDot("a", "b"), "a", "b", "")
	check("only one item", ByDot("a"), "a", "", "")
	check("empty list", ByDot(), "", "", "")
}

func TestFirst4(t *testing.T) {

	check := func(name string, list List, want1 string, want2 string, want3 string, want4 string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			got1, got2, got3, got4 := First4(list)
			require.Equal(t, want1, got1)
			require.Equal(t, want2, got2)
			require.Equal(t, want3, got3)
			require.Equal(t, want4, got4)
		})
	}

	check("five items", ByDot("a", "b", "c", "d", "e"), "a", "b", "c", "d")
	check("exactly four items", ByDot("a", "b", "c", "d"), "a", "b", "c", "d")
	check("only three items", ByDot("a", "b", "c"), "a", "b", "c", "")
	check("only one item", ByDot("a"), "a", "", "", "")
	check("empty list", ByDot(), "", "", "", "")
}

// Last2/Last3/Last4 return their values in REVERSE order: last first, then
// second-to-last, and so on.
//
// Only lists long enough to fill every return value are asserted here. For a list
// SHORTER than the number of requested items, these functions currently pad in the
// MIDDLE rather than at the end -- Last3(["a","b"]) returns ("b", "", "a") -- which
// is almost certainly not intended. That behavior is deliberately left unpinned
// pending a decision on the exported contract; see TestLast_ShortListIsUnspecified.
func TestLast2(t *testing.T) {

	check := func(name string, list List, wantLast string, wantLast2 string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			last, last2 := Last2(list)
			require.Equal(t, wantLast, last)
			require.Equal(t, wantLast2, last2)
		})
	}

	check("four items", ByDot("a", "b", "c", "d"), "d", "c")
	check("three items", ByDot("a", "b", "c"), "c", "b")
	check("exactly two items", ByDot("a", "b"), "b", "a")
	check("empty list", ByDot(), "", "")
}

func TestLast3(t *testing.T) {

	check := func(name string, list List, wantLast string, wantLast2 string, wantLast3 string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			last, last2, last3 := Last3(list)
			require.Equal(t, wantLast, last)
			require.Equal(t, wantLast2, last2)
			require.Equal(t, wantLast3, last3)
		})
	}

	check("five items", ByDot("a", "b", "c", "d", "e"), "e", "d", "c")
	check("four items", ByDot("a", "b", "c", "d"), "d", "c", "b")
	check("exactly three items", ByDot("a", "b", "c"), "c", "b", "a")
	check("empty list", ByDot(), "", "", "")
}

func TestLast4(t *testing.T) {

	check := func(name string, list List, wantLast string, wantLast2 string, wantLast3 string, wantLast4 string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			last, last2, last3, last4 := Last4(list)
			require.Equal(t, wantLast, last)
			require.Equal(t, wantLast2, last2)
			require.Equal(t, wantLast3, last3)
			require.Equal(t, wantLast4, last4)
		})
	}

	check("six items", ByDot("a", "b", "c", "d", "e", "f"), "f", "e", "d", "c")
	check("five items", ByDot("a", "b", "c", "d", "e"), "e", "d", "c", "b")
	check("exactly four items", ByDot("a", "b", "c", "d"), "d", "c", "b", "a")
	check("empty list", ByDot(), "", "", "", "")
}

// Lists shorter than the number of requested items currently pad in the middle instead
// of at the end. This test records what happens today WITHOUT endorsing it, so that a
// deliberate change to the contract shows up here as a failure rather than passing silently.
func TestLast_ShortListIsUnspecified(t *testing.T) {

	// One item: Last2 reports the item as the SECOND-to-last and "" as the last, even
	// though Last() correctly reports "a" as the last item.
	last, last2 := Last2(ByDot("a"))
	require.Equal(t, "", last, "KNOWN ODDITY: the only item is not reported as last")
	require.Equal(t, "a", last2)
	require.Equal(t, "a", ByDot("a").Last(), "Last() disagrees with Last2()")

	// Two items via Last3: the empty value lands in the MIDDLE of the results.
	l1, l2, l3 := Last3(ByDot("a", "b"))
	require.Equal(t, "b", l1)
	require.Equal(t, "", l2, "KNOWN ODDITY: padding lands in the middle, not the end")
	require.Equal(t, "a", l3)
}

// Every algorithm here is delimiter-agnostic, so each one must give the same
// answer no matter which concrete List type carries the items.
func TestAlgorithms_AllDelimiters(t *testing.T) {

	check := func(name string, list List) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			require.Equal(t, "b", Second(list))
			require.Equal(t, "c", Third(list))

			first1, first2 := First2(list)
			require.Equal(t, "a", first1)
			require.Equal(t, "b", first2)

			last, last2 := Last2(list)
			require.Equal(t, "d", last)
			require.Equal(t, "c", last2)
		})
	}

	check("comma", ByComma("a", "b", "c", "d"))
	check("dot", ByDot("a", "b", "c", "d"))
	check("equal", ByEqual("a", "b", "c", "d"))
	check("semicolon", BySemicolon("a", "b", "c", "d"))
	check("slash", BySlash("a", "b", "c", "d"))
	check("space", BySpace("a", "b", "c", "d"))
}
