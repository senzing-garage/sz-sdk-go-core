package szabstractfactory

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
	instanceName      = "SzAbstractFactory Test"
	printResults      = false
	verboseLogging    = sz.SZ_NO_LOGGING
)

var (
	logger logging.LoggingInterface
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

func getTestObject(ctx context.Context, test *testing.T) sz.SzAbstractFactory {
	_ = test
	return getSzAbstractFactory(ctx)
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

func getSzAbstractFactory(ctx context.Context) sz.SzAbstractFactory {
	_ = ctx

	settings, err := getSettings()
	if err != nil {
		return nil
	}
	result := &Szabstractfactory{
		ConfigId:       sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, szAbstractFactory sz.SzAbstractFactory, err error) {
	_ = ctx
	_ = szAbstractFactory
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
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzAbstractFactory_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	config, err := szAbstractFactory.CreateConfig(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer config.Destroy(ctx)
	configHandle, err := config.CreateConfig(ctx)
	testError(test, ctx, szAbstractFactory, err)
	dataSources, err := config.GetDataSources(ctx, configHandle)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, dataSources)
}

func TestSzAbstractFactory_CreateConfigManager(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	configManager, err := szAbstractFactory.CreateConfigManager(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer configManager.Destroy(ctx)
	configList, err := configManager.GetConfigList(ctx)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, configList)
}

func TestSzAbstractFactory_CreateDiagnostic(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	diagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer diagnostic.Destroy(ctx)
	result, err := diagnostic.CheckDatabasePerformance(ctx, 1)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, result)
}

func TestSzAbstractFactory_CreateEngine(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	engine, err := szAbstractFactory.CreateEngine(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer engine.Destroy(ctx)
	stats, err := engine.GetStats(ctx)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, stats)
}

func TestSzAbstractFactory_CreateProduct(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	product, err := szAbstractFactory.CreateProduct(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer product.Destroy(ctx)
	version, err := product.GetVersion(ctx)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, version)
}
