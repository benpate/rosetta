package funcmap

import "github.com/benpate/rosetta/compare"

func addCompareFuncs(result map[string]any) {

	result["isZero"] = func(value any) bool {
		return compare.IsZero(value)
	}

	result["notZero"] = func(value any) bool {
		return compare.NotZero(value)
	}
}
