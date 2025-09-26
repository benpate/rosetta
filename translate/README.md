# translate 🗺️

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/translate)
[![Version](https://img.shields.io/github/v/release/benpate/rosetta?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/rosetta/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/rosetta/go.yml?branch=main&style=flat-square)](https://github.com/benpate/rosetta/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/rosetta?style=flat-square)](https://app.codecov.io/gh/benpate/rosetta)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/rosetta/translate.svg?style=flat-square)](https://codecov.io/gh/benpate/rosetta/translate)

## Object TranslationUtilities for Go

This library maps data from one variable into another.  requires a `schema` for both input and output values, which must be compatable with schema Getter and Setter interfaces.  The best way to use this library is with a `mapof.Any` which already includes these interfaces.


## Mapping Utilities


### Value
`value` sets a static value in the target object.  This does not use any data from the source

`{"value":"FIXED VALUE HERE", "target":"target.path"}`

### Path
`path` maps a value from one path in the source object into a new path in the target object.  This is a direct copy, and no value is changed.

`{"path":"original.path", "target":"target.path"}`

### Expression
`expression` executes a go template expression, using the result to set a value in the target object. This expression can pull from any data in the source object.

`{"expression":"{{go template expression}}", "target":"target.path"}`

### Conditionals
`if` executes a go template expression (using data from the source object) to determine which additional rules to apply.  If the template returns "true",  the `then` rules are executed.  Otherwise, the `else` rules are executed.

`{"if":"{{go template expression}}", "then":[/* additional rules */], "else":[/* additional rules */]}`

### ForEach
`forEach` loops over a map or array, executing a list of rules for each item in the source, prefixed by the target path.

`{"forEach":"original.path", "target":"target.path", "filter":"{{go template expression}}", "rules":[/* additional rules */]}`

### First
`first` executes a list of rules, stopping after the first one sets a non-zero value in the target object

`{"first":"target.path", "rules":[/* additional rules */]}`

## Pull Requests Welcome

This library is a work in progress, and will benefit from your experience reports, use cases, and contributions.  If you have an idea for making this library better, send in a pull request.  We're all in this together! 🗺️
