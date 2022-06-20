# path 🏞

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/rosetta/path)
[![Build Status](https://img.shields.io/github/workflow/status/benpate/rosetta/path/Go/main)](https://github.com/benpate/rosetta/path/actions/workflows/go.yml)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/rosetta/path.svg?style=flat-square)](https://codecov.io/gh/benpate/rosetta/path)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/rosetta/path?style=flat-square)](https://goreportcard.com/report/github.com/benpate/rosetta/path)
[![Version](https://img.shields.io/github/v/release/benpate/rosetta/path?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/rosetta/path/releases)

## Resolve data from arbitrary structures

This is an experimental library for reading/writing values into arbitrary data structures, specifically the `map[string]interface{}` and `[]interface{}` values returned by Go's `json.Unmarshal()` functions.  It is inspired by the JSON-path standard, but has a very simplified syntax -- using a series of strings separated by dots.

## Example Code

```go

s := map[string]interface{}{
    "name":  "John Connor",
    "email": "john@connor.mil",
    "relatives": map[string]interface{}{
        "mom": "Sarah Connor",
        "dad": "Kyle Reese",
    },
    "enemies": []interface{}{"T-1000", "T-3000", "T-5000"},
}

name, err := path.Get(s, "name") // John Connor
email, err := path.Get(s, "email") // john@connor.mil
sarah, err := path.Get(s, "relatives.0") // t-1000

```

## Pull Requests Welcome

Please use GitHub to make suggestions, pull requests, and enhancements.  We're all in this together! 🏞
