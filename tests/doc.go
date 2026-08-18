// Package tests holds cross-package integration tests that exercise schema
// against the mapof and sliceof carrier types together.
//
// It exists to break an import cycle, not because these tests are special.
// Schema defines the Getter and Setter interfaces; mapof and sliceof implement
// them, so schema cannot import either from its own test files. A third package
// can import all three and verify that the implementations really satisfy
// schema's path get/set contract.
package tests
