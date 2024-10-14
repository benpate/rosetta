package convert

import "time"

// Booler interface wraps the Bool() method that enables custom types to convert themselves to bool.
type Booler interface {

	// Bool returns the float64 value of the underlying object
	Bool() bool
}

// Hexer interface wraps the Hex() method that enables a custom type to convert itself into a hexadecimal string
type Hexer interface {
	Hex() string
}

// Inter interface wraps the Int() method that enables custom types to convert themselves to ints.
type Inter interface {

	// Int returns the int value of the underlying object
	Int() int
}

// Int64er interface wraps the Int64() method that enables custom types to convert themselves to int64s.
type Int64er interface {

	// Int64 returns the int64 value of the underlying object
	Int64() int64
}

// Floater interface wraps the Float() method that enables custom types to convert themselves to float64.
type Floater interface {

	// Float returns the float64 value of the underlying object
	Float() float64
}

// Nuller wraps the IsNull interface (implemented by the null.* package) that enables custom types to declare that their value is null (zero)
type Nuller interface {

	// IsNull returns TRUE if the underlying value is null
	IsNull() bool
}

// Stringer interface wraps the String() method that enables a custom type to convert themselves into strings.
type Stringer interface {

	// String returns the string value of the underlying object
	String() string
}

// Timer interface wraps the Time() method that returns the time.Time value of the underlying object
type Timer interface {
	// Time returns the time.Time value of the underlying object
	Time() time.Time
}

// Length interface wraps the Length() method that returns the length of an array or map
type LengthGetter interface {
	Length() int
}

// MapOfAnyGetter wraps the MapOfAny() method that returns a data structure as a MapOfAny
type MapOfAnyGetter interface {

	// MapOfAny returns the underlying data structure as a plain map[string]any
	MapOfAny() map[string]any
}
