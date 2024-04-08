package szconfig

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szconfig"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	instanceName      = "SzConfig Test"
	printResults      = false
	verboseLogging    = sz.SZ_NO_LOGGING
)

var (
	globalSzConfig *Szconfig
	logger         logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
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
	databaseUrl := fmt.Sprintf("sqlite3://na:na@%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseUrl}
	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5902, err)
	}
	return settings, err
}

func getSzConfig(ctx context.Context) sz.SzConfig {
	_ = ctx
	if globalSzConfig == nil {
		settings, err := getSettings()
		if err != nil {
			fmt.Printf("getSettings() Error: %v\n", err)
			return nil
		}
		globalSzConfig = &Szconfig{}
		err = globalSzConfig.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return globalSzConfig
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func getTestObject(ctx context.Context, test *testing.T) sz.SzConfig {
	_ = test
	return getSzConfig(ctx)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, szConfig sz.SzConfig, err error) {
	_ = ctx
	_ = szConfig
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
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
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfig.IdMessages)
	if err != nil {
		return createError(5901, err)
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
	err := globalSzConfig.Destroy(ctx)
	if err != nil {
		return err
	}
	globalSzConfig = nil
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzConfig_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzConfig_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzConfig_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	dataSourceCode := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, szConfig, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
}

func TestSzConfig_AddDataSource_WithLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	testError(test, ctx, szConfig, err)
	dataSourceCode := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle2, dataSourceCode)
	testError(test, ctx, szConfig, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	testError(test, ctx, szConfig, err)
}

func TestSzConfig_CloseConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
}

func TestSzConfig_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	printActual(test, actual)
}

func TestSzConfig_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, szConfig, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	printResult(test, "     Add", actual)
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, szConfig, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
}

func TestSzConfig_DeleteDataSource_WithLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, szConfig, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	printResult(test, "     Add", actual)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	testError(test, ctx, szConfig, err)
	err = szConfig.DeleteDataSource(ctx, configHandle2, dataSourceCode)
	testError(test, ctx, szConfig, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle2)
	testError(test, ctx, szConfig, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	testError(test, ctx, szConfig, err)
}

func TestSzConfig_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
}

func TestSzConfig_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	actual, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	printActual(test, actual)
}

func TestSzConfig_ImportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, szConfig, err)
	jsonConfig, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, szConfig, err)
	actual, err := szConfig.ImportConfig(ctx, jsonConfig)
	testError(test, ctx, szConfig, err)
	printActual(test, actual)
}

func TestSzConfig_Initialize(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	instanceName := "Test name"
	verboseLogging := sz.SZ_NO_LOGGING
	settings, err := getSettings()
	testError(test, ctx, szConfig, err)
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	testError(test, ctx, szConfig, err)
}

func TestSzConfig_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	testError(test, ctx, szConfig, err)
}
