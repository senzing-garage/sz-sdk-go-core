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
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

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
var logger logging.LoggingInterface

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func failOnError(msgID int, err error) {
	logger.Log(msgID, err)
	panic(err.Error())
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("testdata/sqlite/G2C.db")
}

func getLogger(ctx context.Context) (logging.LoggingInterface, error) {
	_ = ctx
	return logging.NewSenzingLogger("my-unique-%04d", Messages)
}

func getSettings(databaseURL string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)
	if err != nil {
		return "", err
	}
	return settings, err
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("target/test/main")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func setupDatabase() (string, error) {
	var err error
	databaseURL, ok := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")
	if !ok {

		// Construct SQLite database URL.

		testDirectoryPath := getTestDirectoryPath()
		dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
		if err != nil {
			err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
				dbTargetPath, err)
			return "", err
		}
		databaseURL = fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

		// Copy template file to test directory.

		databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
		if err != nil {
			err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
				databaseTemplatePath, err.Error())
			return "", err
		}
		_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
		if err != nil {
			return databaseURL, fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
				databaseTemplatePath, testDirectoryPath, err)
		}
	}
	return databaseURL, err
}

// ----------------------------------------------------------------------------
// Demonstrations
// ----------------------------------------------------------------------------

func demonstrateAddRecord(ctx context.Context, szEngine senzing.SzEngine) (string, error) {
	dataSourceCode := "TEST"
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	if err != nil {
		panic(err)
	}
	recordID := randomNumber.String()
	jsonData := fmt.Sprintf(
		"%s%s%s",
		`{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "`,
		recordID,
		`", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "SEAMAN", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`)
	var flags int64 = senzing.SzWithInfo

	// Using SzEngine: Add record and return "withInfo".

	return szEngine.AddRecord(ctx, dataSourceCode, recordID, jsonData, flags)
}

func demonstrateConfigFunctions(ctx context.Context, szAbstractFactory senzing.SzAbstractFactory) error {
	now := time.Now()

	// Create Senzing objects.

	szConfig, err := szAbstractFactory.CreateSzConfig(ctx)
	if err != nil {
		return logger.NewError(5101, err)
	}
	defer func() { handleError(szConfig.Destroy(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateSzConfigManager(ctx)
	if err != nil {
		return logger.NewError(5102, err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	// Using SzConfig: Create a default configuration in memory.

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return logger.NewError(5103, err)
	}

	// Using SzConfig: Add data source to in-memory configuration.

	for testDataSourceCode := range truthset.TruthsetDataSources {
		_, err := szConfig.AddDataSource(ctx, configHandle, testDataSourceCode)
		if err != nil {
			return logger.NewError(5104, err)
		}
	}

	// Using SzConfig: Persist configuration to a string.

	configStr, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return logger.NewError(5105, err)
	}

	// Using SzConfigManager: Persist configuration string to database.

	configComment := fmt.Sprintf("Created by main.go at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configStr, configComment)
	if err != nil {
		return logger.NewError(5106, err)
	}

	// Using SzConfigManager: Set new configuration as the default.

	err = szConfigManager.SetDefaultConfigId(ctx, configID)
	if err != nil {
		return logger.NewError(5107, err)
	}
	return err
}

func demonstrateSenzingFunctions(ctx context.Context, szAbstractFactory senzing.SzAbstractFactory) error {

	// Create Senzing objects.

	szDiagnostic, err := szAbstractFactory.CreateSzDiagnostic(ctx)
	if err != nil {
		return logger.NewError(9999, err)
	}
	defer func() { handleError(szDiagnostic.Destroy(ctx)) }()

	szEngine, err := szAbstractFactory.CreateSzEngine(ctx)
	if err != nil {
		return logger.NewError(9999, err)
	}
	defer func() { handleError(szEngine.Destroy(ctx)) }()

	szProduct, err := szAbstractFactory.CreateSzProduct(ctx)
	if err != nil {
		return logger.NewError(9999, err)
	}
	defer func() { handleError(szProduct.Destroy(ctx)) }()

	// Clean the repository.

	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		failOnError(5301, err)
	}

	// Using SzEngine: Add records with information returned.

	withInfo, err := demonstrateAddRecord(ctx, szEngine)
	if err != nil {
		failOnError(5302, err)
	}
	logger.Log(2003, withInfo)

	// Using SzProduct: Show license metadata.

	license, err := szProduct.GetLicense(ctx)
	if err != nil {
		failOnError(5303, err)
	}
	logger.Log(2004, license)

	// Using SzEngine: Purge repository again.

	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		failOnError(5304, err)
	}

	return err
}

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {
	var err error
	ctx := context.TODO()

	// Create a directory for temporary files.

	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // Cleanup any previous test run.
	if err != nil {
		failOnError(5001, err)
	}

	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // Recreate the test target directory.
	if err != nil {
		failOnError(5002, err)
	}

	// Setup dependencies.

	databaseURL, err := setupDatabase()
	if err != nil {
		failOnError(5003, err)
	}

	log.SetFlags(0)
	logger, err = getLogger(ctx)
	if err != nil {
		failOnError(5004, err)
	}

	// Create a SzAbstractFactory.

	settings, err := getSettings(databaseURL)
	if err != nil {
		failOnError(5005, err)
	}
	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   "Example instance",
		Settings:       settings,
		VerboseLogging: senzing.SzNoLogging,
	}

	// Demonstrate persisting a Senzing configuration to the Senzing repository.

	err = demonstrateConfigFunctions(ctx, szAbstractFactory)
	if err != nil {
		failOnError(5006, err)
	}

	// Demonstrate tests.

	err = demonstrateSenzingFunctions(ctx, szAbstractFactory)
	if err != nil {
		failOnError(5007, err)
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
}
