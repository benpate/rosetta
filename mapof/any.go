package mapof

import "github.com/benpate/rosetta/convert"

type Any map[string]any

/****************************************
 * Path Getters
 ****************************************/

func (x Any) GetAny(key string) (any, bool) {
	if value, ok := x[key]; ok {
		return convert.Interface(value), true
	}
	return nil, false

}

func (x Any) GetBool(key string) (bool, bool) {
	if value, ok := x[key]; ok {
		return convert.BoolOk(value, false)
	}
	return false, false
}

func (x Any) GetFloat(key string) (float64, bool) {
	if value, ok := x[key]; ok {
		return convert.FloatOk(value, 0)
	}
	return 0, false
}

func (x Any) GetInt(key string) (int, bool) {
	if value, ok := x[key]; ok {
		return convert.IntOk(value, 0)
	}
	return 0, false
}

func (x Any) GetInt64(key string) (int64, bool) {
	if value, ok := x[key]; ok {
		return convert.Int64Ok(value, 0)
	}
	return 0, false
}

func (x Any) GetString(key string) (string, bool) {
	if value, ok := x[key]; ok {
		return convert.StringOk(value, "")
	}
	return "", false
}

/****************************************
 * Path Setters
 ****************************************/

func (x *Any) SetAny(key string, value any) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetBool(key string, value bool) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetInt(key string, value int) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetInt64(key string, value Int64) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetString(key string, value string) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) makeNotNil() {
	if *x == nil {
		*x = make(Any)
	}
}

/****************************************
 * Tree Traversal
 ****************************************/

func (x *Any) GetChild(key string) (any, bool) {
	result, ok := (*x)[key]
	return result, ok
}
