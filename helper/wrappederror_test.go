package helper_test

import (
	"errors"
	"testing"

	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func TestHelpers_WrappedError(test *testing.T) {
	var err error
	newError := helper.Errorf(err, "not an error: %w", err)
	assert.NoError(test, newError)
}

func TestHelpers_WrappedError_isError(test *testing.T) {
	err := errors.New("New error") //nolint:err113
	newError := helper.Errorf(err, "is an error: %w", err)
	assert.Error(test, newError)
}
