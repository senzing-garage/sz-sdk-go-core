package szabstractfactory_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/require"
)

const (
	baseCallerSkip    = 4
	defaultTruncation = 76
	instanceName      = "SzAbstractFactory Test"
	printErrors       = false
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badConfigDefinition = "}{"
	badConfigHandle     = uintptr(0)
	badDataSourceCode   = "\n\tGO_TEST"
	badLogLevelName     = "BadLogLevelName"
	badSettings         = "{]"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzAbstractFactory_CreateConfigManager(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)

	defer func() { require.NoError(test, szAbstractFactory.Destroy(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	printDebug(test, err, szConfigManager)
	require.NoError(test, err)
	configList, err := szConfigManager.GetConfigRegistry(ctx)
	printDebug(test, err, configList)
	require.NoError(test, err)
}

func TestSzAbstractFactory_CreateConfigManager_BadConfig(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObjectBadConfig(test)
	actual, err := szAbstractFactory.CreateConfigManager(ctx)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateConfigManager","error":{"function":"szconfigmanager.(*Szconfigmanager).Initialize","error":{"id":"SZSDK60024006","reason":"SENZ0018|Could not process initialization settings"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_CreateDiagnostic(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)

	defer func() { require.NoError(test, szAbstractFactory.Destroy(ctx)) }()

	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	printDebug(test, err, szDiagnostic)
	require.NoError(test, err)
	result, err := szDiagnostic.CheckRepositoryPerformance(ctx, 1)
	printDebug(test, err, result)
	require.NoError(test, err)
}

func TestSzAbstractFactory_CreateDiagnostic_BadConfig(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObjectBadConfig(test)
	actual, err := szAbstractFactory.CreateDiagnostic(ctx)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateDiagnostic","error":{"function":"szdiagnostic.(*Szdiagnostic).Initialize","error":{"id":"SZSDK60034005","reason":"SENZ0018|Could not process initialization settings"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_CreateEngine(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)

	defer func() { require.NoError(test, szAbstractFactory.Destroy(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	printDebug(test, err, szEngine)
	require.NoError(test, err)
	stats, err := szEngine.GetStats(ctx)
	printDebug(test, err, stats)
	require.NoError(test, err)
}

func TestSzAbstractFactory_CreateEngine_BadConfig(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObjectBadConfig(test)
	actual, err := szAbstractFactory.CreateEngine(ctx)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateEngine","error":{"function":"szengine.(*Szengine).Initialize","error":{"id":"SZSDK60044041","reason":"SENZ0018|Could not process initialization settings"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_CreateProduct(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)

	defer func() { require.NoError(test, szAbstractFactory.Destroy(ctx)) }()

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	printDebug(test, err, szProduct)
	require.NoError(test, err)
	version, err := szProduct.GetVersion(ctx)
	printDebug(test, err, version)
	require.NoError(test, err)
}

func TestSzAbstractFactory_CreateProduct_BadConfig(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObjectBadConfig(test)
	actual, err := szAbstractFactory.CreateProduct(ctx)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateProduct","error":{"function":"szproduct.(*Szproduct).Initialize","error":{"id":"SZSDK60064002","reason":"SENZ0018|Could not process initialization settings"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_Destroy(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)

	defer func() { require.NoError(test, szAbstractFactory.Destroy(ctx)) }()
}

func TestSzAbstractFactory_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)

	defer func() { require.NoError(test, szAbstractFactory.Destroy(ctx)) }()

	actual, err := szAbstractFactory.CreateDiagnostic(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	actual2, err := szAbstractFactory.CreateEngine(ctx)
	printDebug(test, err, actual2)
	require.NoError(test, err)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	printDebug(test, err, szConfigManager)
	require.NoError(test, err)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, configID)
	require.NoError(test, err)
	err = szAbstractFactory.Reinitialize(ctx, configID)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getSettings() string {
	var result string

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	panicOnError(err)

	databaseURL := "sqlite3://na:na@nowhere/" + dbTargetPath

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	panicOnError(err)

	return result
}

func getSettingsBadConfig() string {
	return badSettings
}

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var result senzing.SzAbstractFactory

	_ = ctx
	settings := getSettings()
	result = &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	return result
}

func getSzAbstractFactoryBadConfig(ctx context.Context) senzing.SzAbstractFactory {
	var result senzing.SzAbstractFactory

	_ = ctx
	settings := getSettingsBadConfig()
	result = &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func getTestObject(t *testing.T) senzing.SzAbstractFactory {
	t.Helper()
	ctx := t.Context()

	return getSzAbstractFactory(ctx)
}

func getTestObjectBadConfig(t *testing.T) senzing.SzAbstractFactory {
	t.Helper()
	ctx := t.Context()

	return getSzAbstractFactoryBadConfig(ctx)
}

func handleError(err error) {
	if err != nil {
		outputln("Error:", err)
	}
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printDebug(t *testing.T, err error, items ...any) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s\n", err.Error())
		}
	}

	if printResults {
		for _, item := range items {
			outLine := truncator.Truncate(fmt.Sprintf("%v", item), defaultTruncation, "...", truncator.PositionEnd)
			t.Logf("Result: %s\n", outLine)
		}
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	err := teardown()
	if err != nil {
		outputln(err)
	}

	os.Exit(code)
}

func setup() {
	setupDirectories()
	setupDatabase()

	err := setupSenzingConfiguration()
	panicOnError(err)
}

func setupDatabase() {
	testDirectoryPath := getTestDirectoryPath()
	_, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	panicOnError(err)
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	panicOnError(err)

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	panicOnError(err)
}

func setupDirectories() {
	testDirectoryPath := getTestDirectoryPath()
	err := os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	panicOnError(err)
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0o750) // recreate the test target directory
	panicOnError(err)
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	settings := getSettings()

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)

	defer func() { panicOnError(szConfig.Destroy(ctx)) }()

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)

	defer func() { panicOnError(szConfigManager.Destroy(ctx)) }()

	// Create a Senzing configuration.

	err = szConfig.ImportTemplate(ctx)
	panicOnError(err)

	// Add data sources to template Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
		panicOnError(err)
	}

	// Create a string representation of the Senzing configuration.

	configDefinition, err := szConfig.Export(ctx)
	panicOnError(err)

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	_, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	return nil
}

func teardown() error {
	return nil
}
