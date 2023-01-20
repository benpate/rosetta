package mapof

type Float map[string]float64

func (x Float) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Float) GetFloatOK(key string) (float64, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Float) SetFloatOK(key string, value float64) bool {
	if *x == nil {
		*x = make(Float)
	}
	(*x)[key] = value
	return true
}
