package schema

type ElementMap map[string]Element

// Nullable interface wraps the IsNull method, that helps an object
// to identify if it contains a null value or not.  This mirrors
// the null.Nullable interface here, for convenience.
type Nullable interface {
	IsNull() bool
}

type Enumerator interface {
	Enumerate() []string
}

type Getter interface {
	Get(string) (any, error)
}

type Setter interface {
	Set(string, any) error
}

type BoolGetter interface {
	GetBool(string) (bool, error)
}

type BoolSetter interface {
	SetBool(string, bool) error
}

type FloatGetter interface {
	GetFloat(string) (float64, error)
}

type FloatSetter interface {
	SetFloat(string, float64) error
}

type IntGetter interface {
	GetInt(string) (int, error)
}

type IntSetter interface {
	SetInt(string, int) error
}

type Int64Getter interface {
	GetInt64(string) (int64, error)
}

type Int64Setter interface {
	SetInt64(string, int64) error
}

type StringGetter interface {
	GetString(string) (string, error)
}

type StringSetter interface {
	SetString(string, string) error
}
