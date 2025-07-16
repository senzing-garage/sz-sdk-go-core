package szdiagnostic_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
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
	jsonIndentation   = "    "
	observerOrigin    = "SzDiagnostic observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
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
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
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

func TestSzdiagnostic_CheckRepositoryPerformance(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckRepositoryPerformance(ctx, secondsToRun)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_CheckRepositoryPerformance_badSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.CheckRepositoryPerformance(ctx, badSecondsToRun)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_CheckRepositoryPerformance_nilSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.CheckRepositoryPerformance(ctx, nilSecondsToRun)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_GetRepositoryInfo(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.GetRepositoryInfo(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_GetFeature(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szDiagnostic := getTestObject(test)
	featureID := int64(1)
	actual, err := szDiagnostic.GetFeature(ctx, featureID)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_GetFeature_badFeatureID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.GetFeature(ctx, badFeatureID)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"id":"SZSDK60034004","reason":"SENZ0057|Unknown feature ID value '-1'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzdiagnostic_GetFeature_nilFeatureID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.GetFeature(ctx, nilFeatureID)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"id":"SZSDK60034004","reason":"SENZ0057|Unknown feature ID value '0'"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

// PurgeRepository is tested in szdiagnostic_examples_test.go
// func TestSzdiagnostic_PurgeRepository(test *testing.T) {}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzdiagnostic_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzdiagnostic_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	szDiagnostic.SetObserverOrigin(ctx, originMessage)
}

func TestSzdiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	szDiagnostic.SetObserverOrigin(ctx, originMessage)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, originMessage, actual)
}

func TestSzdiagnostic_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzdiagnostic_AsInterface(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getSzDiagnosticAsInterface(ctx)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckRepositoryPerformance(ctx, secondsToRun)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

// func TestSzdiagnostic_Initialize(test *testing.T) {
// 	ctx := test.Context()
// 	szDiagnostic := &szdiagnostic.Szdiagnostic{}
// 	settings := getSettings()
// 	configID := senzing.SzInitializeWithDefaultConfiguration
// 	err := szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
// 	printDebug(test, err)
// 	require.NoError(test, err)
// }

// func TestSzdiagnostic_Initialize_withConfigId(test *testing.T) {
// 	ctx := test.Context()
// 	szDiagnostic := &szdiagnostic.Szdiagnostic{}
// 	settings := getSettings()
// 	configID := getDefaultConfigID()
// 	err := szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
// 	printDebug(test, err)
// 	require.NoError(test, err)
// }

// func TestSzdiagnostic_Initialize_withConfigId_badConfigID(test *testing.T) {
// 	// IMPROVE: Implement TestSzdiagnostic_Initialize_withConfigId_badConfigID
// 	_ = test
// }

func TestSzdiagnostic_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(test)
	configID := getDefaultConfigID()
	err := szDiagnostic.Reinitialize(ctx, configID)
	printDebug(test, err)
	require.NoError(test, err)
}

// func TestSzdiagnostic_Reinitialize_error(test *testing.T) {
// 	// IMPROVE: Implement TestSzdiagnostic_Reinitialize_error
// 	_ = test
// }

func TestSzdiagnostic_Destroy(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzdiagnostic_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

// func TestSzdiagnostic_Destroy_error(test *testing.T) {
// 	// IMPROVE: Implement TestSzdiagnostic_Destroy_error
// 	_ = test
// }

func TestSzdiagnostic_cleanup(test *testing.T) {
	ctx := test.Context()
	destroySzConfigManagers(ctx)
	destroySzDiagnostics(ctx)
	destroySzEngines(ctx)
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

func deleteRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		panicOnError(err)
	}
}

func destroySzConfigManagers(ctx context.Context) {
	szConfigManager := &szconfigmanager.Szconfigmanager{}
	for {
		err := szConfigManager.Destroy(ctx)
		if err != nil {
			break
		}
	}
}

func destroySzDiagnostics(ctx context.Context) {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	for {
		err := szDiagnostic.Destroy(ctx)
		if err != nil {
			break
		}
	}
}

func destroySzEngines(ctx context.Context) {
	szEngine := &szengine.Szengine{}
	for {
		err := szEngine.Destroy(ctx)
		if err != nil {
			break
		}
	}
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getDefaultConfigID() int64 {
	return defaultConfigID
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

func getSzDiagnostic(ctx context.Context) *szdiagnostic.Szdiagnostic {
	if szDiagnosticSingleton == nil {
		settings := getSettings()
		szDiagnosticSingleton = &szdiagnostic.Szdiagnostic{}
		err := szDiagnosticSingleton.SetLogLevel(ctx, logLevel)
		panicOnError(err)

		if logLevel == "TRACE" {
			szDiagnosticSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szDiagnosticSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)
			err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}

		err = szDiagnosticSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		panicOnError(err)
	}

	return szDiagnosticSingleton
}

func getSzDiagnosticAsInterface(ctx context.Context) senzing.SzDiagnostic {
	return getSzDiagnostic(ctx)
}

func getSzEngine(ctx context.Context) senzing.SzEngine {
	if szEngineSingleton == nil {
		settings := getSettings()
		szEngineSingleton = &szengine.Szengine{}
		err := szEngineSingleton.SetLogLevel(ctx, logLevel)
		panicOnError(err)

		if logLevel == "TRACE" {
			szEngineSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szEngineSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)
			err = szEngineSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}

		err = szEngineSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		panicOnError(err)
	}

	return szEngineSingleton
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szdiagnostic")
}

func getTestObject(t *testing.T) *szdiagnostic.Szdiagnostic {
	t.Helper()
	return getSzDiagnostic(t.Context())
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

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	defaultConfigID, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	return nil
}

func teardown() {
	ctx := context.TODO()
	teardownSzDiagnostic(ctx)
	teardownSzEngine(ctx)
}

func teardownSzDiagnostic(ctx context.Context) {
	if szDiagnosticSingleton != nil {
		err := szDiagnosticSingleton.UnregisterObserver(ctx, observerSingleton)
		panicOnError(err)
		szDiagnosticSingleton = nil
	}
}

func teardownSzEngine(ctx context.Context) {
	if szEngineSingleton != nil {
		err := szEngineSingleton.UnregisterObserver(ctx, observerSingleton)
		panicOnError(err)
		szEngineSingleton = nil
	}
}
