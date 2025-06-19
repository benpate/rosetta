package mapof

type Slices[K comparable, V comparable] map[K][]V

func (s Slices[K, V]) Add(key K, value V) {
	if _, exists := s[key]; !exists {
		s[key] = []V{}
	}
	s[key] = append(s[key], value)
}
