package mapof

type Float map[string]float64

func (x Float) GetFloat(key string) float64 {
	return x[key]
}

func (x *Float) SetFloat(key string, value float64) bool {
	(*x)[key] = value
	return true
}
