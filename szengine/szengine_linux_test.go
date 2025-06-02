//go:build linux

package szengine_test

import (
	"strings"
	"testing"

	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/stretchr/testify/require"
)

var expectedExportCsvEntityReport = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`3,0,"","","CUSTOMERS","1001"`,
	`3,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`3,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}

var expectedExportCsvEntityReportIterator = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`20,0,"","","CUSTOMERS","1001"`,
	`20,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`20,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}

var expectedExportCsvEntityReportIteratorNilCsvColumnList = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`26,0,"","","CUSTOMERS","1001"`,
	`26,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`26,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}

func TestSzEngine_ExportCsvEntityReportIterator(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	expected := expectedExportCsvEntityReportIterator
	szEngine := getTestObject(ctx, test)
	csvColumnList := ""
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0

	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		printDebug(test, actual.Error, actual.Value)
		require.NoError(test, actual.Error)
		require.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	require.Equal(test, len(expected), actualCount)
}

func TestSzEngine_ExportCsvEntityReportIterator_nilCsvColumnList(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	expected := expectedExportCsvEntityReportIteratorNilCsvColumnList
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0

	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, nilCsvColumnList, flags) {
		require.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	require.Equal(test, len(expected), actualCount)
}
