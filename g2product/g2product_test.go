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
	printResults      = false
	moduleName        = "Product Test Module"
	verboseLogging    = 0
)

var (
	productInitialized bool      = false
	globalG2product    G2product = G2product{}
	logger             logging.LoggingInterface
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

func setup() error {
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, g2productapi.IdMessages)
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

	// get the database URL and determine if external or a local file just created
	dbUrl, _, err := setupDB(false)
	if err != nil {
		return err
	}

	// get the INI params
	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return err
	}

	// setup the config
	err = setupG2product(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	// get the tem
	return err
}

func teardownG2product(ctx context.Context) error {
	// check if not initialized
	if !productInitialized {
		return nil
	}

	// destroy the G2product
	err := globalG2product.Destroy(ctx)
	if err != nil {
		return err
	}
	productInitialized = false

	return nil
}

func teardown() error {
	ctx := context.TODO()
	err := teardownG2product(ctx)
	return err
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

func setupDB(preserveDB bool) (string, bool, error) {
	var err error = nil

	// get the base directory
	baseDir := baseDirectoryPath()

	// get the template database file path
	dbFilePath := dbTemplatePath()

	dbFilePath, err = filepath.Abs(dbFilePath)
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			dbFilePath, err.Error())
		return "", false, err
	}

	// check the environment for a database URL
	dbUrl, envUrlExists := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")

	dbTargetPath := filepath.Join(baseDirectoryPath(), "G2C.db")

	dbTargetPath, err = filepath.Abs(dbTargetPath)
	if err != nil {
		err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
		return "", false, err
	}

	dbDefaultUrl := fmt.Sprintf("sqlite3://na:na@%s", dbTargetPath)

	dbExternal := envUrlExists && dbDefaultUrl != dbUrl

	if !dbExternal {
		// set the database URL
		dbUrl = dbDefaultUrl

		if !preserveDB {
			// copy the SQLite database file
			_, _, err = futil.CopyFile(dbFilePath, baseDir, true)

			if err != nil {
				err = fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
					dbFilePath, baseDir, err)
				// fall through to return the error
			}
		}
	}

	return dbUrl, dbExternal, err
}

func setupIniParams(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)

	if err != nil {
		err = createError(5902, err)
	}

	return iniParams, err
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

func setupG2product(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
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
	verboseLogging := 0
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

// TODO: Uncomment after fixed
// func TestG2product_ValidateLicenseFile(test *testing.T) {
// 	ctx := context.TODO()
// 	g2product := getTestObject(ctx, test)
// 	licenseFilePath := "testdata/senzing-license/g2.lic"
// 	actual, err := g2product.ValidateLicenseFile(ctx, licenseFilePath)
// 	testErrorNoFail(test, ctx, g2product, err)
// 	printActual(test, actual)
// }

func TestG2product_ValidateLicenseStringBase64(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	licenseString := "AQAAADgCAAAAAAAAU2VuemluZyBQdWJsaWMgVGVzdCBMaWNlbnNlAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARVZBTFVBVElPTiAtIHN1cHBvcnRAc2VuemluZy5jb20AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADIwMjItMTEtMjkAAAAAAAAAAAAARVZBTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNUQU5EQVJEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFDDAAAAAAAAMjAyMy0xMS0yOQAAAAAAAAAAAABNT05USExZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACQfw5e19QAHetkvd+vk0cYHtLaQCLmgx2WUfLorDfLQq15UXmOawNIXc1XguPd8zJtnOaeI6CB2smxVaj10mJE2ndGPZ1JjGk9likrdAj3rw+h6+C/Lyzx/52U8AuaN1kWgErDKdNE9qL6AnnN5LLi7Xs87opP7wbVMOdzsfXx2Xi3H7dSDIam7FitF6brSFoBFtIJac/V/Zc3b8jL/a1o5b1eImQldaYcT4jFrRZkdiVO/SiuLslEb8or3alzT0XsoUJnfQWmh0BjehBK9W74jGw859v/L1SGn1zBYKQ4m8JBiUOytmc9ekLbUKjIg/sCdmGMIYLywKqxb9mZo2TLZBNOpYWVwfaD/6O57jSixfJEHcLx30RPd9PKRO0Nm+4nPdOMMLmd4aAcGPtGMpI6ldTiK9hQyUfrvc9z4gYE3dWhz2Qu3mZFpaAEuZLlKtxaqEtVLWIfKGxwxPargPEfcLsv+30fdjSy8QaHeU638tj67I0uCEgnn5aB8pqZYxLxJx67hvVKOVsnbXQRTSZ00QGX1yTA+fNygqZ5W65wZShhICq5Fz8wPUeSbF7oCcE5VhFfDnSyi5v0YTNlYbF8LOAqXPTi+0KP11Wo24PjLsqYCBVvmOg9ohZ89iOoINwUB32G8VucRfgKKhpXhom47jObq4kSnihxRbTwJRx4o"
	actual, err := g2product.ValidateLicenseStringBase64(ctx, licenseString)
	testErrorNoFail(test, ctx, g2product, err)
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
