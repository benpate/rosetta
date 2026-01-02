package ranges

import "iter"

// Map transforms a rangeFunc of type T into a rangeFunc of type U using the provided transform function.
func Map[IN any, OUT any](iterator iter.Seq[IN], transform func(IN) OUT) iter.Seq[OUT] {

	return func(yield func(OUT) bool) {
		for item := range iterator {
			if mappedItem := transform(item); !yield(mappedItem) {
				return
			}
		}
	}
}
