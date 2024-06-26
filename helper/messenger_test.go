package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestHelpers_GetMessenger(test *testing.T) {
	x := GetMessenger(1, map[int]string{}, 4)
	assert.NotEmpty(test, x)
}
