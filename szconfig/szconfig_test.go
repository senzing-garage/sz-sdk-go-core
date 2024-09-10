package szconfig

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
	nilConfigHandle     uintptr
	nilDataSourceCode   string
)

var (
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigSingleton *Szconfig
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzconfig_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_AddDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	require.NoError(test, err)
	actual, err := szConfig.AddDataSource(ctx, configHandle2, dataSourceCode)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	require.NoError(test, err)
}

func TestSzconfig_AddDataSource_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.AddDataSource(ctx, configHandle, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_nilDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.AddDataSource(ctx, configHandle, nilDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_nilConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.AddDataSource(ctx, nilConfigHandle, dataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzconfig_CloseConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_CloseConfig_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.CloseConfig(ctx, badConfigHandle)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_CloseConfig_nilConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.CloseConfig(ctx, nilConfigHandle)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "Original", actual)
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "     Add", actual)
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_DeleteDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "Original", actual)
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "     Add", actual)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	require.NoError(test, err)
	err = szConfig.DeleteDataSource(ctx, configHandle2, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle2)
	require.NoError(test, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	require.NoError(test, err)
}

func TestSzconfig_DeleteDataSource_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.DeleteDataSource(ctx, badConfigHandle, dataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_DeleteDataSource_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	err = szConfig.DeleteDataSource(ctx, configHandle, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzconfig_DeleteDataSource_nilConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.DeleteDataSource(ctx, nilConfigHandle, dataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_DeleteDataSource_nilDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	err = szConfig.DeleteDataSource(ctx, configHandle, nilDataSourceCode)
	require.NoError(test, err)
}

func TestSzconfig_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_ExportConfig_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.ExportConfig(ctx, badConfigHandle)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_ExportConfig_nilConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.ExportConfig(ctx, nilConfigHandle)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_GetDataSources_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.GetDataSources(ctx, badConfigHandle)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_GetDataSources_nilConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.GetDataSources(ctx, nilConfigHandle)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSz)
}

func TestSzconfig_ImportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	actual, err := szConfig.ImportConfig(ctx, configDefinition)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_ImportConfig_badConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_, err := szConfig.ImportConfig(ctx, badConfigDefinition)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzconfig_ImportConfig_nilConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_, err := szConfig.ImportConfig(ctx, nilConfigDefinition)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func TestSzconfig_getByteArray(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	szProduct.getByteArray(10)
}

func TestSzconfig_getByteArrayC(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	szProduct.getByteArrayC(10)
}

func TestSzconfig_newError(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.newError(ctx, 1)
	require.Error(test, err)
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
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
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
	if err != nil {
		return result, fmt.Errorf("failed to make target database path (%s) absolute. Error: %w", dbTargetPath, err)
	}
	databaseURL := fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	if err != nil {
		return result, fmt.Errorf("failed to BuildSimpleSettingsUsingMap(%s) Error: %w", configAttrMap, err)
	}
	return result, err
}

func getSzConfig(ctx context.Context) (*Szconfig, error) {
	var err error
	if szConfigSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szConfigSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		szConfigSingleton = &Szconfig{}
		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szConfigSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szConfigSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szConfigSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			return szConfigSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szConfigSingleton, err
}

func getSzConfigAsInterface(ctx context.Context) senzing.SzConfig {
	result, err := getSzConfig(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func getTestObject(ctx context.Context, test *testing.T) *Szconfig {
	result, err := getSzConfig(ctx)
	require.NoError(test, err)
	return result
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
	if err != nil {
		return fmt.Errorf("Failed to set up directories. Error: %w", err)
	}
	err = setupDatabase()
	if err != nil {
		return fmt.Errorf("Failed to set up database. Error: %w", err)
	}
	return err
}

func setupDatabase() error {
	var err error

	// Locate source and target paths.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return fmt.Errorf("failed to make target database path (%s) absolute. Error: %w",
			dbTargetPath, err)
	}
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	if err != nil {
		return fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			databaseTemplatePath, err.Error())
	}

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	if err != nil {
		return fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
			databaseTemplatePath, testDirectoryPath, err)
	}
	return err
}

func setupDirectories() error {
	var err error
	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("failed to remove target test directory (%v): %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("failed to recreate target test directory (%v): %w", testDirectoryPath, err)
	}
	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzConfig(ctx)
	return err
}

func teardownSzConfig(ctx context.Context) error {
	err := szConfigSingleton.UnregisterObserver(ctx, observerSingleton)
	if err != nil {
		return err
	}
	err = szConfigSingleton.Destroy(ctx)
	if err != nil {
		return err
	}
	szConfigSingleton = nil
	return err
}
