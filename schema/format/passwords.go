package format

import (
	"strconv"

	"github.com/benpate/derp"
)

func HasLowercase(arg string) StringFormat {

	return func(value string) (string, error) {
		return countCharacters(arg, value, func(ch byte) bool {
			return ((ch >= 'a') && (ch <= 'z'))
		})
	}
}

func HasUppercase(arg string) StringFormat {

	return func(value string) (string, error) {
		return countCharacters(arg, value, func(ch byte) bool {
			return ((ch >= 'A') && (ch <= 'Z'))
		})
	}
}

func HasNumbers(arg string) StringFormat {

	return func(value string) (string, error) {
		return countCharacters(arg, value, func(ch byte) bool {
			return ((ch >= '0') && (ch <= '9'))
		})
	}
}

/*
func HasSymbols(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}
*/

/* TODO
func HasEntropy(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}
*/

func countCharacters(arg string, value string, fn func(byte) bool) (string, error) {

	minCount, err := strconv.Atoi(arg)

	if err != nil {
		minCount = 1
	}

	count := 0

	for index := 0; index < len(value); index = index + 1 {

		if fn(value[index]) {
			count = count + 1
		}
	}

	if minCount > count {
		return "", derp.ValidationError("value does not meet criteria")
	}

	return value, nil
}
