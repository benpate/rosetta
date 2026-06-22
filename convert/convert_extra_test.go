package convert

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBool_Wrapper(t *testing.T) {
	require.True(t, Bool(true))
	require.True(t, Bool("true"))
	require.False(t, Bool("false"))
	require.False(t, Bool(nil))
	require.False(t, BoolDefault("garbage", false))
	require.True(t, BoolDefault("garbage", true))
}

func TestBytes(t *testing.T) {
	require.Equal(t, []byte{}, Bytes(nil))
	require.Equal(t, []byte{0x41}, Bytes(byte(0x41)))
	require.Equal(t, []byte("hello"), Bytes([]byte("hello")))
	require.Equal(t, []byte("world"), Bytes("world"))
	require.Equal(t, []byte("42"), Bytes(42))
}

func TestFloat_Wrapper(t *testing.T) {
	require.Equal(t, 3.5, Float(3.5))
	require.Equal(t, 42.0, Float("42"))
	require.Equal(t, 0.0, Float(nil))
	require.Equal(t, 9.9, FloatDefault("garbage", 9.9))
}

func TestInt64_Wrapper(t *testing.T) {
	require.Equal(t, int64(42), Int64(42))
	require.Equal(t, int64(42), Int64("42"))
	require.Equal(t, int64(0), Int64(nil))
	require.Equal(t, int64(99), Int64Default("garbage", 99))
}

func TestInt32_Wrapper(t *testing.T) {
	require.Equal(t, int32(42), Int32(42))
	require.Equal(t, int32(99), Int32Default("garbage", 99))
}

func TestPointerAndElement(t *testing.T) {

	p := Pointer(42)
	require.Equal(t, 42, *p)

	// Element dereferences pointers
	require.Equal(t, 42, Element(p))
	// Non-pointers are returned unchanged
	require.Equal(t, "hello", Element("hello"))
}

func TestIsMap(t *testing.T) {
	require.True(t, IsMap(map[string]any{"a": 1}))
	require.True(t, IsMap(map[string]string{"a": "b"}))
	require.True(t, IsMap(map[string][]string{"a": {"b"}}))
	require.True(t, IsMap(map[int]int{1: 2}))
	require.False(t, IsMap("not a map"))
	require.False(t, IsMap(42))
}

func TestIsSlice(t *testing.T) {
	require.True(t, IsSlice([]int{1, 2}))
	require.True(t, IsSlice([]string{"a"}))
	require.True(t, IsSlice([]any{1, "two"}))
	require.True(t, IsSlice([2]int{1, 2})) // array
	require.False(t, IsSlice("not a slice"))
	require.False(t, IsSlice(42))

	// Pointer to a slice
	s := []int{1, 2}
	require.True(t, IsSlice(&s))
}

func TestSliceLength(t *testing.T) {
	require.Equal(t, 0, SliceLength(nil))
	require.Equal(t, 3, SliceLength([]int{1, 2, 3}))
	require.Equal(t, 2, SliceLength([]string{"a", "b"}))
	require.Equal(t, 2, SliceLength([]any{1, 2}))
	require.Equal(t, 1, SliceLength([]float64{1.5}))
	require.Equal(t, 1, SliceLength([]int64{1}))
	require.Equal(t, 1, SliceLength([]map[string]any{{}}))
	require.Equal(t, 0, SliceLength("not a slice"))

	// Pointer to a slice
	s := []int{1, 2, 3, 4}
	require.Equal(t, 4, SliceLength(&s))

	// Reflection-based array
	require.Equal(t, 3, SliceLength([3]int{1, 2, 3}))
}

func TestReflectType(t *testing.T) {
	require.Equal(t, reflect.TypeOf(42), ReflectType(42))

	// A reflect.Type is returned as-is
	typeOf := reflect.TypeOf("hello")
	require.Equal(t, typeOf, ReflectType(typeOf))
}

func TestJoinString(t *testing.T) {
	require.Equal(t, "a,b,c", JoinString([]string{"a", "b", "c"}, ","))
	require.Equal(t, "1-2-3", JoinString([]any{1, 2, 3}, "-"))

	// Empty delimiter falls through to String()
	require.Equal(t, "42", JoinString(42, ""))
}

func TestSliceOfAny_Wrapper(t *testing.T) {
	require.Equal(t, []any{1, 2, 3}, SliceOfAny([]int{1, 2, 3}))
	require.Equal(t, []any{"a", "b"}, SliceOfAny([]string{"a", "b"}))
	require.Equal(t, []any{true}, SliceOfAny(true))
	require.Equal(t, []any{"hello"}, SliceOfAny("hello"))

	// Non-convertible types yield an empty slice
	require.Equal(t, []any{}, SliceOfAny(nil))
}

func TestTimeDefault(t *testing.T) {

	fallback := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	parsed := TimeDefault("2026-03-04T00:00:00Z", fallback)
	require.Equal(t, 2026, parsed.Year())

	// Unparseable value returns the default
	require.Equal(t, fallback, TimeDefault("garbage", fallback))
}

func TestMapOfSliceOfString(t *testing.T) {

	{
		result := MapOfSliceOfString(map[string]string{"a": "one"})
		require.Equal(t, []string{"one"}, result["a"])
	}
	{
		result := MapOfSliceOfString(map[string]any{"a": "one"})
		require.Equal(t, []string{"one"}, result["a"])
	}
	{
		input := map[string][]string{"a": {"one", "two"}}
		require.Equal(t, input, MapOfSliceOfString(input))
	}
	{
		// Non-convertible input yields an empty map
		require.Equal(t, map[string][]string{}, MapOfSliceOfString("nope"))
	}
}

func TestURLValuesAndHTTPHeader(t *testing.T) {

	source := map[string][]string{"a": {"one"}}

	require.Equal(t, url.Values(source), URLValues(source))
	require.Equal(t, http.Header(source), HTTPHeader(source))

	urlValues, ok := URLValuesOk(source)
	require.True(t, ok)
	require.Equal(t, url.Values(source), urlValues)

	header, ok := HTTPHeaderOk(source)
	require.True(t, ok)
	require.Equal(t, http.Header(source), header)

	_, ok = URLValuesOk("nope")
	require.False(t, ok)

	_, ok = HTTPHeaderOk("nope")
	require.False(t, ok)
}
