package helper_test

import (
	"testing"

	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestHelpers_GetLogger(test *testing.T) {
	x := helper.GetLogger(1, map[int]string{}, 4)
	assert.NotEmpty(test, x)
}
