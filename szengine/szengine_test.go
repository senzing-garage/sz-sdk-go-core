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
	avoidEntityIDs             = senzing.SzNoAvoidance
	avoidRecordKeys            = senzing.SzNoAvoidance
	baseTen                    = 10
	buildOutDegrees            = int64(2)
	buildOutMaxEntities        = int64(10)
	defaultAttributes          = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	defaultAvoidEntityIDs      = senzing.SzNoAvoidance
	defaultAvoidRecordKeys     = senzing.SzNoAvoidance
	defaultBuildOutDegrees     = int64(2)
	defaultBuildOutMaxEntities = int64(10)
	defaultMaxDegrees          = int64(2)
	defaultSearchAttributes    = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	defaultSearchProfile       = senzing.SzNoSearchProfile
	defaultTruncation          = 76
	defaultVerboseLogging      = senzing.SzNoLogging
	instanceName               = "SzEngine Test"
	jsonIndentation            = "    "
	maxDegrees                 = int64(2)
	observerID                 = "Observer 1"
	observerOrigin             = "SzEngine observer"
	originMessage              = "Machine: nn; Task: UnitTest"
	printErrors                = false
	printResults               = false
	requiredDataSources        = senzing.SzNoRequiredDatasources
	searchAttributes           = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile              = senzing.SzNoSearchProfile
	verboseLogging             = senzing.SzNoLogging
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
	nilSemaphoreString     = "xyzzy"
	nilSemaphoreInt64      = int64(-9999)
)

// Nil/empty parameters

var (
	nilAttributes          = nilSemaphoreString
	nilBuildOutDegrees     = nilSemaphoreInt64
	nilBuildOutMaxEntities = nilSemaphoreInt64
	nilCsvColumnList       string
	nilDataSourceCode      = nilSemaphoreString
	nilEntityID            = nilSemaphoreInt64
	nilExportHandle        uintptr
	nilMaxDegrees          = nilSemaphoreInt64
	nilRecordDefinition    = nilSemaphoreString
	nilRecordID            = nilSemaphoreString
	nilRedoRecord          = nilSemaphoreString
	nilSearchProfile       = nilSemaphoreString
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
	testCases := getTestCasesForAddRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Defaults.
			szEngine := getTestObject(ctx, test)
			record := truthset.CustomerRecords["1001"]

			// Test.

			dataSourceCode := xString(testCase.dataSourceCode, record.DataSource)
			recordID := xString(testCase.recordID, record.ID)

			actual, err := szEngine.AddRecord(ctx,
				dataSourceCode,
				recordID,
				xString(testCase.recordDefinition, record.JSON),
				xInt64(testCase.flags, senzing.SzNoFlags))

			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
				actual, err := szEngine.DeleteRecord(ctx,
					dataSourceCode,
					recordID,
					senzing.SzNoFlags)
				printDebug(test, err, actual)
			}
		})
	}
}

func TestSzEngine_CloseExportReport(test *testing.T) {
	// Tested in:
	//  - TestSzEngine_ExportCsvEntityReport
	//  - TestSzEngine_ExportJSONEntityReport
	_ = test
}

func TestSzEngine_CountRedoRecords(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.CountRedoRecords(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Equal(test, expectedRedoRecordCount, actual)
}

func TestSzEngine_DeleteRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForDeleteRecord()
	szEngine := getTestObject(ctx, test)

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1005 := truthset.CustomerRecords["1005"]

			records := []record.Record{
				record1005,
			}

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Test.

			dataSourceCode := xString(testCase.dataSourceCode, record1005.DataSource)
			recordID := xString(testCase.recordID, record1005.ID)

			actual, err := szEngine.DeleteRecord(ctx,
				dataSourceCode,
				recordID,
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_ExportCsvEntityReport(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	expected := expectedExportCsvEntityReport
	csvColumnList := ""
	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)

	defer func() {
		err := szEngine.CloseExportReport(ctx, exportHandle)
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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, badCsvColumnList, flags)

	defer func() {
		err := szEngine.CloseExportReport(ctx, exportHandle)
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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, nilCsvColumnList, flags)

	defer func() {
		err := szEngine.CloseExportReport(ctx, exportHandle)
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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	expected := expectedExportCsvEntityReportIterator
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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	expected := []string{
		``,
	}
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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	expected := expectedExportCsvEntityReportIteratorNilCsvColumnList
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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

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
		err := szEngine.CloseExportReport(ctx, exportHandle)
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
	szEngine := getTestObject(ctx, test)
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
		err := szEngine.CloseExportReport(ctx, aHandle)
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

func TestSzEngine_ExportJSONEntityReportIterator(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	expected := 1
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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

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
	szEngine := getTestObject(ctx, test)

	defer func() { deleteRecords(ctx, szEngine, records) }()

	addRecords(ctx, szEngine, records)

	actual, err := szEngine.FetchNext(ctx, nilExportHandle)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengine.(*Szengine).FetchNext","error":{"id":"SZSDK60044009","reason":"SENZ3103|Invalid Export Handle [0]"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindInterestingEntitiesByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			entityID := getEntityID(ctx, szEngine, record1001)

			// Test.

			actual, err := szEngine.FindInterestingEntitiesByEntityID(ctx,
				xInt64(testCase.entityID, entityID),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindInterestingEntitiesByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Test.

			actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx,
				xString(testCase.dataSourceCode, record1001.DataSource),
				xString(testCase.recordID, record1001.ID),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_FindNetworkByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindNetworkByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]
			record1002 := truthset.CustomerRecords["1002"]

			records := []record.Record{
				record1001,
				record1002,
				truthset.CustomerRecords["1003"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			entityIDs := entityIDsJSON(
				getEntityID(ctx, szEngine, record1001),
				getEntityID(ctx, szEngine, record1002))

			// Test.

			actual, err := szEngine.FindNetworkByEntityID(
				ctx,
				entityIDs,
				xInt64(testCase.maxDegrees, defaultMaxDegrees),
				xInt64(testCase.buildOutDegrees, defaultBuildOutDegrees),
				xInt64(testCase.buildOutMaxEntities, defaultBuildOutMaxEntities),
				xInt64(testCase.flags, senzing.SzFindNetworkDefaultFlags),
			)
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_FindNetworkByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindNetworkByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			recordKeys := recordKeysFunc()
			if testCase.recordKeys != nil {
				recordKeys = testCase.recordKeys()
			}

			// Test.

			actual, err := szEngine.FindNetworkByRecordID(ctx,
				recordKeys,
				xInt64(testCase.maxDegrees, defaultMaxDegrees),
				xInt64(testCase.buildOutDegrees, defaultBuildOutDegrees),
				xInt64(testCase.buildOutMaxEntities, defaultBuildOutMaxEntities),
				xInt64(testCase.flags, senzing.SzFindNetworkDefaultFlags),
			)
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_FindPathByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindPathByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			startEntityID := getEntityID(ctx, szEngine, truthset.CustomerRecords["1001"])
			endEntityID := getEntityID(ctx, szEngine, truthset.CustomerRecords["1002"])

			avoidEntityIDs := senzing.SzNoAvoidance
			if testCase.avoidEntityIDs != nil {
				avoidEntityIDs = testCase.avoidEntityIDs()
			}

			requiredDataSources := senzing.SzNoRequiredDatasources
			if testCase.requiredDataSources != nil {
				requiredDataSources = testCase.requiredDataSources()
			}

			// Test.

			actual, err := szEngine.FindPathByEntityID(
				ctx,
				xInt64(testCase.startEntityID, startEntityID),
				xInt64(testCase.endEntityID, endEntityID),
				xInt64(testCase.maxDegrees),
				avoidEntityIDs,
				requiredDataSources,
				xInt64(testCase.flags),
			)
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_FindPathByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindPathByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			record1 := truthset.CustomerRecords["1001"]
			startDataSourceCode := record1.DataSource
			startRecordID := record1.ID

			record2 := truthset.CustomerRecords["1002"]
			endDataSourceCode := record2.DataSource
			endRecordID := record2.ID

			avoidRecordKeys := senzing.SzNoAvoidance
			if testCase.avoidRecordKeys != nil {
				avoidRecordKeys = testCase.avoidRecordKeys()
			}

			requiredDataSources := senzing.SzNoRequiredDatasources
			if testCase.requiredDataSources != nil {
				requiredDataSources = testCase.requiredDataSources()
			}

			// Test.

			actual, err := szEngine.FindPathByRecordID(
				ctx,
				xString(testCase.startDataSourceCode, startDataSourceCode),
				xString(testCase.startRecordID, startRecordID),
				xString(testCase.endDataSourceCode, endDataSourceCode),
				xString(testCase.endRecordID, endRecordID),
				xInt64(testCase.maxDegrees, defaultMaxDegrees),
				avoidRecordKeys,
				requiredDataSources,
				xInt64(testCase.flags, senzing.SzFindNetworkDefaultFlags),
			)
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_GetActiveConfigID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetActiveConfigID(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetEntityByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			entityID := getEntityID(ctx, szEngine, truthset.CustomerRecords["1001"])

			// Test.

			actual, err := szEngine.GetEntityByEntityID(ctx,
				xInt64(testCase.entityID, entityID),
				xInt64(testCase.flags, senzing.SzEntityDefaultFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_GetEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetEntityByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			record := truthset.CustomerRecords["1001"]

			// Test.

			actual, err := szEngine.GetEntityByRecordID(ctx,
				xString(testCase.dataSourceCode, record.DataSource),
				xString(testCase.recordID, record.ID),
				xInt64(testCase.flags, senzing.SzEntityDefaultFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_GetRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			record := truthset.CustomerRecords["1001"]

			// Test.

			actual, err := szEngine.GetRecord(ctx,
				xString(testCase.dataSourceCode, record.DataSource),
				xString(testCase.recordID, record.ID),
				xInt64(testCase.flags, senzing.SzRecordDefaultFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_GetRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetRedoRecord(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetStats(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetStats(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetVirtualEntityByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			recordKeys := recordKeysFunc()
			if testCase.recordKeys != nil {
				recordKeys = testCase.recordKeys()
			}

			// Test.

			actual, err := szEngine.GetVirtualEntityByRecordID(ctx,
				recordKeys,
				xInt64(testCase.flags, senzing.SzRecordDefaultFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_HowEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForHowEntityByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
				truthset.CustomerRecords["1002"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			entityID := getEntityID(ctx, szEngine, record1001)

			// Test.

			actual, err := szEngine.HowEntityByEntityID(ctx,
				xInt64(testCase.entityID, entityID),
				xInt64(testCase.flags, senzing.SzRecordDefaultFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_GetRecordPreview(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetRecordPreview()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Test.

			actual, err := szEngine.GetRecordPreview(ctx,
				xString(testCase.recordDefinition, record1001.JSON),
				xInt64(testCase.flags, senzing.SzRecordDefaultFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_PrimeEngine(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	err := szEngine.PrimeEngine(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzEngine_ProcessRedoRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForProcessRedoRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
				truthset.CustomerRecords["1004"],
				truthset.CustomerRecords["1005"],
				truthset.CustomerRecords["1009"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			redoRecord, err := szEngine.GetRedoRecord(ctx)
			printDebug(test, err, redoRecord)
			require.NoError(test, err)

			// Test.

			actual, err := szEngine.ProcessRedoRecord(ctx,
				redoRecord,
				xInt64(testCase.flags, senzing.SzRecordDefaultFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_ReevaluateEntity(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForReevaluateEntity()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			entityID := getEntityID(ctx, szEngine, record1001)

			// Test.

			actual, err := szEngine.ReevaluateEntity(ctx,
				xInt64(testCase.entityID, entityID),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_ReevaluateRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForReevaluateRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Test.

			actual, err := szEngine.ReevaluateRecord(ctx,
				xString(testCase.dataSourceCode, record1001.DataSource),
				xString(testCase.recordID, record1001.ID),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_SearchByAttributes(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForSearchByAttributes()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Test.

			actual, err := szEngine.SearchByAttributes(ctx,
				xString(testCase.attributes, defaultAttributes),
				xString(testCase.searchProfile, defaultSearchProfile),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_WhyEntities(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForWhyEntities()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]
			record1002 := truthset.CustomerRecords["1002"]

			records := []record.Record{
				record1001,
				record1002,
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			entityID1 := getEntityID(ctx, szEngine, record1001)
			entityID2 := getEntityID(ctx, szEngine, record1002)

			// Test.

			actual, err := szEngine.WhyEntities(ctx,
				xInt64(testCase.entityID1, entityID1),
				xInt64(testCase.entityID2, entityID2),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_WhyRecordInEntity(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForWhyRecordInEntity()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Test.

			actual, err := szEngine.WhyRecordInEntity(ctx,
				xString(testCase.dataSourceCode, record1001.DataSource),
				xString(testCase.recordID, record1001.ID),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_WhyRecords(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForWhyRecords()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]
			record1002 := truthset.CustomerRecords["1002"]

			records := []record.Record{
				record1001,
				record1002,
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Test.

			actual, err := szEngine.WhyRecords(
				ctx,
				xString(testCase.dataSourceCode1, record1001.DataSource),
				xString(testCase.recordID1, record1001.ID),
				xString(testCase.dataSourceCode2, record1002.DataSource),
				xString(testCase.recordID2, record1002.ID),
				xInt64(testCase.flags, senzing.SzNoFlags),
			)
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_WhySearch(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForWhySearch()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}
			szEngine := getTestObject(ctx, test)

			defer func() { deleteRecords(ctx, szEngine, records) }()

			addRecords(ctx, szEngine, records)

			// Defaults.

			entityID := getEntityID(ctx, szEngine, truthset.CustomerRecords["1001"])

			// Test.

			actual, err := szEngine.WhySearch(ctx,
				xString(testCase.attributes, defaultAttributes),
				xInt64(testCase.entityID, entityID),
				xString(testCase.searchProfile, defaultSearchProfile),
				xInt64(testCase.flags, senzing.SzNoFlags))
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzEngine_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzEngine_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	szEngine.SetObserverOrigin(ctx, originMessage)
}

func TestSzEngine_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	szEngine.SetObserverOrigin(ctx, originMessage)
	actual := szEngine.GetObserverOrigin(ctx)
	require.Equal(test, originMessage, actual)
	printDebug(test, nil, actual)
}

func TestSzEngine_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
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

// func TestSzEngine_Initialize(test *testing.T) {
// 	ctx := test.Context()
// 	szEngine := getTestObject(ctx, test)
// 	settings := getSettings()

// 	configID := senzing.SzInitializeWithDefaultConfiguration
// 	err := szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
// 	printDebug(test, err)
// 	require.NoError(test, err)
// }

// func TestSzEngine_Initialize_withConfigID(test *testing.T) {
// 	ctx := test.Context()
// 	szEngine := getTestObject(ctx, test)
// 	settings := getSettings()

// 	configID := getDefaultConfigID()
// 	err := szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
// 	printDebug(test, err)
// 	require.NoError(test, err)
// }

func TestSzEngine_Initialize_withConfigID_error(test *testing.T) {
	// IMPROVE: Implement TestSzEngine_Initialize_withConfigID_error
	_ = test
}

func TestSzEngine_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	configID, err := szEngine.GetActiveConfigID(ctx)
	printDebug(test, err, configID)
	require.NoError(test, err)
	err = szEngine.Reinitialize(ctx, configID)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzEngine_Destroy(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	err := szEngine.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
	szEngineSingleton = nil // Reset szEngineSingleton
}

func TestSzEngine_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	err := szEngine.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
	szEngineSingleton = nil // Reset szEngineSingleton
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, szEngine senzing.SzEngine, records []record.Record) {
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		panicOnError(err)
	}
}

func createSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
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

func deleteRecords(ctx context.Context, szEngine senzing.SzEngine, records []record.Record) {
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

func getEntityID(ctx context.Context, szEngine senzing.SzEngine, record record.Record) int64 {
	return getEntityIDForRecord(ctx, szEngine, record.DataSource, record.ID)
}

func getEntityIDForRecord(
	ctx context.Context,
	szEngine senzing.SzEngine,
	datasource string,
	recordID string,
) int64 {
	var (
		err    error
		result int64
	)

	response, err := szEngine.GetEntityByRecordID(ctx, datasource, recordID, senzing.SzWithoutInfo)
	panicOnError(err)

	getEntityByRecordIDResponse := &GetEntityByRecordIDResponse{} //exhaustruct:ignore
	err = json.Unmarshal([]byte(response), &getEntityByRecordIDResponse)
	panicOnError(err)

	result = getEntityByRecordIDResponse.ResolvedEntity.EntityID

	return result
}

func getEntityIDStringForRecord(
	ctx context.Context,
	szEnzine senzing.SzEngine,
	datasource string,
	recordID string,
) string { //nolint
	var result string

	entityID := getEntityIDForRecord(ctx, szEnzine, datasource, recordID)

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

func getTestObject(ctx context.Context, t *testing.T) *szengine.Szengine {
	t.Helper()

	return getSzEngine(ctx)
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

func printDebug(t *testing.T, err error, items ...any) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s\n", err.Error())
		}
	}

	if printResults {
		for _, item := range items {
			outLine := truncator.Truncate(fmt.Sprintf("%v", item), defaultTruncation, "...", truncator.PositionEnd)
			t.Logf("Result: %s\n", outLine)
		}
	}
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
		_, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
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
	teardownSzEngine(ctx)
}

func teardownSzEngine(ctx context.Context) {
	err := szEngineSingleton.UnregisterObserver(ctx, observerSingleton)
	panicOnError(err)

	_ = szEngineSingleton.Destroy(ctx)

	szEngineSingleton = nil
}

// ----------------------------------------------------------------------------
// Test structs
// ----------------------------------------------------------------------------

type TestMetadataForAddRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordDefinition   string
	recordID           string
}

type TestMetadataForDeleteRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForFindInterestingEntitiesByEntityID struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}
type TestMetadataForFindInterestingEntitiesByRecordID struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForFindNetworkByEntityID struct {
	buildOutDegrees     int64
	buildOutMaxEntities int64
	entityIDs           func() string
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
}

type TestMetadataForFindNetworkByRecordID struct {
	buildOutDegrees     int64
	buildOutMaxEntities int64
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
	recordKeys          func() string
}

type TestMetadataForFindPathByEntityID struct {
	avoidEntityIDs      func() string
	endEntityID         int64
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
	requiredDataSources func() string
	startEntityID       int64
}

type TestMetadataForFindPathByRecordID struct {
	avoidRecordKeys     func() string
	endDataSourceCode   string
	endRecordID         string
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
	requiredDataSources func() string
	startDataSourceCode string
	startRecordID       string
}

type TestMetadataForGetEntityByEntityID struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForGetEntityByRecordID struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForGetRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForGetVirtualEntityByRecordID struct {
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordKeys         func() string
}

type TestMetadataForHowEntityByEntityID struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForGetRecordPreview struct {
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordDefinition   string
}

type TestMetadataForProcessRedoRecord struct {
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	redoRecord         string
}

type TestMetadataForReevaluateEntity struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForReevaluateRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForSearchByAttributes struct {
	attributes         string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	searchProfile      string
}

type TestMetadataForWhyEntities struct {
	entityID1          int64
	entityID2          int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForWhyRecordInEntity struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForWhyRecords struct {
	dataSourceCode1    string
	dataSourceCode2    string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID1          string
	recordID2          string
}

type TestMetadataForWhySearch struct {
	attributes         string
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	searchProfile      string
}

// ----------------------------------------------------------------------------
// Test data
// ----------------------------------------------------------------------------

func getTestCasesForAddRecord() []TestMetadataForAddRecord {
	record1002 := truthset.CustomerRecords["1002"]
	result := []TestMetadataForAddRecord{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}`,
		},
		{
			name:               "badDataSourceCode_asErrSz",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}`,
		},
		{
			name:               "badDataSourceCodeInJSON",
			dataSourceCode:     record1002.DataSource,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'CUSTOMERS' and 'BOB'"}}`,
			recordDefinition:   `{"DATA_SOURCE": "BOB", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`,
			recordID:           record1002.ID,
		},
		{
			name:               "badRecordDefinition",
			recordDefinition:   badRecordDefinition,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`,
		},
		{
			name:               "badRecordID",
			recordID:           badRecordID,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0024|Conflicting RECORD_ID values 'BadRecordID' and '1001'"}}`,
		},
		{
			name: "default",
		},
		{
			name:           "nilDataSourceCode",
			dataSourceCode: nilDataSourceCode,
		},
		{
			name:               "nilRecordDefinition",
			recordDefinition:   nilRecordDefinition,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=0]"}}`,
		},
		{
			name:     "nilRecordID",
			recordID: nilRecordID,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:               "withInfo_badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			flags:              senzing.SzWithInfo,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}`,
		},
		{
			name:               "withInfo_badRecordDefinition",
			flags:              senzing.SzWithInfo,
			recordDefinition:   badRecordDefinition,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`,
		},
		{
			name:               "withInfo_badRecordID",
			flags:              senzing.SzWithInfo,
			recordID:           badRecordID,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ0024|Conflicting RECORD_ID values 'BadRecordID' and '1001'"}}`,
		},
		{
			name:           "withInfo_nilDataSourceCode",
			dataSourceCode: nilDataSourceCode,
			flags:          senzing.SzWithInfo,
		},
		{
			name:               "withInfo_nilRecordDefinition",
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=0]"}}`,
			flags:              senzing.SzWithInfo,
			recordDefinition:   nilRecordDefinition,
		},
		{
			name:     "withInfo_nilRecordID",
			flags:    senzing.SzWithInfo,
			recordID: nilRecordID,
		},
	}

	return result
}

func getTestCasesForDeleteRecord() []TestMetadataForDeleteRecord {
	result := []TestMetadataForDeleteRecord{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badDataSourceCode_asSzBadInputError",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badDataSourceCode_asSzErr",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:     "badRecordID",
			recordID: badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzConfiguration,
			expectedErrMessage: `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}`,
		},
		{
			name:     "nilRecordID",
			recordID: nilRecordID,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:               "withInfo_badDataSourceCode",
			flags:              senzing.SzWithInfo,
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044005","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:     "withInfo_badRecordID",
			flags:    senzing.SzWithInfo,
			recordID: badRecordID,
		},
		{
			name:               "withInfo_nilDataSourceCode",
			flags:              senzing.SzWithInfo,
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzConfiguration,
			expectedErrMessage: `{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044005","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}`,
		},
		{
			name:     "withInfo_nilRecordID",
			flags:    senzing.SzWithInfo,
			recordID: nilRecordID,
		},
	}

	return result
}

func getTestCasesForFindInterestingEntitiesByEntityID() []TestMetadataForFindInterestingEntitiesByEntityID {
	result := []TestMetadataForFindInterestingEntitiesByEntityID{
		{
			name:               "badEntityID",
			entityID:           badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEntityID",
			entityID:           nilEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '0'"}}`,
		},
	}

	return result
}

func getTestCasesForFindInterestingEntitiesByRecordID() []TestMetadataForFindInterestingEntitiesByRecordID {
	result := []TestMetadataForFindInterestingEntitiesByRecordID{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badRecordID",
			recordID:           badRecordID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [] does not exist."}}`,
		},
		{
			name:               "nilRecordID",
			recordID:           nilRecordID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`,
		},
	}

	return result
}

func getTestCasesForFindNetworkByEntityID() []TestMetadataForFindNetworkByEntityID {
	result := []TestMetadataForFindNetworkByEntityID{
		{
			name:               "badBuildOutDegrees",
			buildOutDegrees:    badBuildOutDegrees,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0032|Invalid value of build out degree '-1'"}}`,
		},
		{
			name:                "badBuildOutMaxEntities",
			buildOutMaxEntities: badBuildOutMaxEntities,
			expectedErr:         szerror.ErrSz,
			expectedErrMessage:  `{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0029|Invalid value of max entities '-1'"}}`,
		},
		{
			name:      "badEntityIDs",
			entityIDs: badEntityIDsFunc,
			// IMPROVE: Shouldn't this error?
		},
		{
			name:               "badMaxDegrees",
			maxDegrees:         badMaxDegrees,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0031|Invalid value of max degree '-1'"}}`,
		},
		{
			name: "default",
		},
		{
			name:            "nilBuildOutDegrees",
			buildOutDegrees: nilBuildOutDegrees,
		},
		{
			name:                "nilBuildOutMaxEntities",
			buildOutMaxEntities: nilBuildOutMaxEntities,
		},
		{
			name:       "nilMaxDegrees",
			maxDegrees: nilMaxDegrees,
		},
	}

	return result
}

func getTestCasesForFindNetworkByRecordID() []TestMetadataForFindNetworkByRecordID {
	result := []TestMetadataForFindNetworkByRecordID{
		{
			name:               "badBuildOutDegrees",
			buildOutDegrees:    badBuildOutDegrees,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0032|Invalid value of build out degree '-1'"}}`,
		},
		{
			name:               "badBuildOutMaxEntities",
			buildOutDegrees:    badBuildOutMaxEntities,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0032|Invalid value of build out degree '-1'"}}`,
		},
		{
			name:               "badDataSourceCode",
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
			recordKeys:         recordKeysBadDataSourceCodeFunc,
		},
		{
			name:               "badMaxDegree",
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0031|Invalid value of max degree '-1'"}}`,
			maxDegrees:         badMaxDegrees,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			recordKeys:         recordKeysFuncBadRecordIDFunc,
		},
		{
			name: "default",
		},
	}

	return result
}

func getTestCasesForFindPathByEntityID() []TestMetadataForFindPathByEntityID {
	result := []TestMetadataForFindPathByEntityID{
		{
			name:           "avoiding",
			avoidEntityIDs: avoidEntityIDsFunc,
		},
		{
			name:                "avoiding_and_including",
			avoidEntityIDs:      avoidEntityIDsFunc,
			maxDegrees:          1,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:               "avoiding_badStartEntityID",
			avoidEntityIDs:     avoidEntityIDsFunc,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
			startEntityID:      badEntityID,
		},
		{
			name:               "badAvoidEntityIDs",
			avoidEntityIDs:     badAvoidEntityIDsFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`,
		},
		{
			name:               "badEndEntityID",
			endEntityID:        badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
		},
		{
			name:               "badMaxDegrees",
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0031|Invalid value of max degree '-1'"}}`,
			maxDegrees:         badMaxDegrees,
		},
		{
			name:                "badRequiredDataSource",
			expectedErr:         szerror.ErrSzBadInput,
			expectedErrMessage:  `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044025","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`,
			requiredDataSources: badRequiredDataSourcesFunc,
		},
		{
			name:               "badStartEntityID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
			startEntityID:      badEntityID,
		},
		{
			name: "default",
		},
		{
			name:                "including",
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "including_badStartEntityID",
			expectedErr:         szerror.ErrSzNotFound,
			expectedErrMessage:  `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044025","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
			requiredDataSources: requiredDataSourcesFunc,
			startEntityID:       badEntityID,
		},
		{
			name:           "nilAvoidEntityIDs",
			avoidEntityIDs: nilAvoidEntityIDsFunc,
		},
		{
			name: "nilMaxDegrees",
			// expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0031|Invalid value of max degree '-9999'"}}`,
			maxDegrees:         nilMaxDegrees,
		},
		{
			name:                "nilRequiredDataSource",
			requiredDataSources: nilRequiredDataSourcesFunc,
		},
	}

	return result
}

func getTestCasesForFindPathByRecordID() []TestMetadataForFindPathByRecordID {
	result := []TestMetadataForFindPathByRecordID{
		{
			name:            "avoiding",
			avoidRecordKeys: avoidRecordIDsFunc,
			maxDegrees:      1,
		},
		{
			name:                "avoiding_and_including",
			avoidRecordKeys:     avoidRecordIDsFunc,
			maxDegrees:          1,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "avoiding_badStartDataSourceCode",
			avoidRecordKeys:     avoidRecordIDsFunc,
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
			startDataSourceCode: badDataSourceCode,
		},
		{
			name:               "badAvoidRecordKeys",
			avoidRecordKeys:    badAvoidRecordIDsFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`,
		},
		{
			name:               "badDataRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			startRecordID:      badRecordID,
		},
		{
			name:                "badDataSourceCode",
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
			startDataSourceCode: badDataSourceCode,
		},
		{
			name:               "badMaxDegrees",
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0031|Invalid value of max degree '-1'"}}`,
			maxDegrees:         badMaxDegrees,
		},
		{
			name:               "badRequiredDataSources",
			avoidRecordKeys:    badRequiredDataSourcesFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}`,
		},
		{
			name: "default",
		},
		{
			name:                "including",
			avoidRecordKeys:     avoidRecordIDsFunc,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "including_badDataSourceCode",
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044027","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
			requiredDataSources: requiredDataSourcesFunc,
			startDataSourceCode: badDataSourceCode,
		},
		{
			name:                "nilDataSourceCode",
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ2207|Data source code [] does not exist."}}`,
			startDataSourceCode: nilDataSourceCode,
		},
	}

	return result
}

func getTestCasesForGetEntityByEntityID() []TestMetadataForGetEntityByEntityID {
	result := []TestMetadataForGetEntityByEntityID{
		{
			name:               "badEntityID",
			entityID:           badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetEntityByEntityID","error":{"id":"SZSDK60044030","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEntityID",
			entityID:           nilEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetEntityByEntityID","error":{"id":"SZSDK60044030","reason":"SENZ0037|Unknown resolved entity value '0'"}}`,
		},
	}

	return result
}

func getTestCasesForGetEntityByRecordID() []TestMetadataForGetEntityByRecordID {
	result := []TestMetadataForGetEntityByRecordID{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			recordID:           badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ2207|Data source code [] does not exist."}}`,
		},
		{
			name:               "nilRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`,
			recordID:           nilRecordID,
		},
	}

	return result
}

func getTestCasesForGetRecord() []TestMetadataForGetRecord {
	result := []TestMetadataForGetRecord{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			recordID:           badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ2207|Data source code [] does not exist."}}`,
		},
		{
			name:               "nilRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`,
			recordID:           nilRecordID,
		},
	}

	return result
}

func getTestCasesForGetVirtualEntityByRecordID() []TestMetadataForGetVirtualEntityByRecordID {
	result := []TestMetadataForGetVirtualEntityByRecordID{
		{
			name:               "badDataSourceCode",
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetVirtualEntityByRecordID","error":{"id":"SZSDK60044038","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
			recordKeys:         badDataSourcesFunc,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetVirtualEntityByRecordID","error":{"id":"SZSDK60044038","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			recordKeys:         badRecordKeysFunc,
		},
		{
			name: "default",
		},
	}

	return result
}

func getTestCasesForHowEntityByEntityID() []TestMetadataForHowEntityByEntityID {
	result := []TestMetadataForHowEntityByEntityID{
		{
			name:               "badEntityID",
			entityID:           badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).HowEntityByEntityID","error":{"id":"SZSDK60044040","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEntityID",
			entityID:           nilEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).HowEntityByEntityID","error":{"id":"SZSDK60044040","reason":"SENZ0037|Unknown resolved entity value '0'"}}`,
		},
	}

	return result
}

func getTestCasesForGetRecordPreview() []TestMetadataForGetRecordPreview {
	result := []TestMetadataForGetRecordPreview{
		{
			name:               "badRecordDefinition",
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).GetRecordPreview","error":{"id":"SZSDK60044061","reason":"SENZ0002|Invalid Message"}}`,
			recordDefinition:   badRecordDefinition,
		},
		{
			name: "default",
		},
	}

	return result
}

func getTestCasesForProcessRedoRecord() []TestMetadataForProcessRedoRecord {
	result := []TestMetadataForProcessRedoRecord{
		{
			name: "badRedoRecord",
			// expectedErr:        szerror.ErrSzConfiguration,
			// expectedErrMessage: `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}`,
			redoRecord: badRedoRecord,
		},
		{
			name: "default",
		},
		{
			name: "nilRedoRecord",
			// expectedErr:        szerror.ErrSzBadInput,
			// expectedErrMessage: `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ0007|Empty Message"}}`,
			redoRecord: nilRedoRecord,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name: "withInfo_badRedoRecord",
			// expectedErr:        szerror.ErrSzConfiguration,
			// expectedErrMessage: `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}`,
			flags:      senzing.SzWithInfo,
			redoRecord: badRedoRecord,
		},
		{
			name:  "withInfo_nilRedoRecord",
			flags: senzing.SzWithInfo,
			// expectedErr:        szerror.ErrSzBadInput,
			// expectedErrMessage: `{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ0007|Empty Message"}}`,
			redoRecord: nilRedoRecord,
		},
	}

	return result
}

func getTestCasesForReevaluateEntity() []TestMetadataForReevaluateEntity {
	result := []TestMetadataForReevaluateEntity{
		{
			name:     "badEntityID",
			entityID: badEntityID,
		},
		{
			name: "default",
		},
		{
			name:     "nilEntityID",
			entityID: nilEntityID,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:     "withInfo_badEntityID",
			entityID: badEntityID,
			flags:    senzing.SzWithInfo,
		},
		{
			name:     "withInfo_nilEntityID",
			entityID: nilEntityID,
			flags:    senzing.SzWithInfo,
		},
	}

	return result
}

func getTestCasesForReevaluateRecord() []TestMetadataForReevaluateRecord {
	result := []TestMetadataForReevaluateRecord{
		{
			name:               "badDataSourceCode",
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044048","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
			dataSourceCode:     badDataSourceCode,
		},
		{
			name:     "badRecordID",
			recordID: badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044048","reason":"SENZ2207|Data source code [] does not exist."}}`,
		},
		{
			name:     "nilRecordID",
			recordID: nilRecordID,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:               "withInfo_badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044049","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
			flags:              senzing.SzWithInfo,
		},
		{
			name:     "withInfo_badRecordID",
			flags:    senzing.SzWithInfo,
			recordID: badRecordID,
		},
		{
			name:               "withInfo_nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044049","reason":"SENZ2207|Data source code [] does not exist."}}`,
			flags:              senzing.SzWithInfo,
		},
		{
			name:     "withInfo_nilRecordID",
			flags:    senzing.SzWithInfo,
			recordID: nilRecordID,
		},
	}

	return result
}

func getTestCasesForSearchByAttributes() []TestMetadataForSearchByAttributes {
	result := []TestMetadataForSearchByAttributes{
		{
			name:               "badAttributes",
			attributes:         badAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0027|Invalid value for search-attributes"}}`,
		},
		{
			name:               "badSearchProfile",
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0088|Unknown search profile value '}{'"}}`,
			searchProfile:      badSearchProfile,
		},
		{
			name: "default",
		},
		{
			name:               "nilAttributes",
			attributes:         nilAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0027|Invalid value for search-attributes"}}`,
		},
		{
			name:          "nilSearchProfile",
			searchProfile: nilSearchProfile,
		},
	}

	return result
}

func getTestCasesForWhyEntities() []TestMetadataForWhyEntities {
	result := []TestMetadataForWhyEntities{
		{
			name:               "badEnitity1",
			entityID1:          badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
		},
		{
			name:               "badEnitity2",
			entityID2:          badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEnitity1",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '0'"}}`,
			entityID1:          nilEntityID,
		},
		{
			name:               "nilEnitity2",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '0'"}}`,
			entityID2:          nilEntityID,
		},
	}

	return result
}

func getTestCasesForWhyRecordInEntity() []TestMetadataForWhyRecordInEntity {
	result := []TestMetadataForWhyRecordInEntity{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			recordID:           badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ2207|Data source code [] does not exist."}}`,
		},
		{
			name:               "nilRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`,
			recordID:           nilRecordID,
		},
	}

	return result
}

func getTestCasesForWhyRecords() []TestMetadataForWhyRecords {
	result := []TestMetadataForWhyRecords{
		{
			name:               "badDataSourceCode1",
			dataSourceCode1:    badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badDataSourceCode2",
			dataSourceCode2:    badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}`,
		},
		{
			name:               "badRecordID1",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			recordID1:          badRecordID,
		},
		{
			name:               "badRecordID2",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}`,
			recordID2:          badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode1",
			dataSourceCode1:    nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [] does not exist."}}`,
		},
		{
			name:               "nilDataSourceCode2",
			dataSourceCode2:    nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [] does not exist."}}`,
		},
		{
			name:               "nilRecordID1",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`,
			recordID1:          nilRecordID,
		},
		{
			name:               "nilRecordID2",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}`,
			recordID2:          nilRecordID,
		},
	}

	return result
}

func getTestCasesForWhySearch() []TestMetadataForWhySearch {
	result := []TestMetadataForWhySearch{
		{
			name:               "badAttributes",
			attributes:         badAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0027|Invalid value for search-attributes"}}`,
		},
		{
			name:               "badEntityID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0037|Unknown resolved entity value '-1'"}}`,
			entityID:           badEntityID,
		},
		{
			name:               "badSearchProfile",
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0088|Unknown search profile value '}{'"}}`,
			searchProfile:      badSearchProfile,
		},
		{
			name: "default",
		},
		{
			name:               "nilAttributes",
			attributes:         nilAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0027|Invalid value for search-attributes"}}`,
		},
		{
			name:               "nilEntityID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0037|Unknown resolved entity value '0'"}}`,
			entityID:           nilEntityID,
		},
		{
			name:          "searchProfile",
			searchProfile: "SEARCH",
		},
	}

	return result
}

// ----------------------------------------------------------------------------
// Test data helpers
// ----------------------------------------------------------------------------

func avoidEntityIDsFunc() string {
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord(ctx, szEngine, "CUSTOMERS", "1001") + `}]}`

	return result
}

func avoidRecordIDsFunc() string {
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord(ctx, szEngine, "CUSTOMERS", "1001") + `}]}`

	return result
}

func badAvoidEntityIDsFunc() string {
	return "}{"
}

func badAvoidRecordIDsFunc() string {
	return "}{"
}

func badRequiredDataSourcesFunc() string {
	return "}{"
}

func badDataSourcesFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]

	return `{"RECORDS": [{"DATA_SOURCE": "` +
		badDataSourceCode +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}]}`
}

func badEntityIDsFunc() string {
	return `{"ENTITIES": [{"ENTITY_ID": ` +
		strconv.FormatInt(badEntityID, baseTen) +
		`}]}`
}

func badRecordKeysFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]

	return `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		badRecordID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}]}`
}

func entityIDsJSON(entityIDs ...int64) string {
	result := `{"ENTITIES": [`

	for _, entityID := range entityIDs {
		result += `{"ENTITY_ID":` + strconv.FormatInt(entityID, baseTen) + `},`
	}

	result = result[:len(result)-1] // Remove final comma.
	result += `]}`

	return result
}

func nilAvoidEntityIDsFunc() string {
	var result string

	return result
}

func nilRequiredDataSourcesFunc() string {
	var result string

	return result
}

func requiredDataSourcesFunc() string {
	record := truthset.CustomerRecords["1001"]

	return `{"DATA_SOURCES": ["` + record.DataSource + `"]}`
}

func recordKeysJSON(r1DS, r1ID, r2DS, r2ID, r3DS, r3ID string) string {
	return `{"RECORDS": [{"DATA_SOURCE": "` +
		r1DS +
		`", "RECORD_ID": "` +
		r1ID +
		`"}, {"DATA_SOURCE": "` +
		r2DS +
		`", "RECORD_ID": "` +
		r2ID +
		`"}, {"DATA_SOURCE": "` +
		r3DS +
		`", "RECORD_ID": "` +
		r3ID +
		`"}]}`
}

func recordKeysFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]

	return recordKeysJSON(
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		record3.DataSource,
		record3.ID,
	)
}

func recordKeysBadDataSourceCodeFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]

	return recordKeysJSON(
		badDataSourceCode,
		record1.ID,
		record2.DataSource,
		record2.ID,
		record3.DataSource,
		record3.ID,
	)
}

func recordKeysFuncBadRecordIDFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]

	return recordKeysJSON(
		record1.DataSource,
		badRecordID,
		record2.DataSource,
		record2.ID,
		record3.DataSource,
		record3.ID,
	)
}

// Return first non-zero length candidate.  Last candidate is default.
func xString(candidates ...string) string {
	var result string
	for _, result = range candidates {
		if result == nilSemaphoreString {
			return ""
		}

		if len(result) > 0 {
			return result
		}
	}

	return result
}

// Return first non-zero candidate.  Last candidate is default.
func xInt64(candidates ...int64) int64 {
	var result int64
	for _, result = range candidates {
		if result == nilSemaphoreInt64 {
			return 0
		}

		if result != 0 {
			return result
		}
	}

	return result
}
