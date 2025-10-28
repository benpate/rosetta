package sliceof

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestAny_GroupHeader(t *testing.T) {

	values := testAnyGroupValues()

	groupBy := values.GroupBy("group")

	require.True(t, groupBy.IsHeader(-1))
	require.True(t, groupBy.IsHeader(0))
	require.False(t, groupBy.IsHeader(1))
	require.True(t, groupBy.IsHeader(2))
	require.False(t, groupBy.IsHeader(3))
	require.True(t, groupBy.IsHeader(4))
	require.False(t, groupBy.IsHeader(5))
	require.False(t, groupBy.IsHeader(6))
	require.False(t, groupBy.IsHeader(7))
	require.True(t, groupBy.IsHeader(8))
	require.False(t, groupBy.IsHeader(9))
	require.True(t, groupBy.IsHeader(10))
	require.False(t, groupBy.IsHeader(11))
	require.False(t, groupBy.IsHeader(12))
	require.False(t, groupBy.IsHeader(100))
}

func TestAny_GroupFooter(t *testing.T) {

	values := testAnyGroupValues()

	groupBy := values.GroupBy("group")

	require.False(t, groupBy.IsFooter(-1))
	require.False(t, groupBy.IsFooter(0))
	require.True(t, groupBy.IsFooter(1))
	require.False(t, groupBy.IsFooter(2))
	require.True(t, groupBy.IsFooter(3))
	require.False(t, groupBy.IsFooter(4))
	require.False(t, groupBy.IsFooter(5))
	require.False(t, groupBy.IsFooter(6))
	require.True(t, groupBy.IsFooter(7))
	require.False(t, groupBy.IsFooter(8))
	require.True(t, groupBy.IsFooter(9))
	require.False(t, groupBy.IsFooter(10))
	require.True(t, groupBy.IsFooter(11))
	require.True(t, groupBy.IsFooter(12))
	require.True(t, groupBy.IsFooter(100))
}

func testAnyGroupValues() Any {

	return Any{
		mapof.String{
			"value":       "GIPHY",
			"label":       "Giphy",
			"icon":        "film",
			"description": "Embeddable GIF Images",
			"group":       "Images",
		},
		mapof.String{
			"value":       "UNSPLASH",
			"label":       "Unsplash",
			"icon":        "picture",
			"description": "Embeddable Photographs",
			"group":       "Images",
		},
		mapof.String{
			"value":       "ARCGIS",
			"label":       "ArcGIS",
			"icon":        "globe",
			"description": "Geocoding for physical addresses and locations",
			"group":       "Geocoding",
		},
		mapof.String{
			"value":       "GOOGLEMAPS",
			"label":       "Google Maps",
			"icon":        "globe",
			"description": "Geocoding for physical addresses and locations",
			"group":       "Geocoding",
		},
		mapof.String{
			"value":       "FREEIPAPI",
			"label":       "FREEIPAPI.COM",
			"icon":        "globe",
			"description": "Geocoding for IP Addresses",
			"group":       "IP Geocoding",
		},
		mapof.String{
			"value":       "IPAPI",
			"label":       "IPAPI.CO",
			"icon":        "globe",
			"description": "Geocoding for IP Addresses",
			"group":       "IP Geocoding",
		},
		mapof.String{
			"value":       "IP-API",
			"label":       "IP-API.COM",
			"icon":        "globe",
			"description": "Geocoding for IP Addresses",
			"group":       "IP Geocoding",
		},
		mapof.String{
			"value":       "STATIC-GEOCODER-ID",
			"label":       "Static Geocoder",
			"icon":        "globe",
			"description": "Return a fixed location for all IP address geocoding requests.",
			"group":       "IP Geocoding",
		},
		mapof.String{
			"value":       "GEOAPIFY",
			"label":       "Geoapify.com",
			"icon":        "globe",
			"description": "Geosearch / Autocomplete API key",
			"group":       "Geo-Search",
		},
		mapof.String{
			"value":       "NOMINATIM",
			"label":       "Nominatim",
			"icon":        "globe",
			"description": "Address Search Service",
			"group":       "Geo-Search",
		},
		mapof.String{
			"value":       "STRIPE",
			"label":       "Stripe Payments",
			"icon":        "stripe",
			"description": "Users copy/paste API keys from their own Stripe Dashboard.",
			"group":       "User Payments",
		},
		mapof.String{
			"value":       "STRIPE-CONNECT",
			"label":       "Stripe Connect",
			"icon":        "stripe",
			"description": "Users sign in via OAuth. Requires additional setup from admins.",
			"group":       "User Payments",
		},
	}
}
