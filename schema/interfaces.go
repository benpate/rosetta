package schema

import "github.com/benpate/rosetta/list"

// Nullable interface wraps the IsNull method, that helps an object
// to identify if it contains a null value or not.  This mirrors
// the null.Nullable interface here, for convenience.
type Nullable interface {
	IsNull() bool
}

// Enumerator interface wraps the Enumerate method, that helps an object
// to return a list of string values representing its contents.
type Enumerator interface {

	// Enumerate returns a list of string values representing the contents of this object
	Enumerate() []string
}

/******************************************
 * Getter Interfaces
 ******************************************/

// AnyGetter allows an object to get a property of any type element by name
type AnyGetter interface {

	// GetAnyOK gets the property with the specified name
	GetAnyOK(string) (any, bool)
}

// BoolGetter allows an object to get a boolean property by name
type BoolGetter interface {

	// GetBoolOK gets the boolean property with the specified name
	GetBoolOK(string) (bool, bool)
}

// FloatGetter allows an object to get a float64 property by name
type FloatGetter interface {

	// GetFloatOK gets the float64 property with the specified name
	GetFloatOK(string) (float64, bool)
}

// IntGetter allows an object to get an int property by name
type IntGetter interface {

	// GetIntOK gets the int property with the specified name
	GetIntOK(string) (int, bool)
}

// Int64Getter allows an object to get an int64 property by name
type Int64Getter interface {

	// GetInt64OK gets the int64 property with the specified name
	GetInt64OK(string) (int64, bool)
}

// StringGetter allows an object to get a string property by name
type StringGetter interface {

	// GetStringOK gets the string property with the specified name
	GetStringOK(string) (string, bool)
}

// ValueGetter allows an object to get its entire value
type ValueGetter interface {

	// GetValue gets the entire value of the object
	GetValue() any
}

/******************************************
 * Special-Case Getter Interfaces
 ******************************************/

// PointerGetter allows objects to return a pointer to a child object
type PointerGetter interface {

	// GetPointer gets a pointer to the child object with the specified name
	GetPointer(string) (any, bool)
}

// KeysGetter allows an object to return a list of its keys
type KeysGetter interface {

	// Keys returns a list of the object's keys
	Keys() []string
}

// LengthGetter allows an object to return the length of an array
type LengthGetter interface {

	// Length returns the length of the array
	Length() int
}

// ArrayGetter allows an object to get a value at a specific index in an array
type ArrayGetter interface {

	// GetIndex gets the value at the specified index
	GetIndex(int) (any, bool)
}

/******************************************
 * Setter Interfaces
 ******************************************/

// AnySetter allows an object to set a property of any type element by name
type AnySetter interface {

	// SetAny sets the property with the specified name
	SetAny(string, any) bool
}

// BoolSetter allows an object to set a boolean property by name
type BoolSetter interface {

	// SetBool sets the boolean property with the specified name
	SetBool(string, bool) bool
}

// FloatSetter allows an object to set a float64 property by name
type FloatSetter interface {

	// SetFloat sets the float64 property with the specified name
	SetFloat(string, float64) bool
}

// IntSetter allows an object to set an int property by name
type IntSetter interface {

	// SetInt sets the int property with the specified name
	SetInt(string, int) bool
}

// Int64Setter allows an object to set an int64 property by name
type Int64Setter interface {

	// SetInt64 sets the int64 property with the specified name
	SetInt64(string, int64) bool
}

// ObjectSetter allows an object to set a child object by path
type ObjectSetter interface {

	// SetObject sets the child object at the specified path
	SetObject(Element, list.List, any) error
}

// StringSetter allows an object to set a string property by name
type StringSetter interface {

	// SetString sets the string property with the specified name
	SetString(string, string) bool
}

// ValueSetter allows an object to set its entire value
type ValueSetter interface {

	// SetValue sets the entire value of the object
	SetValue(any) error
}

/******************************************
 * Special-Case Setter Interfaces
 ******************************************/

// ArraySetter allows an object to set a value at a specific index in an array
type ArraySetter interface {

	// SetIndex sets the value at the specified index
	SetIndex(int, any) bool

	// Length returns the length of the array
	Length() int
}

/******************************************
 * Remover Interface
 ******************************************/

// Remover allows an object to remove a child element by name
type Remover interface {

	// Remove removes the child element with the specified name
	Remove(string) bool
}
