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
	observerOrigin    = "SzConfig observer"
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
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, nilDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_DeleteDataSource(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printResult(test, "Original", actual)

	_, _ = szConfig.AddDataSource(ctx, dataSourceCode)
	actual, err = szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printResult(test, "     Add", actual)

	_, err = szConfig.DeleteDataSource(ctx, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printResult(test, "  Delete", actual)
}

func TestSzconfig_DeleteDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_, err := szConfig.DeleteDataSource(ctx, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzconfig_DeleteDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_, err := szConfig.DeleteDataSource(ctx, nilDataSourceCode)
	require.NoError(test, err)
}

func TestSzconfig_Export(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.Export(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_GetDataSources(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

func TestSzconfig_Import(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)
	err = szConfig.Import(ctx, configDefinition)
	require.NoError(test, err)
}

func TestSzconfig_Import_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, badConfigDefinition)
	require.NoError(test, err)
}

func TestSzconfig_Import_nilConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, nilConfigDefinition)
	require.NoError(test, err)
}

func TestSzconfig_ImportTemplate(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.ImportTemplate(ctx)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)
	err = szConfig.VerifyConfigDefinition(ctx, configDefinition)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.VerifyConfigDefinition(ctx, badConfigDefinition)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
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
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzconfig_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfig_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfig_AsInterface(test *testing.T) {
	ctx := test.Context()
	szConfig := getSzConfigAsInterface(ctx)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_Initialize(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Initialize_badSettings(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Initialize(ctx, instanceName, badSettings, verboseLogging)
	assert.NoError(test, err)
}

// IMPROVE: Implement TestSzconfig_Initialize_error
// func TestSzconfig_Initialize_error(test *testing.T) {}

func TestSzconfig_Initialize_again(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Destroy(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
	require.NoError(test, err)
}

// IMPROVE: Implement TestSzconfig_Destroy_error
// func TestSzconfig_Destroy_error(test *testing.T) {}

func TestSzconfig_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szConfigSingleton = nil
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
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
	ctx := t.Context()

	return getSzConfig(ctx)
}

func handleError(err error) {
	if err != nil {
		safePrintln("Error:", err)
	}
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

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func safePrintln(message ...any) {
	fmt.Println(message...) //nolint
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
