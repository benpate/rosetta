package translate

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

func TestAppend(t *testing.T) {

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

	// MAPPING RULES
	rules := New(
		Append(mapof.Any{
			"name": "John Wick", "email": "john@wickindustries.com",
		}, ""),
	)

	err := rules.Execute(schema.Wildcard(), &sourceValue, targetSchema, &sourceValue)

	require.Nil(t, err)
	require.Equal(t, 4, len(sourceValue))

	value0 := sourceValue[0]
	require.Equal(t, "Alice", value0["name"])
	require.Equal(t, "alice@wonderland.com", value0["email"])

	value1 := sourceValue[1]
	require.Equal(t, "John Connor", value1["name"])
	require.Equal(t, "john@connor.mil", value1["email"])

	value2 := sourceValue[2]
	require.Equal(t, "Sarah Connor", value2["name"])
	require.Equal(t, "sarah@sky.net", value2["email"])

	value3 := sourceValue[3]
	require.Equal(t, "John Wick", value3["name"])
	require.Equal(t, "john@wickindustries.com", value3["email"])
}
