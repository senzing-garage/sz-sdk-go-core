package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(test *testing.T) {
	_ = test

	main()
}

func TestCopyDatabase(test *testing.T) {
	_, err := copyDatabase()

	require.NoError(test, err)
}
