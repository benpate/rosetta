package bcp47

import "golang.org/x/text/language"

// Parse parses a language tag into its canonical form, keeping its language, script,
// region, and variants, and dropping extensions and private-use subtags.
func Parse(value string) (language.Tag, bool) {

	// RULE: Values are bounded before parsing, because most come from a visitor or an author
	if (value == "") || (len(value) > maxTagLength) {
		return language.Und, false
	}

	tag, err := language.Parse(value)

	if err != nil {
		return language.Und, false
	}

	// Rebuild the tag without extensions ("-u-ca-gregory") or private use ("-x-twain"),
	// which describe formatting preferences, not a different language
	base, script, region := tag.Raw()
	tag, err = language.Compose(base, script, region, tag.Variants())

	if err != nil {
		return language.Und, false
	}

	// RULE: "und" (undetermined, also any private-use tag) and "mul" (Accept-Language's "*")
	// name no language that content can be written in
	switch base.String() {
	case "und", "mul":
		return language.Und, false
	}

	return tag, true
}
