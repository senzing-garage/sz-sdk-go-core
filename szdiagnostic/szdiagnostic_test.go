package szdiagnostic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	instanceName      = "Diagnostic Test"
	printResults      = false
	verboseLogging    = sz.SZ_NO_LOGGING
)

var (
	defaultConfigId    int64
	globalSzDiagnostic *Szdiagnostic
	logger             logging.LoggingInterface
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

func getDefaultConfigId() int64 {
	return defaultConfigId
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

func getSzDiagnostic(ctx context.Context) sz.SzDiagnostic {
	_ = ctx
	if globalSzDiagnostic == nil {
		settings, err := getSettings()
		if err != nil {
			fmt.Printf("getSettings() Error: %v\n", err)
			return nil
		}
		globalSzDiagnostic = &Szdiagnostic{}
		err = globalSzDiagnostic.Initialize(ctx, instanceName, settings, verboseLogging, getDefaultConfigId())
		if err != nil {
			fmt.Println(err)
		}
	}
	return globalSzDiagnostic
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szdiagnostic")
}

func getTestObject(ctx context.Context, test *testing.T) sz.SzDiagnostic {
	_ = test
	return getSzDiagnostic(ctx)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, szDiagnostic sz.SzDiagnostic, err error) {
	_ = ctx
	_ = szDiagnostic
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, szDiagnostic sz.SzDiagnostic, err error) {
	_ = ctx
	_ = szDiagnostic
	if err != nil {
		test.Log("Error:", err.Error())
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
		if szerror.Is(err, szerror.SzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzBadInput) {
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
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szdiagnostic.IdMessages)
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
	err = setupSenzingConfiguration()
	if err != nil {
		return createError(5920, err)
	}
	err = setupAddRecords()
	if err != nil {
		return createError(5922, err)
	}
	return err
}

func setupAddRecords() error {
	ctx := context.TODO()

	settings, err := getSettings()
	if err != nil {
		return createError(9999, err)
	}

	// Create sz objects.

	szEngine := &szengine.Szengine{}
	err = szEngine.Initialize(ctx, instanceName, settings, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
	if err != nil {
		return createError(5916, err)
	}
	defer szEngine.Destroy(ctx)

	szDiagnostic := getSzDiagnostic(ctx)
	err = szDiagnostic.Initialize(ctx, instanceName, settings, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
	if err != nil {
		return createError(5916, err)
	}
	defer szDiagnostic.Destroy(ctx)

	// Add records into Senzing.

	testRecordIds := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	for _, testRecordId := range testRecordIds {
		testRecord := truthset.CustomerRecords[testRecordId]
		_, err := szEngine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Json, sz.SZ_NO_FLAGS)
		if err != nil {
			return createError(5917, err)
		}
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

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	settings, err := getSettings()
	if err != nil {
		return createError(9999, err)
	}

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}
	defer szConfig.Destroy(ctx)

	// Create an in memory Senzing configuration.

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	// Create a string representation of the in-memory configuration.

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	// Close szConfig in-memory object.

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	// Persist the Senzing configuration to the Senzing repository as default.

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}
	defer szConfigManager.Destroy(ctx)

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}
	defaultConfigId = configId

	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzDiagnostic(ctx)
	return err
}

func teardownSzDiagnostic(ctx context.Context) error {
	err := globalSzDiagnostic.Destroy(ctx)
	if err != nil {
		return err
	}
	globalSzDiagnostic = nil
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzDiagnostic_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
}

func TestSzDiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzDiagnostic_CheckDatabasePerformance(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatabasePerformance(ctx, secondsToRun)
	testError(test, ctx, szDiagnostic, err)
	printActual(test, actual)
}

func TestSzDiagnostic_Initialize(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := &Szdiagnostic{}
	instanceName := "Test module name"
	settings, err := getSettings()
	testError(test, ctx, szDiagnostic, err)
	verboseLogging := sz.SZ_NO_LOGGING
	configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
	err = szDiagnostic.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	testError(test, ctx, szDiagnostic, err)
}

func TestSzDiagnostic_Initialize_WithConfigId(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := &Szdiagnostic{}
	instanceName := "Test module name"
	settings, err := getSettings()
	testError(test, ctx, szDiagnostic, err)
	verboseLogging := sz.SZ_NO_LOGGING
	configId := getDefaultConfigId()
	err = szDiagnostic.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	testError(test, ctx, szDiagnostic, err)
}

func TestSzDiagnostic_Reinit(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	configId := getDefaultConfigId()
	err := szDiagnostic.Reinitialize(ctx, configId)
	testErrorNoFail(test, ctx, szDiagnostic, err)
}

func TestSzDiagnostic_Destroy(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.Destroy(ctx)
	testError(test, ctx, szDiagnostic, err)
}
