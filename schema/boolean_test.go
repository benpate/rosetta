package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {

	s := Boolean{}

	assert.Nil(t, s.Validate(true))
	assert.Nil(t, s.Validate(false))

	assert.Nil(t, s.Validate(1))
	assert.Nil(t, s.Validate("string-bad"))

}

func TestBool_Set(t *testing.T) {

	var err error

	s := Schema{Element: Boolean{}}
	value := false

	err = s.Set(&value, "", true)
	require.Nil(t, err)
	require.True(t, value)

	err = s.Set(&value, "", false)
	require.Nil(t, err)
	require.False(t, value)

	err = s.Set(&value, "invalid-sub-property", true)
	require.NotNil(t, err)
}

func TestBool_Set_Array(t *testing.T) {

	var err error
	var value []bool

	s := Schema{Element: Array{Items: Boolean{}}}

	err = s.Set(&value, "5", true)
	require.Nil(t, err)

	err = s.Set(&value, "4", false)
	require.Nil(t, err)

	err = s.Set(&value, "3", true)
	require.Nil(t, err)

	err = s.Set(&value, "2", false)
	require.Nil(t, err)

	err = s.Set(&value, "1", true)
	require.Nil(t, err)

	require.Equal(t, []bool{false, true, false, true, false, true}, value)
}

func TestBool_Set_MapOfBool(t *testing.T) {

	var err error
	value := map[string]bool{}

	s := Schema{Element: Object{Properties: map[string]Element{"read": Boolean{}, "write": Boolean{}}}}

	err = s.Set(&value, "read", true)
	require.Nil(t, err)
	require.True(t, value["read"])

	err = s.Set(&value, "read", false)
	require.Nil(t, err)
	require.False(t, value["read"])

	err = s.Set(&value, "write", true)
	require.Nil(t, err)
	require.True(t, value["write"])

	err = s.Set(&value, "write", false)
	require.Nil(t, err)
	require.False(t, value["write"])
}

func TestBool_Set_MapOfInterface(t *testing.T) {

	var err error
	var value map[string]any

	s := Schema{Element: Object{Properties: map[string]Element{"read": Boolean{}, "write": Boolean{}}}}

	err = s.Set(&value, "read", true)
	require.Nil(t, err)
	require.Equal(t, true, value["read"])

	err = s.Set(&value, "read", false)
	require.Nil(t, err)
	require.Equal(t, false, value["read"])

	err = s.Set(&value, "write", true)
	require.Nil(t, err)
	require.Equal(t, true, value["write"])

	err = s.Set(&value, "write", false)
	require.Nil(t, err)
	require.Equal(t, false, value["write"])
}

func TestBool_Set_Struct(t *testing.T) {

	var err error
	var value struct {
		Read  bool `path:"read"`
		Write bool `path:"write"`
	}

	s := Schema{Element: Object{Properties: map[string]Element{"read": Boolean{}, "write": Boolean{}}}}

	err = s.Set(&value, "read", true)
	require.Nil(t, err)
	require.True(t, value.Read)

	err = s.Set(&value, "read", false)
	require.Nil(t, err)
	require.False(t, value.Read)

	err = s.Set(&value, "write", true)
	require.Nil(t, err)
	require.True(t, value.Write)

	err = s.Set(&value, "write", false)
	require.Nil(t, err)
	require.False(t, value.Write)
}

func TestBool_Required(t *testing.T) {

	s := Schema{Element: Boolean{Required: true}}

	require.Nil(t, s.Validate(true))
	require.NotNil(t, s.Validate(false))
}

func TestBool_Marshal(t *testing.T) {

	b := Boolean{}

	result, err := json.Marshal(b)
	require.Nil(t, err)
	require.Equal(t, `{"type":"boolean"}`, string(result))
}

func TestBool_Unmarshal(t *testing.T) {

	j := []byte(`{"type":"boolean", "required":true}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(Boolean).Required)

	require.Nil(t, s.Validate(true))
	require.NotNil(t, s.Validate(false))
}
