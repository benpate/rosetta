// Package tests holds cross-package integration tests that exercise schema
// against the mapof and sliceof carrier types together.
//
// These tests check that mapof and sliceof satisfy schema's path get/set
// contract. They do not need their own package to avoid an import cycle:
// schema imports mapof and sliceof so it can create their containers, and
// schema's own test files import both as well.
package tests
