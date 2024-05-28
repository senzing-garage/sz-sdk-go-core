package main

import (
	"testing"
)

/*
 * The unit tests in this file simulate command line invocation.
 */
func TestMain(testing *testing.T) {
	_ = testing
	main()
}

func TestCopyDatabase(testing *testing.T) {
	_ = testing
	copyDatabase()
}
