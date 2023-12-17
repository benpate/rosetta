package slice

func Reverse[T any](x []T) []T {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}
