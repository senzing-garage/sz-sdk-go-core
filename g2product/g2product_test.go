package g2product

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/g2-sdk-go/g2error"
	g2productapi "github.com/senzing/g2-sdk-go/g2product"
	futil "github.com/senzing/go-common/fileutil"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logging"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	moduleName        = "Product Test Module"
	printResults      = false
	verboseLogging    = 0
)

var (
	globalG2product    G2product = G2product{}
	logger             logging.LoggingInterface
	productInitialized bool = false
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(logger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) g2api.G2product {
	return &globalG2product
}

func getG2Product(ctx context.Context) g2api.G2product {
	return &globalG2product
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

func testError(test *testing.T, ctx context.Context, g2product g2api.G2product, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2product g2api.G2product, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
	}
}

func baseDirectoryPath() string {
	return filepath.FromSlash("../target/test/g2product")
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

func restoreG2product(ctx context.Context) error {
	iniParams, err := getIniParams()
	if err != nil {
		return err
	}

	err = setupG2product(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}
	return nil
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, g2productapi.IdMessages)
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

	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return err
	}

	err = setupG2product(ctx, moduleName, iniParams, verboseLogging)
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

func setupG2product(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	if productInitialized {
		return fmt.Errorf("G2product is already setup and has not been torn down")
	}
	globalG2product.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalG2product.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	productInitialized = true
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

func teardown() error {
	ctx := context.TODO()
	err := teardownG2product(ctx)
	return err
}

func teardownG2product(ctx context.Context) error {
	if !productInitialized {
		return nil
	}
	err := globalG2product.Destroy(ctx)
	if err != nil {
		return err
	}
	productInitialized = false
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2product_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
}

func TestG2product_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
	actual := g2product.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2product_Init(test *testing.T) {
	ctx := context.TODO()
	g2product := &G2product{}
	moduleName := "Test module name"
	verboseLogging := int64(0)
	iniParams, err := getIniParams()
	testError(test, ctx, g2product, err)
	err = g2product.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2product, err)
}

func TestG2product_License(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	actual, err := g2product.License(ctx)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2product_Version(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	actual, err := g2product.Version(ctx)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2product_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	err := g2product.Destroy(ctx)
	testError(test, ctx, g2product, err)

	// restore the pre-test state
	productInitialized = false
	restoreG2product(ctx)
}
