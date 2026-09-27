// Package mapof provides string-keyed map types with typed values and helper
// methods.
//
// Any, String, Int, Int64, Float, Bool, and the generic Object[T] are all plain
// map types under their methods, so they can be ranged over, indexed, and
// serialized exactly like the maps they wrap. What they add is typed accessors
// — GetString, GetInt, GetBool, and their OK-returning variants — that convert
// a loosely-typed value on the way out, plus the Getter and Setter
// implementations that let the schema package address them by path.
//
// That schema integration is the reason these types exist: they are the
// carriers a JSON-Schema-shaped document is read into and written back out of
// when there is no Go struct to bind to.
//
// Template is the exception to "exactly like the maps they wrap": its
// strings may be stored as loose.Template values, compiled when decoded or
// set, and its getters take the data to render them against. Read it only
// through its methods.
package mapof
