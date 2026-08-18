package sliceof

// stringOKGetter is implemented by objects that can look up a named string property,
// letting this package pull a sortable/searchable string out of an arbitrary item.
type stringOKGetter interface {

	// GetStringOK returns the named string property, and TRUE if it is present
	GetStringOK(string) (string, bool)
}
