package szengine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	futil "github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/testfixtures"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szengineapi "github.com/senzing-garage/sz-sdk-go/szengine"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 158
	instanceName      = "Engine Test"
	printResults      = false
	verboseLogging    = sz.SZ_NO_LOGGING
)

var (
	defaultConfigId     int64
	globalSzDiagnostic  szdiagnostic.Szdiagnostic = szdiagnostic.Szdiagnostic{}
	globalSzEngine      Szengine                  = Szengine{}
	logger              logging.LoggingInterface
	szEngineInitialized bool = false
)

type GetEntityByRecordIdResponse struct {
	ResolvedEntity struct {
		EntityId int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getDefaultConfigId() int64 {
	return defaultConfigId
}

// func getSzDiagnostic(ctx context.Context) sz.SzDiagnostic {
// 	_ = ctx
// 	return &globalSzDiagnostic
// }

func getSzEngine(ctx context.Context) sz.SzEngine {
	_ = ctx
	return &globalSzEngine
}

func getEntityIdForRecord(datasource string, id string) int64 {
	ctx := context.TODO()
	var result int64 = 0
	szEngine := getSzEngine(ctx)
	flags := int64(0)
	response, err := szEngine.GetEntityByRecordId(ctx, datasource, id, flags)
	if err != nil {
		return result
	}
	getEntityByRecordIdResponse := &GetEntityByRecordIdResponse{}
	err = json.Unmarshal([]byte(response), &getEntityByRecordIdResponse)
	if err != nil {
		return result
	}
	return getEntityByRecordIdResponse.ResolvedEntity.EntityId
}

func getEntityIdStringForRecord(datasource string, id string) string {
	entityId := getEntityIdForRecord(datasource, id)
	return strconv.FormatInt(entityId, 10)
}

func getEntityId(record record.Record) int64 {
	return getEntityIdForRecord(record.DataSource, record.Id)
}

func getEntityIdString(record record.Record) string {
	entityId := getEntityId(record)
	return strconv.FormatInt(entityId, 10)
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szengine")
}

func getTestObject(ctx context.Context, test *testing.T) sz.SzEngine {
	_ = ctx
	_ = test
	return &globalSzEngine
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, szEngine sz.SzEngine, err error) {
	_ = ctx
	_ = szEngine
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorBasic(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, szEngine sz.SzEngine, err error) {
	_ = ctx
	_ = szEngine
	if err != nil {
		test.Log("Error:", err.Error())
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
		if szerror.Is(err, szerror.SzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzBadInput) {
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

func createSettings(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}
	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5902, err)
	}
	return settings, err
}

func getSettings() (string, error) {
	dbUrl, _, err := setupDatabase(true)
	if err != nil {
		return "", err
	}
	settings, err := createSettings(dbUrl)
	if err != nil {
		return "", err
	}
	return settings, nil
}

func restoreSzEngine(ctx context.Context) error {
	settings, err := getSettings()
	if err != nil {
		return err
	}
	err = setupSzEngine(ctx, instanceName, settings, verboseLogging, false)
	if err != nil {
		return err
	}
	return nil
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szengineapi.IdMessages)
	if err != nil {
		return createError(5901, err)
	}

	// Cleanup past runs and prepare for current run.

	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath))
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750)
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", testDirectoryPath, err)
	}

	// Get the database URL and determine if external or a local file just created.

	dbUrl, dbPurge, err := setupDatabase(false)
	if err != nil {
		return err
	}

	// Create the Senzing engine configuration JSON.

	settings, err := createSettings(dbUrl)
	if err != nil {
		return err
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfiguration(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5920, err)
	}

	// Setup the engine.

	err = setupSzEngine(ctx, instanceName, settings, verboseLogging, dbPurge)
	if err != nil {
		return err
	}

	// Preload records.

	custRecords := truthset.CustomerRecords
	records := []record.Record{custRecords["1001"], custRecords["1002"], custRecords["1003"]}
	flags := sz.SZ_WITHOUT_INFO
	for _, record := range records {
		_, err = globalSzEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
		if err != nil {
			defer teardownSzEngine(ctx)
			return err
		}
	}

	return err
}

func setupDatabase(preserveDB bool) (string, bool, error) {
	var err error = nil

	// Get paths.

	testDirectoryPath := getTestDirectoryPath()
	dbFilePath, err := filepath.Abs(getDatabaseTemplatePath())
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			dbFilePath, err.Error())
		return "", false, err
	}
	dbTargetPath := filepath.Join(getTestDirectoryPath(), "G2C.db")
	dbTargetPath, err = filepath.Abs(dbTargetPath)
	if err != nil {
		err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
		return "", false, err
	}

	// Check the environment for a database URL.

	dbUrl, envUrlExists := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")
	dbDefaultUrl := fmt.Sprintf("sqlite3://na:na@%s", dbTargetPath)
	dbExternal := envUrlExists && dbDefaultUrl != dbUrl
	if !dbExternal {
		dbUrl = dbDefaultUrl
		if !preserveDB {
			_, _, err = futil.CopyFile(dbFilePath, testDirectoryPath, true) // Copy the SQLite database file.
			if err != nil {
				err = fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
					dbFilePath, testDirectoryPath, err)
				// Fall through to return the error.
			}
		}
	}
	return dbUrl, dbExternal, err
}

func setupSzEngine(ctx context.Context, instanceName string, settings string, verboseLogging int64, purge bool) error {
	if szEngineInitialized {
		return fmt.Errorf("SzEngine is already setup and has not been torn down.")
	}
	globalSzEngine.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)

	err := globalSzEngine.Initialize(ctx, instanceName, settings, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
	if err != nil {
		return createError(5903, err)
	}

	// In case of an external database (e.g.: PostgreSQL) we need to purge since the database
	// may be shared across test suites -- this is not ideal since tests are not isolated.
	// TODO: Look for a way to use external databases while still isolating tests.

	if purge {
		err = globalSzDiagnostic.PurgeRepository(ctx)
		if err != nil {
			// if an error occurred on purge make sure to destroy the engine
			defer globalSzEngine.Destroy(ctx)
			return createError(5904, err)
		}
	}
	szEngineInitialized = true
	return err // Should be nil if we get here.
}

func setupSenzingConfiguration(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	now := time.Now()

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	// Create a string representation of the in-memory configuration.

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	// Create an alternate Senzing configuration.

	configHandle, err = szConfig.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	dataSourceCodes = []string{"CUSTOMERS", "REFERENCE", "WATCHLIST", "EMPLOYEES"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	alternateConfigDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = szConfig.Destroy(ctx)
	if err != nil {
		return createError(5911, err)
	}

	// Persist the Senzing configurations to the Senzing repository.

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComment := fmt.Sprintf("Created by szengine_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}
	defaultConfigId = configId

	configComment = fmt.Sprintf("Alternate config created by szengine_test at %s", now.UTC())
	_, err = szConfigManager.AddConfig(ctx, alternateConfigDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzEngine(ctx)
	return err
}

func teardownSzEngine(ctx context.Context) error {
	if !szEngineInitialized {
		return nil
	}
	err := globalSzEngine.Destroy(ctx)
	if err != nil {
		return err
	}
	szEngineInitialized = false
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzEngine_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
}

func TestSzEngine_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	actual := szEngine.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
	printActual(test, actual)
}

func TestSzEngine_AddRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "ADD_TEST_1", "NAME_FULL": "JIMMITY UNKNOWN"}`)
	testErrorBasic(test, err)

	record2, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "ADD_TEST_2", "NAME_FULL": "SOMEBODY NOTFOUND"}`)
	testErrorBasic(test, err)
	flags := sz.SZ_WITHOUT_INFO

	actual, err := szEngine.AddRecord(ctx, record1.DataSource, record1.Id, record1.Json, flags)
	testError(test, ctx, szEngine, err)
	defer szEngine.DeleteRecord(ctx, record1.DataSource, record1.Id, sz.SZ_NO_FLAGS)
	printActual(test, actual)

	actual, err = szEngine.AddRecord(ctx, record2.DataSource, record2.Id, record2.Json, flags)
	testError(test, ctx, szEngine, err)
	defer szEngine.DeleteRecord(ctx, record2.DataSource, record2.Id, sz.SZ_NO_FLAGS)
	printActual(test, actual)
}

func TestSzEngine_AddRecord_szBadInput(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "ADD_TEST_ERR_1", "NAME_FULL": "NOBODY NOMATCH"}`)
	testErrorBasic(test, err)
	record2, err := record.NewRecord(`{"DATA_SOURCE": "BOB", "RECORD_ID": "ADD_TEST_ERR_2", "NAME_FULL": "ERR BAD SOURCE"}`)
	testErrorBasic(test, err)
	flags := sz.SZ_WITHOUT_INFO

	// This one should succeed.

	actual, err := szEngine.AddRecord(ctx, record1.DataSource, record1.Id, record1.Json, flags)
	testError(test, ctx, szEngine, err)
	defer szEngine.DeleteRecord(ctx, record1.DataSource, record1.Id, sz.SZ_WITHOUT_INFO)
	printActual(test, actual)

	// This one should fail.

	actual, err = szEngine.AddRecord(ctx, "CUSTOMERS", record2.Id, record2.Json, flags)
	assert.True(test, szerror.Is(err, szerror.SzBadInput))
	printActual(test, actual)

	// Clean-up the records we inserted.

	actual, err = szEngine.DeleteRecord(ctx, record1.DataSource, record1.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_AddRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "WITH_INFO_1", "NAME_FULL": "HUBERT WITHINFO"}`)
	testErrorBasic(test, err)
	flags := sz.SZ_WITH_INFO

	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
	testError(test, ctx, szEngine, err)
	defer szEngine.DeleteRecord(ctx, record.DataSource, record.Id, sz.SZ_WITHOUT_INFO)
	printActual(test, actual)

	actual, err = szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.CountRedoRecords(ctx)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_ExportCsvEntityReport(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	expected := []string{
		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
		`1,0,"","","CUSTOMERS","1001"`,
		`1,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
		`1,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
	}
	csvColumnList := ""
	flags := sz.SZ_EXPORT_INCLUDE_ALL_ENTITIES
	aHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, aHandle)
		testError(test, ctx, szEngine, err)
	}()
	testError(test, ctx, szEngine, err)
	actualCount := 0
	for {
		actual, err := szEngine.FetchNext(ctx, aHandle)
		testError(test, ctx, szEngine, err)
		if len(actual) == 0 {
			break
		}
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual))
		actualCount += 1
	}
	assert.Equal(test, len(expected), actualCount)
}

func TestSzEngine_ExportCsvEntityReportIterator(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	expected := []string{
		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
		`1,0,"","","CUSTOMERS","1001"`,
		`1,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
		`1,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
	}
	csvColumnList := ""
	flags := sz.SZ_EXPORT_INCLUDE_ALL_ENTITIES
	actualCount := 0
	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		testError(test, ctx, szEngine, actual.Error)
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))
		actualCount += 1
	}
	assert.Equal(test, len(expected), actualCount)
}

func TestSzEngine_ExportJsonEntityReport(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.AddRecord(ctx, aRecord.DataSource, aRecord.Id, aRecord.Json, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
	defer szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.Id, sz.SZ_WITHOUT_INFO)
	// TODO: Figure out correct flags.
	// flags := sz.Flags(sz.SZ_EXPORT_DEFAULT_FLAGS, sz.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS, sz.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS)
	flags = int64(-1)
	aHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, aHandle)
		testError(test, ctx, szEngine, err)
	}()
	testError(test, ctx, szEngine, err)
	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, aHandle)
		testError(test, ctx, szEngine, err)
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}
	testError(test, ctx, szEngine, err)
	assert.True(test, len(jsonEntityReport) > 65536)
}

func TestSzEngine_ExportJsonEntityReportIterator(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := sz.SZ_EXPORT_INCLUDE_ALL_ENTITIES
	actualCount := 0
	for actual := range szEngine.ExportJsonEntityReportIterator(ctx, flags) {
		testError(test, ctx, szEngine, actual.Error)
		printActual(test, actual.Value)
		actualCount += 1
	}
	assert.Equal(test, 1, actualCount)
}

func TestSzEngine_FindNetworkByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := sz.SZ_FIND_NETWORK_DEFAULT_FLAGS
	actual, err := szEngine.FindNetworkByEntityId(ctx, entityList, maxDegrees, buildOutDegree, maxEntities, flags)
	testErrorNoFail(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindNetworkByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.Id + `"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := sz.SZ_FIND_NETWORK_DEFAULT_FLAGS
	actual, err := szEngine.FindNetworkByRecordId(ctx, recordList, maxDegrees, buildOutDegree, maxEntities, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startEntityId := getEntityId(truthset.CustomerRecords["1001"])
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := sz.SZ_NO_EXCLUSIONS
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByEntityId_excluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(startRecord)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(startRecord) + `}]}`
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByEntityId_excludingAndIncluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(startRecord)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(startRecord) + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByEntityId_including(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(startRecord)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := sz.SZ_NO_EXCLUSIONS
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := sz.SZ_NO_EXCLUSIONS
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByRecordId_excluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	requiredDataSources := sz.SZ_NO_REQUIRED_DATASOURCES
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByRecordId_excludingAndIncluding(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_FindPathByRecordId_including(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := sz.SZ_NO_EXCLUSIONS
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, exclusions, requiredDataSources, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetActiveConfigId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetActiveConfigId(ctx)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetEntityByEntityId(ctx, entityId, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetEntityByRecordId(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetRedoRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetRedoRecord(ctx)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetRepositoryLastModifiedTime(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetRepositoryLastModifiedTime(ctx)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetStats(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetStats(ctx)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_GetVirtualEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.GetVirtualEntityByRecordId(ctx, recordList, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_HowEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.HowEntityByEntityId(ctx, entityId, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	err := szEngine.PrimeEngine(ctx)
	testError(test, ctx, szEngine, err)
}

func TestSzEngine_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_WITHOUT_INFO
	actual, err := szEngine.ReevaluateEntity(ctx, entityId, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.ReevaluateEntity(ctx, entityId, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_WITHOUT_INFO
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_ReevaluateRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_ReplaceRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	recordDefinition := `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1984", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "CUSTOMERS", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "1001", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`
	flags := sz.SZ_WITHOUT_INFO
	actual, err := szEngine.ReplaceRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
	record := truthset.CustomerRecords["1001"]
	actual, err = szEngine.ReplaceRecord(ctx, record.DataSource, record.Id, record.Json, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_ReplaceRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	recordDefinition := `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1985", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "CUSTOMERS", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "1001", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.ReplaceRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
	record := truthset.CustomerRecords["1001"]
	actual, err = szEngine.ReplaceRecord(ctx, record.DataSource, record.Id, record.Json, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := sz.SZ_NO_SEARCH_PROFILE
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_SearchByAttributes_searchProfile(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := sz.SZ_NO_SEARCH_PROFILE // TODO: Figure out the search profile
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.WhyEntities(ctx, entityId1, entityId2, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_WhyRecordInEntity(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.WhyRecordInEntity(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := sz.SZ_NO_FLAGS
	actual, err := szEngine.WhyRecords(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Examples that are not in sorted order.
// ----------------------------------------------------------------------------

func TestSzEngine_Initialize(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	settings, err := getSettings()
	testError(test, ctx, szEngine, err)
	configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
	err = szEngine.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	testError(test, ctx, szEngine, err)
}

func TestSzEngine_Initialize_withConfigId(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	settings, err := getSettings()
	testError(test, ctx, szEngine, err)
	configId := getDefaultConfigId()
	err = szEngine.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	testError(test, ctx, szEngine, err)
}

func TestSzEngine_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	configId, err := szEngine.GetActiveConfigId(ctx)
	testError(test, ctx, szEngine, err)
	err = szEngine.Reinitialize(ctx, configId)
	testError(test, ctx, szEngine, err)
	printActual(test, configId)
}

func TestSzEngine_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "DELETE_TEST", "NAME_FULL": "GONNA B. DELETED"}`)
	testError(test, ctx, szEngine, err)
	flags := sz.SZ_WITHOUT_INFO
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
	printActual(test, actual)
	testError(test, ctx, szEngine, err)
	actual, err = szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

func TestSzEngine_DeleteRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "DELETE_TEST", "NAME_FULL": "DELETE W. INFO"}`)
	testError(test, ctx, szEngine, err)
	flags := sz.SZ_WITH_INFO
	actual, err := szEngine.AddRecord(ctx, record.DataSource, record.Id, record.Json, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
	actual, err = szEngine.DeleteRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, szEngine, err)
	printActual(test, actual)
}

// TODO: Implement TestSzEngine_ProcessRedoRecord
// func TestSzEngine_ProcessRedoRecord(test *testing.T) {
// 	ctx := context.TODO()
// 	szEngine := getTestObject(ctx, test)
//  flags := sz.SZ_WITHOUT_INFO
// 	actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
// 	testError(test, ctx, szEngine, err)
// 	printActual(test, actual)
// }

// TODO: Implement TestSzEngine_ProcessRedoRecord_withInfo
// func TestSzEngine_ProcessRedoRecord_withInfo(test *testing.T) {
// 	ctx := context.TODO()
// 	szEngine := getTestObject(ctx, test)
//  flags := sz.SZ_WITH_INFO
// 	actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
// 	testError(test, ctx, szEngine, err)
// 	printActual(test, actual)
// }

func TestSzEngine_Destroy(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	err := szEngine.Destroy(ctx)
	testError(test, ctx, szEngine, err)
	szEngineInitialized = false
	restoreSzEngine(ctx) // put everything back the way it was
}
