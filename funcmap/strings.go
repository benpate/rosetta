package funcmap

import (
	"strings"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/sliceof"
)

func addStringFuncs(target map[string]any) {

	target["string"] = func(value any) string {
		return convert.String(value)
	}

	target["split"] = func(value string, separator string) sliceof.String {
		if value == "" {
			return sliceof.String{}
		}
		return strings.Split(value, separator)
	}

	target["join"] = func(values ...string) string {
		return strings.Join(values, "")
	}

	target["append"] = func(first []string, second []string) sliceof.String {
		return append(first, second...)
	}

	target["pluralize"] = func(count any, single string, plural string) string {

		if countInt := convert.Int(count); countInt == 1 {
			return single
		}

		return plural
	}

	target["lowerCase"] = func(name any) string {
		return strings.ToLower(convert.String(name))
	}

	target["trim"] = func(value string) string {
		return strings.TrimSpace(value)
	}

	target["hasPrefix"] = func(a string, b string) bool {
		return strings.HasPrefix(a, b)
	}

	target["concat"] = func(values ...any) string {
		result := ""
		for _, value := range values {
			result += convert.String(value)
		}
		return result
	}
}
