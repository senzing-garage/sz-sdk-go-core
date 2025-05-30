package szconfig_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseTen           = 10
	dataSourceCode    = "GO_TEST"
	defaultTruncation = 76
	instanceName      = "SzConfig Test"
	jsonIndentation   = "    "
	observerOrigin    = "SzConfig observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badConfigDefinition = "}{"
	badConfigHandle     = uintptr(0)
	badDataSourceCode   = "\n\tGO_TEST"
	badLogLevelName     = "BadLogLevelName"
	badSettings         = "{]"
)

// Nil/empty parameters

var (
	nilConfigDefinition string
	nilDataSourceCode   string
)

var (
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigSingleton *szconfig.Szconfig
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func TestSzconfig_AddDataSource(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, dataSourceCode)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfig_AddDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, badDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).AddDataSource","error":{"function":"szconfig.(*Szconfig).addDataSourceChoreography","text":"addDataSource: \n\tGO_TEST","error":{"id":"SZSDK60014001","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_AddDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, nilDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).AddDataSource","error":{"function":"szconfig.(*Szconfig).addDataSourceChoreography","text":"addDataSource: ","error":{"id":"SZSDK60014001","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_DeleteDataSource(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSources(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	printResult(test, "Original", actual)

	_, _ = szConfig.AddDataSource(ctx, dataSourceCode)
	actual, err = szConfig.GetDataSources(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	printResult(test, "     Add", actual)

	_, err = szConfig.DeleteDataSource(ctx, dataSourceCode)
	printDebug(test, err, actual)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	printResult(test, "  Delete", actual)
}

func TestSzconfig_DeleteDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.DeleteDataSource(ctx, badDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).DeleteDataSource","error":{"function":"szconfig.(*Szconfig).deleteDataSourceChoreography","text":"deleteDataSource(\n\tGO_TEST)","error":{"id":"SZSDK60014004","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_DeleteDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.DeleteDataSource(ctx, nilDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).DeleteDataSource","error":{"function":"szconfig.(*Szconfig).deleteDataSourceChoreography","text":"deleteDataSource()","error":{"id":"SZSDK60014004","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_Export(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.Export(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfig_GetDataSources(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSources(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

func TestSzconfig_Import(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	printDebug(test, err)
	require.NoError(test, err)
	err = szConfig.Import(ctx, configDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Import_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, badConfigDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Import_nilConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, nilConfigDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_ImportTemplate(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.ImportTemplate(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	printDebug(test, err)
	require.NoError(test, err)
	err = szConfig.VerifyConfigDefinition(ctx, configDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.VerifyConfigDefinition(ctx, badConfigDefinition)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).verifyConfigDefinitionChoreography","text":"load","error":{"id":"SZSDK60014009","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfig_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfig_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	origin := originMessage
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzconfig_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	origin := originMessage
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfig_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfig_AsInterface(test *testing.T) {
	ctx := test.Context()
	szConfig := getSzConfigAsInterface(ctx)
	actual, err := szConfig.GetDataSources(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfig_Initialize(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Initialize_badSettings(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Initialize(ctx, instanceName, badSettings, verboseLogging)
	assert.NoError(test, err)
}

func TestSzconfig_Initialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfig_Initialize_error
	_ = test
}

func TestSzconfig_Initialize_again(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Destroy(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Destroy_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfig_Destroy_error
	_ = test
}

func TestSzconfig_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szConfigSingleton = nil
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
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

func getSzConfig(ctx context.Context) *szconfig.Szconfig {
	if szConfigSingleton == nil {
		settings := getSettings()
		szConfigSingleton = &szconfig.Szconfig{}
		err := szConfigSingleton.SetLogLevel(ctx, logLevel)
		panicOnError(err)

		if logLevel == "TRACE" {
			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)
			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}

		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
		err = szConfigSingleton.ImportTemplate(ctx)
		panicOnError(err)
	}

	return szConfigSingleton
}

func getSzConfigAsInterface(ctx context.Context) senzing.SzConfig {
	return getSzConfig(ctx)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func getTestObject(t *testing.T) *szconfig.Szconfig {
	t.Helper()

	return getSzConfig(t.Context())
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

func teardown() {
	ctx := context.TODO()
	teardownSzConfig(ctx)
}

func teardownSzConfig(ctx context.Context) {
	err := szConfigSingleton.UnregisterObserver(ctx, observerSingleton)
	panicOnError(err)
	err = szConfigSingleton.Destroy(ctx)
	panicOnError(err)

	szConfigSingleton = nil
}
