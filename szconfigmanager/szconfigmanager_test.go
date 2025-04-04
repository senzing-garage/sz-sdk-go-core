package szconfigmanager_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	instanceName      = "SzConfigManager Test"
	observerOrigin    = "SzConfigManager observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badConfigDefinition       = "\n\t"
	badConfigID               = int64(0)
	badCurrentDefaultConfigID = int64(0)
	badLogLevelName           = "BadLogLevelName"
	badNewDefaultConfigID     = int64(0)
	baseTen                   = 10
)

// Nil/empty parameters

var (
	nilConfigComment          string
	nilConfigDefinition       string
	nilConfigID               int64
	nilCurrentDefaultConfigID int64
	nilNewDefaultConfigID     int64
)

var (
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton *szconfigmanager.Szconfigmanager
	szConfigSingleton        *szconfig.Szconfig
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzconfigmanager_CreateConfigFromConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	configID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err1)

	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, configID)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_CreateConfigFromConfigID_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, badConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
	assert.Nil(test, actual)
}

func TestSzconfigmanager_CreateConfigFromConfigID_nilConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, nilConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
	assert.Nil(test, actual)
}

func TestSzconfigmanager_CreateConfigFromString(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)
	szConfig2, err := szConfigManager.CreateConfigFromString(ctx, configDefinition)
	require.NoError(test, err)
	configDefinition2, err := szConfig2.Export(ctx)
	require.NoError(test, err)
	assert.JSONEq(test, configDefinition, configDefinition2)
}

func TestSzconfigmanager_CreateConfigFromTemplate(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromTemplate(ctx)
	require.NoError(test, err)
	assert.NotEmpty(test, actual)
}

func TestSzconfigmanager_GetConfigs(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.GetDefaultConfigID(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_RegisterConfig(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	require.NoError(test, err)

	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	_, err = szConfig.AddDataSource(ctx, dataSourceCode)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)

	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_RegisterConfig_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	_, err := szConfigManager.RegisterConfig(ctx, badConfigDefinition, configComment)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_RegisterConfig_nilConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	_, err := szConfigManager.RegisterConfig(ctx, nilConfigDefinition, configComment)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_RegisterConfig_nilConfigComment(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)
	actual, err := szConfigManager.RegisterConfig(ctx, configDefinition, nilConfigComment)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_ReplaceDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err1)

	// TODO: This is kind of a cheater.

	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err2)

	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	require.NoError(test, err)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badCurrentDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	newDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, badCurrentDefaultConfigID, newDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badNewDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, badNewDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilCurrentDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	newDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, nilCurrentDefaultConfigID, newDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilNewDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, nilNewDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_SetDefaultConfig(test *testing.T) {
	ctx := test.Context()
	now := time.Now()
	szConfigManager := getTestObject(test)
	defaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	require.NoError(test, err)
	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, defaultConfigID)
	require.NoError(test, err)

	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	_, err = szConfig.AddDataSource(ctx, dataSourceCode)
	require.NoError(test, err)
	configDefintion, err := szConfig.Export(ctx)
	require.NoError(test, err)
	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefintion, "Added "+dataSourceCode)
	require.NoError(test, err)
	require.NotZero(test, configID)
}

func TestSzconfigmanager_SetDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	handleErrorWithPanic(err)
	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	require.NoError(test, err)
}

func TestSzconfigmanager_SetDefaultConfigID_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.SetDefaultConfigID(ctx, badConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_SetDefaultConfigID_nilConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.SetDefaultConfigID(ctx, nilConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfigmanager_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	_ = szConfigManager.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfigmanager_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
}

func TestSzconfigmanager_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	actual := szConfigManager.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfigmanager_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfigmanager_AsInterface(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getSzConfigManagerAsInterface(ctx)
	actual, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_Initialize(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	settings := getSettings()
	err := szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfigmanager_Initialize_error
// func TestSzconfigmanager_Initialize_error(test *testing.T) {}

func TestSzconfigmanager_Destroy(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.Destroy(ctx)
	require.NoError(test, err)
}

func TestSzconfigmanager_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szConfigManagerSingleton = nil
	szConfigManager := getTestObject(test)
	err := szConfigManager.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfigmanager_Destroy_error
// func TestSzconfigmanager_Destroy_error(test *testing.T) {}

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
	handleErrorWithPanic(err)

	databaseURL := "sqlite3://na:na@nowhere/" + dbTargetPath

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	handleErrorWithPanic(err)

	return result
}

// func getSzConfig(ctx context.Context) *szconfig.Szconfig {
// 	if szConfigSingleton == nil {
// 		settings, err := getSettings()
// 		handleErrorWithPanic(err)
// 		szConfigSingleton = &szconfig.Szconfig{}
// 		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
// 		handleErrorWithPanic(err)
// 		if logLevel == "TRACE" {
// 			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
// 			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
// 			handleErrorWithPanic(err)
// 			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
// 			handleErrorWithPanic(err)
// 		}
// 		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
// 		handleErrorWithPanic(err)
// 	}
// 	return szConfigSingleton
// }

func getSzConfigManager(ctx context.Context) *szconfigmanager.Szconfigmanager {
	if szConfigManagerSingleton == nil {
		settings := getSettings()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{}
		err := szConfigManagerSingleton.SetLogLevel(ctx, logLevel)
		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szConfigManagerSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigManagerSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)
			err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}

		err = szConfigManagerSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		handleErrorWithPanic(err)
	}

	return szConfigManagerSingleton
}

func getSzConfigManagerAsInterface(ctx context.Context) senzing.SzConfigManager {
	return getSzConfigManager(ctx)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfigmanager")
}

func getTestObject(t *testing.T) *szconfigmanager.Szconfigmanager {
	t.Helper()
	ctx := t.Context()

	return getSzConfigManager(ctx)
}

func handleError(err error) {
	if err != nil {
		safePrintln("Error:", err)
	}
}

func handleErrorWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(t *testing.T, actual interface{}) {
	t.Helper()
	printResult(t, "Actual", actual)
}

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func safePrintln(message ...any) {
	fmt.Println(message...) //nolint
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
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

	err := setupSenzingConfiguration()
	handleErrorWithPanic(err)
}

func setupDatabase() {
	testDirectoryPath := getTestDirectoryPath()
	_, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	handleErrorWithPanic(err)
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	handleErrorWithPanic(err)

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	handleErrorWithPanic(err)
}

func setupDirectories() {
	testDirectoryPath := getTestDirectoryPath()
	err := os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	handleErrorWithPanic(err)
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	handleErrorWithPanic(err)
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()
	settings := getSettings()

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	handleErrorWithPanic(err)

	defer func() { handleErrorWithPanic(szConfig.Destroy(ctx)) }()

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	handleErrorWithPanic(err)

	defer func() { handleErrorWithPanic(szConfigManager.Destroy(ctx)) }()

	// Create a Senzing configuration.

	err = szConfig.ImportTemplate(ctx)
	handleErrorWithPanic(err)

	// Add data sources to template Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		handleErrorWithPanic(err)
	}

	// Create a string representation of the Senzing configuration.

	configDefinition, err := szConfig.Export(ctx)
	handleErrorWithPanic(err)

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szconfigmanager_test at %s", now.UTC())
	_, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	handleErrorWithPanic(err)

	return nil
}

func teardown() {
	ctx := context.TODO()
	teardownSzConfig(ctx)
	teardownSzConfigManager(ctx)
}

func teardownSzConfig(ctx context.Context) {
	if szConfigSingleton != nil {
		err := szConfigSingleton.UnregisterObserver(ctx, observerSingleton)
		handleErrorWithPanic(err)
		err = szConfigSingleton.Destroy(ctx)
		handleErrorWithPanic(err)

		szConfigSingleton = nil
	}
}

func teardownSzConfigManager(ctx context.Context) {
	if szConfigManagerSingleton != nil {
		err := szConfigManagerSingleton.UnregisterObserver(ctx, observerSingleton)
		handleErrorWithPanic(err)
		err = szConfigManagerSingleton.Destroy(ctx)
		handleErrorWithPanic(err)

		szConfigManagerSingleton = nil
	}
}
