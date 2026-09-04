package funcmap

import "github.com/davecgh/go-spew/spew"

// addDebugFuncs registers the debug helpers in the template funcmap.
func addDebugFuncs(target map[string]any) {

	target["dump"] = func(v any) string {
		spew.Dump(v)
		return ""
	}
}
