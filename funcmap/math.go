package funcmap

import (
	"math"

	"github.com/benpate/rosetta/convert"
)

func addMathFuncs(target map[string]any) {

	target["add"] = func(a any, b any) int {
		return convert.Int(a) + convert.Int(b)
	}

	target["subtract"] = func(a any, b any) int {
		return convert.Int(a) - convert.Int(b)
	}

	target["multiply"] = func(a any, b any) int64 {
		return convert.Int64(a) * convert.Int64(b)
	}

	target["divide"] = func(a any, b any) int64 {
		return convert.Int64(a) / convert.Int64(b)
	}

	target["inc"] = func(a any) int {
		return convert.Int(a) + 1
	}

	target["min"] = func(values ...any) int {
		var result = math.MaxInt
		for _, value := range values {
			if value32 := convert.Int(value); value32 < result {
				result = value32
			}
		}
		return result
	}

	target["max"] = func(values ...any) int {
		var result = math.MinInt
		for _, value := range values {
			if value32 := convert.Int(value); value32 > result {
				result = value32
			}
		}
		return result
	}

	target["int"] = func(value any) int {
		return convert.Int(value)
	}

	target["int64"] = func(value any) int64 {
		return convert.Int64(value)
	}
}
