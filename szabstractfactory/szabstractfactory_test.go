package szabstractfactory_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

type DataStores struct {
	Location string `json:"location"` //nolint
}

type GetRepositoryInfoResponse struct {
	DataStores []DataStores `json:"dataStores"` //nolint
}

const (
	baseCallerSkip    = 4
	baseTen           = 10
	defaultTruncation = 76
	instanceName      = "SzAbstractFactory Test"
	location1         = "nowhere1"
	location2         = "nowhere2"
	location3         = "nowhere3"
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
	defer func() { szAbstractFactory.Close(ctx) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	printDebug(test, err, szConfigManager)
	require.NoError(test, err)
	defer func() { szConfigManager.Destroy(ctx) }()

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
	require.Nil(test, actual)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateConfigManager","text":"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzConfigManager]","error":{"function":"szabstractfactory.(*Szabstractfactory).initializeAbstractFactory","error":{"function":"szengine.(*Szengine).Initialize","error":{"id":"SZSDK60044041","reason":"SENZ0018|Could not process initialization settings"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_CreateDiagnostic(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	printDebug(test, err, szDiagnostic)
	require.NoError(test, err)
	defer func() { szDiagnostic.Destroy(ctx) }()

	result, err := szDiagnostic.CheckRepositoryPerformance(ctx, 1)
	printDebug(test, err, result)
	require.NoError(test, err)
}

func TestSzAbstractFactory_CreateDiagnostic_BadConfig(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObjectBadConfig(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	actual, err := szAbstractFactory.CreateDiagnostic(ctx)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)
	require.Nil(test, actual)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateDiagnostic","text":"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzDiagnostic]","error":{"function":"szabstractfactory.(*Szabstractfactory).initializeAbstractFactory","error":{"function":"szengine.(*Szengine).Initialize","error":{"id":"SZSDK60044041","reason":"SENZ0018|Could not process initialization settings"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_CreateEngine(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	printDebug(test, err, szEngine)
	require.NoError(test, err)
	defer func() { szEngine.Destroy(ctx) }()

	stats, err := szEngine.GetStats(ctx)
	printDebug(test, err, stats)
	require.NoError(test, err)
}

func TestSzAbstractFactory_CreateEngine_BadConfig(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObjectBadConfig(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	actual, err := szAbstractFactory.CreateEngine(ctx)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)
	require.Nil(test, actual)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateEngine","text":"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzEngine]","error":{"function":"szabstractfactory.(*Szabstractfactory).initializeAbstractFactory","error":{"function":"szengine.(*Szengine).Initialize","error":{"id":"SZSDK60044041","reason":"SENZ0018|Could not process initialization settings"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_CreateProduct(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	printDebug(test, err, szProduct)
	require.NoError(test, err)
	defer func() { szProduct.Destroy(ctx) }()

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
	require.Nil(test, actual)

	expectedErr := `{"function":"szabstractfactory.(*Szabstractfactory).CreateProduct","text":"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzProduct]","error":{"function":"szabstractfactory.(*Szabstractfactory).initializeAbstractFactory","error":{"function":"szengine.(*Szengine).Initialize","error":{"id":"SZSDK60044041","reason":"SENZ0018|Could not process initialization settings"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzAbstractFactory_Destroy_SzConfigManager(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	for range 3 {
		szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
		require.NoError(test, err)
		err = szConfigManager.Destroy(ctx)
		require.NoError(test, err)
		defer func() { szConfigManager.Destroy(ctx) }()
	}
}

func TestSzAbstractFactory_Destroy_SzDiagnostic(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	for range 3 {
		szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
		require.NoError(test, err)
		err = szDiagnostic.Destroy(ctx)
		require.NoError(test, err)
		defer func() { szDiagnostic.Destroy(ctx) }()
	}
}

func TestSzAbstractFactory_Destroy_SzEngine(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	for range 3 {
		szEngine, err := szAbstractFactory.CreateEngine(ctx)
		require.NoError(test, err)
		err = szEngine.Destroy(ctx)
		require.NoError(test, err)
		defer func() { szEngine.Destroy(ctx) }()
	}
}

func TestSzAbstractFactory_Destroy_SzProduct(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	for range 3 {
		szProduct, err := szAbstractFactory.CreateProduct(ctx)
		require.NoError(test, err)
		err = szProduct.Destroy(ctx)
		require.NoError(test, err)
		defer func() { szProduct.Destroy(ctx) }()
	}
}

func TestSzAbstractFactory_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szAbstractFactory := getTestObject(test)
	defer func() { szAbstractFactory.Close(ctx) }()

	actual, err := szAbstractFactory.CreateDiagnostic(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
	defer func() { actual.Destroy(ctx) }()

	actual2, err := szAbstractFactory.CreateEngine(ctx)
	printDebug(test, err, actual2)
	require.NoError(test, err)
	defer func() { actual2.Destroy(ctx) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	printDebug(test, err, szConfigManager)
	require.NoError(test, err)
	defer func() { szConfigManager.Destroy(ctx) }()

	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	printDebug(test, err, configID)
	require.NoError(test, err)
	err = szAbstractFactory.Reinitialize(ctx, configID)
	printDebug(test, err)
	require.NoError(test, err)
}

/*
Verify that all object can be created and subsequent SzAbstractFactories cannot create
objects until the first SzAbstractFactory is destroyed.
*/
func TestSzAbstractFactory_Multi_CreateAll(test *testing.T) {
	var err error

	ctx := test.Context()

	// Create AbstractFactories.

	szAbstractFactory1 := createSzAbstractFactoryByLocation(ctx, location1)
	szAbstractFactory2 := createSzAbstractFactoryByLocation(ctx, location2)

	// Create artifacts using first AbstractFactory.

	szConfigManager1, err := szAbstractFactory1.CreateConfigManager(ctx)
	require.NoError(test, err, "szAbstractFactory1 should create a SzConfigManager")
	szDiagnostic1, err := szAbstractFactory1.CreateDiagnostic(ctx)
	require.NoError(test, err, "szAbstractFactory1 should create a SzDiagnostic")
	szEngine1, err := szAbstractFactory1.CreateEngine(ctx)
	require.NoError(test, err, "szAbstractFactory1 should create a SzEngine")
	szProduct1, err := szAbstractFactory1.CreateProduct(ctx)
	require.NoError(test, err, "szAbstractFactory1 should create a SzProduct")

	// Second AbstractFactory doesn't created objects when objects are active.

	szConfigManager2, err := szAbstractFactory2.CreateConfigManager(ctx)
	require.Error(test, err, "szAbstractFactory2 should not create a SzConfigManager")
	szDiagnostic2, err := szAbstractFactory2.CreateDiagnostic(ctx)
	require.Error(test, err, "szAbstractFactory2 should not create a SzDiagnostic")
	szEngine2, err := szAbstractFactory2.CreateEngine(ctx)
	require.Error(test, err, "szAbstractFactory2 should not create a SzEngine")
	szProduct2, err := szAbstractFactory2.CreateProduct(ctx)
	require.Error(test, err, "szAbstractFactory2 should not create a SzProduct")

	// Artifacts from first AbstractFactory are destroyed.

	require.NoError(test, szAbstractFactory1.Close(ctx))
	require.NoError(test, szConfigManager1.Destroy(ctx))
	require.NoError(test, szDiagnostic1.Destroy(ctx))
	require.NoError(test, szEngine1.Destroy(ctx))
	require.NoError(test, szProduct1.Destroy(ctx))

	// Now artifacts can be created by second AbstractFactory.

	szConfigManager2, err = szAbstractFactory2.CreateConfigManager(ctx)
	require.NoError(test, err, "szAbstractFactory2 should create a SzConfigManager")
	szDiagnostic2, err = szAbstractFactory2.CreateDiagnostic(ctx)
	require.NoError(test, err, "szAbstractFactory2 should create a SzDiagnostic")
	szEngine2, err = szAbstractFactory2.CreateEngine(ctx)
	require.NoError(test, err, "szAbstractFactory2 should create a SzEngine")
	szProduct2, err = szAbstractFactory2.CreateProduct(ctx)
	require.NoError(test, err, "szAbstractFactory2 should create a SzProduct")

	// Artifacts from second AbstractFactory are destroyed.

	require.NoError(test, szConfigManager2.Destroy(ctx))
	require.NoError(test, szDiagnostic2.Destroy(ctx))
	require.NoError(test, szEngine2.Destroy(ctx))
	require.NoError(test, szProduct2.Destroy(ctx))
	require.NoError(test, szAbstractFactory2.Close(ctx))
}

/*
Verify that a second and third SzAbstractFactories are disabled if the first SzAbstractFactory is active.
*/
func TestSzAbstractFactory_Multi_PreventAdditionalAbstractFactories(test *testing.T) {
	ctx := test.Context()

	// First AbstractFactory without Destroy (it's deferred).

	szAbstractFactory1 := createSzAbstractFactoryByLocation(ctx, location1)
	defer func() { require.NoError(test, szAbstractFactory1.Close(ctx)) }()

	szConfigManager1, err := szAbstractFactory1.CreateConfigManager(ctx)
	require.NoError(test, err)
	defer func() { require.NoError(test, szConfigManager1.Destroy(ctx)) }()

	_, err = szConfigManager1.CreateConfigFromTemplate(ctx)
	require.NoError(test, err)

	// Second AbstractFactory should fail.

	szAbstractFactory2 := createSzAbstractFactoryByLocation(ctx, location2)
	defer func() { require.NoError(test, szAbstractFactory2.Close(ctx)) }()

	_, err = szAbstractFactory2.CreateConfigManager(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzConfigManager")
	_, err = szAbstractFactory2.CreateDiagnostic(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzDiagnostic")
	_, err = szAbstractFactory2.CreateEngine(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzEngine")
	_, err = szAbstractFactory2.CreateProduct(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzProduct")

	// Third AbstractFactory should fail.

	szAbstractFactory3 := createSzAbstractFactoryByLocation(ctx, location3)
	defer func() { require.NoError(test, szAbstractFactory3.Close(ctx)) }()

	_, err = szAbstractFactory3.CreateConfigManager(ctx)
	require.Error(test, err, "szAbstractFactory3 should not create a SzConfigManager")
	_, err = szAbstractFactory3.CreateDiagnostic(ctx)
	require.Error(test, err, "szAbstractFactory3 should not create a SzDiagnostic")
	_, err = szAbstractFactory3.CreateEngine(ctx)
	require.Error(test, err, "szAbstractFactory3 should not create a SzEngine")
	_, err = szAbstractFactory3.CreateProduct(ctx)
	require.Error(test, err, "szAbstractFactory3 should not create a SzProduct")
}

/*
Verify that a second SzAbstractFactory works after the first SzAbstractFactory has been destroyed.
*/
func TestSzAbstractFactory_Multi_PreventSecondAbstractFactory_withRetry(test *testing.T) {
	ctx := test.Context()

	// First AbstractFactory without Destroy (it's deferred).

	szAbstractFactory1 := createSzAbstractFactoryByLocation(ctx, location1)

	szDiagnostic1, err := szAbstractFactory1.CreateDiagnostic(ctx)
	require.NoError(test, err)

	info1, err := szDiagnostic1.GetRepositoryInfo(ctx)
	require.NoError(test, err)
	require.Equal(test, location1, extractLocation(info1), "AbstractFactory1 ")

	// Second AbstractFactory should fail.

	szAbstractFactory2 := createSzAbstractFactoryByLocation(ctx, location2)
	defer func() { require.NoError(test, szAbstractFactory2.Close(ctx)) }()

	_, err = szAbstractFactory2.CreateConfigManager(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzConfigManager")
	_, err = szAbstractFactory2.CreateDiagnostic(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzDiagnostic")
	_, err = szAbstractFactory2.CreateEngine(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzEngine")
	_, err = szAbstractFactory2.CreateProduct(ctx)
	require.Error(test, err, "AbstractFactory2 should not create a SzProduct")

	// Destroy the first AbstractFactory.

	require.NoError(test, szAbstractFactory1.Close(ctx))
	require.NoError(test, szDiagnostic1.Destroy(ctx))

	// Second AbstractFactory should succeed.

	szConfigManager2, err := szAbstractFactory2.CreateConfigManager(ctx)
	require.NoError(test, err, "AbstractFactory2 should create a SzConfigManager")
	defer func() { require.NoError(test, szConfigManager2.Destroy(ctx)) }()

	szDiagnostic2, err := szAbstractFactory2.CreateDiagnostic(ctx)
	require.NoError(test, err, "AbstractFactory2 should create a SzDiagnostic")
	defer func() { require.NoError(test, szDiagnostic2.Destroy(ctx)) }()

	szEngine2, err := szAbstractFactory2.CreateEngine(ctx)
	require.NoError(test, err, "AbstractFactory2 should create a SzEngine")
	defer func() { require.NoError(test, szEngine2.Destroy(ctx)) }()

	szProduct2, err := szAbstractFactory2.CreateProduct(ctx)
	require.NoError(test, err, "AbstractFactory2 should create a SzProduct")
	defer func() { require.NoError(test, szProduct2.Destroy(ctx)) }()

	// Try a third AbstractFactory.

	szAbstractFactory3 := createSzAbstractFactoryByLocation(ctx, location3)
	defer func() { require.NoError(test, szAbstractFactory3.Close(ctx)) }()

	_, err = szAbstractFactory3.CreateDiagnostic(ctx)
	require.Error(test, err, "AbstractFactory3 should not create objects")
}

/*
Verify a Destroy can be called from any SzAbstractFactory.
*/
func TestSzAbstractFactory_Multi_DestroyViaSecondAbstractFactory(test *testing.T) {
	ctx := test.Context()

	// Create AbstractFactories.

	szAbstractFactory1 := createSzAbstractFactoryByLocation(ctx, location1)
	szAbstractFactory2 := createSzAbstractFactoryByLocation(ctx, location2)
	defer func() { require.NoError(test, szAbstractFactory2.Close(ctx)) }()

	// Get object from szAbstractFactory1.

	szDiagnostic1, err := szAbstractFactory1.CreateDiagnostic(ctx)
	require.NoError(test, err)
	info1, err := szDiagnostic1.GetRepositoryInfo(ctx)
	require.NoError(test, err)
	require.Equal(test, location1, extractLocation(info1), "AbstractFactory1")

	// Try to get object from szAbstractFactory2.

	_, err = szAbstractFactory2.CreateDiagnostic(ctx)
	require.Error(test, err, "AbstractFactory2 should create not objects")

	// Destroy via the second AbstractFactory.

	require.NoError(test, szDiagnostic1.Destroy(ctx))
	require.NoError(test, szAbstractFactory1.Close(ctx))

	// Try to get object from szAbstractFactory2, again.

	szDiagnostic2, err := szAbstractFactory2.CreateDiagnostic(ctx)
	require.NoError(test, err)
	defer func() { require.NoError(test, szDiagnostic2.Destroy(ctx)) }()

	info2, err := szDiagnostic2.GetRepositoryInfo(ctx)
	require.NoError(test, err)
	require.Equal(test, location2, extractLocation(info2), "AbstractFactory2")
}

/*
Verify that an "orphaned" Senzing object inherits the configuration of the current SzAbstractFactory.
*/
func TestSzAbstractFactory_Multi_OrphanedObject(test *testing.T) {
	ctx := test.Context()

	// First AbstractFactory with Destroy.

	szAbstractFactory1 := createSzAbstractFactoryByLocation(ctx, location1)

	szDiagnostic1, err := szAbstractFactory1.CreateDiagnostic(ctx)
	require.NoError(test, err)

	info1, err := szDiagnostic1.GetRepositoryInfo(ctx)
	require.NoError(test, err)
	require.Equal(test, location1, extractLocation(info1), "AbstractFactory1 ")

	require.NoError(test, szAbstractFactory1.Close(ctx))
	require.NoError(test, szDiagnostic1.Destroy(ctx))

	// Second AbstractFactory with deferred Destroy.

	szAbstractFactory2 := createSzAbstractFactoryByLocation(ctx, location2)

	szDiagnostic2, err := szAbstractFactory2.CreateDiagnostic(ctx)
	require.NoError(test, err)

	info2, err := szDiagnostic2.GetRepositoryInfo(ctx)
	require.NoError(test, err)
	require.Equal(test, location2, extractLocation(info2))

	require.NoError(test, szAbstractFactory2.Close(ctx))
	require.NoError(test, szDiagnostic2.Destroy(ctx))

	// The orphaned "Go" object is closed, so it will error.

	_, err = szDiagnostic1.GetRepositoryInfo(ctx)
	require.Error(test, err)
}

/*
Verify that an orphaned Senzing object picks up a new configuration.
*/
func TestSzAbstractFactory_Multi_Reinitialize_implicitly(test *testing.T) {
	ctx := test.Context()
	now := time.Now()
	timeSuffix := strconv.FormatInt(now.Unix(), baseTen)
	dataSourceCode := "GO_TEST_" + timeSuffix
	recordID := "RECORD_ID_" + timeSuffix
	recordDefinition := `{"DATA_SOURCE": "` + dataSourceCode + `", "RECORD_ID": "` + recordID + `", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`

	// Create AbstractFactories.

	szAbstractFactory1 := createSzAbstractFactoryByLocation(ctx, location1)
	szAbstractFactory2 := createSzAbstractFactoryByLocation(ctx, location2)
	defer func() { require.NoError(test, szAbstractFactory2.Close(ctx)) }()

	// Get Senzing objects from AbstractFactory1.

	szConfigManager1, err := szAbstractFactory1.CreateConfigManager(ctx)
	require.NoError(test, err)

	szEngine1, err := szAbstractFactory1.CreateEngine(ctx)
	require.NoError(test, err)

	// Add data source to Senzing configuration.

	configID, err := szConfigManager1.GetDefaultConfigID(ctx)
	require.NoError(test, err)

	szConfig, err := szConfigManager1.CreateConfigFromConfigID(ctx, configID)
	require.NoError(test, err)

	_, err = szConfig.RegisterDataSource(ctx, dataSourceCode)
	require.NoError(test, err)

	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)

	newConfigID, err := szConfigManager1.RegisterConfig(ctx, configDefinition, "Add "+dataSourceCode)
	require.NoError(test, err)

	err = szConfigManager1.ReplaceDefaultConfigID(ctx, configID, newConfigID)
	require.NoError(test, err)

	// Inserting record before reinitializing should fail because it hasn't been reinitialized.

	_, err = szEngine1.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, senzing.SzNoFlags)
	require.Error(test, err)

	// Get SzEngine2 fails because AbstractFactory1 hasn't been destroyed.

	_, err = szAbstractFactory2.CreateEngine(ctx)
	require.Error(test, err)

	// Destroy AbstractFactory1 to allow AbstractFactory2 to create Senzing objects.

	err = szAbstractFactory1.Close(ctx)
	require.NoError(test, err)

	// Orphaned szEngine still succeeds.

	_, err = szEngine1.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, senzing.SzNoFlags)
	require.Error(test, err)

	// Create an SzEngine from AbstractFactory2.

	_, err = szAbstractFactory2.CreateEngine(ctx)
	require.Error(test, err)

	err = szEngine1.Destroy(ctx)
	require.NoError(test, err)

	// Orphaned szEngine1 fails.

	_, err = szEngine1.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, senzing.SzNoFlags)
	require.Error(test, err)

	// Try engine2.

	szEngine2, err := szAbstractFactory2.CreateEngine(ctx)
	require.NoError(test, err)
	defer func() { require.NoError(test, szEngine2.Destroy(ctx)) }()

	_, err = szEngine2.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, senzing.SzNoFlags)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	result := createSzAbstractFactoryByLocation(ctx, location1)
	// _ = result.DestroyWithoutClosing(ctx)

	return result
}

func createSzAbstractFactoryByLocation(ctx context.Context, location string) senzing.SzAbstractFactory {
	var result senzing.SzAbstractFactory

	_ = ctx
	settings := getSettings(location)
	result = &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	return result
}

func createSzAbstractFactoryBadConfig(ctx context.Context) senzing.SzAbstractFactory {
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

func extractLocation(location string) string {
	getRepositoryInfoResponse := GetRepositoryInfoResponse{}
	err := json.Unmarshal([]byte(location), &getRepositoryInfoResponse)
	panicOnError(err)

	return getRepositoryInfoResponse.DataStores[0].Location
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getSettings(location string) string {
	var result string

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	panicOnError(err)

	databaseURL := fmt.Sprintf("sqlite3://na:na@%s/%s", location, dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseURL": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	panicOnError(err)

	return result
}

func getSettingsBadConfig() string {
	return badSettings
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szabstractfactory")
}

func getTestObject(t *testing.T) senzing.SzAbstractFactory {
	t.Helper()
	ctx := t.Context()

	return createSzAbstractFactory(ctx)
}

func getTestObjectBadConfig(t *testing.T) senzing.SzAbstractFactory {
	t.Helper()
	ctx := t.Context()

	return createSzAbstractFactoryBadConfig(ctx)
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
	_, _, err = fileutil.CopyFile(
		databaseTemplatePath,
		testDirectoryPath+"/G2C2.db",
		true,
	) // Copy the SQLite database file.
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

	settings := getSettings(location1)

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
