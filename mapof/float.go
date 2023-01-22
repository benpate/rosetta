package mapof

type Float map[string]float64

func (x Float) GetFloatOK(key string) (float64, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Float) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Float) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *Float) makeNotNil() {
	if *x == nil {
		*x = make(Float)
	}
}
