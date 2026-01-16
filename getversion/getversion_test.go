package getversion_test

import (
	"context"
	"testing"

	"github.com/senzing-garage/sz-sdk-go-core/getversion"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestHelpers_GetMessenger(test *testing.T) {
	ctx := context.Background()
	x := getversion.GetSenzingVersion(ctx)
	assert.NotEmpty(test, x)
}
