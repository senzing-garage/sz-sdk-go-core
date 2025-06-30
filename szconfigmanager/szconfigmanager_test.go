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
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
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
	jsonIndentation   = "    "
	observerOrigin    = "SzConfigManager observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
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
	printDebug(test, err1)
	require.NoError(test, err1)

	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, configID)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfigmanager_CreateConfigFromConfigID_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, badConfigID)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromConfigID","error":{"function":"szconfigmanager.(*Szconfigmanager).createConfigFromConfigIDChoreography","text":"getConfig(0)","error":{"id":"SZSDK60024003","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}}`
	require.JSONEq(test, expectedErr, err.Error())
	assert.Nil(test, actual)
}

func TestSzconfigmanager_CreateConfigFromConfigID_nilConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, nilConfigID)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromConfigID","error":{"function":"szconfigmanager.(*Szconfigmanager).createConfigFromConfigIDChoreography","text":"getConfig(0)","error":{"id":"SZSDK60024003","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}}`
	require.JSONEq(test, expectedErr, err.Error())
	assert.Nil(test, actual)
}

func TestSzconfigmanager_CreateConfigFromString(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printDebug(test, err, szConfig)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	printDebug(test, err, configDefinition)
	require.NoError(test, err)
	szConfig2, err := szConfigManager.CreateConfigFromString(ctx, configDefinition)
	printDebug(test, err, szConfig2)
	require.NoError(test, err)
	configDefinition2, err := szConfig2.Export(ctx)
	printDebug(test, err, configDefinition2)
	require.NoError(test, err)
	assert.JSONEq(test, configDefinition, configDefinition2)
}

func TestSzconfigmanager_CreateConfigFromString_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromString(ctx, badConfigDefinition)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromString","error":{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromStringChoreography","text":"VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).verifyConfigDefinitionChoreography","text":"load","error":{"id":"SZSDK60014009","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=2]"}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_CreateConfigFromTemplate(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	assert.NotEmpty(test, actual)
}

func TestSzconfigmanager_GetConfigRegistry(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.GetConfigRegistry(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfigmanager_GetDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfigmanager_RegisterConfig(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printDebug(test, err, szConfig)
	require.NoError(test, err)

	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	actual, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
	printDebug(test, err, actual)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	printDebug(test, err, configDefinition)
	require.NoError(test, err)

	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual2, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	printDebug(test, err, actual2)
	require.NoError(test, err)
}

func TestSzconfigmanager_RegisterConfig_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szConfigManager.RegisterConfig(ctx, badConfigDefinition, configComment)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).RegisterConfig","error":{"id":"SZSDK60024001","reason":"SENZ0028|Invalid JSON config document"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_RegisterConfig_nilConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szConfigManager.RegisterConfig(ctx, nilConfigDefinition, configComment)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).RegisterConfig","error":{"id":"SZSDK60024001","reason":"SENZ0028|Invalid JSON config document"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_RegisterConfig_nilConfigComment(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printDebug(test, err, szConfig)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	printDebug(test, err, configDefinition)
	require.NoError(test, err)
	actual, err := szConfigManager.RegisterConfig(ctx, configDefinition, nilConfigComment)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfigmanager_ReplaceDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err1)
	require.NoError(test, err1)

	// IMPROVE: This is kind of a cheater.

	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err2)
	require.NoError(test, err2)

	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badCurrentDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	newDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, newDefaultConfigID)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, badCurrentDefaultConfigID, newDefaultConfigID)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7245|Current configuration ID does not match specified data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badNewDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, currentDefaultConfigID)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, badNewDefaultConfigID)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilCurrentDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	newDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, newDefaultConfigID)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, nilCurrentDefaultConfigID, newDefaultConfigID)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7245|Current configuration ID does not match specified data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilNewDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, currentDefaultConfigID)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, nilNewDefaultConfigID)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_SetDefaultConfig(test *testing.T) {
	ctx := test.Context()
	now := time.Now()
	szConfigManager := getTestObject(test)
	defaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, defaultConfigID)
	require.NoError(test, err)
	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, defaultConfigID)
	printDebug(test, err, szConfig)
	require.NoError(test, err)

	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	actual, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
	printDebug(test, err, actual)
	require.NoError(test, err)
	configDefintion, err := szConfig.Export(ctx)
	printDebug(test, err, configDefintion)
	require.NoError(test, err)
	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefintion, "Added "+dataSourceCode)
	printDebug(test, err, configID)
	require.NoError(test, err)
	require.NotZero(test, configID)
}

func TestSzconfigmanager_SetDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, configID)
	require.NoError(test, err)
	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_SetDefaultConfigID_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.SetDefaultConfigID(ctx, badConfigID)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).SetDefaultConfigID","error":{"id":"SZSDK60024008","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_SetDefaultConfigID_nilConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.SetDefaultConfigID(ctx, nilConfigID)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).SetDefaultConfigID","error":{"id":"SZSDK60024008","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
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
	szConfigManager.SetObserverOrigin(ctx, originMessage)
}

func TestSzconfigmanager_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfigManager.SetObserverOrigin(ctx, originMessage)
	actual := szConfigManager.GetObserverOrigin(ctx)
	assert.Equal(test, originMessage, actual)
}

func TestSzconfigmanager_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfigmanager_AsInterface(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getSzConfigManagerAsInterface(ctx)
	actual, err := szConfigManager.GetConfigRegistry(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfigmanager_Initialize(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	settings := getSettings()
	err := szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_Initialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfigmanager_Initialize_error
	_ = test
}

func TestSzconfigmanager_Destroy(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szConfigManagerSingleton = nil
	szConfigManager := getTestObject(test)
	err := szConfigManager.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_Destroy_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfigmanager_Destroy_error
	_ = test
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

// func getSzConfig(ctx context.Context) *szconfig.Szconfig {
// 	if szConfigSingleton == nil {
// 		settings, err := getSettings()
// 		panicOnError(err)
// 		szConfigSingleton = &szconfig.Szconfig{}
// 		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
// 		panicOnError(err)
// 		if logLevel == "TRACE" {
// 			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
// 			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
// 			panicOnError(err)
// 			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
// 			panicOnError(err)
// 		}
// 		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
// 		panicOnError(err)
// 	}
// 	return szConfigSingleton
// }

func getSzConfigManager(ctx context.Context) *szconfigmanager.Szconfigmanager {
	if szConfigManagerSingleton == nil {
		settings := getSettings()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{}
		err := szConfigManagerSingleton.SetLogLevel(ctx, logLevel)
		panicOnError(err)

		if logLevel == "TRACE" {
			szConfigManagerSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigManagerSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)
			err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}

		err = szConfigManagerSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
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

	return getSzConfigManager(t.Context())
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

	configComment := fmt.Sprintf("Created by szconfigmanager_test at %s", now.UTC())
	_, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	panicOnError(err)

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
		panicOnError(err)
		err = szConfigSingleton.Destroy(ctx)
		panicOnError(err)

		szConfigSingleton = nil
	}
}

func teardownSzConfigManager(ctx context.Context) {
	if szConfigManagerSingleton != nil {
		err := szConfigManagerSingleton.UnregisterObserver(ctx, observerSingleton)
		panicOnError(err)
		err = szConfigManagerSingleton.Destroy(ctx)
		panicOnError(err)

		szConfigManagerSingleton = nil
	}
}
