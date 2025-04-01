package szconfig_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
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
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
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
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.AddDataSource(ctx, dataSourceCode)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.AddDataSource(ctx, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_nilDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.AddDataSource(ctx, nilDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
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
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_, err := szConfig.DeleteDataSource(ctx, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzconfig_DeleteDataSource_nilDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_, err := szConfig.DeleteDataSource(ctx, nilDataSourceCode)
	require.NoError(test, err)
}

func TestSzconfig_Export(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.Export(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

func TestSzconfig_Import(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)
	err = szConfig.Import(ctx, configDefinition)
	require.NoError(test, err)
}

func TestSzconfig_Import_badConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Import(ctx, badConfigDefinition)
	require.NoError(test, err)
}

func TestSzconfig_Import_nilConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Import(ctx, nilConfigDefinition)
	require.NoError(test, err)
}

func TestSzconfig_ImportTemplate(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.ImportTemplate(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfig_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfig_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzconfig_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfig_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfig_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szConfig := getSzConfigAsInterface(ctx)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_Initialize(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	settings, err := getSettings()
	require.NoError(test, err)
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Initialize_badSettings(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Initialize(ctx, instanceName, badSettings, verboseLogging)
	assert.NoError(test, err)
}

// TODO: Implement TestSzconfig_Initialize_error
// func TestSzconfig_Initialize_error(test *testing.T) {}

func TestSzconfig_Initialize_again(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	settings, err := getSettings()
	require.NoError(test, err)
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfig_Destroy_error
// func TestSzconfig_Destroy_error(test *testing.T) {}

func TestSzconfig_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szConfigSingleton = nil
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
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

func getSzConfig(ctx context.Context) *szconfig.Szconfig {
	if szConfigSingleton == nil {
		settings, err := getSettings()
		handleErrorWithPanic(err)
		szConfigSingleton = &szconfig.Szconfig{}
		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
		handleErrorWithPanic(err)
		if logLevel == "TRACE" {
			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)
			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}
		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		handleErrorWithPanic(err)
		szConfigSingleton.ImportTemplate(ctx)
	}
	return szConfigSingleton
}

func getSzConfigAsInterface(ctx context.Context) senzing.SzConfig {
	return getSzConfig(ctx)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func getTestObject(ctx context.Context, test *testing.T) *szconfig.Szconfig {
	_ = test
	return getSzConfig(ctx)
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
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error
	err = setupDirectories()
	handleErrorWithPanic(err)
	err = setupDatabase()
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

func teardown() error {
	ctx := context.TODO()
	err := teardownSzConfig(ctx)
	return err
}

func teardownSzConfig(ctx context.Context) error {
	err := szConfigSingleton.UnregisterObserver(ctx, observerSingleton)
	handleErrorWithPanic(err)
	err = szConfigSingleton.Destroy(ctx)
	handleErrorWithPanic(err)
	szConfigSingleton = nil
	return err
}
