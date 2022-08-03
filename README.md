# Rosetta 💐

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta)
[![Build Status](https://img.shields.io/github/workflow/status/benpate/rosetta/Go/main)](https://github.com/benpate/rosetta/actions/workflows/go.yml)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/rosetta.svg?style=flat-square)](https://codecov.io/gh/benpate/rosetta)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/rosetta?style=flat-square)](https://goreportcard.com/report/github.com/benpate/rosetta)
[![Version](https://img.shields.io/github/v/release/benpate/rosetta?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/rosetta/releases)

## A Collection of Data Mapping Tools

Rosetta combines several different data mapping tools into a single module.  While each was useful on its own, the dependencies between them made updates hellish, so this collection was born.  

All of these packages pre-date Go generics, and many are being refactored to take advantage of this new capability.

### What's Included

* [Compare](compare) values of unknown data types.
* [Convert](convert) between arbitrary data types with sensible, configurable defaults.
* [HTML](html) conversion tools
* [Nullable](null) values with strong type
* [Path getter/setter](path) for generic, complex data structures
* [Schema](schema) validation based on JSON Schema

### Enhanced Data Types

* [List](list) parsing library
* [Map](maps) data type with type safe getters/setters
* [Null](null)-able data types
* [Slice](slice) manipulation library

## Pull Requests Welcome

While many parts of this module have been used for years in production environments, it is still a work in progress and will benefit from your experience reports, use cases, and contributions.  If you have an idea for making Rosetta better, send in a pull request.  We're all in this together! 👍
