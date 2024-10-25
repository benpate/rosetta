# Rosetta 🌹

<img alt="Plate 22 in volume I: five women in blue robes carrying white cloths in the centre, men with turbans rowing and sailing boats behind the women, many sailing boats moored at a bank in the left background; after Mayer. 1801-1803 Hand-coloured aquatint with etching © The Trustees of the British Museum" src="https://github.com/benpate/rosetta/raw/main/meta/city-of-rosetta.webp" style="width:100%; display:block; margin-bottom:20px;">


[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta)
[![Version](https://img.shields.io/github/v/release/benpate/rosetta?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/rosetta/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/rosetta/go.yml?branch=main&style=flat-square)](https://github.com/benpate/rosetta/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/rosetta?style=flat-square)](https://goreportcard.com/report/github.com/benpate/rosetta)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/rosetta.svg?style=flat-square)](https://codecov.io/gh/benpate/rosetta)

## A Collection of Translation and Data Manipulation Tools

Rosetta combines several different data manipulation tools into a single module.  While each was useful on its own, the dependencies between them made updates hellish, so this collection was born.  

Many of these packages pre-date Go generics, and many are being refactored to take advantage of this new capability.

### What's Included

* [Channel](https://github.com/benpate/rosetta/tree/main/channel) tools for manipulating Go channels
* [Compare](https://github.com/benpate/rosetta/tree/main/compare) values of unknown data types.
* [Convert](https://github.com/benpate/rosetta/tree/main/convert) between arbitrary data types with sensible, configurable defaults.
* [HTML](https://github.com/benpate/rosetta/tree/main/html) conversion tools
* [Iterator](https://github.com/benpate/rosetta/tree/main/iterator) tools for iterating over data structures
* [Maps](https://github.com/benpate/rosetta/tree/main/maps) tools for working with map types
* [Schema](https://github.com/benpate/rosetta/tree/main/schema) validation based on JSON Schema
* [Slice](https://github.com/benpate/rosetta/tree/main/slice) manipulation library
* [Translate](https://github.com/benpate/rosetta/tree/main/translate) data mapping library

### Enhanced Data Types

* [List](https://github.com/benpate/rosetta/tree/main/list) parsing library
* [MapOf](https://github.com/benpate/rosetta/tree/main/mapof) data type with type safe getters/setters
* [Null](https://github.com/benpate/rosetta/tree/main/null)-able data types
* [SliceOf](https://github.com/benpate/rosetta/tree/main/sliceof) data types


## Image Used With Permission

**© The Trustees of the British Museum.** Shared under a [Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0) licence](http://creativecommons.org/licenses/by-nc-sa/4.0/).

## Pull Requests Welcome

While many parts of this module have been used for years in production environments, it is still a work in progress and will benefit from your experience reports, use cases, and contributions.  If you have an idea for making Rosetta better, send in a pull request.  We're all in this together! 🌹

