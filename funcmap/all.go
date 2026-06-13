package funcmap

// All returns a map of every template function provided by this package, keyed by name.
func All() map[string]any {

	// TODO: Consider using https://github.com/Masterminds/sprig

	result := make(map[string]any)

	addArraysFuncs(result)
	addCompareFuncs(result)
	addCurrencyFuncs(result)
	addDateFuncs(result)
	addHTMLFuncs(result)
	addLogicFuncs(result)
	addMathFuncs(result)
	addStringFuncs(result)

	return result
}
