package mapof

type Int64 map[string]int64

func NewInt64() Int64 {
	return make(Int64)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x Int64) Keys() []string {
	keys := make([]string, 0, len(x))
	for key := range x {
		keys = append(keys, key)
	}
	return keys
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (x Int64) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Int64) GetInt64OK(key string) (int64, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Int64) SetInt64(key string, value int64) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Int64) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *Int64) makeNotNil() {
	if *x == nil {
		*x = make(Int64)
	}
}
