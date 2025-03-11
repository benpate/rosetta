package maps

// Equal returns TRUE if the two maps are identical, having the same items in the same order, with no alterations.
func Equal[T comparable](map1 map[string]T, map2 map[string]T) bool {

	// Lengths must be identical
	if len(map1) != len(map2) {
		return false
	}

	// Items at each index must be identical
	for key := range map1 {
		if map1[key] != map2[key] {
			return false
		}
	}

	return true
}

// NotEqual returns TRUE if the two maps are NOT identical
func NotEqual[T comparable](map1 map[string]T, map2 map[string]T) bool {
	return !Equal(map1, map2)
}
