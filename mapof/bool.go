package mapof

type Bool map[string]bool

func NewBool() Bool {
	return make(Bool)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x Bool) Keys() []string {
	keys := make([]string, 0, len(x))
	for key := range x {
		keys = append(keys, key)
	}
	return keys
}

func (x Bool) Equal(value Bool) bool {
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

func (x Bool) IsEmpty() bool {
	return len(x) == 0
}

func (x Bool) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (x Bool) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

func (x Bool) GetBoolOK(key string) (bool, bool) {
	if result, ok := x[key]; ok {
		return result, true
	}
	return false, false
}

func (x *Bool) SetBool(key string, value bool) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Bool) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *Bool) makeNotNil() {
	if *x == nil {
		*x = make(Bool)
	}
}
