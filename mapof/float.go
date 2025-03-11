package mapof

import "github.com/benpate/rosetta/maps"

type Float map[string]float64

func NewFloat() Float {
	return make(Float)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x Float) Keys() []string {
	return maps.KeysSorted(x)
}

func (x Float) Equal(value map[string]float64) bool {
	return maps.Equal(x, value)
}

func (x Float) NotEqual(value map[string]float64) bool {
	return maps.NotEqual(x, value)
}

func (x Float) IsEmpty() bool {
	return len(x) == 0
}

func (x Float) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (x Float) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Float) GetFloatOK(key string) (float64, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Float) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
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
