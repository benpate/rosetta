package schema

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// childKind names the shape a path segment needs: a map, a list, or (under Any) either
type childKind int

const (
	kindAny childKind = iota
	kindMap
	kindList
	kindScalar
)

// childRef is a child container that SetProperty writes into, and the value to store back in
// its parent once the write succeeds
type childRef struct {
	pointer any
	value   func() any
}

// convertibleContainers are the rosetta types that a plain Go map or slice is written through
var convertibleContainers = []reflect.Type{
	reflect.TypeFor[mapof.Any](),
	reflect.TypeFor[sliceof.Any](),
	reflect.TypeFor[sliceof.Object[mapof.Any]](),
}

// setMapChild sets the value at head.tail within a map, creating or entering the child at head
func setMapChild(
	setter KeySetter, childElement Element, object any, head string, tail string, value any,
) error {

	const location = "schema.setMapChild"

	// A single key is stored as given
	if tail == "" {
		if err := setter.SetKey(head, value); err != nil {
			return derp.Wrap(err, location, "Setting value", head, typeName(value))
		}
		return nil
	}

	// A map's GetPointer returns the child itself, or the zero value of its child type
	var current any
	if getter, ok := object.(PointerGetter); ok {
		current, _ = getter.GetPointer(head)
	}

	child, err := enterChild(childElement, current, nextSegment(tail))

	if err != nil {
		return derp.Wrap(err, location, "Entering child", head)
	}

	if err := SetProperty(childElement, child.pointer, tail, value); err != nil {
		return derp.Wrap(err, location, "Setting value", head)
	}

	// Store the child back, because a copied map or a grown list is a new value
	if err := setter.SetKey(head, child.value()); err != nil {
		return derp.Wrap(err, location, "Storing child", head, typeName(child.value()))
	}

	return nil
}

// setTypedMapKey stores a single key in a map that declares itself through MapTyper but has no
// KeySetter, using SetAny, or SetString for a string.  It returns FALSE if neither applies.
func setTypedMapKey(object any, key string, value any) bool {

	if typer, ok := object.(MapTyper); !ok || !typer.IsMap() {
		return false
	}

	if setter, ok := object.(AnySetter); ok {
		return setter.SetAny(key, value)
	}

	if text, isString := value.(string); isString {
		if setter, ok := object.(StringSetter); ok {
			return setter.SetString(key, text)
		}
	}

	return false
}

// setOtherChild sets the value at head.tail within a struct or a list, which hand back each
// child through PointerGetter
func setOtherChild(
	parentElement Element, childElement Element, object any,
	path string, head string, tail string, value any,
) error {

	const location = "schema.setOtherChild"

	// RULE: Under an Any element, a list grows by at most one item per write (BUG-116)
	if err := checkAnyIndex(parentElement, object, head); err != nil {
		return derp.Wrap(err, location, "Setting list item", path)
	}

	// A list stores a whole item directly, because its GetPointer may return a copy
	if (tail == "") && setListItem(object, head, value) {
		return nil
	}

	// PointerGetter works for Structs, Slices, and Arrays. We have already
	// descended to "head", so continue with the child element and the tail.
	if getter, ok := object.(PointerGetter); ok {
		if subPointer, ok := getter.GetPointer(head); ok {
			return setPointerChild(object, childElement, subPointer, head, tail, value)
		}
	}

	// Cannot set the value
	return derp.Internal(
		location, "Target Object must be a KeySetter or PointerGetter", path, typeName(object),
	)
}

// enterChild returns the existing child to write into, or creates one when it is missing.  A
// child of the wrong kind for its element is an error, never replaced.
func enterChild(element Element, current any, next string) (childRef, error) {

	const location = "schema.enterChild"

	if isMissing(current) {
		return newChild(element, current, next)
	}

	have := kindOf(reflect.TypeOf(current))

	// RULE: Never replace a value of the wrong kind. Report it instead
	if have == kindScalar {
		return childRef{}, derp.Validation("Expected a map or a list", location, typeName(current))
	}

	if want := wantedKind(element); (want != kindAny) && (want != have) {
		return childRef{}, kindError(location, want, have, typeName(current))
	}

	// A child held by pointer is written in place
	if reflect.TypeOf(current).Kind() == reflect.Pointer {
		return childRef{pointer: current, value: func() any { return current }}, nil
	}

	return copyChild(current)
}

// nextSegment returns the first segment of a dotted path
func nextSegment(path string) string {
	segment, _, _ := strings.Cut(path, ".")
	return segment
}

// checkAnyIndex requires a write into a list beneath an Any element to use a plain index no
// greater than the list's length, so one path cannot allocate a list of any size
func checkAnyIndex(parentElement Element, object any, head string) error {

	const location = "schema.checkAnyIndex"

	if _, isAny := parentElement.(Any); !isAny {
		return nil
	}

	lister, isList := object.(ArrayGetter)

	if !isList {
		return nil
	}

	if !isIndex(head) {
		return derp.Validation("List index must be a number: "+head, location)
	}

	if index, err := strconv.Atoi(head); (err != nil) || (index > lister.Length()) {
		return derp.Validation("List index is past the end of the list: "+head, location)
	}

	return nil
}

// setListItem stores a whole item in a list through SetAny.  It returns FALSE if the object is
// not a list, or cannot store the item.
func setListItem(object any, head string, value any) bool {

	if _, isList := object.(ArrayGetter); !isList {
		return false
	}

	if setter, ok := object.(AnySetter); ok {
		return setter.SetAny(head, value)
	}

	return false
}

// setPointerChild sets the value at tail within a child that a struct or list returned through
// GetPointer, creating a missing list item first
func setPointerChild(
	object any, childElement Element, subPointer any, head string, tail string, value any,
) error {

	// A missing list item is created from its element, then stored in the list
	if isMissing(subPointer) && (tail != "") {
		return createListItem(object, childElement, head, tail, value)
	}

	if err := SetProperty(childElement, subPointer, tail, value); err != nil {
		return err
	}

	// A list returns a copy of a list item, so the item's new length is stored back
	if tail != "" {
		storeListItem(object, subPointer, head)
	}

	return nil
}

// isMissing returns TRUE for nil, and for a nil map, slice, pointer, or interface
func isMissing(value any) bool {

	if value == nil {
		return true
	}

	switch reflected := reflect.ValueOf(value); reflected.Kind() {

	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface:
		return reflected.IsNil()
	}

	return false
}

// newChild creates a missing child: an empty value of the parent's own child type when the
// parent supplied one, otherwise the container its element names
func newChild(element Element, current any, next string) (childRef, error) {

	// A typed parent, such as mapof.Object[T], supplies the type of its children
	if current != nil {
		return newTypedChild(wantedKind(element), reflect.TypeOf(current))
	}

	// Under Any, the next path segment decides between a list and a map
	if want := wantedKind(element); (want == kindList) || ((want == kindAny) && isIndex(next)) {

		// A new list is stored by pointer, because sliceof.Any reads its items through a pointer
		created := sliceof.NewAny()
		return childRef{pointer: &created, value: func() any { return &created }}, nil
	}

	return copyChild(mapof.NewAny())
}

// kindOf returns whether a type (or the type it points to) is a map, a list, or a scalar.
// A struct counts as a map, because its fields are reached by name.
func kindOf(valueType reflect.Type) childKind {

	if valueType.Kind() == reflect.Pointer {
		valueType = valueType.Elem()
	}

	switch valueType.Kind() {

	case reflect.Map, reflect.Struct:
		return kindMap

	case reflect.Slice, reflect.Array:
		return kindList
	}

	return kindScalar
}

// wantedKind returns the kind of value an element describes
func wantedKind(element Element) childKind {

	switch element.(type) {

	case Object:
		return kindMap

	case Array:
		return kindList
	}

	return kindAny
}

// kindError reports a child, named by its type, whose kind does not match its element
func kindError(location string, want childKind, have childKind, childType string) error {
	message := "Expected " + want.String() + ", found " + have.String()
	return derp.Validation(message, location, childType)
}

// copyChild returns a pointer to a copy of a child value, which is stored back in the parent
// afterwards.  A plain Go map or slice is written through the rosetta type with its layout.
func copyChild(value any) (childRef, error) {

	const location = "schema.copyChild"

	original := reflect.ValueOf(value)
	copied := reflect.New(original.Type())
	copied.Elem().Set(original)

	if isWritable(copied.Interface()) {
		return childRef{
			pointer: copied.Interface(),
			value:   func() any { return copied.Elem().Interface() },
		}, nil
	}

	for _, containerType := range convertibleContainers {

		if !original.Type().ConvertibleTo(containerType) {
			continue
		}

		converted := reflect.New(containerType)
		converted.Elem().Set(original.Convert(containerType))

		// Store the child back as its original type, so callers see the type they stored
		return childRef{
			pointer: converted.Interface(),
			value:   func() any { return converted.Elem().Convert(original.Type()).Interface() },
		}, nil
	}

	return childRef{}, derp.Internal(location, "Cannot write into this value", typeName(value))
}

// isIndex returns TRUE for a list index written as plain digits, with no sign or leading zero
func isIndex(segment string) bool {

	if (len(segment) > 1) && (segment[0] == '0') {
		return false
	}

	return isDigits(segment)
}

// createListItem creates a missing list item from its element, sets the value at tail within
// it, and stores it in the list
func createListItem(object any, childElement Element, head string, tail string, value any) error {

	const location = "schema.createListItem"

	setter, ok := object.(AnySetter)

	if !ok {
		return derp.Internal(location, "Cannot create a missing item", head, typeName(object))
	}

	// Unreachable today: newChild fails only for a typed parent's child, and this passes none
	child, err := newChild(childElement, nil, nextSegment(tail))

	if err != nil {
		return derp.Wrap(err, location, "Creating item", head)
	}

	if err := SetProperty(childElement, child.pointer, tail, value); err != nil {
		return derp.Wrap(err, location, "Setting value", head)
	}

	if !setter.SetAny(head, child.value()) {
		return derp.Internal(location, "Unable to store item", head, typeName(object))
	}

	return nil
}

// storeListItem stores a list item that is itself a list back into its parent list, because
// the parent's GetPointer returned a copy that may since have grown
func storeListItem(object any, subPointer any, head string) {

	if _, isList := object.(ArrayGetter); !isList || !isSlicePointer(subPointer) {
		return
	}

	if setter, ok := object.(AnySetter); ok {
		setter.SetAny(head, reflect.ValueOf(subPointer).Elem().Interface())
	}
}

// newTypedChild creates an empty value of the child type that a typed parent supplies
func newTypedChild(want childKind, childType reflect.Type) (childRef, error) {

	const location = "schema.newTypedChild"

	if have := kindOf(childType); (want != kindAny) && (want != have) {
		return childRef{}, kindError(location, want, have, childType.String())
	}

	// A nil pointer is replaced by a pointer to a new, empty value
	if childType.Kind() == reflect.Pointer {
		fresh := reflect.New(childType.Elem()).Interface()
		return childRef{pointer: fresh, value: func() any { return fresh }}, nil
	}

	return copyChild(reflect.New(childType).Elem().Interface())
}

// String returns the name of a kind, for error messages
func (kind childKind) String() string {

	switch kind {

	case kindMap:
		return "a map"

	case kindList:
		return "a list"

	case kindScalar:
		return "a single value"
	}

	return "any value"
}

// isWritable returns TRUE if schema can write into the value through one of its interfaces
func isWritable(pointer any) bool {

	switch pointer.(type) {

	case KeySetter, PointerGetter, AnySetter, StringSetter:
		return true
	}

	return false
}

// isSlicePointer returns TRUE if the value is a non-nil pointer to a slice
func isSlicePointer(value any) bool {

	reflected := reflect.ValueOf(value)

	if (reflected.Kind() != reflect.Pointer) || reflected.IsNil() {
		return false
	}

	return reflected.Elem().Kind() == reflect.Slice
}
