package format

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/bcp47"
)

// Language accepts an empty value or any usable BCP 47 language tag, rewriting it to its
// canonical form so that one language is always stored with one spelling ("EN-us" is "en-US").
func Language(_ string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		if canonical, ok := bcp47.Canonical(value); ok {
			return canonical, nil
		}

		return "", derp.Validation("Value must be a language tag, such as en-US", value)
	}
}
