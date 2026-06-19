package pointer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPointerTo(t *testing.T) {

	value := map[string]string{
		"key": "value",
		"foo": "bar",
	}

	result := To(value)

	value["other"] = "thing"

	require.Equal(t, &value, result)
	require.Equal(t, "thing", (*result.(*map[string]string))["other"])
	require.Equal(t, "thing", value["other"])
}

func TestTo_Scalar(t *testing.T) {

	result := To(42)

	pointer, ok := result.(*int)
	require.True(t, ok)
	require.Equal(t, 42, *pointer)
}

func TestTo_String(t *testing.T) {

	result := To("hello")

	pointer, ok := result.(*string)
	require.True(t, ok)
	require.Equal(t, "hello", *pointer)
}

func TestTo_AlreadyPointer(t *testing.T) {

	original := 42
	p := &original

	// A value that is already a pointer is returned as-is
	result := To(p)
	require.Same(t, p, result.(*int))
}

// NOTE: the `reflect.Interface` branch inside To() is not exercised here. A
// value passed through an `any` parameter always reflects as its concrete kind,
// so reflect.ValueOf(value).Kind() never reports reflect.Interface. That branch
// is effectively unreachable from this entry point.
