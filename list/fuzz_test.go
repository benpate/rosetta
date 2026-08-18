package list

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzList_Decompose feeds arbitrary values and delimiters to the read-side list functions
// and asserts that they never panic, and that they decompose the value exactly the way
// strings.Split does.  strings.Split is the oracle here: for a single-byte separator it
// splits on raw bytes, which is precisely what this package's Index/Head/Tail do.
func FuzzList_Decompose(f *testing.F) {

	f.Add("a,b,c", byte(','))
	f.Add("", byte(','))
	f.Add(",", byte(','))
	f.Add(",,,", byte(','))
	f.Add("a,,b", byte(','))
	f.Add("no-delimiter-here", byte(','))
	f.Add("a/b/c", byte('/'))
	f.Add("héllo,wörld", byte(','))
	f.Add("\x00,\xff", byte(','))
	f.Add("trailing,", byte(','))
	f.Add(",leading", byte(','))

	f.Fuzz(func(t *testing.T, value string, delimiter byte) {

		// A single-byte separator makes strings.Split a byte-exact oracle
		separator := string([]byte{delimiter})
		parts := strings.Split(value, separator)

		require.Equal(t, parts[0], Head(value, delimiter))
		require.Equal(t, parts[0], First(value, delimiter))
		require.Equal(t, parts[len(parts)-1], Last(value, delimiter))

		// Tail is everything after the first delimiter, or empty when there is none
		if len(parts) == 1 {
			require.Equal(t, "", Tail(value, delimiter))
			require.Equal(t, "", RemoveLast(value, delimiter))
		} else {
			require.Equal(t, strings.Join(parts[1:], separator), Tail(value, delimiter))
			require.Equal(t, strings.Join(parts[:len(parts)-1], separator), RemoveLast(value, delimiter))
		}

		// Split and SplitTail must agree with their single-value counterparts
		head, tail := Split(value, delimiter)
		require.Equal(t, Head(value, delimiter), head)
		require.Equal(t, Tail(value, delimiter), tail)

		if len(parts) > 1 {
			removed, last := SplitTail(value, delimiter)
			require.Equal(t, RemoveLast(value, delimiter), removed)
			require.Equal(t, Last(value, delimiter), last)
		}

		// Every index must round-trip, and out-of-range indexes must be empty (never a panic)
		if value != "" {
			for index, part := range parts {
				require.Equal(t, part, At(value, delimiter, index), "index %d of %q", index, value)
			}
		}

		require.Equal(t, "", At(value, delimiter, len(parts)))
		require.Equal(t, "", At(value, delimiter, len(parts)+1))

		// IsEmpty / IsEmptyTail must agree with the decomposition
		require.Equal(t, value == "", IsEmpty(value))
		require.Equal(t, (len(parts) == 1) || ((len(parts) == 2) && (parts[1] == "")), IsEmptyTail(value, delimiter))
	})
}

// FuzzList_ByteParity asserts that the []byte instantiation of every generic list function
// returns the same answer as the string instantiation for the same bytes.
func FuzzList_ByteParity(f *testing.F) {

	f.Add("a,b,c", byte(','))
	f.Add("", byte(','))
	f.Add("a,,b", byte(','))
	f.Add("\xff\xfe,\x00", byte(','))

	f.Fuzz(func(t *testing.T, value string, delimiter byte) {

		bytes := []byte(value)

		require.Equal(t, Index(value, delimiter), Index(bytes, delimiter))
		require.Equal(t, LastIndex(value, delimiter), LastIndex(bytes, delimiter))
		require.Equal(t, Head(value, delimiter), Head(bytes, delimiter))
		require.Equal(t, Last(value, delimiter), Last(bytes, delimiter))
		require.Equal(t, Tail(value, delimiter), string(Tail(bytes, delimiter)))
		require.Equal(t, RemoveLast(value, delimiter), string(RemoveLast(bytes, delimiter)))
		require.Equal(t, IsEmpty(value), IsEmpty(bytes))
		require.Equal(t, IsEmptyTail(value, delimiter), IsEmptyTail(bytes, delimiter))
	})
}

// FuzzList_Push asserts that a value pushed onto either end of a list reads back from that
// same end.  The delimiter is restricted to ASCII: PushHead/PushTail join with
// string(delimiter), which UTF-8-encodes any byte above 0x7f into two bytes, so the
// delimiter they write is not the byte Index() later searches for.
func FuzzList_Push(f *testing.F) {

	f.Add("b,c", "a", byte(','))
	f.Add("", "a", byte(','))
	f.Add("b", "", byte(','))
	f.Add("", "", byte(','))
	f.Add("b,c", "a,x", byte(','))
	f.Add("héllo", "wörld", byte('/'))

	f.Fuzz(func(t *testing.T, value string, item string, delimiter byte) {

		if delimiter > 0x7f {
			t.Skip("PushHead/PushTail cannot write a non-ASCII delimiter -- see the doc comment above")
		}

		pushedHead := PushHead(value, item, delimiter)
		pushedTail := PushTail(value, item, delimiter)

		// Pushing an empty item is a no-op on both ends
		if item == "" {
			require.Equal(t, value, pushedHead)
			require.Equal(t, value, pushedTail)
			return
		}

		// The pushed item reads back from the end it was pushed onto.  An item that
		// itself contains the delimiter arrives as multiple elements, so only its own
		// first (or last) element comes back.
		require.Equal(t, Head(item, delimiter), Head(pushedHead, delimiter))
		require.Equal(t, Last(item, delimiter), Last(pushedTail, delimiter))

		// Pushing onto a non-empty list leaves the original list reachable
		if value != "" {
			require.Equal(t, Last(value, delimiter), Last(pushedHead, delimiter))
			require.Equal(t, Head(value, delimiter), Head(pushedTail, delimiter))
		}
	})
}
