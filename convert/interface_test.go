package convert

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterface_Int(t *testing.T) {
	v := Interface(6)
	require.Equal(t, 6, v)
}

func TestInterface_String(t *testing.T) {
	v := Interface("hello there")
	require.Equal(t, "hello there", v)
}

func TestInterface_Nil(t *testing.T) {
	v := Interface(nil)
	require.Equal(t, nil, v)
}

func TestInterface_Map(t *testing.T) {
	var value map[string]any
	v := Interface(value)
	require.Equal(t, value, v)
}

func TestInterface_Int_Reflect(t *testing.T) {
	value := 6
	v := Interface(reflect.ValueOf(value))
	require.Equal(t, 6, v)
}

func TestInterface_String_Reflect(t *testing.T) {
	value := "general kenobi"
	v := Interface(reflect.ValueOf(value))
	require.Equal(t, "general kenobi", v)
}

func TestInterface_Nil_Reflect(t *testing.T) {
	var value *int
	v := Interface(reflect.ValueOf(value))
	require.Equal(t, value, v)
}

func TestInterface_Map_Reflect(t *testing.T) {
	var value map[string]any
	v := Interface(reflect.ValueOf(value))
	require.Equal(t, value, v)
}

func TestInterface_Invalid(t *testing.T) {
	var value map[string]any
	v := Interface(reflect.ValueOf(value["missing"]))
	require.Equal(t, nil, v)
}
