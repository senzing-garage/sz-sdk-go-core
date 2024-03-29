package szconfig

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/g2-sdk-go/szapi"
	szconfigapi "github.com/senzing-garage/g2-sdk-go/szconfig"
	"github.com/senzing-garage/g2-sdk-go/szerror"
	futil "github.com/senzing-garage/go-common/fileutil"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	moduleName        = "Config Test Module"
	printResults      = false
	verboseLogging    = 0
)

var (
	configInitialized bool     = false
	globalSzconfig    Szconfig = Szconfig{}
	logger            logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) szapi.G2config {
	_ = ctx
	_ = test
	return &globalSzconfig
}

func getSzconfig(ctx context.Context) szapi.G2config {
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

func testError(test *testing.T, ctx context.Context, szconfig szapi.G2config, err error) {
	_ = ctx
	_ = szconfig
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func baseDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
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
	settings, err := createSettings(dbUrl)
	if err != nil {
		return "", err
	}
	return settings, nil
}

func restoreG2config(ctx context.Context) error {
	settings, err := getSettings()
	if err != nil {
		return err
	}

	err = setupG2config(ctx, moduleName, settings, verboseLogging)
	if err != nil {
		return err
	}
	return nil
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfigapi.IdMessages)
	if err != nil {
		return createError(5901, err)
	}

	// Cleanup past runs and prepare for current run.

	baseDir := baseDirectoryPath()
	err = os.RemoveAll(filepath.Clean(baseDir))
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", baseDir, err)
	}
	err = os.MkdirAll(filepath.Clean(baseDir), 0750)
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", baseDir, err)
	}

	// Get the database URL and determine if external or a local file just created.

	dbUrl, _, err := setupDB(false)
	if err != nil {
		return err
	}

	// Create the Senzing engine configuration JSON.

	settings, err := createSettings(dbUrl)
	if err != nil {
		return err
	}

	err = setupG2config(ctx, moduleName, settings, verboseLogging)
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

func setupG2config(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	if configInitialized {
		return fmt.Errorf("G2config is already setup and has not been torn down")
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

func createSettings(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}
	settings, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5902, err)
	}
	return settings, err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownG2config(ctx)
	return err
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

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2config_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szconfig.SetObserverOrigin(ctx, origin)
}

func TestG2config_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szconfig.SetObserverOrigin(ctx, origin)
	actual := szconfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2config_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	dataSourceDefinition := `{"DSRC_CODE": "GO_TEST"}`
	actual, err := szconfig.AddDataSource(ctx, configHandle, dataSourceDefinition)
	testError(test, ctx, szconfig, err)
	printActual(test, actual)
	err = szconfig.Close(ctx, configHandle)
	testError(test, ctx, szconfig, err)
}

func TestG2config_AddDataSource_WithLoad(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	configDefinition, err := szconfig.GetJsonString(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	err = szconfig.Close(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	configHandle2, err := szconfig.Load(ctx, configDefinition)
	testError(test, ctx, szconfig, err)
	dataSourceDefinition := `{"DSRC_CODE": "GO_TEST"}`
	actual, err := szconfig.AddDataSource(ctx, configHandle2, dataSourceDefinition)
	testError(test, ctx, szconfig, err)
	printActual(test, actual)
	err = szconfig.Close(ctx, configHandle2)
	testError(test, ctx, szconfig, err)
}

func TestG2config_Close(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	err = szconfig.Close(ctx, configHandle)
	testError(test, ctx, szconfig, err)
}

func TestG2config_Create(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	actual, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	printActual(test, actual)
}

func TestG2config_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	actual, err := szconfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	dataSourceDefinition := `{"DSRC_CODE": "` + dataSourceCode + `"}`
	_, err = szconfig.AddDataSource(ctx, configHandle, dataSourceDefinition)
	testError(test, ctx, szconfig, err)
	actual, err = szconfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	printResult(test, "     Add", actual)
	err = szconfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, szconfig, err)
	actual, err = szconfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	printResult(test, "  Delete", actual)
	err = szconfig.Close(ctx, configHandle)
	testError(test, ctx, szconfig, err)
}

func TestG2config_DeleteDataSource_WithLoad(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	actual, err := szconfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	dataSourceDefinition := `{"DSRC_CODE": "` + dataSourceCode + `"}`
	_, err = szconfig.AddDataSource(ctx, configHandle, dataSourceDefinition)
	testError(test, ctx, szconfig, err)
	actual, err = szconfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	printResult(test, "     Add", actual)
	configDefinition, err := szconfig.GetJsonString(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	err = szconfig.Close(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	configHandle2, err := szconfig.Load(ctx, configDefinition)
	testError(test, ctx, szconfig, err)
	err = szconfig.DeleteDataSource(ctx, configHandle2, dataSourceCode)
	testError(test, ctx, szconfig, err)
	actual, err = szconfig.GetDataSources(ctx, configHandle2)
	testError(test, ctx, szconfig, err)
	printResult(test, "  Delete", actual)
	err = szconfig.Close(ctx, configHandle2)
	testError(test, ctx, szconfig, err)
}

func TestG2config_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	actual, err := szconfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	printActual(test, actual)
	err = szconfig.Close(ctx, configHandle)
	testError(test, ctx, szconfig, err)
}

func TestG2config_GetJsonString(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	actual, err := szconfig.GetJsonString(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	printActual(test, actual)
}

func TestG2config_Load(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	configHandle, err := szconfig.Create(ctx)
	testError(test, ctx, szconfig, err)
	jsonConfig, err := szconfig.GetJsonString(ctx, configHandle)
	testError(test, ctx, szconfig, err)
	actual, err := szconfig.Load(ctx, jsonConfig)
	testError(test, ctx, szconfig, err)
	printActual(test, actual)
}

func TestG2config_Initialize(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	instanceName := "Test module name"
	verboseLogging := int64(0)
	settings, err := getSettings()
	testError(test, ctx, szconfig, err)
	err = szconfig.Initialize(ctx, instanceName, settings, verboseLogging)
	testError(test, ctx, szconfig, err)
}

func TestG2config_Destroy(test *testing.T) {
	ctx := context.TODO()
	szconfig := getTestObject(ctx, test)
	err := szconfig.Destroy(ctx)
	testError(test, ctx, szconfig, err)

	// restore the state that existed prior to this test
	configInitialized = false
	restoreG2config(ctx)
}
