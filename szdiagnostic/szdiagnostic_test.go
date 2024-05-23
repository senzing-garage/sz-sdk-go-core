package szdiagnostic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	instanceName      = "SzDiagnostic Test"
	observerOrigin    = "SzDiagnostic observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
	badLogLevelName   = "BadLogLevelName"
)

var (
	defaultConfigID   int64
	logger            logging.LoggingInterface
	logLevel          = "INFO"
	observerSingleton = &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	szDiagnosticSingleton *Szdiagnostic
	szEngineSingleton     *szengine.Szengine
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzdiagnostic_CheckDatastorePerformance(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_CheckDatastorePerformance_badSecondsToRun(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	badSecondsToRun := -1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, badSecondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetDatastoreInfo(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.GetDatastoreInfo(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szDiagnostic := getTestObject(ctx, test)
	featureID := int64(1)
	actual, err := szDiagnostic.GetFeature(ctx, int64(featureID))
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature_badFeatureID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szDiagnostic := getTestObject(ctx, test)
	featureID := int64(-1)
	actual, err := szDiagnostic.GetFeature(ctx, int64(featureID))
	require.ErrorIs(test, err, szerror.ErrSzBase)
	printActual(test, actual)
}

// TODO:  Determine if PurgeRepository can be tested without disturbing other testcases
// func TestSzdiagnostic_PurgeRepository(test *testing.T) {
// 	ctx := context.TODO()
// 	szDiagnostic := getTestObject(ctx, test)
// 	err := szDiagnostic.PurgeRepository(ctx)
// 	require.NoError(test, err)
// }

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfig_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzdiagnostic_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
}

func TestSzdiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzdiagnostic_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzdiagnostic_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getSzDiagnosticAsInterface(ctx)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_Initialize(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := &Szdiagnostic{}
	instanceName := "Test name"
	settings, err := getSettings()
	require.NoError(test, err)
	verboseLogging := senzing.SzNoLogging
	configID := senzing.SzInitializeWithDefaultConfiguration
	err = szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

func TestSzdiagnostic_Initialize_withConfigId(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := &Szdiagnostic{}
	instanceName := "Test name"
	settings, err := getSettings()
	require.NoError(test, err)
	verboseLogging := senzing.SzNoLogging
	configID := getDefaultConfigID()
	err = szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

func TestSzdiagnostic_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	configID := getDefaultConfigID()
	err := szDiagnostic.Reinitialize(ctx, configID)
	require.NoError(test, err)
}

// func TestSzdiagnostic_Reinitialize_badConfigID_scenario(test *testing.T) {
// 	ctx := context.TODO()
// 	szDiagnostic := getTestObject(ctx, test)
// 	// goodConfigID := getDefaultConfigID()
// 	badConfigID := int64(-1)
// 	err := szDiagnostic.Reinitialize(ctx, badConfigID)
// 	if err != nil {
// 		test.Logf("\n>>>> Reinitialize(ctx, badConfigID) error %v", err)
// 	}
// 	// require.ErrorIs(test, err, szerror.ErrSzConfiguration)

// 	// err = szDiagnostic.Reinitialize(ctx, goodConfigID)
// 	// if err != nil {
// 	// 	test.Logf("\n>>>> Reinitialize(ctx, goodConfigID) error %v", err)
// 	// }
// 	// require.NoError(test, err)

// 	szDiagnosticSingleton = nil
// 	szDiagnostic = getTestObject(ctx, test)

// 	secondsToRun := 1
// 	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
// 	require.NoError(test, err)
// 	printActual(test, actual)
// }

// func TestSzdiagnostic_Reinitialize_badConfigID(test *testing.T) {
// 	ctx := context.TODO()
// 	szDiagnostic := getTestObject(ctx, test)
// 	badConfigID := int64(-1)
// 	err := szDiagnostic.Reinitialize(ctx, badConfigID)
// 	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
// }

// func TestSzdiagnostic_Reinitialize_again(test *testing.T) {
// 	ctx := context.TODO()
// 	// szDiagnosticSingleton = nil
// 	szDiagnostic := getTestObject(ctx, test)
// 	configID := getDefaultConfigID()
// 	test.Logf("\n>>>>>> configID-2: %d", configID)

// 	err := szDiagnostic.Reinitialize(ctx, configID)
// 	require.NoError(test, err)

// 	secondsToRun := 1
// 	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
// 	require.NoError(test, err)
// 	printActual(test, actual)
// }

func TestSzdiagnostic_Destroy(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.Destroy(ctx)
	require.NoError(test, err)
}

func TestSzdiagnostic_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.Destroy(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, records []record.Record) error {
	var err error
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo
	for _, record := range records {
		_, err = szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
		if err != nil {
			return err
		}
	}
	return err
}

func createError(errorID int, err error) error {
	// return errors.Cast(logger.NewError(errorId, err), err)
	return logger.NewError(errorID, err)
}

func deleteRecords(ctx context.Context, records []record.Record) error {
	var err error
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo
	for _, record := range records {
		_, err = szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
		if err != nil {
			return err
		}
	}
	return err
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getDefaultConfigID() int64 {
	return defaultConfigID
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
	databaseURL := fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5900, err)
	}
	return settings, err
}

func getSzDiagnostic(ctx context.Context) *Szdiagnostic {
	_ = ctx
	if szDiagnosticSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			fmt.Printf("getSettings() Error: %v\n", err)
			return nil
		}
		szDiagnosticSingleton = &Szdiagnostic{}
		err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if logLevel == "TRACE" {
			szDiagnosticSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szDiagnosticSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				fmt.Println(err)
				return nil
			}
		}
		err = szDiagnosticSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return szDiagnosticSingleton
}

func getSzEngine(ctx context.Context) *szengine.Szengine {
	_ = ctx
	if szEngineSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			fmt.Printf("getSettings() Error: %v\n", err)
			return nil
		}
		szEngineSingleton = &szengine.Szengine{}
		err = szEngineSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if logLevel == "TRACE" {
			szEngineSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szEngineSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			err = szEngineSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				fmt.Println(err)
				return nil
			}
		}
		err = szEngineSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return szEngineSingleton
}

func getSzDiagnosticAsInterface(ctx context.Context) senzing.SzDiagnostic {
	return getSzDiagnostic(ctx)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szdiagnostic")
}

func getTestObject(ctx context.Context, test *testing.T) *Szdiagnostic {
	_ = test
	return getSzDiagnostic(ctx)
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
	logger, err = logging.NewSenzingSdkLogger(ComponentID, szdiagnostic.IDMessages)
	if err != nil {
		return createError(5901, err)
	}
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
	err = setupSenzingConfiguration()
	if err != nil {
		return createError(5920, err)
	}
	return err
}

func setupDatabase() error {
	var err error

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
	var err error
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

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	settings, err := getSettings()
	if err != nil {
		return createError(5901, err)
	}

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5902, err)
	}
	defer func() { handleError(szConfig.Destroy(ctx)) }()

	// Create an in memory Senzing configuration.

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return createError(5903, err)
	}

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5904, err)
		}
	}

	// Create a string representation of the in-memory configuration.

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5905, err)
	}

	// Close szConfig in-memory object.

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5906, err)
	}

	// Persist the Senzing configuration to the Senzing repository as default.

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5907, err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5908, err)
	}

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return createError(5909, err)
	}
	defaultConfigID = configID

	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzDiagnostic(ctx)
	return err
}

func teardownSzDiagnostic(ctx context.Context) error {
	err := szDiagnosticSingleton.UnregisterObserver(ctx, observerSingleton)
	if err != nil {
		return err
	}
	err = szDiagnosticSingleton.Destroy(ctx)
	if err != nil {
		return err
	}
	szDiagnosticSingleton = nil
	err = szEngineSingleton.UnregisterObserver(ctx, observerSingleton)
	if err != nil {
		return err
	}
	err = szEngineSingleton.Destroy(ctx)
	if err != nil {
		return err
	}
	szEngineSingleton = nil
	return nil
}
