package g2config

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logging"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configSingleton g2api.G2config
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) g2api.G2config {
	if g2configSingleton == nil {
		g2configSingleton = &G2config{}
		g2configSingleton.SetLogLevel(ctx, logging.LevelTraceName)
		log.SetFlags(0)
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = g2configSingleton.Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return g2configSingleton
}

func getG2Config(ctx context.Context) g2api.G2config {
	if g2configSingleton == nil {
		g2configSingleton := &G2config{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if err != nil {
			fmt.Println(err)
		}
		err = g2configSingleton.Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return g2configSingleton
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

func testError(test *testing.T, ctx context.Context, g2config g2api.G2config, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
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

func setup() error {
	var err error = nil
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

func TestG2config_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	actual, err := g2config.AddDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_Close(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_Create(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	actual, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
}

func TestG2config_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "Original", actual)
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	_, err = g2config.AddDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	actual, err = g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "     Add", actual)
	err = g2config.DeleteDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	actual, err = g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "  Delete", actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_ListDataSources(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_Load(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	jsonConfig, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	err = g2config.Load(ctx, configHandle, jsonConfig)
	testError(test, ctx, g2config, err)
}

func TestG2config_Save(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
}

func TestG2config_Init(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, g2config, err)
	err = g2config.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2config, err)
}

func TestG2config_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	err := g2config.Destroy(ctx)
	testError(test, ctx, g2config, err)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2config_AddDataSource() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	result, err := g2config.AddDataSource(ctx, configHandle, inputJson)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DSRC_ID":1001}
}

func ExampleG2config_Close() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = g2config.Close(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Create() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2config_DeleteDataSource() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	inputJson := `{"DSRC_CODE": "TEST"}`
	err = g2config.DeleteDataSource(ctx, configHandle, inputJson)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_ListDataSources() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	result, err := g2config.ListDataSources(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleG2config_Load() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.Save(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	err = g2config.Load(ctx, configHandle, jsonConfig)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Save() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.Save(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfig, 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

func ExampleG2config_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	err := g2config.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := 0
	err = g2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	err := g2config.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
