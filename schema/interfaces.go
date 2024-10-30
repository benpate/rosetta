package schema

import "github.com/benpate/rosetta/list"

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

type AnyGetter interface {
	GetAnyOK(string) (any, bool)
}

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

type ValueGetter interface {
	GetValue() any
}

/******************************************
 * Special-Case Getter Interfaces
 ******************************************/

// PointerGetter allows objects to return a pointer to a child object
type PointerGetter interface {
	GetPointer(string) (any, bool)
}

type KeysGetter interface {
	Keys() []string
}

type LengthGetter interface {
	Length() int
}

type ArrayGetter interface {
	GetIndex(int) (any, bool)
}

/******************************************
 * Setter Interfaces
 ******************************************/

type AnySetter interface {
	SetAny(string, any) bool
}

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

type ObjectSetter interface {
	SetObject(Element, list.List, any) error
}

type StringSetter interface {
	SetString(string, string) bool
}

type ValueSetter interface {
	SetValue(any) error
}

/******************************************
 * Special-Case Setter Interfaces
 ******************************************/

type ArraySetter interface {
	SetIndex(int, any) bool
	Length() int
}

/******************************************
 * Remover Interface
 ******************************************/
type Remover interface {
	Remove(string) bool
}
