package g2diagnostic

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/g2-sdk-go-base/g2config"
	"github.com/senzing-garage/g2-sdk-go-base/g2configmgr"
	"github.com/senzing-garage/g2-sdk-go-base/g2engine"
	"github.com/senzing-garage/g2-sdk-go/g2api"
	g2diagnosticapi "github.com/senzing-garage/g2-sdk-go/g2diagnostic"
	"github.com/senzing-garage/g2-sdk-go/g2error"
	futil "github.com/senzing-garage/go-common/fileutil"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
	"github.com/senzing-garage/go-common/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	moduleName        = "Diagnostic Test Module"
	printResults      = false
	verboseLogging    = 0
)

var (
	defaultConfigID       int64
	diagnosticInitialized bool         = false
	globalG2diagnostic    G2diagnostic = G2diagnostic{}
	logger                logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(logger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) g2api.G2diagnostic {
	return &globalG2diagnostic
}

func getG2Diagnostic(ctx context.Context) g2api.G2diagnostic {
	return &globalG2diagnostic
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

func testError(test *testing.T, ctx context.Context, g2diagnostic g2api.G2diagnostic, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2diagnostic g2api.G2diagnostic, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
	}
}

func baseDirectoryPath() string {
	return filepath.FromSlash("../target/test/g2diagnostic")
}

func dbTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getDefaultConfigID() int64 {
	return defaultConfigID
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
		if g2error.Is(err, g2error.G2BadInput) {
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

func restoreG2diagnostic(ctx context.Context) error {
	iniParams, err := getIniParams()
	if err != nil {
		return err
	}
	err = setupG2diagnostic(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}
	return nil
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	moduleName := "Test module name"
	verboseLogging := int64(0)
	logger, err = logging.NewSenzingSdkLogger(ComponentId, g2diagnosticapi.IdMessages)
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

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5920, err)
	}

	// Add records.

	err = setupAddRecords(ctx, moduleName, iniParams, verboseLogging, dbPurge)
	if err != nil {
		return createError(5922, err)
	}

	// Setup the G2 diagnostic object.

	err = setupG2diagnostic(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	return err
}

func setupAddRecords(ctx context.Context, moduleName string, iniParams string, verboseLogging int64, purge bool) error {

	aG2engine := &g2engine.G2engine{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5916, err)
	}

	// If requested, purge existing database.

	if purge {
		err = aG2engine.PurgeRepository(ctx)
		if err != nil {
			aG2engine.Destroy(ctx) // If an error occurred on purge make sure to destroy the engine.
			return createError(5904, err)
		}
	}

	// Add records into Senzing.

	testRecordIds := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	for _, testRecordId := range testRecordIds {
		testRecord := truthset.CustomerRecords[testRecordId]
		err := aG2engine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Json, "G2Diagnostic_test")
		if err != nil {
			return createError(5917, err)
		}
	}

	// All done. Destroy Senzing engine

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return createError(5905, err)
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

func setupG2diagnostic(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	if diagnosticInitialized {
		return fmt.Errorf("G2diagnostic is already setup and has not been torn down.")
	}
	globalG2diagnostic.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalG2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5903, err)
	}

	diagnosticInitialized = true
	return err
}

func setupIniParams(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		err = createError(5902, err)
	}
	return iniParams, err
}

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
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
	defaultConfigID = configID

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownG2diagnostic(ctx)
	return err
}

func teardownG2diagnostic(ctx context.Context) error {
	if !diagnosticInitialized {
		return nil
	}
	err := globalG2diagnostic.Destroy(ctx)
	if err != nil {
		return err
	}
	diagnosticInitialized = false
	return nil
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2diagnostic_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
}

func TestG2diagnostic_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	actual := g2diagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2diagnostic_CheckDBPerf(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := g2diagnostic.CheckDBPerf(ctx, secondsToRun)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

// func TestG2diagnostic_EntityListBySize(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	aSize := 1000
// 	aHandle, err := g2diagnostic.GetEntityListBySize(ctx, aSize)
// 	testError(test, ctx, g2diagnostic, err)
// 	anEntity, err := g2diagnostic.FetchNextEntityBySize(ctx, aHandle)
// 	testError(test, ctx, g2diagnostic, err)
// 	printResult(test, "Entity", anEntity)
// 	err = g2diagnostic.CloseEntityListBySize(ctx, aHandle)
// 	testError(test, ctx, g2diagnostic, err)
// }

// func TestG2diagnostic_FindEntitiesByFeatureIDs(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	features := "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}"
// 	actual, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, features)
// 	testError(test, ctx, g2diagnostic, err)
// 	printResult(test, "len(Actual)", len(actual))
// }

func TestG2diagnostic_GetAvailableMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetAvailableMemory(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, int64(0))
	printActual(test, actual)
}

// func TestG2diagnostic_GetDataSourceCounts(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	actual, err := g2diagnostic.GetDataSourceCounts(ctx)
// 	testError(test, ctx, g2diagnostic, err)
// 	printResult(test, "Data Source counts", actual)
// }

func TestG2diagnostic_GetDBInfo(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetDBInfo(ctx)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

// func TestG2diagnostic_GetEntityDetails(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	entityID := int64(1)
// 	includeInternalFeatures := 1
// 	actual, err := g2diagnostic.GetEntityDetails(ctx, entityID, includeInternalFeatures)
// 	testErrorNoFail(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

// func TestG2diagnostic_GetEntityResume(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	entityID := int64(1)
// 	actual, err := g2diagnostic.GetEntityResume(ctx, entityID)
// 	testErrorNoFail(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

// func TestG2diagnostic_GetEntitySizeBreakdown(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	minimumEntitySize := 1
// 	includeInternalFeatures := 1
// 	actual, err := g2diagnostic.GetEntitySizeBreakdown(ctx, minimumEntitySize, includeInternalFeatures)
// 	testError(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

// func TestG2diagnostic_GetFeature(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	libFeatID := int64(1)
// 	actual, err := g2diagnostic.GetFeature(ctx, libFeatID)
// 	testErrorNoFail(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

// func TestG2diagnostic_GetGenericFeatures(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	featureType := "PHONE"
// 	maximumEstimatedCount := 10
// 	actual, err := g2diagnostic.GetGenericFeatures(ctx, featureType, maximumEstimatedCount)
// 	testError(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

func TestG2diagnostic_GetLogicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetLogicalCores(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, 0)
	printActual(test, actual)
}

// func TestG2diagnostic_GetMappingStatistics(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	includeInternalFeatures := 1
// 	actual, err := g2diagnostic.GetMappingStatistics(ctx, includeInternalFeatures)
// 	testError(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

func TestG2diagnostic_GetPhysicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetPhysicalCores(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, 0)
	printActual(test, actual)
}

// func TestG2diagnostic_GetRelationshipDetails(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	relationshipID := int64(1)
// 	includeInternalFeatures := 1
// 	actual, err := g2diagnostic.GetRelationshipDetails(ctx, relationshipID, includeInternalFeatures)
// 	testErrorNoFail(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

// func TestG2diagnostic_GetResolutionStatistics(test *testing.T) {
// 	ctx := context.TODO()
// 	g2diagnostic := getTestObject(ctx, test)
// 	actual, err := g2diagnostic.GetResolutionStatistics(ctx)
// 	testError(test, ctx, g2diagnostic, err)
// 	printActual(test, actual)
// }

func TestG2diagnostic_GetTotalSystemMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetTotalSystemMemory(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, int64(0))
	printActual(test, actual)
}

func TestG2diagnostic_Init(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := &G2diagnostic{}
	moduleName := "Test module name"
	verboseLogging := int64(0)
	iniParams, err := getIniParams()
	testError(test, ctx, g2diagnostic, err)
	err = g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2diagnostic, err)
}

func TestG2diagnostic_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := &G2diagnostic{}
	moduleName := "Test module name"
	initConfigID := int64(1)
	verboseLogging := int64(0)
	iniParams, err := getIniParams()
	testError(test, ctx, g2diagnostic, err)
	err = g2diagnostic.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	testError(test, ctx, g2diagnostic, err)
}

func TestG2diagnostic_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	initConfigID := getDefaultConfigID()
	err := g2diagnostic.Reinit(ctx, initConfigID)
	testErrorNoFail(test, ctx, g2diagnostic, err)
}

func TestG2diagnostic_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	err := g2diagnostic.Destroy(ctx)
	testError(test, ctx, g2diagnostic, err)

	// restore the state that existed prior to this test
	diagnosticInitialized = false
	restoreG2diagnostic(ctx)
}
