package path

// Getter interface allows other objects to make it easy to trace through their property trees, and get values from them.
type Getter interface {
	GetPath(string) (any, bool)
}

// Setter interface allows other objects to make it easy to trace through their property trees, and set values into them.
type Setter interface {
	SetPath(string, any) error
}

// Deleter interface allows other objects to make it easy to trace through their property trees, and delete values from them.
type Deleter interface {
	DeletePath(string) error
}

/**************************
 * New Style Getters
 **************************/

type BoolGetter interface {
	GetBool(string) bool
}

type FloatGetter interface {
	GetFloat(string) float64
}

type IntGetter interface {
	GetInt(string) int
}

type Int64Getter interface {
	GetInt64(string) int64
}

type StringGetter interface {
	GetString(string) string
}

type ChildGetter interface {
	GetChild(string) (any, bool)
}

/**************************
 * New Style Setters
 **************************/

type BoolSetter interface {
	SetBool(string, bool) bool
}

type FloatSetter interface {
	SetFloat(string, float64) bool
}

type IntSetter interface {
	SetInt(string, int) bool
}

type Int64Setter interface {
	SetInt64(string, int64) bool
}

type StringSetter interface {
	SetString(string, string) bool
}

/**************************
 * New Style Getter/Setters
 **************************/

type BoolGetterSetter interface {
	BoolGetter
	BoolSetter
}

type FloatGetterSetter interface {
	FloatGetter
	FloatSetter
}

type IntGetterSetter interface {
	IntGetter
	IntSetter
}

type Int64GetterSetter interface {
	Int64Getter
	Int64Setter
}

type StringGetterSetter interface {
	StringGetter
	StringSetter
}
