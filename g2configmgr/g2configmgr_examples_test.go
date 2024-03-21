//go:build linux

package g2configmgr

import (
	"context"
	"fmt"

	"github.com/senzing-garage/g2-sdk-go-base/g2config"
	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2configmgr_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2configmgr_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	result := g2configmgr.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2configmgr_AddConfig() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	g2configmgr := getG2Configmgr(ctx)
	configStr, err := g2config.GetJsonString(ctx, configHandle)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	configComments := "Example configuration"
	configId, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	fmt.Println(configId > 0) // Dummy output.
	// Output: true
}

func ExampleG2configmgr_GetConfig() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configStr, err := g2configmgr.GetConfig(ctx, configId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configStr, defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleG2configmgr_GetConfigList() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	jsonConfigList, err := g2configmgr.GetConfigList(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfigList, 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleG2configmgr_GetDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configId > 0) // Dummy output.
	// Output: true
}

func ExampleG2configmgr_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	oldConfigId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	g2config := &g2config.G2config{}
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configStr, err := g2config.GetJsonString(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	configComments := "Example configuration"
	newConfigId, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		fmt.Println(err)
	}
	err = g2configmgr.ReplaceDefaultConfigId(ctx, oldConfigId, newConfigId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgr_SetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configId, err := g2configmgr.GetDefaultConfigId(ctx) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		fmt.Println(err)
	}
	err = g2configmgr.SetDefaultConfigId(ctx, configId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgr_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	err := g2configmgr.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgr_Init() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := &G2configmgr{}
	moduleName := "Test module name"
	iniParams, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	err = g2configmgr.Initialize(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgr_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	err := g2configmgr.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
