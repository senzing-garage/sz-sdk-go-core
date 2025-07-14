package szproduct_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szproduct"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	instanceName      = "SzProduct Test"
	jsonIndentation   = "    "
	observerOrigin    = "SzProduct observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badLogLevelName = "BadLogLevelName"
)

var (
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szProductSingleton *szproduct.Szproduct
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzproduct_GetLicense(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	actual, err := szProduct.GetLicense(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzproduct_GetVersion(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	actual, err := szProduct.GetVersion(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzproduct_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzproduct_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	szProduct.SetObserverOrigin(ctx, originMessage)
}

func TestSzproduct_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	szProduct.SetObserverOrigin(ctx, originMessage)
	actual := szProduct.GetObserverOrigin(ctx)
	assert.Equal(test, originMessage, actual)
}

func TestSzproduct_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	err := szProduct.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzproduct_AsInterface(test *testing.T) {
	ctx := test.Context()
	szProduct := getSzProductAsInterface(ctx)
	actual, err := szProduct.GetLicense(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzproduct_Initialize(test *testing.T) {
	ctx := test.Context()
	szProduct := &szproduct.Szproduct{}
	settings := getSettings()
	err := szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzproduct_Initialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzengine_Initialize_error
	_ = test
}

func TestSzproduct_Destroy(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	err := szProduct.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzproduct_Destroy_error(test *testing.T) {
	// IMPROVE: Implement TestSzengine_Destroy_error
	_ = test
}

func TestSzproduct_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szProductSingleton = nil
	szProduct := getTestObject(test)
	err := szProduct.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
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

func getSzProduct(ctx context.Context) *szproduct.Szproduct {
	if szProductSingleton == nil {
		settings := getSettings()

		szProductSingleton = &szproduct.Szproduct{}
		err := szProductSingleton.SetLogLevel(ctx, logLevel)
		panicOnError(err)

		if logLevel == "TRACE" {
			szProductSingleton.SetObserverOrigin(ctx, observerOrigin)

			err = szProductSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)

			err = szProductSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}

		err = szProductSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
	}

	return szProductSingleton
}

func getSzProductAsInterface(ctx context.Context) senzing.SzProduct {
	return getSzProduct(ctx)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szproduct")
}

func getTestObject(t *testing.T) *szproduct.Szproduct {
	t.Helper()

	return getSzProduct(t.Context())
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

	teardown()
	os.Exit(code)
}

func setup() {
	setupDirectories()
	setupDatabase()
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

func teardown() {
	ctx := context.TODO()
	teardownSzProduct(ctx)
}

func teardownSzProduct(ctx context.Context) {
	err := szProductSingleton.UnregisterObserver(ctx, observerSingleton)
	panicOnError(err)

	_ = szProductSingleton.Destroy(ctx)

	szProductSingleton = nil
}
