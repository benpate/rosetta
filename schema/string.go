package schema

import (
	"reflect"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/null"
	"github.com/benpate/rosetta/schema/format"
)

// String represents a string data type within a JSON-Schema.
type String struct {
	Default   string
	MinLength null.Int
	MaxLength null.Int
	Enum      []string
	Pattern   string
	Format    string
	Required  bool
}

// Enumerate implements the "Enumerator" interface
func (element String) Enumerate() []string {
	return element.Enum
}

// Type returns the data type of this Element
func (element String) Type() reflect.Type {
	return reflect.TypeOf("")
}

// IsRequired returns TRUE if this element is a required field
func (element String) IsRequired() bool {
	return element.Required
}

// Find locates a child of this element
func (element String) Get(object reflect.Value, path string) (any, Element, error) {

	if path != "" {
		return nil, element, derp.NewInternalError("schema.String.Find", "Can't find sub-properties on a 'string' type", path)
	}

	return convert.StringDefault(object, element.Default), element, nil
}

// Set validates/formats a generic value using this schema
func (element String) Set(object reflect.Value, path string, value any) error {

	if path != "" {
		return derp.NewInternalError("schema.String.Set", "Can't set sub-properties on a string", path, value)
	}

	return setWithReflection(object, convert.StringDefault(value, element.Default))
}

// Validate compares a generic data value using this Schema
func (element String) Validate(value any) error {

	var errorReport error

	stringValue := convert.StringDefault(value, element.Default)

	// Validate against all formatting functions
	for _, format := range element.formatFunctions() {
		if _, err := format(stringValue); err != nil {
			errorReport = derp.Append(errorReport, err)
		}
	}

	// Verify required fields (after format functions are applied)
	if element.Required {
		if stringValue == "" {
			errorReport = derp.Append(errorReport, ValidationError{Message: "field is required"})
		}
	}

	// Validate minimum length
	if element.MinLength.IsPresent() {
		if len(stringValue) < element.MinLength.Int() {
			errorReport = derp.Append(errorReport, ValidationError{Message: "minimum length is " + element.MinLength.String()})
		}
	}

	// Validate maximum length
	if element.MaxLength.IsPresent() {
		if len(stringValue) > element.MaxLength.Int() {
			errorReport = derp.Append(errorReport, ValidationError{Message: "maximum length is " + element.MaxLength.String()})
		}
	}

	// Validate enumerated values
	if len(element.Enum) > 0 {
		if !compare.Contains(element.Enum, stringValue) {
			errorReport = derp.Append(errorReport, ValidationError{Message: "must match one of the required values."})
		}
	}

	return errorReport
}

// MarshalMap populates object data into a map[string]any
func (element String) MarshalMap() map[string]any {

	result := map[string]any{
		"type":     TypeString,
		"required": element.Required,
	}

	if element.Default != "" {
		result["default"] = element.Default
	}

	if element.MinLength.IsPresent() {
		result["minLength"] = element.MinLength.Int()
	}

	if element.MaxLength.IsPresent() {
		result["maxLength"] = element.MaxLength.Int()
	}

	if element.Pattern != "" {
		result["pattern"] = element.Pattern
	}

	if element.Format != "" {
		result["format"] = element.Format
	}

	if len(element.Enum) > 0 {
		result["enum"] = element.Enum
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *String) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "string" {
		return derp.New(500, "schema.String.UnmarshalMap", "Data is not type 'string'", data)
	}

	element.Default = convert.String(data["default"])
	element.MinLength = convert.NullInt(data["minLength"])
	element.MaxLength = convert.NullInt(data["maxLength"])
	element.Pattern = convert.String(data["pattern"])
	element.Format = convert.String(data["format"])
	element.Required = convert.Bool(data["required"])
	element.Enum = convert.SliceOfString(data["enum"])

	return err
}

// formatFunctions parses all of the formatting functions
// used by this string value.
func (element String) formatFunctions() []format.StringFormat {

	result := make([]format.StringFormat, 0)

	// Split multiple formats into individual function calls
	params := strings.Split(element.Format, " ")

	for _, arg := range params {

		name, value := list.Equal(arg).Split()

		if formatFunction, ok := formats[name]; ok {
			result = append(result, formatFunction(value.String()))
		}
	}

	// If there are no valid formats defined, then default to
	// no-html, which strictly removes all HTML tags from the value.
	if len(result) == 0 {
		result = []format.StringFormat{formats["no-html"]("")}
	}

	return result
}
