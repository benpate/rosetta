package translate

import (
	"testing"

	"github.com/benpate/derp"
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
				// "comment" contains "Name <email>"; opt out of the default
				// no-html format so the angle-bracketed text is preserved.
				"comment": schema.String{Format: "unsafe-any"},
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

// activityStreamSchema matches the schema Emissary uses for both sides of its social rules
func activityStreamSchema() schema.Schema {
	return schema.New(schema.Object{
		Properties: schema.ElementMap{"@context": schema.Array{Items: schema.Any{}}},
		Wildcard:   schema.Any{},
	})
}

// TestForEach_Source pins forEach with each shape of source: a missing or nil source has no
// items, so the rule after it still runs
func TestForEach_Source(t *testing.T) {

	// BUG-234: a missing source used to stop every rule after it

	rules, err := NewFromJSON(`[
		{"target": "before", "value": "b"},
		{"target": "links", "forEach": "data.links", "rules": [{"target": "href", "path": "value"}]},
		{"target": "after", "value": "a"}
	]`)
	require.NoError(t, err)

	run := func(source mapof.Any) (mapof.Any, error) {
		target := mapof.Any{}
		err := rules.Execute(activityStreamSchema(), source, activityStreamSchema(), &target)
		return target, err
	}

	t.Run("missing parent, missing source, and nil source write nothing", func(t *testing.T) {
		for _, source := range []mapof.Any{
			{},
			{"data": mapof.Any{}},
			{"data": mapof.Any{"links": nil}},
		} {
			target, err := run(source)
			require.NoError(t, err, "%#v", source)
			require.Equal(t, mapof.Any{"before": "b", "after": "a"}, target, "%#v", source)
		}
	})

	t.Run("a source that is not a list is still an error", func(t *testing.T) {
		target, err := run(mapof.Any{"data": mapof.Any{"links": "not a list"}})
		require.Error(t, err)
		require.Equal(t, "Source value must implement schema.KeysGetter", derp.RootMessage(err))
		require.Equal(t, mapof.Any{"before": "b"}, target)
	})

	t.Run("empty source writes nothing, and the next rule runs", func(t *testing.T) {
		target, err := run(mapof.Any{"data": mapof.Any{"links": mapof.Any{}}})
		require.NoError(t, err)
		require.Equal(t, mapof.Any{"before": "b", "after": "a"}, target)
	})

	t.Run("each item is written as a pointer to a list of maps", func(t *testing.T) {
		target, err := run(mapof.Any{"data": mapof.Any{"links": mapof.Any{"SPOTIFY": "https://spotify.example/a"}}})
		require.NoError(t, err)
		require.Equal(t, mapof.Any{
			"before": "b",
			"links":  &sliceof.Object[mapof.Any]{{"href": "https://spotify.example/a"}},
			"after":  "a",
		}, target)
	})
}

// TestForEach_TypedSources requires forEach to read through a typed source schema: a map
// declared as an Object, and a list declared as an Array of a single type
func TestForEach_TypedSources(t *testing.T) {

	rules, err := NewFromJSON(`[
		{"target": "links", "forEach": "links", "rules": [
			{"target": "label", "path": "key"},
			{"target": "href", "path": "value"}
		]}
	]`)
	require.NoError(t, err)

	run := func(sourceSchema schema.Schema, source mapof.Any) (mapof.Any, error) {
		target := mapof.Any{}
		err := rules.Execute(sourceSchema, source, activityStreamSchema(), &target)
		return target, err
	}

	t.Run("a map declared as an Object iterates its keys", func(t *testing.T) {

		// FUNKWHALE task 1.1 step 4: Bandwagon's data.links is an Object with a wildcard
		sourceSchema := schema.New(schema.Object{Properties: schema.ElementMap{
			"links": schema.Object{Wildcard: schema.String{Format: "url"}},
		}})

		source := mapof.Any{"links": mapof.Any{"SPOTIFY": "https://spotify.example/a"}}
		target, err := run(sourceSchema, source)
		require.NoError(t, err)
		require.Equal(t, mapof.Any{
			"links": &sliceof.Object[mapof.Any]{
				{"label": "SPOTIFY", "href": "https://spotify.example/a"},
			},
		}, target)
	})

	t.Run("a list of strings gives item rules its key and value", func(t *testing.T) {

		// FUNKWHALE task 1.1 step 4: item rules used to read through the item schema, so
		// "key" and "value" were unreadable beneath a String and became ""
		sourceSchema := schema.New(schema.Object{Properties: schema.ElementMap{
			"links": schema.Array{Items: schema.String{}},
		}})

		source := mapof.Any{"links": sliceof.String{"https://a.example", "https://b.example"}}
		target, err := run(sourceSchema, source)
		require.NoError(t, err)
		require.Equal(t, mapof.Any{
			"links": &sliceof.Object[mapof.Any]{
				{"label": "0", "href": "https://a.example"},
				{"label": "1", "href": "https://b.example"},
			},
		}, target)
	})

	t.Run("an Object with named properties iterates them", func(t *testing.T) {
		sourceSchema := schema.New(schema.Object{Properties: schema.ElementMap{
			"links": schema.Object{Properties: schema.ElementMap{"home": schema.String{}}},
		}})

		source := mapof.Any{"links": mapof.Any{"home": "https://home.example"}}
		target, err := run(sourceSchema, source)
		require.NoError(t, err)
		require.Equal(t, mapof.Any{
			"links": &sliceof.Object[mapof.Any]{{"label": "home", "href": "https://home.example"}},
		}, target)
	})

	t.Run("a scalar source element is still an error", func(t *testing.T) {
		sourceSchema := schema.New(schema.Object{Properties: schema.ElementMap{
			"links": schema.String{},
		}})

		_, err := run(sourceSchema, mapof.Any{"links": "https://a.example"})
		require.Error(t, err)
		require.Equal(t, "Source element must be a list or a map", derp.RootMessage(err))
	})
}
