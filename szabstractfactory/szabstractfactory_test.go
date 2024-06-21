package szabstractfactory

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szconfig"
	"github.com/stretchr/testify/require"
)

const (
	baseCallerSkip    = 4
	defaultTruncation = 76
	instanceName      = "SzAbstractFactory Test"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

var (
	logger logging.Logging
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzAbstractFactory_CreateSzConfig(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szConfig, err := szAbstractFactory.CreateSzConfig(ctx)
	require.NoError(test, err)
	defer func() { handleError(szConfig.Destroy(ctx)) }()
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	dataSources, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, dataSources)
}

func TestSzAbstractFactory_CreateSzConfigManager(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szConfigManager, err := szAbstractFactory.CreateSzConfigManager(ctx)
	require.NoError(test, err)
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()
	configList, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, configList)
}

func TestSzAbstractFactory_CreateSzDiagnostic(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szDiagnostic, err := szAbstractFactory.CreateSzDiagnostic(ctx)
	require.NoError(test, err)
	defer func() { handleError(szDiagnostic.Destroy(ctx)) }()
	result, err := szDiagnostic.CheckDatastorePerformance(ctx, 1)
	require.NoError(test, err)
	printActual(test, result)
}

func TestSzAbstractFactory_CreateSzEngine(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szEngine, err := szAbstractFactory.CreateSzEngine(ctx)
	require.NoError(test, err)
	defer func() { handleError(szEngine.Destroy(ctx)) }()
	stats, err := szEngine.GetStats(ctx)
	require.NoError(test, err)
	printActual(test, stats)
}

func TestSzAbstractFactory_CreateSzProduct(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szProduct, err := szAbstractFactory.CreateSzProduct(ctx)
	require.NoError(test, err)
	defer func() { handleError(szProduct.Destroy(ctx)) }()
	version, err := szProduct.GetVersion(ctx)
	require.NoError(test, err)
	printActual(test, version)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getSettings() (string, error) {
	var result string

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return result, fmt.Errorf("failed to make target database path (%s) absolute. Error: %w", dbTargetPath, err)
	}
	databaseURL := fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	if err != nil {
		return result, fmt.Errorf("failed to BuildSimpleSettingsUsingMap(%s) Error: %w", configAttrMap, err)
	}
	return result, err
}

func getSzAbstractFactory(ctx context.Context) (senzing.SzAbstractFactory, error) {
	var err error
	var result senzing.SzAbstractFactory
	_ = ctx
	settings, err := getSettings()
	if err != nil {
		return result, err
	}
	result = &Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}
	return result, err
}

func getTestObject(ctx context.Context, test *testing.T) senzing.SzAbstractFactory {
	result, err := getSzAbstractFactory(ctx)
	require.NoError(test, err)
	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
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
	var err error
	logger = helper.GetLogger(ComponentID, szconfig.IDMessages, baseCallerSkip)
	err = setupDirectories()
	if err != nil {
		return fmt.Errorf("failed to set up directories. Error: %w", err)
	}
	err = setupDatabase()
	if err != nil {
		return fmt.Errorf("failed to set up database. Error: %w", err)
	}
	return err
}

func setupDatabase() error {
	var err error

	// Locate source and target paths.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return fmt.Errorf("failed to make target database path (%s) absolute. Error: %w", dbTargetPath, err)
	}
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	if err != nil {
		return fmt.Errorf("failed to obtain absolute path to database file (%s). Error: %w", databaseTemplatePath, err)
	}

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	if err != nil {
		return fmt.Errorf("setup failed to copy template database (%v) to target path (%v). Error: %w", databaseTemplatePath, testDirectoryPath, err)
	}
	return err
}

func setupDirectories() error {
	var err error
	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("failed to remove target test directory (%v). Error: %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("failed to recreate target test directory (%v). Error: %w", testDirectoryPath, err)
	}
	return err
}

func teardown() error {
	return nil
}
