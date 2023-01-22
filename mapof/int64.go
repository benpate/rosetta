package mapof

type Int64 map[string]int64

func (x Int64) GetInt64(key string) (int64, bool) {
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
