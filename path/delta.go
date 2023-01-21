package path

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

// Delta tracks changes to an object
type Delta struct {
	object  any
	changed bool
}

// NewDelta returns a fully initialized Delta object
func NewDelta(object any) Delta {
	return Delta{
		object: object,
	}
}

// SetBool tracks changes to a bool value and collects errors
func (d *Delta) SetBool(path string, value bool) error {

	if leaf, last, ok := getLeaf(d.object, list.Dot(path)); ok {

		if getter, ok := leaf.(BoolGetterSetter); ok {

			if getter.GetBool(last) != value {
				getter.SetBool(last, value)
				d.changed = true
			}

			return nil
		}
	}

	return derp.NewInternalError("delta.SetBool", "Unable to set bool", path, value)
}

// SetFloat tracks changes to a float value and collects errors
func (d *Delta) SetFloat(path string, value float64) error {

	if leaf, last, ok := getLeaf(d.object, list.Dot(path)); ok {

		if getter, ok := leaf.(FloatGetterSetter); ok {

			if getter.GetFloat(last) != value {
				getter.SetFloat(last, value)
				d.changed = true
			}

			return nil
		}
	}

	return derp.NewInternalError("delta.SetBool", "Unable to set float", path, value)
}

// SetInt tracks changes to an int value and collects errors
func (d *Delta) SetInt(path string, value int) error {

	if leaf, last, ok := getLeaf(d.object, list.Dot(path)); ok {

		if getter, ok := leaf.(IntGetterSetter); ok {

			if getter.GetInt(last) != value {
				getter.SetInt(last, value)
				d.changed = true
			}

			return nil
		}
	}

	return derp.NewInternalError("delta.SetBool", "Unable to set int", path, value)
}

// SetInt64 tracks changes to an int64 value and collects errors
func (d *Delta) SetInt64(path string, value int64) error {

	if leaf, last, ok := getLeaf(d.object, list.Dot(path)); ok {

		if getter, ok := leaf.(Int64GetterSetter); ok {

			if getter.GetInt64(last) != value {
				getter.SetInt64(last, value)
				d.changed = true
			}

			return nil
		}
	}

	return derp.NewInternalError("delta.SetBool", "Unable to set int64", path, value)
}

// SetString tracks changes to a string value and collects errors
func (d *Delta) SetString(path string, value string) error {

	if leaf, last, ok := getLeaf(d.object, list.Dot(path)); ok {

		if getter, ok := leaf.(StringGetterSetter); ok {

			if getter.GetString(last) != value {
				getter.SetString(last, value)
				d.changed = true
			}

			return nil
		}
	}

	return derp.NewInternalError("delta.SetBool", "Unable to set string", path, value)
}

// HasChanged returns TRUE if any of the values have been changed
func (d *Delta) HasChanged() bool {
	return d.changed
}
