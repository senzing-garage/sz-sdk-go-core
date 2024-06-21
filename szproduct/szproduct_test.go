package szproduct

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/senzing-garage/sz-sdk-go/szproduct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	badLogLevelName   = "BadLogLevelName"
	defaultTruncation = 76
	instanceName      = "SzProduct Test"
	observerOrigin    = "SzProduct observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

var (
	logger            logging.Logging
	logLevel          = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szProductSingleton *Szproduct
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzproduct_GetLicense(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetLicense(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_GetVersion(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetVersion(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func TestSzproduct_getByteArray(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	szProduct.getByteArray(10)
}

func TestSzproduct_getByteArrayC(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	szProduct.getByteArrayC(10)
}

func TestSzproduct_newError(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.newError(ctx, 1)
	require.Error(test, err)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzproduct_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzproduct_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
}

func TestSzproduct_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	actual := szProduct.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzproduct_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzproduct_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szProduct := getSzProductAsInterface(ctx)
	actual, err := szProduct.GetLicense(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_Initialize(test *testing.T) {
	ctx := context.TODO()
	szProduct := &Szproduct{}
	settings, err := getSettings()
	require.NoError(test, err)
	err = szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Initialize_error
// func TestSzproduct_Initialize_error(test *testing.T) {}

func TestSzproduct_Destroy(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Destroy_error
// func TestSzproduct_Destroy_error(test *testing.T) {}

func TestSzproduct_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szProductSingleton = nil
	szProduct := getTestObject(ctx, test)
	err := szProduct.Destroy(ctx)
	require.NoError(test, err)
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

func getSzProduct(ctx context.Context) (*Szproduct, error) {
	var err error
	_ = ctx
	if szProductSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szProductSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		szProductSingleton = &Szproduct{}
		err = szProductSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szProductSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szProductSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szProductSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szProductSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szProductSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szProductSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szProductSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			return szProductSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szProductSingleton, err
}

func getSzProductAsInterface(ctx context.Context) senzing.SzProduct {
	result, err := getSzProduct(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szproduct")
}

func getTestObject(ctx context.Context, test *testing.T) *Szproduct {
	result, err := getSzProduct(ctx)
	require.NoError(test, err)
	return result
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
		if errors.Is(err, szerror.ErrSzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzBadInput) {
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

func setup() error {
	var err error
	logger = helper.GetLogger(ComponentID, szproduct.IDMessages, baseCallerSkip)
	osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
	if len(osenvLogLevel) > 0 {
		logLevel = osenvLogLevel
	}
	err = setupDirectories()
	if err != nil {
		return fmt.Errorf("Failed to set up directories. Error: %w", err)
	}
	err = setupDatabase()
	if err != nil {
		return fmt.Errorf("Failed to set up database. Error: %w", err)
	}
	return err
}

func setupDatabase() error {
	var err error

	// Locate source and target paths.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return fmt.Errorf("failed to make target database path (%s) absolute. Error: %w",
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
	var err error
	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("failed to remove target test directory (%v): %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("failed to recreate target test directory (%v): %w", testDirectoryPath, err)
	}
	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzProduct(ctx)
	return err
}

func teardownSzProduct(ctx context.Context) error {
	err := szProductSingleton.UnregisterObserver(ctx, observerSingleton)
	if err != nil {
		return err
	}
	err = szProductSingleton.Destroy(ctx)
	if err != nil {
		return err
	}
	szProductSingleton = nil
	return nil
}
