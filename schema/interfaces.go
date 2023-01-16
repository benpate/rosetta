package schema

// Nullable interface wraps the IsNull method, that helps an object
// to identify if it contains a null value or not.  This mirrors
// the null.Nullable interface here, for convenience.
type Nullable interface {
	IsNull() bool
}

type Enumerator interface {
	Enumerate() []string
}

/******************************************
 * Getter Interfaces
 ******************************************/

type BoolGetter interface {
	GetBoolOK(string) (bool, bool)
}

type FloatGetter interface {
	GetFloatOK(string) (float64, bool)
}

type IntGetter interface {
	GetIntOK(string) (int, bool)
}

type Int64Getter interface {
	GetInt64OK(string) (int64, bool)
}

type StringGetter interface {
	GetStringOK(string) (string, bool)
}

/******************************************
 * Special-Case Getter Interfaces
 ******************************************/

type ObjectGetter interface {
	GetObjectOK(string) (any, bool)
}

// LengthGetter allows arrays to report their total length
type LengthGetter interface {
	Length() int
}

/******************************************
 * Setter Interfaces
 ******************************************/

type BoolSetter interface {
	SetBoolOK(string, bool) bool
}

type FloatSetter interface {
	SetFloatOK(string, float64) bool
}

type IntSetter interface {
	SetIntOK(string, int) bool
}

type Int64Setter interface {
	SetInt64OK(string, int64) bool
}

type StringSetter interface {
	SetStringOK(string, string) bool
}

/******************************************
 * Remover Interface
 ******************************************/
type Remover interface {
	Remove(string) bool
}
