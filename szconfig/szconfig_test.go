package szconfig

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szconfig"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	instanceName      = "SzConfig Test"
	observerOrigin    = "SzConfig observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

var (
	logger            logging.LoggingInterface
	logLevel          = "INFO"
	observerSingleton = &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	szConfigSingleton *Szconfig
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzconfig_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	dataSourceCode := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	assert.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	assert.NoError(test, err)
}

func TestSzconfig_AddDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	assert.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	assert.NoError(test, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	assert.NoError(test, err)
	dataSourceCode := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle2, dataSourceCode)
	assert.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	assert.NoError(test, err)
}

func TestSzconfig_AddDataSource_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	badDataSourceCode := "\n\tGO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle, badDataSourceCode)
	assert.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_CloseConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	assert.NoError(test, err)
}

// TODO: Implement TestSzconfig_CloseConfig_badConfigHandle
// func TestSzconfig_CloseConfig_badConfigHandle(test *testing.T) {
// 	ctx := context.TODO()
// 	szConfig := getTestObject(ctx, test)
// 	badConfigHandle := uintptr(0)
// 	err := szConfig.CloseConfig(ctx, badConfigHandle)
// 	expectError(test, szerror.SzUnhandledError{}, err)
// }

func TestSzconfig_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzconfig_CreateConfig_badXxx
// func TestSzconfig_CreateConfig_badXxx(test *testing.T) {}

func TestSzconfig_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	assert.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	assert.NoError(test, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	assert.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	assert.NoError(test, err)
	printResult(test, "     Add", actual)
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	assert.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	assert.NoError(test, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	assert.NoError(test, err)
}

func TestSzconfig_DeleteDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	assert.NoError(test, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	assert.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	assert.NoError(test, err)
	printResult(test, "     Add", actual)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	assert.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	assert.NoError(test, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	assert.NoError(test, err)
	err = szConfig.DeleteDataSource(ctx, configHandle2, dataSourceCode)
	assert.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle2)
	assert.NoError(test, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	assert.NoError(test, err)
}

func TestSzconfig_DeleteDataSource_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	badConfigHandle := uintptr(0)
	dataSourceCode := "GO_TEST"
	err := szConfig.DeleteDataSource(ctx, badConfigHandle, dataSourceCode)
	assert.ErrorIs(test, err, szerror.ErrSzBase)
}

func TestSzconfig_DeleteDataSource_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	badDataSourceCode := "\n\tGO_TEST"
	err = szConfig.DeleteDataSource(ctx, configHandle, badDataSourceCode)
	assert.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzconfig_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	actual, err := szConfig.ExportConfig(ctx, configHandle)
	assert.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_ExportConfig_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	badConfigHandle := uintptr(0)
	actual, err := szConfig.ExportConfig(ctx, badConfigHandle)
	assert.Equal(test, "", actual)
	assert.ErrorIs(test, err, szerror.ErrSzBase)
}

func TestSzconfig_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	assert.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	assert.NoError(test, err)
}

func TestSzconfig_GetDataSources_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	badConfigHandle := uintptr(0)
	actual, err := szConfig.GetDataSources(ctx, badConfigHandle)
	assert.Equal(test, "", actual)
	assert.ErrorIs(test, err, szerror.ErrSzBase)
}

func TestSzconfig_ImportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	assert.NoError(test, err)
	actual, err := szConfig.ImportConfig(ctx, configDefinition)
	assert.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_ImportConfig_badConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	badConfigDefinition := "}{"
	_, err := szConfig.ImportConfig(ctx, badConfigDefinition)
	assert.ErrorIs(test, err, szerror.ErrSzBadInput)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfig_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	badLogLevelName := "BadLogLevelName"
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
	assert.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfig_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szConfig := getSzConfigAsInterface(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	assert.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	assert.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	assert.NoError(test, err)
}

func TestSzconfig_Initialize(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	instanceName := "Test name"
	verboseLogging := senzing.SzNoLogging
	settings, err := getSettings()
	assert.NoError(test, err)
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	assert.NoError(test, err)
}

// TODO: Implement TestSzconfig_Initialize_badSettings
// func TestSzconfig_Initialize_badSettings(test *testing.T) {
// 	ctx := context.TODO()
// 	szConfig := getTestObject(ctx, test)
// 	instanceName := "Test name"
// 	verboseLogging := senzing.SzNoLogging
// 	badSettings := "\n\t"
// 	err := szConfig.Initialize(ctx, instanceName, badSettings, verboseLogging)
// 	expectError(test, szerror.ErrSzBase, err)
// }

func TestSzconfig_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	assert.NoError(test, err)
}

func TestSzconfig_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szConfigSingleton = nil
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	assert.NoError(test, err)
}

// TODO: Implement TestSzconfig_Destroy_badXxx
// func TestSzconfig_Destroy_badXxx(test *testing.T) {}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	// return errors.Cast(logger.NewError(errorId, err), err)
	return logger.NewError(errorId, err)
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getSettings() (string, error) {

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
		return "", err
	}
	databaseUrl := fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseUrl}
	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5900, err)
	}
	return settings, err
}

func getSzConfig(ctx context.Context) *Szconfig {
	_ = ctx
	if szConfigSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			fmt.Printf("getSettings() Error: %v\n", err)
			return nil
		}
		szConfigSingleton = &Szconfig{}
		szConfigSingleton.SetLogLevel(ctx, logLevel)
		if logLevel == "TRACE" {
			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
			szConfigSingleton.RegisterObserver(ctx, observerSingleton)
			szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
		}
		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return szConfigSingleton
}

func getSzConfigAsInterface(ctx context.Context) senzing.SzConfig {
	return getSzConfig(ctx)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func getTestObject(ctx context.Context, test *testing.T) *Szconfig {
	_ = test
	return getSzConfig(ctx)
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
	var err error = nil
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfig.IDMessages)
	if err != nil {
		return createError(5901, err)
	}
	osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
	if len(osenvLogLevel) > 0 {
		logLevel = osenvLogLevel
	}
	err = setupDirectories()
	if err != nil {
		return fmt.Errorf("Failed to set up directories. Error: %v", err)
	}
	err = setupDatabase()
	if err != nil {
		return fmt.Errorf("Failed to set up database. Error: %v", err)
	}
	return err
}

func setupDatabase() error {
	var err error = nil

	// Locate source and target paths.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return fmt.Errorf("failed to make target database path (%s) absolute: %w",
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
	var err error = nil
	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", testDirectoryPath, err)
	}
	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzConfig(ctx)
	return err
}

func teardownSzConfig(ctx context.Context) error {
	err := szConfigSingleton.Destroy(ctx)
	if err != nil {
		return err
	}
	szConfigSingleton = nil
	return nil
}
