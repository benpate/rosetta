package format

import (
	"errors"
	"regexp"
)

// MatchRegex returns a StringFormat that requires the value to match the regular expression arg.
func MatchRegex(arg string) StringFormat {

	rx, err := regexp.Compile(arg)

	return func(value string) (string, error) {

		if err != nil {
			return "", err
		}

		if rx.Match([]byte(value)) {
			return value, nil
		}

		return value, errors.New("Does not match Regular Expression: " + arg)
	}
}
