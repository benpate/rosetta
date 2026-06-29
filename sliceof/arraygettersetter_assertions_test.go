package sliceof_test

import (
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
)

// Compile-time assertions that every concrete sliceof type satisfies schema.ArrayGetterSetter
// through its POINTER. schema.validate_Array requires this interface for any value used as an
// Array property; a type that silently loses it (e.g. a receiver changed from pointer to value)
// would only fail at runtime during validation. These assertions move that failure to compile time.
//
// RULE: the assertion must be on the POINTER (&T), because SetIndex has a pointer receiver.
var (
	_ schema.ArrayGetterSetter = (*sliceof.Any)(nil)
	_ schema.ArrayGetterSetter = (*sliceof.Int)(nil)
	_ schema.ArrayGetterSetter = (*sliceof.Float)(nil)
	_ schema.ArrayGetterSetter = (*sliceof.String)(nil)
	_ schema.ArrayGetterSetter = (*sliceof.Object[string])(nil)
	_ schema.ArrayGetterSetter = (*sliceof.MapOfAny)(nil)
	_ schema.ArrayGetterSetter = (*sliceof.MapOfString)(nil)
)
