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

	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
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

var logger logging.Logging

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {
	var err error

	ctx := context.TODO()

	// Create a directory for temporary files.

	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // Cleanup any previous test run.
	failOnError(5001, err)

	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // Recreate the test target directory.
	failOnError(5002, err)

	// Setup dependencies.

	databaseURL, err := setupDatabase()
	failOnError(5003, err)

	log.SetFlags(0)

	logger, err = getLogger(ctx)
	failOnError(5004, err)

	// Create a SzAbstractFactory.

	settings, err := getSettings(databaseURL)
	failOnError(5005, err)

	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   "Example instance",
		Settings:       settings,
		VerboseLogging: senzing.SzNoLogging,
	}

	// Demonstrate persisting a Senzing configuration to the Senzing repository.

	err = demonstrateConfigFunctions(ctx, szAbstractFactory)
	failOnError(5006, err)

	// Demonstrate tests.

	err = demonstrateSenzingFunctions(ctx, szAbstractFactory)
	failOnError(5007, err)

	err = szAbstractFactory.Destroy(ctx)
	failOnError(5008, err)

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
}

// ----------------------------------------------------------------------------
// Demonstrations
// ----------------------------------------------------------------------------

func demonstrateAddRecord(ctx context.Context, szEngine senzing.SzEngine) (string, error) {
	var (
		flags  = senzing.SzWithInfo
		result string
	)

	dataSourceCode := "TEST"
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	failOnError(5010, err)

	recordID := randomNumber.String()
	jsonData := fmt.Sprintf(
		"%s%s%s",
		`{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "`,
		recordID,
		`", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "SEAMAN", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`)

	// Using SzEngine: Add record and return "withInfo".

	result, err = szEngine.AddRecord(ctx, dataSourceCode, recordID, jsonData, flags)

	if err != nil {
		return "", fmt.Errorf("demonstrateAddRecord error: %w", err)
	}

	return result, nil
}

func demonstrateConfigFunctions(ctx context.Context, szAbstractFactory senzing.SzAbstractFactory) error {
	now := time.Now()

	// Create Senzing objects.

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	failOnError(5101, err)

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	failOnError(5102, err)

	// Using SzConfig: Add data source to in-memory configuration.

	for testDataSourceCode := range truthset.TruthsetDataSources {
		_, err := szConfig.AddDataSource(ctx, testDataSourceCode)
		failOnError(5104, err)
	}

	// Using SzConfig: Persist configuration to a string.

	configStr, err := szConfig.Export(ctx)
	failOnError(5105, err)

	// Using SzConfigManager: Persist configuration string to database.

	configComment := fmt.Sprintf("Created by main.go at %s", now.UTC())
	_, err = szConfigManager.SetDefaultConfig(ctx, configStr, configComment)
	failOnError(5106, err)

	return nil
}

func demonstrateSenzingFunctions(ctx context.Context, szAbstractFactory senzing.SzAbstractFactory) error {
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	failOnError(5201, err)

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	failOnError(5202, err)

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	failOnError(5203, err)

	// Clean the repository.

	err = szDiagnostic.PurgeRepository(ctx)
	failOnError(5204, err)

	// Using SzEngine: Add records with information returned.

	withInfo, err := demonstrateAddRecord(ctx, szEngine)
	failOnError(5205, err)
	logger.Log(2003, withInfo)

	// Using SzProduct: Show license metadata.

	license, err := szProduct.GetLicense(ctx)
	failOnError(5206, err)
	logger.Log(2004, license)

	// Using SzEngine: Purge repository again.

	err = szDiagnostic.PurgeRepository(ctx)
	failOnError(5207, err)

	return nil
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func copyDatabase() (string, error) {
	var result string

	// Construct SQLite database URL.

	testDirectoryPath := getTestDirectoryPath()

	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return result, fmt.Errorf("failed to make target database path (%s) absolute. Error: %w", dbTargetPath, err)
	}

	result = "sqlite3://na:na@nowhere/" + dbTargetPath

	// Copy template file to test directory.

	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	if err != nil {
		return result, fmt.Errorf("failed to obtain absolute path to database file (%s): %s", databaseTemplatePath, err.Error())
	}

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	if err != nil {
		return result, fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w", databaseTemplatePath, testDirectoryPath, err)
	}

	return result, nil
}

func failOnError(msgID int, err error) {
	if err != nil {
		logger.Log(msgID, err)
		panic(err.Error())
	}
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("testdata/sqlite/G2C.db")
}

func getLogger(ctx context.Context) (logging.Logging, error) {
	_ = ctx
	result, err := logging.NewSenzingLogger(9999, Messages)

	if err != nil {
		return nil, fmt.Errorf("getLogger error: %w", err)
	}

	return result, nil
}

func getSettings(databaseURL string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": databaseURL}

	result, err := settings.BuildSimpleSettingsUsingMap(configAttrMap)
	if err != nil {
		return result, fmt.Errorf("getSettings error: %w", err)
	}

	return result, nil
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("target/test/main")
}

func setupDatabase() (string, error) {
	databaseURL, ok := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")
	if ok {
		return databaseURL, nil
	}

	return copyDatabase()
}
