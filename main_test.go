package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
 * The unit tests in this file simulate command line invocation.
 */
func TestMain(test *testing.T) {
	_ = test
	main()
}

func TestCopyDatabase(test *testing.T) {
	_, err := copyDatabase()
	require.NoError(test, err)
}
