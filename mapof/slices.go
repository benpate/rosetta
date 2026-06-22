package mapof

// Slices is a map from keys to slices of values, i.e. a one-to-many multimap.
type Slices[K comparable, V comparable] map[K][]V

// Add appends a value to the slice stored at the key, creating the slice if necessary.
func (s Slices[K, V]) Add(key K, value V) {
	if _, exists := s[key]; !exists {
		s[key] = []V{}
	}
	s[key] = append(s[key], value)
}

// Flatten returns all values across every key as a single slice (in unspecified order).
func (s Slices[K, V]) Flatten() []V {
	var result []V
	for _, values := range s {
		result = append(result, values...)
	}
	return result
}
