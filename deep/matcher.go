package deep

import (
	"slices"
	"strings"
)

// Matcher returns TRUE for the name of a property that Remove should delete
type Matcher func(name string) bool

// Named returns a Matcher for properties with one of the given names. Names are matched
// whole, and upper and lower case letters are different.
func Named(names ...string) Matcher {

	// The names are copied, so that a caller who reuses its slice does not change the Matcher
	return matchNames(slices.Clone(names))
}

// matchNames returns a Matcher for the given names, without copying them
func matchNames(names []string) Matcher {
	return func(name string) bool {
		return slices.Contains(names, name)
	}
}

// HasPrefix returns a Matcher for properties whose names begin with one of the given
// prefixes. An empty prefix is ignored.
func HasPrefix(prefixes ...string) Matcher {

	// prefixes are copied, so a caller who reuses its slice does not change the Matcher
	prefixes = slices.Clone(prefixes)

	// RULE: An empty prefix begins every name, so keeping one would delete every property
	prefixes = slices.DeleteFunc(prefixes, func(prefix string) bool {
		return prefix == ""
	})

	return func(name string) bool {

		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				return true
			}
		}

		return false
	}
}

// AnyOf returns a Matcher for properties that at least one of the given Matchers selects.
// A nil Matcher is ignored.
func AnyOf(matchers ...Matcher) Matcher {

	matchers = slices.DeleteFunc(slices.Clone(matchers), func(matcher Matcher) bool {
		return matcher == nil
	})

	return func(name string) bool {

		for _, matcher := range matchers {
			if matcher(name) {
				return true
			}
		}

		return false
	}
}
