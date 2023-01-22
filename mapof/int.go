package mapof

type Int map[string]int

func (x Int) GetInt(key string) (int, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Int) SetInt(key string, value int) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Int) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *Int) makeNotNil() {
	if *x == nil {
		*x = make(Int)
	}
}
