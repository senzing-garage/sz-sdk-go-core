//go:build linux

package g2config

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2config_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2config_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
	result := g2config.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2config_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
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
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
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
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
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
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
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
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
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
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	mockConfigHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.Save(ctx, mockConfigHandle)
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

func ExampleG2config_Save() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
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
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	err := g2config.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Init() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	moduleName := "Test module name"
	iniParams, err := getIniParams()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	err = g2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2config_examples_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	err := g2config.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
