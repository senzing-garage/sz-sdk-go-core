package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/observerpb"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go-core/szproduct"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const MessageIdTemplate = "senzing-9999%04d"

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var Messages = map[int]string{
	1:    "%s",
	2:    "WithInfo: %s",
	2001: "Testing %s.",
	2002: "Physical cores: %d.",
	2003: "withInfo",
	2004: "License",
	2999: "Cannot retrieve last error message.",
}

// Values updated via "go install -ldflags" parameters.

var programName string = "unknown"
var buildVersion string = "0.0.0"
var buildIteration string = "0"
var logger logging.LoggingInterface

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func getTestDirectoryPath() string {
	return filepath.FromSlash("target/test/main")
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("testdata/sqlite/G2C.db")
}

func setupDatabase(preserveDB bool) (string, bool, error) {
	var err error = nil
	testDirectoryPath := getTestDirectoryPath()
	databaseTemplatePath := getDatabaseTemplatePath()
	databaseTemplatePath, err = filepath.Abs(databaseTemplatePath)
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			databaseTemplatePath, err.Error())
		return "", false, err
	}

	// check the environment for a database URL
	dbUrl, envUrlExists := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")

	dbTargetPath := filepath.Join(getTestDirectoryPath(), "G2C.db")

	dbTargetPath, err = filepath.Abs(dbTargetPath)
	if err != nil {
		err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
		return "", false, err
	}

	dbDefaultUrl := fmt.Sprintf("sqlite3://na:na@%s", dbTargetPath)

	dbExternal := envUrlExists && dbDefaultUrl != dbUrl

	if !dbExternal {
		// set the database URL
		dbUrl = dbDefaultUrl

		if !preserveDB {
			// copy the SQLite database file
			_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true)

			if err != nil {
				err = fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
					databaseTemplatePath, testDirectoryPath, err)
				// fall through to return the error
			}
		}
	}

	return dbUrl, dbExternal, err
}

func createSettings(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}
	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		return "", err
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

func getSzConfig(ctx context.Context) (sz.SzConfig, error) {
	result := szconfig.Szconfig{}
	instanceName := "Test name"
	verboseLogging := int64(0) // 0 for no Senzing logging; 1 for logging
	settings, err := getSettings()
	if err != nil {
		return &result, err
	}
	err = result.Initialize(ctx, instanceName, settings, verboseLogging)
	return &result, err
}

func getSzConfigManager(ctx context.Context) (sz.SzConfigManager, error) {
	result := szconfigmanager.Szconfigmanager{}
	instanceName := "Test name"
	verboseLogging := int64(0)
	settings, err := getSettings()
	if err != nil {
		return &result, err
	}
	err = result.Initialize(ctx, instanceName, settings, verboseLogging)
	return &result, err
}

func getSzDiagnostic(ctx context.Context) (sz.SzDiagnostic, error) {
	result := szdiagnostic.Szdiagnostic{}
	instanceName := "Test name"
	verboseLogging := int64(0)
	settings, err := getSettings()
	if err != nil {
		return &result, err
	}
	err = result.Initialize(ctx, instanceName, settings, verboseLogging, 0)
	return &result, err
}

func getSzEngine(ctx context.Context) (sz.SzEngine, error) {
	result := szengine.Szengine{}
	moduleName := "Test name"
	verboseLogging := int64(0)
	settings, err := getSettings()
	if err != nil {
		return &result, err
	}
	err = result.Initialize(ctx, moduleName, settings, verboseLogging, 0)
	return &result, err
}

func getSzProduct(ctx context.Context) (sz.SzProduct, error) {
	result := szproduct.Szproduct{}
	moduleName := "Test module name"
	verboseLogging := int64(0)
	settings, err := getSettings()
	if err != nil {
		return &result, err
	}
	err = result.Initialize(ctx, moduleName, settings, verboseLogging)
	return &result, err
}

func getLogger(ctx context.Context) (logging.LoggingInterface, error) {
	_ = ctx
	logger, err := logging.NewSenzingLogger("my-unique-%04d", Messages)
	if err != nil {
		fmt.Println(err)
	}
	return logger, err
}

func demonstrateConfigFunctions(ctx context.Context, szConfig sz.SzConfig, szConfigManager sz.SzConfigManager) error {
	now := time.Now()

	// Using SzConfig: Create a default configuration in memory.

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return logger.NewError(5100, err)
	}

	// Using SzConfig: Add data source to in-memory configuration.

	for testDataSourceCode, _ := range truthset.TruthsetDataSources {
		_, err := szConfig.AddDataSource(ctx, configHandle, testDataSourceCode)
		if err != nil {
			return logger.NewError(5101, err)
		}
	}

	// Using SzConfig: Persist configuration to a string.

	configStr, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return logger.NewError(5102, err)
	}

	// Using SzConfigManager: Persist configuration string to database.

	configComment := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configStr, configComment)
	if err != nil {
		return logger.NewError(5103, err)
	}

	// Using SzConfigManager: Set new configuration as the default.

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return logger.NewError(5104, err)
	}

	return err
}

func demonstrateAddRecord(ctx context.Context, g2Engine sz.SzEngine) (string, error) {
	dataSourceCode := "TEST"
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	if err != nil {
		panic(err)
	}
	recordId := randomNumber.String()
	jsonData := fmt.Sprintf(
		"%s%s%s",
		`{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "`,
		recordId,
		`", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "SEAMAN", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`)
	var flags int64 = sz.SZ_WITH_INFO

	// Using G2Engine: Add record and return "withInfo".

	return g2Engine.AddRecord(ctx, dataSourceCode, recordId, jsonData, flags)
}

func demonstrateAdditionalFunctions(ctx context.Context, g2Diagnostic sz.SzDiagnostic, g2Engine sz.SzEngine, g2Product sz.SzProduct) error {

	err := g2Diagnostic.PurgeRepository(ctx)
	if err != nil {
		failOnError(5301, err)
	}

	// Using G2Engine: Add records with information returned.

	withInfo, err := demonstrateAddRecord(ctx, g2Engine)
	if err != nil {
		failOnError(5302, err)
	}
	logger.Log(2003, withInfo)

	// Using G2Product: Show license metadata.

	license, err := g2Product.GetLicense(ctx)
	if err != nil {
		failOnError(5303, err)
	}
	logger.Log(2004, license)

	// Using G2Engine: Purge repository again.

	err = g2Diagnostic.PurgeRepository(ctx)
	if err != nil {
		failOnError(5304, err)
	}

	return err
}

func destroyObjects(ctx context.Context, g2Config sz.SzConfig, g2Configmgr sz.SzConfigManager, g2Diagnostic sz.SzDiagnostic, g2Engine sz.SzEngine, g2Product sz.SzProduct) error {

	err := g2Config.Destroy(ctx)
	if err != nil {
		failOnError(5401, err)
	}

	err = g2Configmgr.Destroy(ctx)
	if err != nil {
		failOnError(5402, err)
	}

	err = g2Diagnostic.Destroy(ctx)
	if err != nil {
		failOnError(5403, err)
	}

	err = g2Engine.Destroy(ctx)
	if err != nil {
		failOnError(5404, err)
	}

	err = g2Product.Destroy(ctx)
	if err != nil {
		failOnError(5405, err)
	}

	return err
}

func failOnError(msgId int, err error) {
	logger.Log(msgId, err)
	panic(err.Error())
}

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {
	var err error = nil
	ctx := context.TODO()

	fmt.Printf(">>>>>> Step 1.0\n")

	// get the base directory for temporary files
	baseDir := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(baseDir)) // cleanup any previous test run
	if err != nil {
		fmt.Printf("Failed to remove target test directory: %v\n", baseDir)
		fmt.Println(err)
		return
	}

	fmt.Printf(">>>>>> Step 2.0\n")

	err = os.MkdirAll(filepath.Clean(baseDir), 0750) // recreate the test target directory
	if err != nil {
		fmt.Printf("Failed to recreate target test directory: %v\n", baseDir)
		fmt.Println(err)
		return
	}

	fmt.Printf(">>>>>> Step 3.0\n")

	// setup the database
	_, _, err = setupDatabase(false)
	if err != nil {
		fmt.Println("Failed to setup database")
		fmt.Println(err)
		return
	}

	fmt.Printf(">>>>>> Step 4.0\n")

	// Configure the "log" standard library.
	log.SetFlags(0)
	logger, err = getLogger(ctx)
	if err != nil {
		failOnError(5000, err)
	}

	// Test logger.

	programmMetadataMap := map[string]interface{}{
		"ProgramName":    programName,
		"BuildVersion":   buildVersion,
		"BuildIteration": buildIteration,
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
	logger.Log(2001, "Just a test of logging", programmMetadataMap)

	// Create observers.

	observer1 := &observer.ObserverNull{
		Id: "Observer 1",
	}
	observer2 := &observer.ObserverNull{
		Id: "Observer 2",
	}

	grpcConnection, err := grpc.Dial("localhost:8260", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Did not connect: %v\n", err)
	}

	observer3 := &observer.ObserverGrpc{
		Id:         "Observer 3",
		GrpcClient: observerpb.NewObserverClient(grpcConnection),
	}

	// Get Senzing objects for installing a Senzing Engine configuration.

	g2Config, err := getSzConfig(ctx)
	if err != nil {
		failOnError(5001, err)
	}
	err = g2Config.RegisterObserver(ctx, observer1)
	if err != nil {
		panic(err)
	}
	err = g2Config.RegisterObserver(ctx, observer2)
	if err != nil {
		panic(err)
	}
	err = g2Config.RegisterObserver(ctx, observer3)
	if err != nil {
		panic(err)
	}
	g2Config.SetObserverOrigin(ctx, "sz-sdk-go-core main.go")

	g2Configmgr, err := getSzConfigManager(ctx)
	if err != nil {
		failOnError(5005, err)
	}
	err = g2Configmgr.RegisterObserver(ctx, observer1)
	if err != nil {
		panic(err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	err = demonstrateConfigFunctions(ctx, g2Config, g2Configmgr)
	if err != nil {
		failOnError(5008, err)
	}

	// Now that a Senzing configuration is installed, get the remainder of the Senzing objects.

	g2Diagnostic, err := getSzDiagnostic(ctx)
	if err != nil {
		failOnError(5009, err)
	}
	err = g2Diagnostic.RegisterObserver(ctx, observer1)
	if err != nil {
		panic(err)
	}

	g2Engine, err := getSzEngine(ctx)
	if err != nil {
		failOnError(5010, err)
	}
	err = g2Engine.RegisterObserver(ctx, observer1)
	if err != nil {
		panic(err)
	}

	g2Product, err := getSzProduct(ctx)
	if err != nil {
		failOnError(5011, err)
	}
	err = g2Product.RegisterObserver(ctx, observer1)
	if err != nil {
		panic(err)
	}

	// Demonstrate tests.

	err = demonstrateAdditionalFunctions(ctx, g2Diagnostic, g2Engine, g2Product)
	if err != nil {
		failOnError(5015, err)
	}

	// Destroy Senzing objects.

	err = destroyObjects(ctx, g2Config, g2Configmgr, g2Diagnostic, g2Engine, g2Product)
	if err != nil {
		failOnError(5016, err)
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
}
