package mapof

type String map[string]string

func NewString() String {
	return make(String)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x String) Keys() []string {
	keys := make([]string, 0, len(x))
	for key := range x {
		keys = append(keys, key)
	}
	return keys
}

func (x String) Equal(value String) bool {
	// Lengths must be identical
	if len(x) != len(value) {
		return false
	}

	// Items at each index must be identical
	for key := range x {
		if x[key] != value[key] {
			return false
		}
	}

	return true
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (x String) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x String) GetStringOK(key string) (string, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *String) SetString(key string, value string) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *String) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *String) makeNotNil() {
	if *x == nil {
		*x = make(String)
	}
}
