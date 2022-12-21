package path

import (
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
)

// Match returns TRUE if the provided Getter matches the provided expression
func Match(object any, criteria exp.Expression) bool {

	// Call criteria.Match using a custom matcher function that knows how to resolve Path values.
	return criteria.Match(func(predicate exp.Predicate) bool {

		value, ok := GetOK(object, predicate.Field)

		if !ok {
			return false
		}

		result, err := compare.WithOperator(value, predicate.Operator, predicate.Value)

		if err != nil {
			return false
		}

		return result
	})
}
