package funcmap

func All() map[string]any {

	// TODO: Consider using https://github.com/Masterminds/sprig

	result := make(map[string]any)

	addArraysFuncs(result)
	addCurrencyFuncs(result)
	addDateFuncs(result)
	addHTMLFuncs(result)
	addLogicFuncs(result)
	addMathFuncs(result)
	addStringFuncs(result)

	return result
}
