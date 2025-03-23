package funcmap

import "github.com/benpate/rosetta/convert"

func addCurrencyFuncs(target map[string]any) {

	target["dollarFormat"] = func(value any) string {

		var unitAmount int64

		switch value := value.(type) {
		case float32:
			unitAmount = int64(value * 100)
		case float64:
			unitAmount = int64(value * 100)
		default:
			unitAmount = convert.Int64(value)
		}

		stringValue := convert.String(unitAmount)
		length := len(stringValue)
		for length < 3 {
			stringValue = "0" + stringValue
			length = len(stringValue)
		}
		return "$" + stringValue[:length-2] + "." + stringValue[length-2:]
	}

}
