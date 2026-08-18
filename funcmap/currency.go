package funcmap

import (
	"strings"

	"github.com/benpate/rosetta/convert"
)

// addCurrencyFuncs registers the currency formatting helpers in the template funcmap.
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

		// Lift the sign off the digits before padding.  A "-" left in place counts as a
		// digit, which pushes the decimal point into the wrong spot ("-5" => "$0.-5"), and
		// the sign belongs OUTSIDE the currency symbol anyway ("-$0.05", not "$-0.05").
		// Splitting the string (instead of negating the number) also avoids the overflow
		// that negating math.MinInt64 would cause.
		sign := ""

		if unsigned, isNegative := strings.CutPrefix(stringValue, "-"); isNegative {
			sign = "-"
			stringValue = unsigned
		}

		// Pad to at least "0dd" so there is always a digit on each side of the decimal point
		for len(stringValue) < 3 {
			stringValue = "0" + stringValue
		}

		length := len(stringValue)
		return sign + "$" + stringValue[:length-2] + "." + stringValue[length-2:]
	}

}
