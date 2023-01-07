package mapof

import "github.com/benpate/rosetta/convert"

type Any map[string]any

/****************************************
 * Path Getters
 ****************************************/

func (x Any) GetAny(key string) any {
	return convert.Interface(x[key])
}

func (x Any) GetBool(key string) bool {
	return convert.Bool(x[key])
}

func (x Any) GetFloat(key string) float64 {
	return convert.Float(x[key])
}

func (x Any) GetInt(key string) int {
	return convert.Int(x[key])
}

func (x Any) GetInt64(key string) int64 {
	return convert.Int64(x[key])
}

func (x Any) GetString(key string) string {
	return convert.String(x[key])
}

/****************************************
 * Path Setters
 ****************************************/

func (x *Any) SetAny(key string, value any) bool {
	(*x)[key] = value
	return true
}

func (x *Any) SetBool(key string, value bool) bool {
	(*x)[key] = value
	return true
}

func (x *Any) SetFloat(key string, value float64) bool {
	(*x)[key] = value
	return true
}

func (x *Any) SetInt(key string, value int) bool {
	(*x)[key] = value
	return true
}

func (x *Any) SetInt64(key string, value Int64) bool {
	(*x)[key] = value
	return true
}

func (x *Any) SetString(key string, value string) bool {
	(*x)[key] = value
	return true
}

/****************************************
 * Tree Traversal
 ****************************************/

func (x *Any) GetChild(key string) (any, bool) {
	result, ok := (*x)[key]
	return result, ok
}
