# translate 🗺️

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/translate)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/rosetta/go.yml?branch=main&style=flat-square)](https://github.com/benpate/rosetta/actions/workflows/go.yml)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/rosetta/translate.svg?style=flat-square)](https://codecov.io/gh/benpate/rosetta/translate)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/rosetta?style=flat-square)](https://goreportcard.com/report/github.com/benpate/rosetta)
[![Version](https://img.shields.io/github/v/release/benpate/rosetta?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/rosetta/releases)

## Object TranslationUtilities for Go

This library maps data from one variable into another.  requires a `schema` for both input and output values, which must be compatable with schema Getter and Setter interfaces.  The best way to use this library is with a `mapof.Any` which already includes these interfaces.


## Mapping Utilities


### Path

`{"path":"original.path", "target":"target.path"}`

### Value

`{"value":"FIXED VALUE HERE", "target":"target.path"}`

### Expression

`{"path":"{{go template expression}}", "target":"target.path"}`

### Conditionals

`{"if":"{{go template expression}}", "then":[/* additional rules */], "else":[/* additional rules */]}`

### ForEach

`{"path":"original.path", "target":"target.path", "filter":"{{go template expression}}", "rules":[/* additional rules*/]}`

## Pull Requests Welcome

This library is a work in progress, and will benefit from your experience reports, use cases, and contributions.  If you have an idea for making this library better, send in a pull request.  We're all in this together! 🗺️
