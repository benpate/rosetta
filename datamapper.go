package rosetta

import "github.com/benpate/path"

// DataMapping contains all of the rules for mapping from one data format to another.
type DataMapping struct {
	From         path.Path
	To           path.Path
	DefaultValue interface{}
}
