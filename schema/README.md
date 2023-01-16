# schema 👍

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/schema)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/rosetta/schema.svg?style=flat-square)](https://codecov.io/gh/benpate/rosetta/schema)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/rosetta?style=flat-square)](https://goreportcard.com/report/github.com/benpate/rosetta)
[![Version](https://img.shields.io/github/v/release/benpate/rosetta?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/rosetta/releases)

## Simplified JSON-Schema

This is a simplified, minimal implementation of JSON-Schema that is fast and has an easy API.  If you're looking for a complete, rigorous implementation of JSON-Schema, you should try another tool.

## What it Does

This library implements a sub-set of the [JSON-Schema specification](http://json-schema.org)

### What's Included

* Unmarshal schema from JSON
* Array, Boolean, Integer, Number, Object, and String type validators.
* Custom Format rules
* Happy API for accessing schema information, and walking a schema tree with a [JSON-Pointer](https://tools.ietf.org/html/rfc6901)

Schemas can also get/set data from a compliant data structure.  This is helpful when dealing with dynamic data (like form inputs to a dynamic CMS) that may not be known at compile time.

### What's Left Out

* References
* Loading remote Schemas by URI.

## How It Works

Schema now requires data structures to implement Getter/Setter interfaces to expose the data inside of them.  While this is a larger code requirement on client libraries, it means that schema has a type-safe way of getting to dynamic data.  There are several examples in the test cases, and the `mapof` and `sliceof` libraries in Rosetta also implement these interfaces fully.  

Here's a quick example:

```go
type MyStruct struct {
	Name string
	Email string
	Age int
}

func (m MyStruct) GetStringOK(path string) (string, bool) {
	switch path {
	case "name":
		return m.Name, true
	case "email":
		return m.Email, true
	default:
		return "", false
	}
}

func (m *MyStruct) SetStringOK(path string, value string) bool {
	switch path {
	case "name":
		m.Name = value
		return true

	case "email":
		m.Email = value
		return true
	default:
		return false
	}
}

```

Additional interfaces enable schemas to traverse nested structures and arrays.

## Pull Requests Welcome

This library is a work in progress, and will benefit from your experience reports, use cases, and contributions.  If you have an idea for making Rosetta better, send in a pull request.  We're all in this together! 👍
 
