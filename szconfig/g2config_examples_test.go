//go:build linux

package szconfig

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfig_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfig_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
	result := g2config.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzconfig_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
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

func ExampleSzconfig_Close() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
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

func ExampleSzconfig_Create() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfig_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
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

func ExampleSzconfig_GetDataSources() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	result, err := g2config.GetDataSources(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleSzconfig_GetJsonString() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.GetJsonString(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfig, 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

func ExampleSzconfig_Load() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	mockConfigHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.GetJsonString(ctx, mockConfigHandle)
	if err != nil {
		fmt.Println(err)
	}
	configHandle, err := g2config.Load(ctx, jsonConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfig_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	err := g2config.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzconfig_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	moduleName := "Test module name"
	iniParams, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	err = g2config.Initialize(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzconfig_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	err := g2config.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
