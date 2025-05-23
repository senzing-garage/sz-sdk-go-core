package szengine_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/testfixtures"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/require"
)

const (
	avoidEntityIDs      = senzing.SzNoAvoidance
	avoidRecordKeys     = senzing.SzNoAvoidance
	baseTen             = 10
	buildOutDegrees     = int64(2)
	buildOutMaxEntities = int64(10)
	defaultTruncation   = 76
	instanceName        = "SzEngine Test"
	jsonIndentation     = "    "
	maxDegrees          = int64(2)
	observerOrigin      = "SzEngine observer"
	originMessage       = "Machine: nn; Task: UnitTest"
	printErrors         = false
	printResults        = false
	requiredDataSources = senzing.SzNoRequiredDatasources
	searchAttributes    = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile       = senzing.SzNoSearchProfile
	verboseLogging      = senzing.SzNoLogging
)

// Bad parameters

const (
	badAttributes          = "}{"
	badAvoidEntityIDs      = "}{"
	badAvoidRecordKeys     = "}{"
	badBuildOutDegrees     = int64(-1)
	badBuildOutMaxEntities = int64(-1)
	badCsvColumnList       = "BAD, CSV, COLUMN, LIST"
	badDataSourceCode      = "BadDataSourceCode"
	badEntityID            = int64(-1)
	badExportHandle        = uintptr(0)
	badLogLevelName        = "BadLogLevelName"
	badMaxDegrees          = int64(-1)
	badRecordDefinition    = "}{"
	badRecordID            = "BadRecordID"
	badRedoRecord          = "{}"
	badRequiredDataSources = "}{"
	badSearchProfile       = "}{"
)

// Nil/empty parameters

var (
	nilAttributes          string
	nilAvoidEntityIDs      string
	nilBuildOutDegrees     int64
	nilBuildOutMaxEntities int64
	nilCsvColumnList       string
	nilDataSourceCode      string
	nilEntityID            int64
	nilExportHandle        uintptr
	nilMaxDegrees          int64
	nilRecordDefinition    string
	nilRecordID            string
	nilRedoRecord          string
	nilRequiredDataSources string
	nilSearchProfile       string
)

var (
	defaultConfigID   int64
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szEngineSingleton *szengine.Szengine
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityID int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func TestSzEngine_AddRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.Empty(test, actual)
	}

	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.Empty(test, actual)
	}
}

func TestSzEngine_AddRecord_badDataSourceCodeInJSON(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	record := truthset.CustomerRecords["1002"]
	badRecordJSON := `{"DATA_SOURCE": "BOB", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, badRecordJSON, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'CUSTOMERS' and 'BOB'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, badDataSourceCode, record.ID, record.JSON, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, badRecordID, record.JSON, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0024|Conflicting RECORD_ID values 'BadRecordID' and '1001'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_badRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, badRecordDefinition, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, nilDataSourceCode, record.ID, record.JSON, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_AddRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, nilRecordID, record.JSON, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_AddRecord_nilRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, nilRecordDefinition, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}

	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
	}

	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
	}
}

func TestSzEngine_AddRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}

	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, badDataSourceCode, record.ID, record.JSON, flags)
		printDebug(test, err, actual)
		require.ErrorIs(test, err, szerror.ErrSzBadInput)

		expectedErr := `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}`
		require.JSONEq(test, expectedErr, err.Error())
	}

	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
	}
}

func TestSzEngine_CloseExport(test *testing.T) {
	// Tested in:
	//  - TestSzEngine_ExportCsvEntityReport
	//  - TestSzEngine_ExportJSONEntityReport
	_ = test
}

func TestSzEngine_CountRedoRecords(test *testing.T) {
	ctx := test.Context()
	expected := int64(2)
	szEngine := getTestObject(test)
	actual, err := szEngine.CountRedoRecords(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Equal(test, expected, actual)
}

func TestSzEngine_DeleteRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_DeleteRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_DeleteRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, badRecordID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_DeleteRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, nilDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_DeleteRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, nilRecordID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_DeleteRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1009"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.NotEmpty(test, actual)
}

func TestSzEngine_DeleteRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1009"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.DeleteRecord(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044005","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_DeleteRecord_withInfo_badDataSourceCode_fix(test *testing.T) {
	_ = test
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	deleteRecords(ctx, records)
}

func TestSzEngine_ExportCsvEntityReport(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	expected := expectedExportCsvEntityReport
	szEngine := getTestObject(test)
	csvColumnList := ""
	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)

	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.NoError(test, err)
	}()

	printDebug(test, err)
	require.NoError(test, err)

	actualCount := 0

	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		require.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	require.Equal(test, len(expected), actualCount)
}

func TestSzEngine_ExportCsvEntityReport_badCsvColumnList(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, badCsvColumnList, flags)

	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.ErrorIs(test, err, szerror.ErrSz)
	}()

	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).ExportCsvEntityReport","error":{"id":"SZSDK60044007","reason":"SENZ3131|Invalid column [BAD] requested for CSV export."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ExportCsvEntityReport_nilCsvColumnList(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, nilCsvColumnList, flags)

	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.NoError(test, err)
	}()

	printDebug(test, err)
	require.NoError(test, err)
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
	szEngine := getTestObject(test)
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

func TestSzEngine_ExportCsvEntityReportIterator_badCsvColumnList(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	expected := []string{
		``,
	}
	szEngine := getTestObject(test)
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0

	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, badCsvColumnList, flags) {
		printDebug(test, actual.Error, actual.Value)
		require.ErrorIs(test, actual.Error, szerror.ErrSzBadInput)
		require.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		expectedErr := `{"function":"szengine.(*Szengine).ExportCsvEntityReport","error":{"id":"SZSDK60044007","reason":"SENZ3131|Invalid column [BAD] requested for CSV export."}}`
		require.JSONEq(test, expectedErr, actual.Error.Error())

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
	szEngine := getTestObject(test)
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0

	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, nilCsvColumnList, flags) {
		require.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	require.Equal(test, len(expected), actualCount)
}

func TestSzEngine_ExportJSONEntityReport(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.AddRecord(ctx, aRecord.DataSource, aRecord.ID, aRecord.JSON, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)

	defer func() {
		panicOnErrorWithString(szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.ID, senzing.SzWithoutInfo))
	}()

	flags = senzing.SzExportDefaultFlags
	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)

	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.NoError(test, err)
	}()

	printDebug(test, err, actual)
	require.NoError(test, err)

	jsonEntityReport := ""

	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, exportHandle)
		printDebug(test, err, actual)
		require.NoError(test, err)

		if len(jsonEntityReportFragment) == 0 {
			break
		}

		jsonEntityReport += jsonEntityReportFragment
	}

	printDebug(test, err, actual)
	require.NoError(test, err)
	require.NotEmpty(test, jsonEntityReport)
}

func TestSzEngine_ExportJSONEntityReport_65536(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.AddRecord(ctx, aRecord.DataSource, aRecord.ID, aRecord.JSON, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.NotEmpty(test, actual)

	defer func() { _, _ = szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.ID, senzing.SzWithoutInfo) }()

	flags = getFlagsForEntityReport()
	// flags = int64(-1)
	aHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)

	defer func() {
		err := szEngine.CloseExport(ctx, aHandle)
		require.NoError(test, err)
	}()

	printDebug(test, err, actual)
	require.NoError(test, err)

	jsonEntityReport := ""

	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, aHandle)
		printDebug(test, err, actual)
		require.NoError(test, err)

		if len(jsonEntityReportFragment) == 0 {
			break
		}

		jsonEntityReport += jsonEntityReportFragment
	}

	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Greater(test, len(jsonEntityReport), 65536)
}

// IMPROVE: Implement TestSzEngine_ExportJSONEntityReport_error
// func TestSzEngine_ExportJSONEntityReport_error(test *testing.T) {}

func TestSzEngine_ExportJSONEntityReportIterator(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	expected := 1
	szEngine := getTestObject(test)
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0

	for actual := range szEngine.ExportJSONEntityReportIterator(ctx, flags) {
		printDebug(test, actual.Error, actual.Value)
		require.NoError(test, actual.Error)

		actualCount++
	}

	require.Equal(test, expected, actualCount)
}

func TestSzEngine_FetchNext(test *testing.T) {
	// Tested in:
	//  - TestSzEngine_ExportJSONEntityReport
	_ = test
}

func TestSzEngine_FetchNext_badExportHandle(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	actual, err := szEngine.FetchNext(ctx, badExportHandle)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FetchNext","error":{"id":"SZSDK60044009","reason":"SENZ3103|Invalid Export Handle [0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FetchNext_nilExportHandle(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	actual, err := szEngine.FetchNext(ctx, nilExportHandle)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FetchNext","error":{"id":"SZSDK60044009","reason":"SENZ3103|Invalid Export Handle [0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindInterestingEntitiesByEntityID_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByEntityID(ctx, badEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByEntityID_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByEntityID(ctx, nilEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '0'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindInterestingEntitiesByRecordID_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, record.DataSource, badRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, nilDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID_nilRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, record.DataSource, nilRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		entityIDs,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByEntityID_badEntityIDs(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	badEntityID1 := 0
	badEntityID2 := 1
	badEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + strconv.Itoa(
		badEntityID1,
	) + `}, {"ENTITY_ID": ` + strconv.Itoa(
		badEntityID2,
	) + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		badEntityIDs,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0037|Unknown resolved entity value '0'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_badMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		entityIDs,
		badMaxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0031|Invalid value of max degree '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_badBuildOutDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		entityIDs,
		maxDegrees,
		badBuildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0032|Invalid value of build out degree '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_badBuildOutMaxEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		entityIDs,
		maxDegrees,
		buildOutDegrees,
		badBuildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0029|Invalid value of max entities '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_nilMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		entityIDs,
		nilMaxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByEntityID_nilBuildOutDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		entityIDs,
		maxDegrees,
		nilBuildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByEntityID_nilBuildOutMaxEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(
		ctx,
		entityIDs,
		maxDegrees,
		buildOutDegrees,
		nilBuildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(
		ctx,
		recordKeys,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByRecordID_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	badRecordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		badDataSourceCode +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(
		ctx,
		badRecordKeys,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByRecordID_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	badRecordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		badRecordID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(
		ctx,
		badRecordKeys,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByRecordID_nilMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(
		ctx,
		recordKeys,
		nilMaxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByRecordID_nilBuildOutDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(
		ctx,
		recordKeys,
		maxDegrees,
		nilBuildOutDegrees,
		buildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByRecordID_nilBuildOutMaxEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(
		ctx,
		recordKeys,
		maxDegrees,
		buildOutDegrees,
		nilBuildOutMaxEntities,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID_badStartEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	badStartEntityID := badEntityID
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		badStartEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByEntityID_badEndEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])

	badEndEntityID := badEntityID
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		badEndEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByEntityID_badMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		badMaxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0031|Invalid value of max degree '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByEntityID_badAvoidEntityIDs(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		badAvoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByEntityID_badRequiredDataSource(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		badRequiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044025","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByEntityID_nilMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		nilMaxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID_nilAvoidEntityIDs(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		nilAvoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID_nilRequiredDataSource(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		nilRequiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID_avoiding(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID := getEntityID(startRecord)
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	startEntityIDString := getEntityIDStringForRecord("CUSTOMERS", "1001")

	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + startEntityIDString + `}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID_avoiding_badStartEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID := badEntityID
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	startRecordEntityIDString := getEntityIDString(startRecord)

	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + startRecordEntityIDString + `}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByEntityID_avoidingAndIncluding(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID := getEntityID(startRecord)
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	startRecordEntityIDString := getEntityIDString(startRecord)

	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + startRecordEntityIDString + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID_including(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID := getEntityID(startRecord)
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByEntityID_including_badStartEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	badStartEntityID := badEntityID
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])

	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(
		ctx,
		badStartEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044025","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByRecordID_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		badDataSourceCode,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByRecordID_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		badRecordID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByRecordID_badMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		badMaxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0031|Invalid value of max degree '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByRecordID_badAvoidRecordKeys(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		badAvoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByRecordID_badRequiredDataSources(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		badRequiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044027","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByRecordID_avoiding(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	avoidRecordKeys := `{"RECORDS": [{ "DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByRecordID_avoiding_badStartDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	avoidRecordKeys := `{"RECORDS": [{ "DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		badDataSourceCode,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindPathByRecordID_avoidingAndIncluding(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	avoidRecordKeys := `{"RECORDS": [{ "DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByRecordID_including(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record1EntityID := getEntityIDString(record1)

	avoidRecordKeys := `{"ENTITIES": [{"ENTITY_ID": ` + record1EntityID + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindPathByRecordID_including_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record1EntityID := getEntityIDString(record1)

	avoidRecordKeys := `{"ENTITIES": [{"ENTITY_ID": ` + record1EntityID + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(
		ctx,
		badDataSourceCode,
		record1.ID,
		record2.DataSource,
		record2.ID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044027","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetActiveConfigID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	actual, err := szEngine.GetActiveConfigID(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, entityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByEntityID_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, badEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).GetEntityByEntityID","error":{"id":"SZSDK60044030","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetEntityByEntityID_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, nilEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).GetEntityByEntityID","error":{"id":"SZSDK60044030","reason":"SENZ0037|Unknown resolved entity value '0'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByRecordID_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetEntityByRecordID_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, record.DataSource, badRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetEntityByRecordID_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, nilDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ2207|Data source code [] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetEntityByRecordID_nilRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, record.DataSource, nilRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetRecord(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, record.DataSource, badRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, nilDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ2207|Data source code [] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, record.DataSource, nilRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	actual, err := szEngine.GetRedoRecord(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetStats(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	actual, err := szEngine.GetStats(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetVirtualEntityByRecordID_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` +
		badDataSourceCode +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).GetVirtualEntityByRecordID","error":{"id":"SZSDK60044038","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_GetVirtualEntityByRecordID_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		badRecordID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).GetVirtualEntityByRecordID","error":{"id":"SZSDK60044038","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_HowEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, entityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_HowEntityByEntityID_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, badEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).HowEntityByEntityID","error":{"id":"SZSDK60044040","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_HowEntityByEntityID_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, nilEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).HowEntityByEntityID","error":{"id":"SZSDK60044040","reason":"SENZ0037|Unknown resolved entity value '0'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_PreprocessRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	for _, record := range records {
		actual, err := szEngine.PreprocessRecord(ctx, record.JSON, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
	}
}

func TestSzEngine_PreprocessRecord_badRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.PreprocessRecord(ctx, badRecordDefinition, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).PreprocessRecord","error":{"id":"SZSDK60044061","reason":"SENZ0002|Invalid Message"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_PrimeEngine(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	err := szEngine.PrimeEngine(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzEngine_ProcessRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	redoRecord, err := szEngine.GetRedoRecord(ctx)
	printDebug(test, err, redoRecord)
	require.NoError(test, err)

	if len(redoRecord) > 0 {
		flags := senzing.SzWithoutInfo
		actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.Empty(test, actual)
	}
}

func TestSzEngine_ProcessRedoRecord_badRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, badRedoRecord, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ProcessRedoRecord_nilRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, nilRedoRecord, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ0007|Empty Message"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ProcessRedoRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
		truthset.CustomerRecords["1005"],
		truthset.CustomerRecords["1009"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	redoRecord, err := szEngine.GetRedoRecord(ctx)
	printDebug(test, err, redoRecord)
	require.NoError(test, err)

	if len(redoRecord) > 0 {
		flags := senzing.SzWithInfo
		actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
	}
}

func TestSzEngine_ProcessRedoRecord_withInfo_badRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, badRedoRecord, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ProcessRedoRecord_withInfo_nilRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, nilRedoRecord, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ0007|Empty Message"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ReevaluateEntity(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_ReevaluateEntity_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, badEntityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_ReevaluateEntity_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, nilEntityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ReevaluateEntity_withInfo_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, badEntityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ReevaluateEntity_withInfo_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, nilEntityID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ReevaluateRecord(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_ReevaluateRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044048","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ReevaluateRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, badRecordID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_ReevaluateRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, nilDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044048","reason":"SENZ2207|Data source code [] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ReevaluateRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, nilRecordID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual)
}

func TestSzEngine_ReevaluateRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ReevaluateRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateRecord(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044049","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ReevaluateRecord_withInfo_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateRecord(ctx, nilDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044049","reason":"SENZ2207|Data source code [] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_SearchByAttributes(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, searchAttributes, searchProfile, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_SearchByAttributes_badAttributes(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, badAttributes, searchProfile, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0027|Invalid value for search-attributes"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_SearchByAttributes_badSearchProfile(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, searchAttributes, badSearchProfile, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0088|Unknown search profile value '}{'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_SearchByAttributes_nilAttributes(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, nilAttributes, searchProfile, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0027|Invalid value for search-attributes"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_SearchByAttributes_nilSearchProfile(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, searchAttributes, nilSearchProfile, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_SearchByAttributes_withSearchProfile(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	// searchProfile := "SEARCH"
	searchProfile := "INGEST"
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, searchAttributes, searchProfile, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_SearchByAttributes_searchProfile(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, searchAttributes, searchProfile, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_WhyEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID1 := getEntityID(truthset.CustomerRecords["1001"])
	entityID2 := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID1, entityID2, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_WhyEntities_badEnitity1(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, badEntityID, entityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyEntities_badEnitity2(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID, badEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyEntities_nilEnitity1(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1002"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, nilEntityID, entityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '0'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyEntities_nilEnitity2(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID, nilEntityID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '0'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecordInEntity(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, record.DataSource, record.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_WhyRecordInEntity_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, badDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecordInEntity_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, record.DataSource, badRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecordInEntity_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, nilDataSourceCode, record.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ2207|Data source code [] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecordInEntity_nilRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, record.DataSource, nilRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecords(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, record1.DataSource, record1.ID, record2.DataSource, record2.ID, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_WhyRecords_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, badDataSourceCode, record1.ID, record2.DataSource, record2.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecords_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, record1.DataSource, record1.ID, record2.DataSource, badRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecords_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, nilDataSourceCode, record1.ID, record2.DataSource, record2.ID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [] does not exist."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhyRecords_nilRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, record1.DataSource, record1.ID, record2.DataSource, nilRecordID, flags)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_WhySearch(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhySearch(ctx, searchAttributes, entityID, searchProfile, flags)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzEngine_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzEngine_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	szEngine.SetObserverOrigin(ctx, originMessage)
}

func TestSzEngine_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	szEngine.SetObserverOrigin(ctx, originMessage)
	actual := szEngine.GetObserverOrigin(ctx)
	require.Equal(test, originMessage, actual)
	printDebug(test, nil, actual)
}

func TestSzEngine_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	err := szEngine.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzEngine_AsInterface(test *testing.T) {
	expected := int64(4)
	ctx := test.Context()
	szEngine := getSzEngineAsInterface(ctx)
	actual, err := szEngine.CountRedoRecords(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Equal(test, expected, actual)
}

func TestSzEngine_Initialize(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	settings := getSettings()

	configID := senzing.SzInitializeWithDefaultConfiguration
	err := szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzEngine_Initialize_error
// func TestSzEngine_Initialize_error(test *testing.T) {}

func TestSzEngine_Initialize_withConfigID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	settings := getSettings()

	configID := getDefaultConfigID()
	err := szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzEngine_Initialize_withConfigID_error
// func TestSzEngine_Initialize_withConfigID_error(test *testing.T) {}

func TestSzEngine_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	configID, err := szEngine.GetActiveConfigID(ctx)
	printDebug(test, err, configID)
	require.NoError(test, err)
	err = szEngine.Reinitialize(ctx, configID)
	printDebug(test, err)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzEngine_Reinitialize_badConfigID
// func TestSzEngine_Reinitialize_badConfigID(test *testing.T) {}

func TestSzEngine_Destroy(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	err := szEngine.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzEngine_Destroy_error
// func TestSzEngine_Destroy_error(test *testing.T) {}

func TestSzEngine_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szEngineSingleton = nil
	szEngine := getTestObject(test)
	err := szEngine.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		panicOnError(err)
	}
}

func deleteRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		panicOnError(err)
	}
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getDefaultConfigID() int64 {
	return defaultConfigID
}

func getEntityID(record record.Record) int64 {
	return getEntityIDForRecord(record.DataSource, record.ID)
}

func getEntityIDForRecord(datasource string, recordID string) int64 {
	var (
		err    error
		result int64
	)

	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	response, err := szEngine.GetEntityByRecordID(ctx, datasource, recordID, senzing.SzWithoutInfo)
	panicOnError(err)

	getEntityByRecordIDResponse := &GetEntityByRecordIDResponse{} //exhaustruct:ignore
	err = json.Unmarshal([]byte(response), &getEntityByRecordIDResponse)
	panicOnError(err)

	result = getEntityByRecordIDResponse.ResolvedEntity.EntityID

	return result
}

func getEntityIDString(record record.Record) string {
	entityID := getEntityID(record)

	result := strconv.FormatInt(entityID, baseTen)

	return result
}

func getEntityIDStringForRecord(datasource string, recordID string) string {
	var result string

	entityID := getEntityIDForRecord(datasource, recordID)

	result = strconv.FormatInt(entityID, baseTen)

	return result
}

func getFlagsForEntityReport() int64 {
	return senzing.Flags(
		senzing.SzEntityIncludeAllFeatures,
		senzing.SzEntityIncludeDisclosedRelations,
		senzing.SzEntityIncludeEntityName,
		senzing.SzEntityIncludeFeatureStats,
		senzing.SzEntityIncludeInternalFeatures,
		senzing.SzEntityIncludeNameOnlyRelations,
		senzing.SzEntityIncludePossiblyRelatedRelations,
		senzing.SzEntityIncludePossiblySameRelations,
		senzing.SzEntityIncludeRecordData,
		senzing.SzEntityIncludeRecordFeatureDetails,
		senzing.SzEntityIncludeRecordFeatures,
		senzing.SzEntityIncludeRecordFeatureStats,
		senzing.SzEntityIncludeRecordJSONData,
		senzing.SzEntityIncludeRecordMatchingInfo,
		senzing.SzEntityIncludeRecordSummary,
		senzing.SzEntityIncludeRecordTypes,
		senzing.SzEntityIncludeRecordUnmappedData,
		senzing.SzEntityIncludeRelatedEntityName,
		senzing.SzEntityIncludeRelatedMatchingInfo,
		senzing.SzEntityIncludeRelatedRecordData,
		senzing.SzEntityIncludeRelatedRecordSummary,
		senzing.SzEntityIncludeRelatedRecordTypes,
		senzing.SzEntityIncludeRepresentativeFeatures,
		senzing.SzExportIncludeDisclosed,
		senzing.SzExportIncludeDisclosed,
		senzing.SzExportIncludeMultiRecordEntities,
		senzing.SzExportIncludeMultiRecordEntities,
		senzing.SzExportIncludeNameOnly,
		senzing.SzExportIncludeNameOnly,
		senzing.SzExportIncludePossiblyRelated,
		senzing.SzExportIncludePossiblyRelated,
		senzing.SzExportIncludePossiblySame,
		senzing.SzExportIncludePossiblySame,
		senzing.SzExportIncludeSingleRecordEntities,
		senzing.SzExportIncludeSingleRecordEntities,
	)
}

func getNewSzEngine(ctx context.Context) *szengine.Szengine {
	var (
		err         error
		newSzEngine *szengine.Szengine
	)

	settings := getSettings()

	newSzEngine = &szengine.Szengine{}
	err = newSzEngine.SetLogLevel(ctx, logLevel)
	panicOnError(err)

	if logLevel == "TRACE" {
		newSzEngine.SetObserverOrigin(ctx, observerOrigin)
		err = newSzEngine.RegisterObserver(ctx, observerSingleton)
		panicOnError(err)
		err = newSzEngine.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
		panicOnError(err)
	}

	err = newSzEngine.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
	panicOnError(err)

	return newSzEngine
}

func getSettings() string {
	var result string

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	panicOnError(err)

	databaseURL := "sqlite3://na:na@nowhere/" + dbTargetPath

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	panicOnError(err)

	return result
}

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var result senzing.SzAbstractFactory

	_ = ctx

	settings := getSettings()

	result = &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	return result
}

func getSzEngine(ctx context.Context) *szengine.Szengine {
	if szEngineSingleton == nil {
		szEngineSingleton = getNewSzEngine(ctx)
	}

	return szEngineSingleton
}

func getSzEngineAsInterface(ctx context.Context) senzing.SzEngine {
	return getSzEngine(ctx)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szengine")
}

func getTestObject(t *testing.T) *szengine.Szengine {
	t.Helper()

	return getSzEngine(t.Context())
}

func handleError(err error) {
	if err != nil {
		outputln("Error:", err)
	}
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicOnErrorWithString(aString string, err error) {
	_ = aString

	panicOnError(err)
}

func printActual(t *testing.T, actual interface{}) {
	t.Helper()
	printResult(t, "Actual", actual)
}

func printDebug(t *testing.T, err error, items ...any) {
	t.Helper()
	printError(t, err)

	for item := range items {
		printActual(t, item)
	}
}

func printError(t *testing.T, err error) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}
}

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

// func ramCheck(test *testing.T, iteration int) {
// 	sysInfo := &syscall.Sysinfo_t{}
// 	printer := message.NewPrinter(language.English)
// 	err := syscall.Sysinfo(sysInfo)
// 	require.NoError(test, err)
// 	usedRAM := sysInfo.Totalram - sysInfo.Freeram
// 	printer.Printf(">>> iteration: %d,  Used memory: %d\n", iteration, usedRAM)
// }

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func setup() {
	setupDirectories()
	setupDatabase()

	err := setupSenzingConfiguration()
	panicOnError(err)
}

func setupDatabase() {
	testDirectoryPath := getTestDirectoryPath()
	_, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	panicOnError(err)
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	panicOnError(err)

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	panicOnError(err)
}

func setupDirectories() {
	testDirectoryPath := getTestDirectoryPath()
	err := os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	panicOnError(err)
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0o750) // recreate the test target directory
	panicOnError(err)
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	settings := getSettings()

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)

	defer func() { panicOnError(szConfig.Destroy(ctx)) }()

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)

	defer func() { panicOnError(szConfigManager.Destroy(ctx)) }()

	// Create a Senzing configuration.

	err = szConfig.ImportTemplate(ctx)
	panicOnError(err)

	// Add data sources to template Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		panicOnError(err)
	}

	// Create a string representation of the Senzing configuration.

	configDefinition, err := szConfig.Export(ctx)
	panicOnError(err)

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szengine_test at %s", now.UTC())
	defaultConfigID, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	return nil
}

func teardown() {
	ctx := context.TODO()
	_ = getSzEngine(ctx)
	teardownSzEngine(ctx)
}

func teardownSzEngine(ctx context.Context) {
	err := szEngineSingleton.UnregisterObserver(ctx, observerSingleton)
	panicOnError(err)
	err = szEngineSingleton.Destroy(ctx)
	panicOnError(err)

	szEngineSingleton = nil
}
