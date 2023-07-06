package g2engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go-base/g2config"
	"github.com/senzing/g2-sdk-go-base/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2api"
	g2engineapi "github.com/senzing/g2-sdk-go/g2engine"
	"github.com/senzing/g2-sdk-go/g2error"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-common/record"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-logging/logging"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	loadId            = "G2Engine_test"
	printResults      = false
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityId int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

var (
	g2engineSingleton g2api.G2engine
	localLogger       logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(localLogger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) g2api.G2engine {
	if g2engineSingleton == nil {
		g2engineSingleton = &G2engine{}
		g2engineSingleton.SetLogLevel(ctx, logging.LevelInfoName)
		log.SetFlags(0)
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = g2engineSingleton.Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return g2engineSingleton
}

func getG2Engine(ctx context.Context) g2api.G2engine {
	if g2engineSingleton == nil {
		g2engineSingleton = &G2engine{}
		g2engineSingleton.SetLogLevel(ctx, logging.LevelInfoName)
		log.SetFlags(0)
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if err != nil {
			fmt.Println(err)
		}
		g2engineSingleton.Init(ctx, moduleName, iniParams, verboseLogging)
	}
	return g2engineSingleton
}

func getEntityIdForRecord(datasource string, id string) int64 {
	ctx := context.TODO()
	var result int64 = 0
	g2engine := getG2Engine(ctx)
	response, err := g2engine.GetEntityByRecordID(ctx, datasource, id)
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

func testError(test *testing.T, ctx context.Context, g2engine g2api.G2engine, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2engine g2api.G2engine, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if g2error.Is(err, g2error.G2Unrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2Retryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2BadUserInput) {
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	now := time.Now()

	aG2config := &g2config.G2config{}
	err := aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return createError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := aG2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return createError(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return createError(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgr{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return createError(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	aG2engine := &G2engine{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return createError(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return createError(5905, err)
	}
	return err
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	moduleName := "Test module name"
	verboseLogging := 0
	localLogger, err = logging.NewSenzingSdkLogger(ComponentId, g2engineapi.IdMessages)
	if err != nil {
		return createError(5901, err)
	}

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		return createError(5902, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5921, err)
	}
	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestBuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
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

func TestG2engine_AddRecord_G2Unrecoverable(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record2Json := `{"DATA_SOURCE": "BOB", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	err := g2engine.AddRecord(ctx, record1.DataSource, record1.Id, record1.Json, loadId)
	testError(test, ctx, g2engine, err)
	err = g2engine.AddRecord(ctx, record2.DataSource, record2.Id, record2Json, loadId)
	assert.True(test, g2error.Is(err, g2error.G2Unrecoverable))
}

func TestG2engine_AddRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	err := g2engine.AddRecord(ctx, record1.DataSource, record1.Id, record1.Json, loadId)
	testError(test, ctx, g2engine, err)
	err = g2engine.AddRecord(ctx, record2.DataSource, record2.Id, record2.Json, loadId)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_AddRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	var flags int64 = 0
	actual, err := g2engine.AddRecordWithInfo(ctx, record.DataSource, record.Id, record.Json, loadId, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_AddRecordWithInfoWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	//	record := truthset.TestRecordsWithoutRecordId[0]
	var record_dsrc string = "TEST"
	var record_json string = "{\"ADDR_LINE1\":\"123 Commerce Street, Las Vegas NV 89101\",\"ADDR_TYPE\":\"MAILING\",\"AMOUNT\":\"100\",\"DATE\":\"1/2/18\",\"DATE_OF_BIRTH\":\"12/11/1978\",\"EMAIL_ADDRESS\":\"bsmith@work.com\",\"PHONE_NUMBER\":\"702-919-1300\",\"PHONE_TYPE\":\"HOME\",\"PRIMARY_NAME_FIRST\":\"Robert\",\"PRIMARY_NAME_LAST\":\"Smith\",\"RECORD_TYPE\":\"PERSON\",\"STATUS\":\"Active\"}"
	var flags int64 = 0
	//	actual, actualRecordID, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, record.DataSource, record.Json, loadId, flags)
	actual, actualRecordID, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, record_dsrc, record_json, loadId, flags)
	testError(test, ctx, g2engine, err)
	printResult(test, "Actual RecordID", actualRecordID)
	printActual(test, actual)
}

func TestG2engine_AddRecordWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.TestRecordsWithoutRecordId[1]
	actual, err := g2engine.AddRecordWithReturnedRecordID(ctx, record.DataSource, record.Json, loadId)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_CheckRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	recordQueryList := `{"RECORDS": [{"DATA_SOURCE": "` + record.DataSource + `","RECORD_ID": "` + record.Id + `"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "123456789"}]}`
	actual, err := g2engine.CheckRecord(ctx, record.Json, recordQueryList)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.CountRedoRecords(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

// FAIL:
func TestG2engine_ExportJSONEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	flags := int64(0)
	aHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	testError(test, ctx, g2engine, err)
	anEntity, err := g2engine.FetchNext(ctx, aHandle)
	testError(test, ctx, g2engine, err)
	printResult(test, "Entity", anEntity)
	err = g2engine.CloseExport(ctx, aHandle)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_ExportConfigAndConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actualConfig, actualConfigId, err := g2engine.ExportConfigAndConfigID(ctx)
	testError(test, ctx, g2engine, err)
	printResult(test, "Actual Config", actualConfig)
	printResult(test, "Actual Config ID", actualConfigId)
}

func TestG2engine_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.ExportConfig(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

//func TestG2engine_ExportCSVEntityReport(test *testing.T) {
//	ctx := context.TODO()
//	g2engine := getTestObject(ctx, test)
//	csvColumnList := ""
//	var flags int64 = 0
//	actual, err := g2engine.ExportCSVEntityReport(ctx, csvColumnList, flags)
//	testError(test, ctx, g2engine, err)
//	test.Log("Actual:", actual)
//}

func TestG2engine_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	actual, err := g2engine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	actual, err := g2engine.FindInterestingEntitiesByRecordID(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindNetworkByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegree := 2
	buildOutDegree := 1
	maxEntities := 10
	actual, err := g2engine.FindNetworkByEntityID(ctx, entityList, maxDegree, buildOutDegree, maxEntities)
	testErrorNoFail(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindNetworkByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegree := 2
	buildOutDegree := 1
	maxEntities := 10
	var flags int64 = 0
	actual, err := g2engine.FindNetworkByEntityID_V2(ctx, entityList, maxDegree, buildOutDegree, maxEntities, flags)
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
	maxDegree := 1
	buildOutDegree := 2
	maxEntities := 10
	actual, err := g2engine.FindNetworkByRecordID(ctx, recordList, maxDegree, buildOutDegree, maxEntities)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindNetworkByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.Id + `"}]}`
	maxDegree := 1
	buildOutDegree := 2
	maxEntities := 10
	var flags int64 = 0
	actual, err := g2engine.FindNetworkByRecordID_V2(ctx, recordList, maxDegree, buildOutDegree, maxEntities, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := 1
	actual, err := g2engine.FindPathByEntityID(ctx, entityID1, entityID2, maxDegree)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := 1
	var flags int64 = 0
	actual, err := g2engine.FindPathByEntityID_V2(ctx, entityID1, entityID2, maxDegree, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := 1
	actual, err := g2engine.FindPathByRecordID(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := 1
	var flags int64 = 0
	actual, err := g2engine.FindPathByRecordID_V2(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathExcludingByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	actual, err := g2engine.FindPathExcludingByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathExcludingByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	var flags int64 = 0
	actual, err := g2engine.FindPathExcludingByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathExcludingByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := 1
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	actual, err := g2engine.FindPathExcludingByRecordID(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedRecords)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathExcludingByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := 1
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	var flags int64 = 0
	actual, err := g2engine.FindPathExcludingByRecordID_V2(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedRecords, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathIncludingSourceByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	actual, err := g2engine.FindPathIncludingSourceByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathIncludingSourceByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	var flags int64 = 0
	actual, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathIncludingSourceByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	actual, err := g2engine.FindPathIncludingSourceByRecordID(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedEntities, requiredDsrcs)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_FindPathIncludingSourceByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	var flags int64 = 0
	actual, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, maxDegree, excludedEntities, requiredDsrcs, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetActiveConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.GetActiveConfigID(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	actual, err := g2engine.GetEntityByEntityID(ctx, entityID)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	actual, err := g2engine.GetEntityByEntityID_V2(ctx, entityID, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.GetEntityByRecordID(ctx, record.DataSource, record.Id)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	actual, err := g2engine.GetEntityByRecordID_V2(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.GetRecord(ctx, record.DataSource, record.Id)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetRecord_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	actual, err := g2engine.GetRecord_V2(ctx, record.DataSource, record.Id, flags)
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

func TestG2engine_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	actual, err := g2engine.GetVirtualEntityByRecordID(ctx, recordList)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_GetVirtualEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	var flags int64 = 0
	actual, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, recordList, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_HowEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	actual, err := g2engine.HowEntityByEntityID(ctx, entityID)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_HowEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	actual, err := g2engine.HowEntityByEntityID_V2(ctx, entityID, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	err := g2engine.PrimeEngine(ctx)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_Process(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	err := g2engine.Process(ctx, record.Json)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_ProcessRedoRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.ProcessRedoRecord(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_ProcessRedoRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	var flags int64 = 0
	actual, actualInfo, err := g2engine.ProcessRedoRecordWithInfo(ctx, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
	printResult(test, "Actual Info", actualInfo)
}

func TestG2engine_ProcessWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	actual, err := g2engine.ProcessWithInfo(ctx, record.Json, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_ProcessWithResponse(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.ProcessWithResponse(ctx, record.Json)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_ProcessWithResponseResize(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.ProcessWithResponseResize(ctx, record.Json)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	err := g2engine.ReevaluateEntity(ctx, entityID, flags)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_ReevaluateEntityWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	actual, err := g2engine.ReevaluateEntityWithInfo(ctx, entityID, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	err := g2engine.ReevaluateRecord(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_ReevaluateRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	actual, err := g2engine.ReevaluateRecordWithInfo(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

// FIXME: Remove after GDEV-3576 is fixed
// func TestG2engine_ReplaceRecord(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	dataSourceCode := "CUSTOMERS"
// 	recordID := "1001"
// 	jsonData := `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1984", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "CUSTOMERS", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "1001", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`
// 	loadID := "CUSTOMERS"
// 	err := g2engine.ReplaceRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
// 	testError(test, ctx, g2engine, err)
// }

// FIXME: Remove after GDEV-3576 is fixed
// func TestG2engine_ReplaceRecordWithInfo(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	dataSourceCode := "CUSTOMERS"
// 	recordID := "1001"
// 	jsonData := `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1985", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "CUSTOMERS", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "1001", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`
// 	loadID := "CUSTOMERS"
// 	var flags int64 = 0
// 	actual, err := g2engine.ReplaceRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
// 	testError(test, ctx, g2engine, err)
// 	printActual(test, actual)
// }

func TestG2engine_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	actual, err := g2engine.SearchByAttributes(ctx, jsonData)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_SearchByAttributes_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	var flags int64 = 0
	actual, err := g2engine.SearchByAttributes_V2(ctx, jsonData, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_Stats(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	actual, err := g2engine.Stats(ctx)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	actual, err := g2engine.WhyEntities(ctx, entityID1, entityID2)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyEntities_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	var flags int64 = 0
	actual, err := g2engine.WhyEntities_V2(ctx, entityID1, entityID2, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	actual, err := g2engine.WhyEntityByEntityID(ctx, entityID)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	actual, err := g2engine.WhyEntityByEntityID_V2(ctx, entityID, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	actual, err := g2engine.WhyEntityByRecordID(ctx, record.DataSource, record.Id)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	actual, err := g2engine.WhyEntityByRecordID_V2(ctx, record.DataSource, record.Id, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	actual, err := g2engine.WhyRecords(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_WhyRecords_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	var flags int64 = 0
	actual, err := g2engine.WhyRecords_V2(ctx, record1.DataSource, record1.Id, record2.DataSource, record2.Id, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_Init(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, g2engine, err)
	err = g2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	moduleName := "Test module name"
	var initConfigID int64 = 1
	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, g2engine, err)
	err = g2engine.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	initConfigID, err := g2engine.GetActiveConfigID(ctx)
	testError(test, ctx, g2engine, err)
	err = g2engine.Reinit(ctx, initConfigID)
	testError(test, ctx, g2engine, err)
	printActual(test, initConfigID)
}

func TestG2engine_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	err := g2engine.DeleteRecord(ctx, record.DataSource, record.Id, loadId)
	testError(test, ctx, g2engine, err)
}

func TestG2engine_DeleteRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	var flags int64 = 0
	actual, err := g2engine.DeleteRecordWithInfo(ctx, record.DataSource, record.Id, record.Json, flags)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engine_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	err := g2engine.Destroy(ctx)
	testError(test, ctx, g2engine, err)
	g2engineSingleton = nil
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2engine_SetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2engine_GetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	result := g2engine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2engine_AddRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	err := g2engine.AddRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1002"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	loadID := "G2Engine_test"
	err := g2engine.AddRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_AddRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`
	loadID := "G2Engine_test"
	var flags int64 = 0
	result, err := g2engine.AddRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_AddRecordWithInfoWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Kellar", "PRIMARY_NAME_FIRST": "Candace", "ADDR_LINE1": "1824 AspenOak Way", "ADDR_CITY": "Elmwood Park", "ADDR_STATE": "CA", "ADDR_POSTAL_CODE": "95865", "EMAIL_ADDRESS": "info@ca-state.gov"}`
	loadID := "G2Engine_test"
	var flags int64 = 0
	result, _, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, dataSourceCode, jsonData, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 42))
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_AddRecordWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Sanders", "PRIMARY_NAME_FIRST": "Sandy", "ADDR_LINE1": "1376 BlueBell Rd", "ADDR_CITY": "Sacramento", "ADDR_STATE": "CA", "ADDR_POSTAL_CODE": "95823", "EMAIL_ADDRESS": "info@ca-state.gov"}`
	loadID := "G2Engine_test"
	result, err := g2engine.AddRecordWithReturnedRecordID(ctx, dataSourceCode, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Length of record identifier is %d hexadecimal characters.\n", len(result))
	// Output: Length of record identifier is 40 hexadecimal characters.
}

func ExampleG2engine_CheckRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	recordQueryList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "123456789"}]}`
	result, err := g2engine.CheckRecord(ctx, record, recordQueryList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"CHECK_RECORD_RESPONSE":[{"DSRC_CODE":"CUSTOMERS","RECORD_ID":"1001","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","MATCH_KEY":"","ERRULE_CODE":"","ERRULE_ID":0,"CANDIDATE_MATCH":"N","NON_GENERIC_CANDIDATE_MATCH":"N"}]}
}

func ExampleG2engine_CloseExport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	g2engine.CloseExport(ctx, responseHandle)
	// Output:
}

func ExampleG2engine_CountRedoRecords() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 1
}

func ExampleG2engine_ExportCSVEntityReport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	csvColumnList := ""
	var flags int64 = 0
	responseHandle, err := g2engine.ExportCSVEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportConfig() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.ExportConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 42))
	// Output: {"G2_CONFIG":{"CFG_ETYPE":[{"ETYPE_ID":...
}

func ExampleG2engine_ExportConfigAndConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	_, configId, err := g2engine.ExportConfigAndConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configId > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportJSONEntityReport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_FetchNext() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	anEntity, _ := g2engine.FetchNext(ctx, responseHandle)
	fmt.Println(len(anEntity) >= 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	var flags int64 = 0
	result, err := g2engine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_FindInterestingEntitiesByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	var flags int64 = 0
	result, err := g2engine.FindInterestingEntitiesByRecordID(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_FindNetworkByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegree := 2
	buildOutDegree := 1
	maxEntities := 10
	result, err := g2engine.FindNetworkByEntityID(ctx, entityList, maxDegree, buildOutDegree, maxEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2023...
}

func ExampleG2engine_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegree := 2
	buildOutDegree := 1
	maxEntities := 10
	var flags int64 = 0
	result, err := g2engine.FindNetworkByEntityID_V2(ctx, entityList, maxDegree, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindNetworkByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegree := 1
	buildOutDegree := 2
	maxEntities := 10
	result, err := g2engine.FindNetworkByRecordID(ctx, recordList, maxDegree, buildOutDegree, maxEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2023...
}

func ExampleG2engine_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegree := 1
	buildOutDegree := 2
	maxEntities := 10
	var flags int64 = 0
	result, err := g2engine.FindNetworkByRecordID_V2(ctx, recordList, maxDegree, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := 1
	result, err := g2engine.FindPathByEntityID(ctx, entityID1, entityID2, maxDegree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := 1
	var flags int64 = 0
	result, err := g2engine.FindPathByEntityID_V2(ctx, entityID1, entityID2, maxDegree, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := 1
	result, err := g2engine.FindPathByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleG2engine_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := 1
	var flags int64 = 0
	result, err := g2engine.FindPathByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	result, err := g2engine.FindPathExcludingByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	var flags int64 = 0
	result, err := g2engine.FindPathExcludingByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := 1
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	result, err := g2engine.FindPathExcludingByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := 1
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	var flags int64 = 0
	result, err := g2engine.FindPathExcludingByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	result, err := g2engine.FindPathIncludingSourceByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	var flags int64 = 0
	result, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	result, err := g2engine.FindPathIncludingSourceByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedEntities, requiredDsrcs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2engine_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	var flags int64 = 0
	result, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedEntities, requiredDsrcs, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_GetActiveConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetActiveConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_GetEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.GetEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2engine_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	var flags int64 = 0
	result, err := g2engine.GetEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_GetEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.GetEntityByRecordID(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2engine_GetEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	var flags int64 = 0
	result, err := g2engine.GetEntityByRecordID_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_GetRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.GetRecord(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","JSON_DATA":{"RECORD_TYPE":"PERSON","PRIMARY_NAME_LAST":"Smith","PRIMARY_NAME_FIRST":"Robert","DATE_OF_BIRTH":"12/11/1978","ADDR_TYPE":"MAILING","ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","PHONE_TYPE":"HOME","PHONE_NUMBER":"702-919-1300","EMAIL_ADDRESS":"bsmith@work.com","DATE":"1/2/18","STATUS":"Active","AMOUNT":"100","DATA_SOURCE":"CUSTOMERS","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"A","RECORD_ID":"1001"}}
}

func ExampleG2engine_GetRecord_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	var flags int64 = 0
	result, err := g2engine.GetRecord_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleG2engine_GetRedoRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"X"}
}

func ExampleG2engine_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_GetVirtualEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, recordList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2engine_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	var flags int64 = 0
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_HowEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.HowEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"INBOUND_FEAT_ID":20,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_USAGE_TYPE":"HOME","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":19,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"FMES"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":97,"GNR_SN":100,"GNR_GN":95,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_USAGE_TYPE":"MOBILE","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V4","MEMBER_RECORDS":[{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"2E2C94D1772FBC0E2D435DFDE28A3943495EF2D0"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","RESULT_VIRTUAL_ENTITY_ID":"V1-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE+EMAIL","ERRULE_CODE":"SF1_SNAME_CFF_CSTAB","FEATURE_SCORES":{"ADDRESS":[{"INBOUND_FEAT_ID":20,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_USAGE_TYPE":"HOME","CANDIDATE_FEAT_ID":38,"CANDIDATE_FEAT":"123 Commerce Street, Las Vegas NV 89101","CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":30,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FF"},{"INBOUND_FEAT_ID":3,"INBOUND_FEAT":"123 Main Street, Las Vegas NV 89132","INBOUND_FEAT_USAGE_TYPE":"MAILING","CANDIDATE_FEAT_ID":38,"CANDIDATE_FEAT":"123 Commerce Street, Las Vegas NV 89101","CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":80,"SCORE_BUCKET":"LIKELY","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FMES"}],"EMAIL":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":100,"GNR_SN":100,"GNR_GN":100,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_USAGE_TYPE":"MOBILE","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}},{"STEP":3,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S2","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"2E2C94D1772FBC0E2D435DFDE28A3943495EF2D0"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S2","RESULT_VIRTUAL_ENTITY_ID":"V1-S3","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FMES"}],"EMAIL":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":93,"GNR_SN":100,"GNR_GN":93,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"},{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":90,"GNR_SN":100,"GNR_GN":88,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}}],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1-S3","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"2E2C94D1772FBC0E2D435DFDE28A3943495EF2D0"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2engine_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	var flags int64 = 0
	result, err := g2engine.HowEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL"}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V4","MEMBER_RECORDS":[{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"2E2C94D1772FBC0E2D435DFDE28A3943495EF2D0"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","RESULT_VIRTUAL_ENTITY_ID":"V1-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE+EMAIL","ERRULE_CODE":"SF1_SNAME_CFF_CSTAB"}},{"STEP":3,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S2","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"2E2C94D1772FBC0E2D435DFDE28A3943495EF2D0"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S2","RESULT_VIRTUAL_ENTITY_ID":"V1-S3","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB"}}],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1-S3","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"2E2C94D1772FBC0E2D435DFDE28A3943495EF2D0"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2engine_PrimeEngine() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_SearchByAttributes() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	result, err := g2engine.SearchByAttributes(ctx, jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 1962))
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"INBOUND_FEAT":"bsmith@work.com","CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100}],"NAME":[{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Bob J Smith","GNR_FN":83,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1},{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Robert Smith","GNR_FN":88,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20}]},{"FEAT_DESC":"123 Commerce Street, Las Vegas NV 89101","LIB_FEAT_ID":38,"USAGE_TYPE":"MAILING","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Commerce Street, Las Vegas NV 89101","LIB_FEAT_ID":38}]},{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3,"USAGE_TYPE":"MAILING","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3}]}],"DOB":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2},{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":19}]}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5}]}],"NAME":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1,"USAGE_TYPE":"PRIMARY","FEAT_DESC_VALUES":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1},{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":32},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":18}]}],"PHONE":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]},{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"MOBILE","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]}],"RECOR...
}

func ExampleG2engine_SearchByAttributes_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	var flags int64 = 0
	result, err := g2engine.SearchByAttributes_V2(ctx, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1"},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}}}]}
}

func ExampleG2engine_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Stats() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.Stats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

// FIXME: Remove after GDEV-3576 is fixed
// func ExampleG2engine_WhyEntities() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := &G2engine{}
// 	var entityID1 int64 = 1
// 	var entityID2 int64 = 4
// 	result, err := g2engine.WhyEntities(ctx, entityID1, entityID2)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(truncate(result, 74))
// 	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":4,"MATCH_INFO":{"WHY_KEY":...
// }

// FIXME: Remove after GDEV-3576 is fixed
// func ExampleG2engine_WhyEntities_V2() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := &G2engine{}
// 	var entityID1 int64 = 1
// 	var entityID2 int64 = 4
// 	var flags int64 = 0
// 	result, err := g2engine.WhyEntities_V2(ctx, entityID1, entityID2, flags)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":4,"MATCH_INFO":{"WHY_KEY":"","WHY_ERRULE_CODE":"","MATCH_LEVEL_CODE":""}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":4}}]}
// }

func ExampleG2engine_WhyEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.WhyEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	var flags int64 = 0
	result, err := g2engine.WhyEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.WhyEntityByRecordID(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	var flags int64 = 0
	result, err := g2engine.WhyEntityByRecordID_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyRecords() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	result, err := g2engine.WhyRecords(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleG2engine_WhyRecords_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	var flags int64 = 0
	result, err := g2engine.WhyRecords_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":1,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_Process() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	err := g2engine.Process(ctx, record)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_ProcessRedoRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.ProcessRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2engine_ProcessRedoRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	var flags int64 = 0
	_, result, err := g2engine.ProcessRedoRecordWithInfo(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2engine_ProcessWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	var flags int64 = 0
	result, err := g2engine.ProcessWithInfo(ctx, record, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ProcessWithResponse() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	result, err := g2engine.ProcessWithResponse(ctx, record)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2engine_ProcessWithResponseResize() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	result, err := g2engine.ProcessWithResponseResize(ctx, record)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2engine_ReevaluateEntity() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	var flags int64 = 0
	err := g2engine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
func ExampleG2engine_ReevaluateEntityWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	var flags int64 = 0
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 37))
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2engine_ReevaluateRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	var flags int64 = 0
	err := g2engine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_ReevaluateRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	var flags int64 = 0
	result, err := g2engine.ReevaluateRecordWithInfo(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ReplaceRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	err := g2engine.ReplaceRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_ReplaceRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	var flags int64 = 0
	result, err := g2engine.ReplaceRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_DeleteRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	loadID := "G2Engine_test"
	err := g2engine.DeleteRecord(ctx, dataSourceCode, recordID, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_DeleteRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	loadID := "G2Engine_test"
	var flags int64 = 0
	result, err := g2engine.DeleteRecordWithInfo(ctx, dataSourceCode, recordID, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := 0
	err = g2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_InitWithConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	initConfigID := int64(1)
	verboseLogging := 0
	err = g2engine.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Reinit() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	initConfigID, _ := g2engine.GetActiveConfigID(ctx) // Example initConfigID.
	err := g2engine.Reinit(ctx, initConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_PurgeRepository() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.PurgeRepository(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
