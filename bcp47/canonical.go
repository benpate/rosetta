package bcp47

// Canonical returns the canonical BCP 47 form of a language tag, and FALSE if the
// value is not a usable language.
func Canonical(value string) (string, bool) {

	tag, ok := Parse(value)

	if !ok {
		return "", false
	}

	return tag.String(), true
}
