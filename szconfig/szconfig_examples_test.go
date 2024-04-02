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

func ExampleSzConfig_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzConfig_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	result := szConfig.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzConfig_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	dataSourceCode := "GO_TEST"
	result, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DSRC_ID":1001}
}

func ExampleSzConfig_Close() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = szConfig.Close(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfig_Create() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfig_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	dataSourceCode := "TEST"
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfig_GetDataSources() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	result, err := szConfig.GetDataSources(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleSzConfig_GetJsonString() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := szConfig.GetJsonString(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfig, 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"INTERNAL":"No"},{"ATTR_ID":1003,"...
}

func ExampleSzConfig_Load() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	mockConfigHandle, err := szConfig.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := szConfig.GetJsonString(ctx, mockConfigHandle)
	if err != nil {
		fmt.Println(err)
	}
	configHandle, err := szConfig.Load(ctx, jsonConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfig_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	err := szConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// func ExampleSzConfig_Initialize() {
// 	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
// 	ctx := context.TODO()
// 	szConfig := getSzconfig(ctx)
// 	moduleName := "Test module name"
// 	iniParams, err := getSettings()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	verboseLogging := int64(0)
// 	err = szConfig.Initialize(ctx, moduleName, iniParams, verboseLogging)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// Output:
// }

func ExampleSzConfig_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	err := szConfig.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
