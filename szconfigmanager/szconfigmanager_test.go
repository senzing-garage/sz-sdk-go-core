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
	szconfigmanagerapi "github.com/senzing-garage/sz-sdk-go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/senzing-garage/sz-sdk-go/szinterface"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	moduleName        = "Config Manager Test Module"
	printResults      = false
	verboseLogging    = 0
)

var (
	configInitialized     bool              = false
	configMgrInitialized  bool              = false
	globalSzconfig        szconfig.Szconfig = szconfig.Szconfig{}
	globalSzConfigManager SzConfigManager   = SzConfigManager{}
	logger                logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) szinterface.SzConfigManager {
	_ = ctx
	_ = test
	return &globalSzConfigManager
}

func getSzConfigManager(ctx context.Context) szinterface.SzConfigManager {
	_ = ctx
	return &globalSzConfigManager
}

func getSzconfig(ctx context.Context) szinterface.SzConfig {
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

func testError(test *testing.T, ctx context.Context, szconfigmanager szinterface.SzConfigManager, err error) {
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

func getSettings() (string, error) {
	dbUrl, _, err := setupDB(true)
	if err != nil {
		return "", err
	}
	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return "", err
	}
	return iniParams, nil
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfigmanagerapi.IdMessages)
	if err != nil {
		return createError(5901, err)
	}

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

	dbUrl, _, err := setupDB(false)
	if err != nil {
		return err
	}

	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return err
	}

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	err = setupG2config(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	err = setupG2configmgr(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	return err
}

func setupDB(preserveDB bool) (string, bool, error) {
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

func setupG2configmgr(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	if configMgrInitialized {
		return fmt.Errorf("G2configmgr is already setup and has not been torn down")
	}

	globalSzConfigManager.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalSzConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	configMgrInitialized = true
	return err
}

func setupG2config(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	if configInitialized {
		return fmt.Errorf("G2configmgr is already setup and has not been torn down")
	}
	globalSzconfig.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalSzconfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	configInitialized = true
	return err
}

func setupIniParams(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5902, err)
	}
	return iniParams, err
}

func setupSenzingConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	now := time.Now()
	aG2config := &szconfig.Szconfig{}
	err := aG2config.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return createError(5907, err)
	}
	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := aG2config.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	configStr, err := aG2config.GetJsonString(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return createError(5911, err)

	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &SzConfigManager{}
	err = aG2configmgr.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configId, err := aG2configmgr.AddConfig(ctx, configStr, configComments)

	if err != nil {
		return createError(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func restoreG2configmgr(ctx context.Context) error {
	iniParams, err := getSettings()
	if err != nil {
		return err
	}

	err = setupG2configmgr(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	return nil
}

func teardown() error {
	var resultErr error = nil
	ctx := context.TODO()
	err := teardownG2config(ctx)
	if err != nil {
		fmt.Println(err)
		resultErr = err
	}
	teardownG2configmgr(ctx)
	if err != nil {
		fmt.Println(err)
		resultErr = err
	}
	return resultErr
}

func teardownG2config(ctx context.Context) error {
	if !configInitialized {
		return nil
	}
	err := globalSzconfig.Destroy(ctx)
	if err != nil {
		return err
	}
	configInitialized = false
	return nil
}

func teardownG2configmgr(ctx context.Context) error {
	if !configMgrInitialized {
		return nil
	}
	err := globalSzConfigManager.Destroy(ctx)
	if err != nil {
		return err
	}
	configMgrInitialized = false
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
	g2config := getSzconfig(ctx)
	configHandle, err1 := g2config.Create(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2config.Create()")
	}
	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), 10)
	_, err2 := g2config.AddDataSource(ctx, configHandle, dataSourceCode)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "g2config.AddDataSource()")
	}
	configDefinition, err3 := g2config.GetJsonString(ctx, configHandle)
	if err3 != nil {
		test.Log("Error:", err3.Error())
		assert.FailNow(test, configDefinition)
	}
	configComments := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szconfigmanager.AddConfig(ctx, configDefinition, configComments)
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
	configMgrInitialized = false
	restoreG2configmgr(ctx)
}