package szengine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	futil "github.com/senzing-garage/go-common/fileutil"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
	"github.com/senzing-garage/go-common/record"
	"github.com/senzing-garage/go-common/testfixtures"
	"github.com/senzing-garage/go-common/truthset"
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
	defaultTruncation = 76
	instanceName      = "Engine Test Module"
	printResults      = false
	verboseLogging    = 0
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityId int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

var (
	engineInitialized  bool                      = false
	globalG2engine     Szengine                  = Szengine{}
	globalG2Diagnostic szdiagnostic.Szdiagnostic = szdiagnostic.Szdiagnostic{}
	logger             logging.LoggingInterface
	senzingConfigId    int64 = 0
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) sz.SzEngine {
	_ = ctx
	_ = test
	return &globalG2engine
}

func getG2Diagnostic(ctx context.Context) sz.SzDiagnostic {
	_ = ctx
	return &globalG2Diagnostic
}

func getG2Engine(ctx context.Context) sz.SzEngine {
	_ = ctx
	return &globalG2engine
}

func getEntityIdForRecord(datasource string, id string) int64 {
	ctx := context.TODO()
	var result int64 = 0
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	response, err := g2engine.GetEntityByRecordId(ctx, datasource, id, flags)
	if err != nil {
		return result
	}
	getEntityByRecordIDResponse := &GetEntityByRecordIDResponse{}
	err = json.Unmarshal([]byte(response), &getEntityByRecordIDResponse)
	if err != nil {
		return result
	}
	return getEntityByRecordIDResponse.ResolvedEntity.EntityId
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

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testErrorBasic(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testError(test *testing.T, ctx context.Context, g2engine sz.SzEngine, err error) {
	_ = ctx
	_ = g2engine
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2engine sz.SzEngine, err error) {
	_ = ctx
	_ = g2engine
	if err != nil {
		test.Log("Error:", err.Error())
	}
}

func baseDirectoryPath() string {
	return filepath.FromSlash("../target/test/g2engine")
}

func dbTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
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

func getIniParams() (string, error) {
	dbUrl, _, err := setupDB(true)
	if err != nil {
		return "", err
	}
	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return "", err
	}
	return iniParams, nil
}

func restoreG2engine(ctx context.Context) error {
	iniParams, err := getIniParams()
	if err != nil {
		return err
	}
	err = setupG2engine(ctx, instanceName, iniParams, verboseLogging, false)
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

	baseDir := baseDirectoryPath()
	err = os.RemoveAll(filepath.Clean(baseDir))
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", baseDir, err)
	}
	err = os.MkdirAll(filepath.Clean(baseDir), 0750)
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", baseDir, err)
	}

	// Get the database URL and determine if external or a local file just created.

	dbUrl, dbPurge, err := setupDB(false)
	if err != nil {
		return err
	}

	// Create the Senzing engine configuration JSON.

	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return err
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfiguration(ctx, instanceName, iniParams, verboseLogging)
	if err != nil {
		return createError(5920, err)
	}

	// Setup the engine.

	err = setupG2engine(ctx, instanceName, iniParams, verboseLogging, dbPurge)
	if err != nil {
		return err
	}

	// Preload records.

	custRecords := truthset.CustomerRecords
	records := []record.Record{custRecords["1001"], custRecords["1002"], custRecords["1003"]}
	for _, record := range records {
		_, err = globalG2engine.AddRecord(ctx, record.DataSource, record.Id, record.Json, sz.SZ_NO_FLAGS)
		if err != nil {
			defer teardownG2engine(ctx)
			return err
		}
	}

	return err
}

func setupDB(preserveDB bool) (string, bool, error) {
	var err error = nil

	// Get paths.

	baseDir := baseDirectoryPath()
	dbFilePath, err := filepath.Abs(dbTemplatePath())
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			dbFilePath, err.Error())
		return "", false, err
	}
	dbTargetPath := filepath.Join(baseDirectoryPath(), "G2C.db")
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
			_, _, err = futil.CopyFile(dbFilePath, baseDir, true) // Copy the SQLite database file.
			if err != nil {
				err = fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
					dbFilePath, baseDir, err)
				// Fall through to return the error.
			}
		}
	}
	return dbUrl, dbExternal, err
}

func setupG2engine(ctx context.Context, moduleName string, iniParams string, verboseLogging int64, purge bool) error {
	if engineInitialized {
		return fmt.Errorf("G2engine is already setup and has not been torn down.")
	}
	globalG2engine.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)

	err := globalG2engine.Initialize(ctx, moduleName, iniParams, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
	if err != nil {
		return createError(5903, err)
	}

	// In case of an external database (e.g.: PostgreSQL) we need to purge since the database
	// may be shared across test suites -- this is not ideal since tests are not isolated.
	// TODO: look for a way to use external databases while still isolating tests.

	if purge {
		err = globalG2Diagnostic.PurgeRepository(ctx)
		if err != nil {
			// if an error occurred on purge make sure to destroy the engine
			defer globalG2engine.Destroy(ctx)
			return createError(5904, err)
		}
	}
	engineInitialized = true
	return err // Should be nil if we get here.
}

func setupIniParams(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5902, err)
	}
	return iniParams, err
}

func setupSenzingConfiguration(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	now := time.Now()

	aG2config := &szconfig.Szconfig{}
	err := aG2config.Initialize(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := aG2config.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	// Add data sources to in-memory Senzing configuration.

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := aG2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return createError(5908, err)
		}
	}

	// Create a string representation of the in-memory configuration.

	configStr, err := aG2config.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = aG2config.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	// Create an alternate Senzing configuration.

	configHandle, err = aG2config.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	datasourceNames = []string{"CUSTOMERS", "REFERENCE", "WATCHLIST", "EMPLOYEES"}
	for _, datasourceName := range datasourceNames {
		datasourceJson := fmt.Sprintf(`{"DSRC_CODE": "%v"}`, datasourceName)
		_, err := aG2config.AddDataSource(ctx, configHandle, datasourceJson)
		if err != nil {
			return createError(5908, err)
		}
	}

	altConfigStr, err := aG2config.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = aG2config.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return createError(5911, err)
	}

	// Persist the Senzing configurations to the Senzing repository.

	aG2configmgr := &szconfigmanager.Szconfigmanager{}
	err = aG2configmgr.Initialize(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComment := fmt.Sprintf("Created by g2engine_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComment)
	if err != nil {
		return createError(5913, err)
	}

	senzingConfigId = configID

	err = aG2configmgr.SetDefaultConfigId(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	configComment = fmt.Sprintf("Alternate config created by g2engine_test at %s", now.UTC())
	_, err = aG2configmgr.AddConfig(ctx, altConfigStr, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownG2engine(ctx)
	return err
}

func teardownG2engine(ctx context.Context) error {
	if !engineInitialized {
		return nil
	}
	err := globalG2engine.Destroy(ctx)
	if err != nil {
		return err
	}
	engineInitialized = false
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2engine_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
}

func TestG2engine_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	actual := g2engine.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2engine_AddRecord_G2BadInput(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "ADD_TEST_ERR_1", "NAME_FULL": "NOBODY NOMATCH"}`)
	testErrorBasic(test, err)
	record2, err := record.NewRecord(`{"DATA_SOURCE": "BOB", "RECORD_ID": "ADD_TEST_ERR_2", "NAME_FULL": "ERR BAD SOURCE"}`)
	testErrorBasic(test, err)

	// this one should succeed
	_, err = g2engine.AddRecord(ctx, record1.DataSource, record1.Id, record1.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	defer g2engine.DeleteRecord(ctx, record1.DataSource, record1.Id, sz.SZ_NO_FLAGS)

	// this one should fail
	_, err = g2engine.AddRecord(ctx, "CUSTOMERS", record2.Id, record2.Json, sz.SZ_NO_FLAGS)
	assert.True(test, szerror.Is(err, szerror.SzBadInput))

	// clean-up the records we inserted
	_, err = g2engine.DeleteRecord(ctx, record1.DataSource, record1.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_AddRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "ADD_TEST_1", "NAME_FULL": "JIMMITY UNKNOWN"}`)
	testErrorBasic(test, err)

	record2, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "ADD_TEST_2", "NAME_FULL": "SOMEBODY NOTFOUND"}`)
	testErrorBasic(test, err)

	_, err = g2engine.AddRecord(ctx, record1.DataSource, record1.Id, record1.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	defer g2engine.DeleteRecord(ctx, record1.DataSource, record1.Id, sz.SZ_NO_FLAGS)

	_, err = g2engine.AddRecord(ctx, record2.DataSource, record2.Id, record2.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	defer g2engine.DeleteRecord(ctx, record2.DataSource, record2.Id, sz.SZ_NO_FLAGS)

	// Clean-up the records we inserted.

	_, err = g2engine.DeleteRecord(ctx, record1.DataSource, record1.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)

	_, err = g2engine.DeleteRecord(ctx, record2.DataSource, record2.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
}

// func TestG2engine_AddRecordWithInfo(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "WITH_INFO_1", "NAME_FULL": "HUBERT WITHINFO"}`)
// 	testErrorBasic(test, err)

// 	flags := int64(0)
// 	actual, err := g2engine.AddRecordWithInfo(ctx, record.DataSource, record.Id, record.Json, sz.SZ_NO_FLAGS, flags)
// 	testError(test, ctx, g2engine, err)
// 	defer g2engine.DeleteRecord(ctx, record.DataSource, record.Id, sz.SZ_NO_FLAGS)
// 	printActual(test, actual)

// 	_, err = g2engine.DeleteRecord(ctx, record.DataSource, record.Id, sz.SZ_NO_FLAGS)
// 	testError(test, ctx, g2engine, err)
// }

func TestG2engine_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.CountRedoRecords(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

// func TestG2engine_ExportCSVEntityReport(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	expected := []string{
// 		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
// 		`1,0,0,"","CUSTOMERS","1001"`,
// 		`1,0,1,"+NAME+DOB+PHONE","CUSTOMERS","1002"`,
// 		`1,0,1,"+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
// 	}
// 	csvColumnList := ""
// 	flags := int64(-1)
// 	aHandle, err := g2engine.ExportCSVEntityReport(ctx, csvColumnList, flags)
// 	defer func() {
// 		err := g2engine.CloseExport(ctx, aHandle)
// 		testError(test, ctx, g2engine, err)
// 	}()
// 	testError(test, ctx, g2engine, err)
// 	actualCount := 0
// 	for {
// 		actual, err := g2engine.FetchNext(ctx, aHandle)
// 		testError(test, ctx, g2engine, err)
// 		if len(actual) == 0 {
// 			break
// 		}
// 		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual))
// 		actualCount += 1
// 	}
// 	assert.Equal(test, len(expected), actualCount)
// }

// func TestG2engine_ExportCSVEntityReportIterator(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	expected := []string{
// 		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
// 		`1,0,0,"","CUSTOMERS","1001"`,
// 		`1,0,1,"+NAME+DOB+PHONE","CUSTOMERS","1002"`,
// 		`1,0,1,"+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
// 	}
// 	csvColumnList := ""
// 	flags := int64(-1)
// 	actualCount := 0
// 	for actual := range g2engine.ExportCSVEntityReportIterator(ctx, csvColumnList, flags) {
// 		testError(test, ctx, g2engine, actual.Error)
// 		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))
// 		actualCount += 1
// 	}
// 	assert.Equal(test, len(expected), actualCount)
// }

func TestG2engine_ExportJSONEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	_, err := g2engine.AddRecord(ctx, aRecord.DataSource, aRecord.Id, aRecord.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	defer g2engine.DeleteRecord(ctx, aRecord.DataSource, aRecord.Id, sz.SZ_NO_FLAGS)
	flags := int64(-1)
	aHandle, err := g2engine.ExportJsonEntityReport(ctx, flags)
	defer func() {
		err := g2engine.CloseExport(ctx, aHandle)
		testError(test, ctx, g2engine, err)
	}()
	testError(test, ctx, g2engine, err)
	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := g2engine.FetchNext(ctx, aHandle)
		testError(test, ctx, g2engine, err)
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}
	testError(test, ctx, g2engine, err)
	assert.True(test, len(jsonEntityReport) > 65536)
}

func TestG2engine_ExportJSONEntityReportIterator(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	flags := int64(-1)
	actualCount := 0
	for actual := range g2engine.ExportJsonEntityReportIterator(ctx, flags) {
		testError(test, ctx, g2engine, actual.Error)
		printActual(test, actual.Value)
		actualCount += 1
	}
	assert.Equal(test, 1, actualCount)
}

func TestG2engine_FindNetworkByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	actual, err := g2engine.FindNetworkByEntityId(ctx, entityList, maxDegree, buildOutDegree, maxEntities, sz.SZ_NO_FLAGS)
	testErrorNoFail(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindNetworkByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.Id + `"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	actual, err := g2engine.FindNetworkByRecordId(ctx, recordList, maxDegree, buildOutDegree, maxEntities, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	actual, err := g2engine.FindPathByEntityId(ctx, entityID1, entityID2, maxDegree, sz.SZ_NO_EXCLUSIONS, sz.SZ_NO_REQUIRED_DATASOURCES, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	actual, err := g2engine.FindPathByRecordId(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, sz.SZ_NO_EXCLUSIONS, sz.SZ_NO_REQUIRED_DATASOURCES, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

// func TestG2engine_FindPathExcludingByEntityID(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	entityID1 := getEntityId(record1)
// 	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
// 	maxDegree := int64(1)
// 	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
// 	actual, err := g2engine.FindPathExcludingByEntityId(ctx, entityID1, entityID2, maxDegree, excludedEntities)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_FindPathExcludingByEntityID_V2(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	entityID1 := getEntityId(record1)
// 	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
// 	maxDegree := int64(1)
// 	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
// 	flags := int64(0)
// 	actual, err := g2engine.FindPathExcludingByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, flags)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_FindPathExcludingByRecordID(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	record2 := truthset.CustomerRecords["1002"]
// 	maxDegree := int64(1)
// 	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
// 	actual, err := g2engine.FindPathExcludingByRecordID(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedRecords)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_FindPathExcludingByRecordID_V2(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	record2 := truthset.CustomerRecords["1002"]
// 	maxDegree := int64(1)
// 	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
// 	flags := int64(0)
// 	actual, err := g2engine.FindPathExcludingByRecordID_V2(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedRecords, flags)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_FindPathIncludingSourceByEntityID(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	entityID1 := getEntityId(record1)
// 	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
// 	maxDegree := int64(1)
// 	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
// 	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
// 	actual, err := g2engine.FindPathIncludingSourceByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_FindPathIncludingSourceByEntityID_V2(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	entityID1 := getEntityId(record1)
// 	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
// 	maxDegree := int64(1)
// 	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
// 	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
// 	flags := int64(0)
// 	actual, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs, flags)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_FindPathIncludingSourceByRecordID(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	record2 := truthset.CustomerRecords["1002"]
// 	maxDegree := int64(1)
// 	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
// 	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
// 	actual, err := g2engine.FindPathIncludingSourceByRecordID(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedEntities, requiredDsrcs)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_FindPathIncludingSourceByRecordID_V2(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	record1 := truthset.CustomerRecords["1001"]
// 	record2 := truthset.CustomerRecords["1002"]
// 	maxDegree := int64(1)
// 	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
// 	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
// 	flags := int64(0)
// 	actual, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedEntities, requiredDsrcs, flags)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

func TestG2engine_GetActiveConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.GetActiveConfigId(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	actual, err := g2engine.GetEntityByEntityId(ctx, entityID, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.GetEntityByRecordId(ctx, record.DataSource, record.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.GetRecord(ctx, record.DataSource, record.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetRedoRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.GetRedoRecord(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetRepositoryLastModifiedTime(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetVirtualEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	actual, err := g2engine.GetVirtualEntityByRecordId(ctx, recordList, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_HowEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	actual, err := g2engine.HowEntityByEntityId(ctx, entityID, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	err := g2engine.PrimeEngine(ctx)
	testError(test, ctx, g2engine, err)
}

// func TestG2engine_ProcessRedoRecord(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	actual, err := g2engine.ProcessRedoRecord(ctx)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

// func TestG2engine_ProcessRedoRecordWithInfo(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	flags := int64(0)
// 	actual, actualInfo, err := g2engine.ProcessRedoRecordWithInfo(ctx, flags)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// 	printResult(test, "Actual Info", actualInfo)
// }

func TestG2engine_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	_, err := g2engine.ReevaluateEntity(ctx, entityID, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_ReevaluateEntity_WithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	actual, err := g2engine.ReevaluateEntity(ctx, entityID, sz.SZ_WITH_INFO)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	_, err := g2engine.ReevaluateRecord(ctx, record.DataSource, record.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_ReevaluateRecord_WithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.ReevaluateRecord(ctx, record.DataSource, record.Id, sz.SZ_WITH_INFO)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_ReplaceRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1984", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "CUSTOMERS", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "1001", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`
	_, err := g2engine.ReplaceRecord(ctx, dataSourceCode, recordID, jsonData, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)

	record := truthset.CustomerRecords["1001"]
	_, err = g2engine.ReplaceRecord(ctx, record.DataSource, record.Id, record.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
}

// FIXME: Remove after GDEV-3576 is fixed
func TestG2engine_ReplaceRecord_WithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1985", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "CUSTOMERS", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "1001", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`
	actual, err := g2engine.ReplaceRecord(ctx, dataSourceCode, recordID, jsonData, sz.SZ_WITH_INFO)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
	record := truthset.CustomerRecords["1001"]
	_, err = g2engine.ReplaceRecord(ctx, record.DataSource, record.Id, record.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	actual, err := g2engine.SearchByAttributes(ctx, attributes, sz.SZ_NO_SEARCH_PROFILE, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetStats(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.GetStats(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	actual, err := g2engine.WhyEntities(ctx, entityID1, entityID2, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	actual, err := g2engine.WhyRecords(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_Initialize(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	settings, err := getIniParams()
	testError(test, ctx, g2engine, err)
	err = g2engine.Initialize(ctx, instanceName, settings, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	initConfigID, err := g2engine.GetActiveConfigId(ctx)
	testError(test, ctx, g2engine, err)
	err = g2engine.Reinitialize(ctx, initConfigID)
	testError(test, ctx, g2engine, err)
	printActual(test, initConfigID)
}

func TestG2engine_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)

	// first create and add the record to be deleted
	record, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "DELETE_TEST", "NAME_FULL": "GONNA B. DELETED"}`)
	testError(test, ctx, g2engine, err)

	_, err = g2engine.AddRecord(ctx, record.DataSource, record.Id, record.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)

	// now delete the record
	_, err = g2engine.DeleteRecord(ctx, record.DataSource, record.Id, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_DeleteRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)

	// first create and add the record to be deleted
	record, err := record.NewRecord(`{"DATA_SOURCE": "TEST", "RECORD_ID": "DELETE_TEST", "NAME_FULL": "DELETE W. INFO"}`)
	testError(test, ctx, g2engine, err)

	_, err = g2engine.AddRecord(ctx, record.DataSource, record.Id, record.Json, sz.SZ_NO_FLAGS)
	testError(test, ctx, g2engine, err)

	// now delete the record
	actual, err := g2engine.DeleteRecord(ctx, record.DataSource, record.Id, sz.SZ_WITH_INFO)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	err := g2engine.Destroy(ctx)
	testError(test, ctx, g2engine, err)
	engineInitialized = false
	restoreG2engine(ctx) // put everything back the way it was
}
