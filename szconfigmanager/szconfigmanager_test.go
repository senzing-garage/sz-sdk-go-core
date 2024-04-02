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
	globalSzconfig             szconfig.Szconfig = szconfig.Szconfig{}
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

func getTestObject(ctx context.Context, test *testing.T) sz.SzConfigManager {
	_ = ctx
	_ = test
	return getSzConfigManager(ctx)
}

func getSzConfigManager(ctx context.Context) sz.SzConfigManager {
	_ = ctx
	return &globalSzConfigManager
}

func getSzConfig(ctx context.Context) sz.SzConfig {
	_ = ctx
	return &globalSzconfig
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, szconfigmanager sz.SzConfigManager, err error) {
	_ = ctx
	_ = szconfigmanager
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func baseDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfigmanager")
}

func dbTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
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

	baseDir := baseDirectoryPath()
	err = os.RemoveAll(filepath.Clean(baseDir)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", baseDir, err)
	}
	err = os.MkdirAll(filepath.Clean(baseDir), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", baseDir, err)
	}

	// Get the database URL and determine if external or a local file just created.

	dbUrl, _, err := setupDatabase(false)
	if err != nil {
		return err
	}

	iniParams, err := createSettings(dbUrl)
	if err != nil {
		return err
	}

	err = setupSenzingConfiguration(ctx, instanceName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	err = setupSzConfig(ctx, instanceName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	err = setupSzConfigManager(ctx, instanceName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	return err
}

func setupDatabase(preserveDB bool) (string, bool, error) {
	var err error = nil

	// Get paths.

	baseDir := baseDirectoryPath()
	dbFilePath, err := filepath.Abs(dbTemplatePath())
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			dbFilePath, err.Error())
		return "", false, err
	}
	dbTargetPath := filepath.Join(baseDirectoryPath(), "G2C.db")
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
			_, _, err = futil.CopyFile(dbFilePath, baseDir, true) // Copy the SQLite database file.
			if err != nil {
				err = fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
					dbFilePath, baseDir, err)
				// Fall through to return the error.
			}
		}
	}
	return dbUrl, dbExternal, err
}

func setupSzConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	if szConfigInitialized {
		return fmt.Errorf("SzConfigManager is already setup and has not been torn down")
	}
	globalSzconfig.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalSzconfig.Initialize(ctx, instanceName, settings, verboseLogging)
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
	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := szConfig.Create(ctx)
	if err != nil {
		return createError(5907, err)
	}
	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	configDefinition, err := szConfig.GetJsonString(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = szConfig.Close(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = szConfig.Destroy(ctx)
	if err != nil {
		return createError(5911, err)

	}

	// Persist the Senzing configuration to the Senzing repository.

	szConfigManager := &Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComment := fmt.Sprintf("Created by szconfigmanager_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)

	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}

	err = szConfigManager.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func restoreSzConfigManager(ctx context.Context) error {
	iniParams, err := getSettings()
	if err != nil {
		return err
	}

	err = setupSzConfigManager(ctx, instanceName, iniParams, verboseLogging)
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
	err := globalSzconfig.Destroy(ctx)
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
	configHandle, err1 := szConfig.Create(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfig.Create()")
	}
	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), 10)
	_, err2 := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "scConfig.AddDataSource()")
	}
	configDefinition, err3 := szConfig.GetJsonString(ctx, configHandle)
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
