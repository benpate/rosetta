package format

// Function is a function that takes an optional parameter and generates a StringFormat function
type Function func(string) StringFormat

// StringFormat tries to force a value to fit the desired format.  If it cannot
// safely convert the value into the specified format, then it returns empty string and an error message.
type StringFormat func(string) (string, error)
