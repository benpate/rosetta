package translate

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

func TestForEach(t *testing.T) {

	// SOURCE DATA
	sourceValue := sliceof.Object[mapof.Any]{
		{"name": "Alice", "email": "alice@wonderland.com"},
		{"name": "John Connor", "email": "john@connor.mil"},
		{"name": "Sarah Connor", "email": "sarah@sky.net"},
	}

	// TARGET CONFIGURATION
	targetSchema := schema.New(schema.Array{
		Items: schema.Object{
			Properties: schema.ElementMap{
				"fullName":     schema.String{},
				"emailAddress": schema.String{Format: "email"},
				"type":         schema.String{},
				"comment":      schema.String{},
			},
		},
	})

	targetValue := make(sliceof.Object[mapof.Any], 0, 3)

	// MAPPING RULES
	rules := New(
		ForEach("", "", "", []mapof.Any{
			{"path": "name", "target": "fullName"},
			{"path": "email", "target": "emailAddress"},
			{"value": "person", "target": "type"},
			{"expression": "{{.name}} <{{.email}}>", "target": "comment"},
		}),
	)

	err := rules.Execute(schema.Wildcard(), &sourceValue, targetSchema, &targetValue)
	require.Nil(t, err)

	require.Equal(t, 3, len(targetValue))

	require.Equal(t, "Alice", targetValue[0]["fullName"])
	require.Equal(t, "alice@wonderland.com", targetValue[0]["emailAddress"])
	require.Equal(t, "person", targetValue[0]["type"])
	require.Equal(t, "Alice <alice@wonderland.com>", targetValue[0]["comment"])

	require.Equal(t, "John Connor", targetValue[1]["fullName"])
	require.Equal(t, "john@connor.mil", targetValue[1]["emailAddress"])
	require.Equal(t, "person", targetValue[1]["type"])
	require.Equal(t, "John Connor <john@connor.mil>", targetValue[1]["comment"])

	require.Equal(t, "Sarah Connor", targetValue[2]["fullName"])
	require.Equal(t, "sarah@sky.net", targetValue[2]["emailAddress"])
	require.Equal(t, "person", targetValue[2]["type"])
	require.Equal(t, "Sarah Connor <sarah@sky.net>", targetValue[2]["comment"])
}
