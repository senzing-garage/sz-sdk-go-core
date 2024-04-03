package szconfigmanager

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	futil "github.com/senzing-garage/go-common/fileutil"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szconfigmanagerapi "github.com/senzing-garage/sz-sdk-go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	instanceName      = "Config Manager Test Module"
	printResults      = false
	verboseLogging    = sz.SZ_NO_LOGGING
)

var (
	globalSzConfig             szconfig.Szconfig = szconfig.Szconfig{}
	globalSzConfigManager      Szconfigmanager   = Szconfigmanager{}
	logger                     logging.LoggingInterface
	szConfigInitialized        bool = false
	szConfigManagerInitialized bool = false
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

func getSzConfigManager(ctx context.Context) sz.SzConfigManager {
	_ = ctx
	return &globalSzConfigManager
}

func getSzConfig(ctx context.Context) sz.SzConfig {
	_ = ctx
	return &globalSzConfig
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfigmanager")
}

func getTestObject(ctx context.Context, test *testing.T) sz.SzConfigManager {
	_ = ctx
	_ = test
	return getSzConfigManager(ctx)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, szConfigManager sz.SzConfigManager, err error) {
	_ = ctx
	_ = szConfigManager
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

func createSettings(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}
	settings, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5902, err)
	}
	return settings, err
}

func getSettings() (string, error) {
	dbUrl, _, err := setupDatabase(true)
	if err != nil {
		return "", err
	}
	settings, err := createSettings(dbUrl)
	if err != nil {
		return "", err
	}
	return settings, nil
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfigmanagerapi.IdMessages)
	if err != nil {
		return createError(5901, err)
	}

	// Cleanup past runs and prepare for current run.

	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", testDirectoryPath, err)
	}

	// Get the database URL and determine if external or a local file just created.

	dbUrl, _, err := setupDatabase(false)
	if err != nil {
		return err
	}

	// Create the Senzing engine configuration JSON.

	settings, err := createSettings(dbUrl)
	if err != nil {
		return err
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfiguration(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return err
	}

	err = setupSzConfig(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return err
	}

	err = setupSzConfigManager(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return err
	}

	return err
}

func setupDatabase(preserveDB bool) (string, bool, error) {
	var err error = nil

	// Get paths.

	testDirectoryPath := getTestDirectoryPath()
	dbFilePath, err := filepath.Abs(getDatabaseTemplatePath())
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			dbFilePath, err.Error())
		return "", false, err
	}
	dbTargetPath := filepath.Join(getTestDirectoryPath(), "G2C.db")
	dbTargetPath, err = filepath.Abs(dbTargetPath)
	if err != nil {
		err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
		return "", false, err
	}

	// Check the environment for a database URL.

	dbUrl, envUrlExists := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")
	dbDefaultUrl := fmt.Sprintf("sqlite3://na:na@%s", dbTargetPath)
	dbExternal := envUrlExists && dbDefaultUrl != dbUrl
	if !dbExternal {
		dbUrl = dbDefaultUrl
		if !preserveDB {
			_, _, err = futil.CopyFile(dbFilePath, testDirectoryPath, true) // Copy the SQLite database file.
			if err != nil {
				err = fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
					dbFilePath, testDirectoryPath, err)
				// Fall through to return the error.
			}
		}
	}
	return dbUrl, dbExternal, err
}

func setupSzConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	if szConfigInitialized {
		return fmt.Errorf("SzConfig is already setup and has not been torn down")
	}
	globalSzConfig.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalSzConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	szConfigInitialized = true
	return err
}

func setupSzConfigManager(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	if szConfigManagerInitialized {
		return fmt.Errorf("SzConfigManager is already setup and has not been torn down")
	}

	globalSzConfigManager.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalSzConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	szConfigManagerInitialized = true
	return err
}

func setupSenzingConfiguration(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	now := time.Now()

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
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

	szConfigManager := &Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}
	defer szConfigManager.Destroy(ctx)

	configComment := fmt.Sprintf("Created by szconfigmanager_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}

	return err
}

func restoreSzConfigManager(ctx context.Context) error {
	settings, err := getSettings()
	if err != nil {
		return err
	}

	err = setupSzConfigManager(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return err
	}

	return nil
}

func teardown() error {
	var resultErr error = nil
	ctx := context.TODO()
	err := teardownSzConfig(ctx)
	if err != nil {
		fmt.Println(err)
		resultErr = err
	}
	teardownSzConfigManager(ctx)
	if err != nil {
		fmt.Println(err)
		resultErr = err
	}
	return resultErr
}

func teardownSzConfig(ctx context.Context) error {
	if !szConfigInitialized {
		return nil
	}
	err := globalSzConfig.Destroy(ctx)
	if err != nil {
		return err
	}
	szConfigInitialized = false
	return nil
}

func teardownSzConfigManager(ctx context.Context) error {
	if !szConfigManagerInitialized {
		return nil
	}
	err := globalSzConfigManager.Destroy(ctx)
	if err != nil {
		return err
	}
	szConfigManagerInitialized = false
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzConfigManager_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szconfigmanager.SetObserverOrigin(ctx, origin)
}

func TestSzConfigManager_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szconfigmanager.SetObserverOrigin(ctx, origin)
	actual := szconfigmanager.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzConfigManager_AddConfig(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	now := time.Now()
	szConfig := getSzConfig(ctx)
	configHandle, err1 := szConfig.CreateConfig(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfig.CreateConfig()")
	}
	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), 10)
	_, err2 := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "scConfig.AddDataSource()")
	}
	configDefinition, err3 := szConfig.ExportConfig(ctx, configHandle)
	if err3 != nil {
		test.Log("Error:", err3.Error())
		assert.FailNow(test, configDefinition)
	}
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szconfigmanager.AddConfig(ctx, configDefinition, configComment)
	testError(test, ctx, szconfigmanager, err)
	printActual(test, actual)
}

func TestSzConfigManager_GetConfig(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	configId, err1 := szconfigmanager.GetDefaultConfigId(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szconfigmanager.GetDefaultConfigId()")
	}
	actual, err := szconfigmanager.GetConfig(ctx, configId)
	testError(test, ctx, szconfigmanager, err)
	printActual(test, actual)
}

func TestSzConfigManager_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	actual, err := szconfigmanager.GetConfigList(ctx)
	testError(test, ctx, szconfigmanager, err)
	printActual(test, actual)
}

func TestSzConfigManager_GetDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	actual, err := szconfigmanager.GetDefaultConfigId(ctx)
	testError(test, ctx, szconfigmanager, err)
	printActual(test, actual)
}

func TestSzConfigManager_ReplaceDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	currentDefaultConfigId, err1 := szconfigmanager.GetDefaultConfigId(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szconfigmanager.GetDefaultConfigId()")
	}

	// FIXME: This is kind of a cheater.

	newDefaultConfigId, err2 := szconfigmanager.GetDefaultConfigId(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szconfigmanager.GetDefaultConfigId()-2")
	}

	err := szconfigmanager.ReplaceDefaultConfigId(ctx, currentDefaultConfigId, newDefaultConfigId)
	testError(test, ctx, szconfigmanager, err)
}

func TestSzConfigManager_SetDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	configId, err1 := szconfigmanager.GetDefaultConfigId(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szconfigmanager.GetDefaultConfigId()")
	}
	err := szconfigmanager.SetDefaultConfigId(ctx, configId)
	testError(test, ctx, szconfigmanager, err)
}

func TestSzConfigManager_Initialize(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	instanceName := "Test module name"
	verboseLogging := int64(0)
	settings, err := getSettings()
	testError(test, ctx, szconfigmanager, err)
	err = szconfigmanager.Initialize(ctx, instanceName, settings, verboseLogging)
	testError(test, ctx, szconfigmanager, err)
}

func TestSzConfigManager_Destroy(test *testing.T) {
	ctx := context.TODO()
	szconfigmanager := getTestObject(ctx, test)
	err := szconfigmanager.Destroy(ctx)
	testError(test, ctx, szconfigmanager, err)

	// restore the state that existed prior to this test
	szConfigManagerInitialized = false
	restoreSzConfigManager(ctx)
}
