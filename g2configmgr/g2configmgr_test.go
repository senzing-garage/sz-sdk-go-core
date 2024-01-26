package g2configmgr

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
	"github.com/senzing-garage/g2-sdk-go-base/g2config"
	"github.com/senzing-garage/g2-sdk-go/g2api"
	g2configmgrapi "github.com/senzing-garage/g2-sdk-go/g2configmgr"
	"github.com/senzing-garage/g2-sdk-go/g2error"
	futil "github.com/senzing-garage/go-common/fileutil"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
	"github.com/senzing-garage/go-common/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	moduleName        = "Config Manager Test Module"
	printResults      = false
	verboseLogging    = 0
)

var (
	configInitialized    bool              = false
	configMgrInitialized bool              = false
	globalG2config       g2config.G2config = g2config.G2config{}
	globalG2configmgr    G2configmgr       = G2configmgr{}
	logger               logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(logger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) g2api.G2configmgr {
	return &globalG2configmgr
}

func getG2Configmgr(ctx context.Context) g2api.G2configmgr {
	return &globalG2configmgr
}

func getG2Config(ctx context.Context) g2api.G2config {
	return &globalG2config
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

func testError(test *testing.T, ctx context.Context, g2configmgr g2api.G2configmgr, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func baseDirectoryPath() string {
	return filepath.FromSlash("../target/test/g2configmgr")
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
		if g2error.Is(err, g2error.G2Unrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2Retryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2BadInput) {
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

func dbUrl() (string, bool, error) {
	var err error = nil

	// Get paths.

	baseDir := baseDirectoryPath()
	dbFilePath, err := filepath.Abs(dbTemplatePath())
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			dbFilePath, err.Error())
		return "", false, err
	}
	dbTargetPath := filepath.Join(baseDir, "G2C.db")
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
	}

	return dbUrl, dbExternal, err
}

func getIniParams() (string, error) {
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
	logger, err = logging.NewSenzingSdkLogger(ComponentId, g2configmgrapi.IdMessages)
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

func setupG2configmgr(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	if configMgrInitialized {
		return fmt.Errorf("G2configmgr is already setup and has not been torn down")
	}

	globalG2configmgr.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	configMgrInitialized = true
	return err
}

func setupG2config(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	if configInitialized {
		return fmt.Errorf("G2configmgr is already setup and has not been torn down")
	}
	globalG2config.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalG2config.Init(ctx, moduleName, iniParams, verboseLogging)
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	now := time.Now()
	aG2config := &g2config.G2config{}
	err := aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return createError(5907, err)
	}
	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := aG2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return createError(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
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

	aG2configmgr := &G2configmgr{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)

	if err != nil {
		return createError(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
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
	iniParams, err := getIniParams()
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
	err := globalG2config.Destroy(ctx)
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
	err := globalG2configmgr.Destroy(ctx)
	if err != nil {
		return err
	}
	configMgrInitialized = false
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2configmgr_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
}

func TestG2configmgr_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	actual := g2configmgr.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2configmgr_AddConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	now := time.Now()
	g2config := getG2Config(ctx)
	configHandle, err1 := g2config.Create(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2config.Create()")
	}
	inputJson := `{"DSRC_CODE": "GO_TEST_` + strconv.FormatInt(now.Unix(), 10) + `"}`
	_, err2 := g2config.AddDataSource(ctx, configHandle, inputJson)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "g2config.AddDataSource()")
	}
	configStr, err3 := g2config.Save(ctx, configHandle)
	if err3 != nil {
		test.Log("Error:", err3.Error())
		assert.FailNow(test, configStr)
	}
	configComments := fmt.Sprintf("g2configmgr_test at %s", now.UTC())
	actual, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_GetConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	configID, err1 := g2configmgr.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()")
	}
	actual, err := g2configmgr.GetConfig(ctx, configID)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	actual, err := g2configmgr.GetConfigList(ctx)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	actual, err := g2configmgr.GetDefaultConfigID(ctx)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	oldConfigID, err1 := g2configmgr.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()")
	}

	// FIXME: This is kind of a cheater.

	newConfigID, err2 := g2configmgr.GetDefaultConfigID(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()-2")
	}

	err := g2configmgr.ReplaceDefaultConfigID(ctx, oldConfigID, newConfigID)
	testError(test, ctx, g2configmgr, err)
}

func TestG2configmgr_SetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	configID, err1 := g2configmgr.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()")
	}
	err := g2configmgr.SetDefaultConfigID(ctx, configID)
	testError(test, ctx, g2configmgr, err)
}

func TestG2configmgr_Init(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := int64(0)
	iniParams, err := getIniParams()
	testError(test, ctx, g2configmgr, err)
	err = g2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2configmgr, err)
}

func TestG2configmgr_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	err := g2configmgr.Destroy(ctx)
	testError(test, ctx, g2configmgr, err)

	// restore the state that existed prior to this test
	configMgrInitialized = false
	restoreG2configmgr(ctx)
}
