package helper_test

import (
	"testing"

	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestHelpers_GetEnv(test *testing.T) {
	expected := "EXPECTED_VALUE"
	test.Setenv("TEST_ENV_VAR", expected)

	actual := helper.GetEnv("TEST_ENV_VAR", "DEFAULT_VALUE")
	assert.Equal(test, expected, actual)
}

func TestHelpers_GetEnv_default(test *testing.T) {
	expected := "DEFAULT_VALUE"
	actual := helper.GetEnv("NO_ENV_VAR", expected)
	assert.Equal(test, expected, actual)
}
