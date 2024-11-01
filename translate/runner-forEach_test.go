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
		ForEach("", "", "", []map[string]any{
			{"path": "value.name", "target": "fullName"},
			{"path": "value.email", "target": "emailAddress"},
			{"value": "person", "target": "type"},
			{"expression": "{{.value.name}} <{{.value.email}}>", "target": "comment"},
		}),
	)

	err := rules.Execute(schema.Wildcard(), &sourceValue, targetSchema, &targetValue)
	require.Nil(t, err)
	require.Equal(t, 3, len(targetValue))

	value0 := targetValue[0]
	require.Equal(t, "Alice", value0["fullName"])
	require.Equal(t, "alice@wonderland.com", value0["emailAddress"])
	require.Equal(t, "person", value0["type"])
	require.Equal(t, "Alice <alice@wonderland.com>", value0["comment"])

	value1 := targetValue[1]
	require.Equal(t, "John Connor", value1["fullName"])
	require.Equal(t, "john@connor.mil", value1["emailAddress"])
	require.Equal(t, "person", value1["type"])
	require.Equal(t, "John Connor <john@connor.mil>", value1["comment"])

	value2 := targetValue[2]
	require.Equal(t, "Sarah Connor", value2["fullName"])
	require.Equal(t, "sarah@sky.net", value2["emailAddress"])
	require.Equal(t, "person", value2["type"])
	require.Equal(t, "Sarah Connor <sarah@sky.net>", value2["comment"])
}
