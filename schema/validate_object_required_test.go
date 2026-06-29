package schema_test

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

// requiredOptionalSchema is a map-backed Object schema with one REQUIRED property and one OPTIONAL
// property. It is the fixture for the MapTyper required/optional behavior tests below.
func requiredOptionalSchema() schema.Schema {
	return schema.New(schema.Object{
		Properties: schema.ElementMap{
			"required": schema.String{Required: true},
			"optional": schema.String{},
		},
	})
}

// TestValidateObject_Map_OptionalAbsent_RequiredPresent: the happy path. A map (IsMap() == true)
// that contains the required key but omits the optional key must validate cleanly -- an absent
// optional property is legitimately empty.
func TestValidateObject_Map_OptionalAbsent_RequiredPresent(t *testing.T) {

	s := requiredOptionalSchema()
	value := mapof.String{"required": "present"}

	_, err := s.Validate(value)
	require.Nil(t, err)
}

// TestValidateObject_Map_RequiredAbsent: the rule under test. A map that omits the REQUIRED key
// must FAIL validation, even though absent optional keys are tolerated. (If skipAbsentKeys causes
// the required check to be skipped, this test fails -- which is the bug it guards against.)
func TestValidateObject_Map_RequiredAbsent(t *testing.T) {

	s := requiredOptionalSchema()
	value := mapof.String{"optional": "present"} // required key is missing

	_, err := s.Validate(value)
	require.NotNil(t, err)
}

// TestValidateObject_Map_RequiredPresentButEmpty: a present-but-empty required value must also fail,
// matching the non-map contract (Required + "" is invalid).
func TestValidateObject_Map_RequiredPresentButEmpty(t *testing.T) {

	s := requiredOptionalSchema()
	value := mapof.String{"required": ""}

	_, err := s.Validate(value)
	require.NotNil(t, err)
}

// TestValidateObject_Map_BothPresent: both keys present and non-empty validates cleanly.
func TestValidateObject_Map_BothPresent(t *testing.T) {

	s := requiredOptionalSchema()
	value := mapof.String{"required": "here", "optional": "also here"}

	_, err := s.Validate(value)
	require.Nil(t, err)
}
