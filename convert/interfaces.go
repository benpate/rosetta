package convert

// Booler interface wraps the Bool() function that enables custom types to convert themselves to bool.
type Booler interface {

	// Bool returns the float64 value of the underlying object
	Bool() bool
}

// Inter interface wraps the Int() function that enables custom types to convert themselves to ints.
type Inter interface {

	// Int returns the int value of the underlying object
	Int() int
}

// Floater interface wraps the Float() function that enables custom types to convert themselves to float64.
type Floater interface {

	// Float returns the float64 value of the underlying object
	Float() float64
}

// Nuller wraps the IsNull interface (implemented by the null.* package) that enables custom types to declare that their value is null (zero)
type Nuller interface {

	// IsNull returns TRUE if the underlying value is null
	IsNull() bool
}

// Stringer interface wraps the String() function that enables a custom type to convert themselves into strings.
type Stringer interface {

	// String returns the string value of the underlying object
	String() string
}

// Hexer interface wraps the Hex() function that enables a custom type to convert itself into a hexadecimal string
type Hexer interface {
	Hex() string
}
