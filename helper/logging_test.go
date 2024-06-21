package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestHelpers_GetLogger(test *testing.T) {
	x := GetLogger(1, map[int]string{}, 4)
	assert.NotEmpty(test, x)
}
