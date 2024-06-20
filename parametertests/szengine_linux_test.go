//go:build linux

package szengine

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/sz-sdk-go/senzing"

	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func TestParameters_Szengine_AddRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getVerboseTestObject(ctx, test)
	flags := senzing.SzWithoutInfo
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	for _, record := range records {
		stdOut, _, err := captureStdoutReturningString(func() (string, error) {
			return szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		})
		require.NoError(test, err)
		inspectStdout(stdOut)
	}
	for _, record := range records {
		stdOut, _, err := captureStdoutReturningString(func() (string, error) {
			return szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		})
		require.NoError(test, err)
		inspectStdout(stdOut)
	}
}

// ----------------------------------------------------------------------------
// utility functions
// ----------------------------------------------------------------------------

func inspectStdout(stdout string) {
	_ = stdout
}
