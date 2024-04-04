package szproduct

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	futil "github.com/senzing-garage/go-common/fileutil"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szproductapi "github.com/senzing-garage/sz-sdk-go/szproduct"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	instanceName      = "Product Test Module"
	printResults      = false
	verboseLogging    = sz.SZ_NO_LOGGING
)

var (
	globalSzProduct      Szproduct = Szproduct{}
	logger               logging.LoggingInterface
	szProductInitialized bool = false
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

func getSzProduct(ctx context.Context) sz.SzProduct {
	_ = ctx
	return &globalSzProduct
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szproduct")
}

func getTestObject(ctx context.Context, test *testing.T) sz.SzProduct {
	_ = ctx
	_ = test
	return getSzProduct(ctx)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, szProduct sz.SzProduct, err error) {
	_ = ctx
	_ = szProduct
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

// func testErrorNoFail(test *testing.T, ctx context.Context, szProduct sz.SzProduct, err error) {
// 	if err != nil {
// 		test.Log("Error:", err.Error())
// 	}
// }

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

func restoreSzProduct(ctx context.Context) error {
	settings, err := getSettings()
	if err != nil {
		return err
	}

	err = setupSzProduct(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return err
	}
	return nil
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szproductapi.IdMessages)
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

	err = setupSzProduct(ctx, instanceName, settings, verboseLogging)
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

func setupSzProduct(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	if szProductInitialized {
		return fmt.Errorf("SzProduct is already setup and has not been torn down")
	}
	globalSzProduct.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalSzProduct.Initialize(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	szProductInitialized = true
	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzProduct(ctx)
	return err
}

func teardownSzProduct(ctx context.Context) error {
	if !szProductInitialized {
		return nil
	}
	err := globalSzProduct.Destroy(ctx)
	if err != nil {
		return err
	}
	szProductInitialized = false
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzProduct_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
}

func TestSzProduct_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	actual := szProduct.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzProduct_Initialize(test *testing.T) {
	ctx := context.TODO()
	szProduct := &Szproduct{}
	instanceName := "Test name"
	settings, err := getSettings()
	testError(test, ctx, szProduct, err)
	err = szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	testError(test, ctx, szProduct, err)
}

func TestSzProduct_GetLicense(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetLicense(ctx)
	testError(test, ctx, szProduct, err)
	printActual(test, actual)
}

func TestSzProduct_GetVersion(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetVersion(ctx)
	testError(test, ctx, szProduct, err)
	printActual(test, actual)
}

func TestSzProduct_Destroy(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.Destroy(ctx)
	testError(test, ctx, szProduct, err)

	// restore the pre-test state
	szProductInitialized = false
	restoreSzProduct(ctx)
}
