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
	configDefinition, err := g2config.GetJsonString(ctx, configHandle)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	configComments := "Example configuration"
	configId, err := g2configmgr.AddConfig(ctx, configDefinition, configComments)
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
	defaultConfigId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configDefinition, err := g2configmgr.GetConfig(ctx, defaultConfigId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configDefinition, defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleG2configmgr_GetConfigList() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configList, err := g2configmgr.GetConfigList(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configList, 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleG2configmgr_GetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	defaultConfigId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(defaultConfigId > 0) // Dummy output.
	// Output: true
}

func ExampleG2configmgr_ReplaceDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	currentDefaultConfigId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	g2config := &g2config.G2config{}
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configDefinition, err := g2config.GetJsonString(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	configComments := "Example configuration"
	newDefaultConfigId, err := g2configmgr.AddConfig(ctx, configDefinition, configComments)
	if err != nil {
		fmt.Println(err)
	}
	err = g2configmgr.ReplaceDefaultConfigId(ctx, currentDefaultConfigId, newDefaultConfigId)
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

func ExampleG2configmgr_Initialize() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := &G2configmgr{}
	instanceName := "Test module name"
	settings, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	err = g2configmgr.Initialize(ctx, instanceName, settings, verboseLogging)
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
