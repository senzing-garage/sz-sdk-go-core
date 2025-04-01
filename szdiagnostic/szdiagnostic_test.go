package szdiagnostic_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	instanceName      = "SzDiagnostic Test"
	observerOrigin    = "SzDiagnostic observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badFeatureID    = int64(-1)
	badLogLevelName = "BadLogLevelName"
	badSecondsToRun = -1
)

// Nil/empty parameters

var (
	nilSecondsToRun int
	nilFeatureID    int64
)

var (
	defaultConfigID   int64
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szDiagnosticSingleton *szdiagnostic.Szdiagnostic
	szEngineSingleton     *szengine.Szengine
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzdiagnostic_CheckDatastorePerformance(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_CheckDatastorePerformance_badSecondsToRun(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, badSecondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_CheckDatastorePerformance_nilSecondsToRun(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, nilSecondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetDatastoreInfo(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.GetDatastoreInfo(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { handleError(deleteRecords(ctx, records)) }()

	err := addRecords(ctx, records)
	require.NoError(test, err)
	szDiagnostic := getTestObject(ctx, test)
	featureID := int64(1)
	actual, err := szDiagnostic.GetFeature(ctx, featureID)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature_badFeatureID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { handleError(deleteRecords(ctx, records)) }()

	err := addRecords(ctx, records)
	require.NoError(test, err)
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.GetFeature(ctx, badFeatureID)
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature_nilFeatureID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { handleError(deleteRecords(ctx, records)) }()

	err := addRecords(ctx, records)
	require.NoError(test, err)
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.GetFeature(ctx, nilFeatureID)
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

// PurgeRepository is tested in szdiagnostic_examples_test.go
// func TestSzdiagnostic_PurgeRepository(test *testing.T) {}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzdiagnostic_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzdiagnostic_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
}

func TestSzdiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzdiagnostic_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzdiagnostic_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getSzDiagnosticAsInterface(ctx)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_Initialize(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	settings, err := getSettings()
	require.NoError(test, err)

	configID := senzing.SzInitializeWithDefaultConfiguration
	err = szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

func TestSzdiagnostic_Initialize_withConfigId(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	settings, err := getSettings()
	require.NoError(test, err)

	configID := getDefaultConfigID()
	err = szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzdiagnostic_Initialize_withConfigId_badConfigID
// func TestSzdiagnostic_Initialize_withConfigId_badConfigID(test *testing.T) {}

func TestSzdiagnostic_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(ctx, test)
	configID := getDefaultConfigID()
	err := szDiagnostic.Reinitialize(ctx, configID)
	require.NoError(test, err)
}

// TODO: Implement TestSzdiagnostic_Reinitialize_error
// func TestSzdiagnostic_Reinitialize_error(test *testing.T) {}

func TestSzdiagnostic_Destroy(test *testing.T) {
	ctx := context.TODO()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.Destroy(ctx)
	require.NoError(test, err)
}

func TestSzdiagnostic_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzdiagnostic_Destroy_error
// func TestSzdiagnostic_Destroy_error(test *testing.T) {}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, records []record.Record) error {
	var err error

	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err = szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		handleErrorWithPanic(err)
	}

	return err
}

func deleteRecords(ctx context.Context, records []record.Record) error {
	var err error

	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err = szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		handleErrorWithPanic(err)
	}

	return err
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getDefaultConfigID() int64 {
	return defaultConfigID
}

func getSettings() (string, error) {
	var result string

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	handleErrorWithPanic(err)

	databaseURL := fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	handleErrorWithPanic(err)

	return result, err
}

func getSzDiagnostic(ctx context.Context) *szdiagnostic.Szdiagnostic {
	if szDiagnosticSingleton == nil {
		settings, err := getSettings()
		handleErrorWithPanic(err)

		szDiagnosticSingleton = &szdiagnostic.Szdiagnostic{}
		err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel)
		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szDiagnosticSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szDiagnosticSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)
			err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}

		err = szDiagnosticSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		handleErrorWithPanic(err)
	}

	return szDiagnosticSingleton
}

func getSzDiagnosticAsInterface(ctx context.Context) senzing.SzDiagnostic {
	return getSzDiagnostic(ctx)
}

func getSzEngine(ctx context.Context) senzing.SzEngine {
	if szEngineSingleton == nil {
		settings, err := getSettings()
		handleErrorWithPanic(err)

		szEngineSingleton = &szengine.Szengine{}
		err = szEngineSingleton.SetLogLevel(ctx, logLevel)
		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szEngineSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szEngineSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)
			err = szEngineSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}

		err = szEngineSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		handleErrorWithPanic(err)
	}

	return szEngineSingleton
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szdiagnostic")
}

func getTestObject(ctx context.Context, test *testing.T) *szdiagnostic.Szdiagnostic {
	_ = test
	return getSzDiagnostic(ctx)
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func handleErrorWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if errors.Is(err, szerror.ErrSzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}

		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}

		if errors.Is(err, szerror.ErrSzBadInput) {
			fmt.Printf("\nBad user input error detected. \n\n")
		}

		fmt.Print(err)

		os.Exit(1)
	}

	code := m.Run()
	err = teardown()
	handleErrorWithPanic(err)

	os.Exit(code)
}

func setup() error {
	var err error

	err = setupDirectories()
	handleErrorWithPanic(err)
	err = setupDatabase()
	handleErrorWithPanic(err)
	err = setupSenzingConfiguration()
	handleErrorWithPanic(err)

	return err
}

func setupDatabase() error {
	var err error

	// Locate source and target paths.

	testDirectoryPath := getTestDirectoryPath()
	_, err = filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	handleErrorWithPanic(err)
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	handleErrorWithPanic(err)

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	handleErrorWithPanic(err)

	return err
}

func setupDirectories() error {
	var err error

	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	handleErrorWithPanic(err)
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	handleErrorWithPanic(err)

	return err
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	settings, err := getSettings()
	handleErrorWithPanic(err)

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	handleErrorWithPanic(err)

	defer func() { handleErrorWithPanic(szConfig.Destroy(ctx)) }()

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	handleErrorWithPanic(err)

	defer func() { handleErrorWithPanic(szConfigManager.Destroy(ctx)) }()

	// Create a Senzing configuration.

	err = szConfig.ImportTemplate(ctx)
	handleErrorWithPanic(err)

	// Add data sources to template Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		handleErrorWithPanic(err)
	}

	// Create a string representation of the Senzing configuration.

	configDefinition, err := szConfig.Export(ctx)
	handleErrorWithPanic(err)

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	defaultConfigID, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	handleErrorWithPanic(err)

	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzDiagnostic(ctx)
	handleErrorWithPanic(err)
	err = teardownSzEngine(ctx)
	handleErrorWithPanic(err)

	return err
}

func teardownSzDiagnostic(ctx context.Context) error {
	var err error

	err = szDiagnosticSingleton.UnregisterObserver(ctx, observerSingleton)
	handleErrorWithPanic(err)
	err = szDiagnosticSingleton.Destroy(ctx)
	handleErrorWithPanic(err)

	szDiagnosticSingleton = nil

	return err
}

func teardownSzEngine(ctx context.Context) error {
	err := szEngineSingleton.UnregisterObserver(ctx, observerSingleton)
	handleErrorWithPanic(err)
	err = szEngineSingleton.Destroy(ctx)
	handleErrorWithPanic(err)

	szEngineSingleton = nil

	return err
}
