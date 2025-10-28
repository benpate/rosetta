package slice

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {

	values := testGroupValues()
	grouper := NewGrouper(values, "group")

	require.True(t, grouper.IsHeader(-1))
	require.True(t, grouper.IsHeader(0))
	require.False(t, grouper.IsHeader(1))
	require.True(t, grouper.IsHeader(2))
	require.False(t, grouper.IsHeader(3))
	require.True(t, grouper.IsHeader(4))
	require.False(t, grouper.IsHeader(5))
	require.False(t, grouper.IsHeader(6))
	require.False(t, grouper.IsHeader(7))
	require.True(t, grouper.IsHeader(8))
	require.False(t, grouper.IsHeader(9))
	require.True(t, grouper.IsHeader(10))
	require.False(t, grouper.IsHeader(11))
	require.False(t, grouper.IsHeader(1000))
}

func TestGroupFooter(t *testing.T) {

	values := testGroupValues()
	grouper := NewGrouper(values, "group")

	require.False(t, grouper.IsFooter(-1))
	require.False(t, grouper.IsFooter(0))
	require.True(t, grouper.IsFooter(1))
	require.False(t, grouper.IsFooter(2))
	require.True(t, grouper.IsFooter(3))
	require.False(t, grouper.IsFooter(4))
	require.False(t, grouper.IsFooter(5))
	require.False(t, grouper.IsFooter(6))
	require.True(t, grouper.IsFooter(7))
	require.False(t, grouper.IsFooter(8))
	require.True(t, grouper.IsFooter(9))
	require.False(t, grouper.IsFooter(10))
	require.True(t, grouper.IsFooter(11))
	require.True(t, grouper.IsFooter(1000))
}

func testGroupValues() []mapof.String {

	return []mapof.String{
		{
			"value":       "GIPHY",
			"label":       "Giphy",
			"icon":        "film",
			"description": "Embeddable GIF Images",
			"group":       "Images",
		},
		{
			"value":       "UNSPLASH",
			"label":       "Unsplash",
			"icon":        "picture",
			"description": "Embeddable Photographs",
			"group":       "Images",
		},
		{
			"value":       "ARCGIS",
			"label":       "ArcGIS",
			"icon":        "globe",
			"description": "Geocoding for physical addresses and locations",
			"group":       "Geocoding",
		},
		{
			"value":       "GOOGLEMAPS",
			"label":       "Google Maps",
			"icon":        "globe",
			"description": "Geocoding for physical addresses and locations",
			"group":       "Geocoding",
		},
		{
			"value":       "FREEIPAPI",
			"label":       "FREEIPAPI.COM",
			"icon":        "globe",
			"description": "Geocoding for IP Addresses",
			"group":       "IP Geocoding",
		},
		{
			"value":       "IPAPI",
			"label":       "IPAPI.CO",
			"icon":        "globe",
			"description": "Geocoding for IP Addresses",
			"group":       "IP Geocoding",
		},
		{
			"value":       "IP-API",
			"label":       "IP-API.COM",
			"icon":        "globe",
			"description": "Geocoding for IP Addresses",
			"group":       "IP Geocoding",
		},
		{
			"value":       "STATIC-GEOCODER-ID",
			"label":       "Static Geocoder",
			"icon":        "globe",
			"description": "Return a fixed location for all IP address geocoding requests.",
			"group":       "IP Geocoding",
		},
		{
			"value":       "GEOAPIFY",
			"label":       "Geoapify.com",
			"icon":        "globe",
			"description": "Geosearch / Autocomplete API key",
			"group":       "Geo-Search",
		},
		{
			"value":       "NOMINATIM",
			"label":       "Nominatim",
			"icon":        "globe",
			"description": "Address Search Service",
			"group":       "Geo-Search",
		},
		{
			"value":       "STRIPE",
			"label":       "Stripe Payments",
			"icon":        "stripe",
			"description": "Users copy/paste API keys from their own Stripe Dashboard.",
			"group":       "User Payments",
		},
		{
			"value":       "STRIPE-CONNECT",
			"label":       "Stripe Connect",
			"icon":        "stripe",
			"description": "Users sign in via OAuth. Requires additional setup from admins.",
			"group":       "User Payments",
		},
	}
}
