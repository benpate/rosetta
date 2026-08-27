package convert

// SliceOf converts the value into a slice of T.
// It works with any value that the matching SliceOfXxx function accepts, plus any slice, array,
// pointer, or ArrayGetter whose items are already Ts.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOf[T any](value any) []T {
	if result, _ := SliceOfOk[T](value); result != nil {
		return result
	}

	return make([]T, 0)
}

// SliceOfOk converts the value into a slice of T.
// It returns TRUE if the conversion was lossless, and FALSE otherwise.
func SliceOfOk[T any](value any) ([]T, bool) {

	// A value that is already the right shape passes straight through, with no
	// intermediate slice and no per-item boxing
	if typed, ok := value.([]T); ok {
		return typed, true
	}

	// A value that is not a collection at all converts to nothing
	items, isCollection := SliceOfAnyOk(value)

	if !isCollection {
		return make([]T, 0), false
	}

	// Hand the scalar element types to the concrete converter that knows their coercion rules.
	// Dispatching on a nil *T rather than a zero T keeps an interface element type (T = any)
	// from vanishing into a nil interface, which would match no case at all.
	switch any((*T)(nil)).(type) {

	case *any:
		return castSliceOf[T](items), true

	case *string:
		result, lossless := SliceOfStringOk(value)
		return castSliceOf[T](result), lossless

	case *int:
		result, lossless := SliceOfIntOk(value)
		return castSliceOf[T](result), lossless

	case *int64:
		result, lossless := SliceOfInt64Ok(value)
		return castSliceOf[T](result), lossless

	case *float64:
		result, lossless := SliceOfFloatOk(value)
		return castSliceOf[T](result), lossless

	case *map[string]any:
		result, lossless := SliceOfMapOk(value)
		return castSliceOf[T](result), lossless
	}

	// Reshape every other element type structurally
	return structuralSliceOf[T](items)
}

// structuralSliceOf reshapes already-unwrapped items into a []T without converting the items
// themselves, and returns FALSE if it cannot.
func structuralSliceOf[T any](items []any) ([]T, bool) {

	result := make([]T, 0, len(items))

	// RULE: There is no value conversion for an unknown element type. This package cannot
	// know how a caller wants an arbitrary value rendered as a T, so each item must be one.
	for _, item := range items {

		typed, ok := item.(T)

		if !ok {
			return make([]T, 0), false
		}

		result = append(result, typed)
	}

	// Slice, slice, baby
	return result, true
}

// castSliceOf re-types a slice produced by one of the concrete SliceOfXxx converters
// into the caller's []T.
func castSliceOf[T any](value any) []T {

	// The dispatch in SliceOfOk has already established that these are the same type, so this
	// costs an interface assertion rather than a copy.
	if result, ok := value.([]T); ok {
		return result
	}

	// UNREACHABLE: every caller is a branch of that type switch, each of which has proven T to
	// be the element type of the slice it passes in. The comma-ok form has to be handled, and
	// an empty slice beats a panic if the invariant is ever broken.
	return make([]T, 0)
}
