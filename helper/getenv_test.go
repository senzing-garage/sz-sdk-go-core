package helper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestHelpers_GetEnv(test *testing.T) {
	expected := "EXPECTED_VALUE"
	os.Setenv("TEST_ENV_VAR", expected)
	actual := GetEnv("TEST_ENV_VAR", "DEFAULT_VALUE")
	assert.Equal(test, expected, actual)
}

func TestHelpers_GetEnv_default(test *testing.T) {
	expected := "DEFAULT_VALUE"
	actual := GetEnv("NO_ENV_VAR", expected)
	assert.Equal(test, expected, actual)
}
