# null 🚫

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/null)
[![Build Status](https://img.shields.io/github/workflow/status/benpate/null/Go/master)](https://github.com/benpate/rosetta/null/actions/workflows/go.yml)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/null.svg?style=flat-square)](https://codecov.io/gh/benpate/null)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/rosetta/null?style=flat-square)](https://goreportcard.com/report/github.com/benpate/rosetta/null)
[![Version](https://img.shields.io/github/v/release/benpate/null?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/rosetta/null/releases)

## Simple library for null values in Go

This library provides simple, (mostly) idiomatic primitives for nullable values in Go.  It supports Int, Bool, and Float types.

```go
// "b" is null, and ready to use
var b null.Bool

// Or, create a nullable value that's not null
b := null.NewBool(true)

// Set value to false
b.Set(false)

// Set value to true
b.Set(true)

// Get the value
b.Bool()

// Make the value null again
b.Unset()
```

## Pull Requests Welcome

Please use GitHub to make suggestions, pull requests, and enhancements.  We're all in this together! 🚫
