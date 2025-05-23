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
	"github.com/stretchr/testify/assert"
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
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzengine_AddRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		require.NoError(test, err)
		require.Empty(test, actual)
		printActual(test, actual)
	}

	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		require.NoError(test, err)
		require.Empty(test, actual)
		printActual(test, actual)
	}
}

func TestG2engine_AddRecord_badDataSourceCodeInJSON(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record2Json := `{"DATA_SOURCE": "BOB", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}` //nolint
	actual, err := szEngine.AddRecord(ctx, record1.DataSource, record1.ID, record1.JSON, flags)
	require.NoError(test, err)
	require.Empty(test, actual)

	_, err = szEngine.AddRecord(ctx, record2.DataSource, record2.ID, record2Json, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzengine_AddRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, badDataSourceCode, record.ID, record.JSON, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_AddRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, badRecordID, record.JSON, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_AddRecord_badRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, badRecordDefinition, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_AddRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, nilDataSourceCode, record.ID, record.JSON, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_AddRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, nilRecordID, record.JSON, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_AddRecord_nilRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, nilRecordDefinition, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_AddRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}

	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
		printActual(test, actual)
	}

	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
		printActual(test, actual)
	}
}

func TestSzengine_AddRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}

	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, badDataSourceCode, record.ID, record.JSON, flags)
		require.ErrorIs(test, err, szerror.ErrSzBadInput)
		printActual(test, actual)
	}

	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
		printActual(test, actual)
	}
}

func TestSzengine_CloseExport(test *testing.T) {
	// Tested in:
	//  - TestSzengine_ExportCsvEntityReport
	//  - TestSzengine_ExportJSONEntityReport
	_ = test
}

func TestSzengine_CountRedoRecords(test *testing.T) {
	ctx := test.Context()
	expected := int64(2)
	szEngine := getTestObject(test)
	actual, err := szEngine.CountRedoRecords(ctx)
	require.NoError(test, err)
	printActual(test, actual)
	assert.Equal(test, expected, actual)
}

func TestSzengine_DeleteRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, badRecordID, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, nilDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, nilRecordID, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1009"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	require.NotEmpty(test, actual)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1009"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.DeleteRecord(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_withInfo_badDataSourceCode_fix(test *testing.T) {
	_ = test
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	deleteRecords(ctx, records)
}

func TestSzengine_ExportCsvEntityReport(test *testing.T) {
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

	require.NoError(test, err)

	actualCount := 0

	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportCsvEntityReport_badCsvColumnList(test *testing.T) {
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

	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzengine_ExportCsvEntityReport_nilCsvColumnList(test *testing.T) {
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

	require.NoError(test, err)
}

func TestSzengine_ExportCsvEntityReportIterator(test *testing.T) {
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
		require.NoError(test, actual.Error)
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportCsvEntityReportIterator_badCsvColumnList(test *testing.T) {
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
		require.ErrorIs(test, actual.Error, szerror.ErrSzBadInput)
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportCsvEntityReportIterator_nilCsvColumnList(test *testing.T) {
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
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))

		actualCount++
	}

	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportJSONEntityReport(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)

	defer func() {
		panicOnErrorWithString(szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.ID, senzing.SzWithoutInfo))
	}()

	flags = senzing.SzExportDefaultFlags
	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)

	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.NoError(test, err)
	}()

	require.NoError(test, err)

	jsonEntityReport := ""

	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, exportHandle)
		require.NoError(test, err)

		if len(jsonEntityReportFragment) == 0 {
			break
		}

		jsonEntityReport += jsonEntityReportFragment
	}

	require.NoError(test, err)
	assert.NotEmpty(test, jsonEntityReport)
}

func TestSzengine_ExportJSONEntityReport_65536(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.AddRecord(ctx, aRecord.DataSource, aRecord.ID, aRecord.JSON, flags)
	require.NoError(test, err)
	require.NotEmpty(test, actual)
	printActual(test, actual)

	defer func() { _, _ = szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.ID, senzing.SzWithoutInfo) }()

	flags = getFlagsForEntityReport()
	// flags = int64(-1)
	aHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)

	defer func() {
		err := szEngine.CloseExport(ctx, aHandle)
		require.NoError(test, err)
	}()

	require.NoError(test, err)

	jsonEntityReport := ""

	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, aHandle)
		require.NoError(test, err)

		if len(jsonEntityReportFragment) == 0 {
			break
		}

		jsonEntityReport += jsonEntityReportFragment
	}

	require.NoError(test, err)
	assert.Greater(test, len(jsonEntityReport), 65536)
}

// IMPROVE: Implement TestSzengine_ExportJSONEntityReport_error
// func TestSzengine_ExportJSONEntityReport_error(test *testing.T) {}

func TestSzengine_ExportJSONEntityReportIterator(test *testing.T) {
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
		require.NoError(test, actual.Error)
		printActual(test, actual.Value)

		actualCount++
	}

	assert.Equal(test, expected, actualCount)
}

func TestSzengine_FetchNext(test *testing.T) {
	// Tested in:
	//  - TestSzengine_ExportJSONEntityReport
	_ = test
}

func TestSzengine_FetchNext_badExportHandle(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	actual, err := szEngine.FetchNext(ctx, badExportHandle)
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_FetchNext_nilExportHandle(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	actual, err := szEngine.FetchNext(ctx, nilExportHandle)
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByEntityID_badEntityID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByEntityID_nilEntityID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID_nilDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID_nilRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badEntityIDs(test *testing.T) {
	ctx := test.Context()
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badMaxDegrees(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badBuildOutDegrees(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badBuildOutMaxEntities(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_nilMaxDegrees(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_nilBuildOutDegrees(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_nilBuildOutMaxEntities(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID_nilMaxDegrees(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID_nilBuildOutDegrees(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID_nilBuildOutMaxEntities(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badStartEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	badStartEntityID := badEntityID
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badEndEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badAvoidEntityIDs(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badRequiredDataSource(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_nilMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_nilAvoidEntityIDs(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_nilRequiredDataSource(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_avoiding(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID, err := getEntityID(startRecord)
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_avoiding_badStartEntityID(test *testing.T) {
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
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_avoidingAndIncluding(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID, err := getEntityID(startRecord)
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_including(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID, err := getEntityID(startRecord)
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_including_badStartEntityID(test *testing.T) {
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
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_badMaxDegrees(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_badAvoidRecordKeys(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_badRequiredDataSources(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_avoiding(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_avoiding_badStartDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_avoidingAndIncluding(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_including(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_including_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetActiveConfigID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	actual, err := szEngine.GetActiveConfigID(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetEntityByEntityID_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, badEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetEntityByEntityID_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, nilEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID_nilDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID_nilRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetRecord(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetRecord_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetRecord_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetRecord_nilDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetRecord_nilRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	actual, err := szEngine.GetRedoRecord(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetStats(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	actual, err := szEngine.GetStats(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetVirtualEntityByRecordID(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetVirtualEntityByRecordID_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetVirtualEntityByRecordID_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_HowEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_HowEntityByEntityID_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, badEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_HowEntityByEntityID_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, nilEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_PreprocessRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	for _, record := range records {
		actual, err := szEngine.PreprocessRecord(ctx, record.JSON, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
}

func TestSzengine_PreprocessRecord_badRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.PreprocessRecord(ctx, badRecordDefinition, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_PrimeEngine(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	err := szEngine.PrimeEngine(ctx)
	require.NoError(test, err)
}

func TestSzengine_ProcessRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	redoRecord, err := szEngine.GetRedoRecord(ctx)
	require.NoError(test, err)

	if len(redoRecord) > 0 {
		flags := senzing.SzWithoutInfo
		actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
		require.NoError(test, err)
		require.Empty(test, actual)
		printActual(test, actual)
	}
}

func TestSzengine_ProcessRedoRecord_badRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, badRedoRecord, flags)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
	printActual(test, actual)
}

func TestSzengine_ProcessRedoRecord_nilRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, nilRedoRecord, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_ProcessRedoRecord_withInfo(test *testing.T) {
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
	require.NoError(test, err)

	if len(redoRecord) > 0 {
		flags := senzing.SzWithInfo
		actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
		require.NoError(test, err)
		require.NotEmpty(test, actual)
		printActual(test, actual)
	}
}

func TestSzengine_ProcessRedoRecord_withInfo_badRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, badRedoRecord, flags)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
	printActual(test, actual)
}

func TestSzengine_ProcessRedoRecord_withInfo_nilRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, nilRedoRecord, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, badEntityID, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, nilEntityID, flags)
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_withInfo_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, badEntityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_withInfo_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, nilEntityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord(test *testing.T) {
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
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_badRecordID(test *testing.T) {
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
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_nilDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_nilRecordID(test *testing.T) {
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
	require.NoError(test, err)
	require.Empty(test, actual)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_withInfo(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_withInfo_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_withInfo_nilDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_badAttributes(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_badSearchProfile(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_nilAttributes(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_nilSearchProfile(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_withSearchProfile(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_searchProfile(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, searchAttributes, searchProfile, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID1, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	entityID2, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID1, entityID2, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyEntities_badEnitity1(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID2, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, badEntityID, entityID2, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyEntities_badEnitity2(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID1, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID1, badEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyEntities_nilEnitity1(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID2, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, nilEntityID, entityID2, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyEntities_nilEnitity2(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID1, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID1, nilEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity_nilDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity_nilRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyRecords(test *testing.T) {
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyRecords_badDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_WhyRecords_badRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyRecords_nilDataSourceCode(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_WhyRecords_nilRecordID(test *testing.T) {
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
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhySearch(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)

	flags := senzing.SzNoFlags
	actual, err := szEngine.WhySearch(ctx, searchAttributes, entityID, searchProfile, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzengine_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzengine_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	szEngine.SetObserverOrigin(ctx, originMessage)
}

func TestSzengine_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	szEngine.SetObserverOrigin(ctx, originMessage)
	actual := szEngine.GetObserverOrigin(ctx)
	assert.Equal(test, originMessage, actual)
	printActual(test, actual)
}

func TestSzengine_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	err := szEngine.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzengine_AsInterface(test *testing.T) {
	expected := int64(4)
	ctx := test.Context()
	szEngine := getSzEngineAsInterface(ctx)
	actual, err := szEngine.CountRedoRecords(ctx)
	require.NoError(test, err)
	printActual(test, actual)
	assert.Equal(test, expected, actual)
}

func TestSzengine_Initialize(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	settings := getSettings()

	configID := senzing.SzInitializeWithDefaultConfiguration
	err := szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzengine_Initialize_error
// func TestSzengine_Initialize_error(test *testing.T) {}

func TestSzengine_Initialize_withConfigID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	settings := getSettings()

	configID := getDefaultConfigID()
	err := szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzengine_Initialize_withConfigID_error
// func TestSzengine_Initialize_withConfigID_error(test *testing.T) {}

func TestSzengine_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	configID, err := szEngine.GetActiveConfigID(ctx)
	require.NoError(test, err)
	err = szEngine.Reinitialize(ctx, configID)
	require.NoError(test, err)
	printActual(test, configID)
}

// IMPROVE: Implement TestSzengine_Reinitialize_badConfigID
// func TestSzengine_Reinitialize_badConfigID(test *testing.T) {}

func TestSzengine_Destroy(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	err := szEngine.Destroy(ctx)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzengine_Destroy_error
// func TestSzengine_Destroy_error(test *testing.T) {}

func TestSzengine_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szEngineSingleton = nil
	szEngine := getTestObject(test)
	err := szEngine.Destroy(ctx)
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

func getEntityID(record record.Record) (int64, error) {
	return getEntityIDForRecord(record.DataSource, record.ID)
}

func getEntityIDForRecord(datasource string, recordID string) (int64, error) {
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

	return result, nil
}

func getEntityIDString(record record.Record) string {
	entityID, err := getEntityID(record)
	panicOnError(err)

	result := strconv.FormatInt(entityID, baseTen)

	return result
}

func getEntityIDStringForRecord(datasource string, recordID string) string {
	var result string

	entityID, err := getEntityIDForRecord(datasource, recordID)
	panicOnError(err)

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
