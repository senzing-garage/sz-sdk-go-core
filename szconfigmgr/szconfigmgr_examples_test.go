//go:build linux

package szconfigmgr

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfigmgr_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfigmgr_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	result := g2configmgr.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzconfigmgr_AddConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2config := getSzconfig(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		text := err.Error()
		fmt.Println(text[len(text)-40:])
	}
	g2configmgr := getSzconfigmgr(ctx)
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

func ExampleSzconfigmgr_GetConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
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

func ExampleSzconfigmgr_GetConfigList() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
	configList, err := g2configmgr.GetConfigList(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configList, 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleSzconfigmgr_GetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
	defaultConfigId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(defaultConfigId > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmgr_ReplaceDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
	currentDefaultConfigId, err := g2configmgr.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	g2config := &szconfig.Szconfig{}
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

func ExampleSzconfigmgr_SetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
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

func ExampleSzconfigmgr_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
	err := g2configmgr.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzconfigmgr_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := &Szconfigmgr{}
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

func ExampleSzconfigmgr_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmgr/szconfigmgr_examples_test.go
	ctx := context.TODO()
	g2configmgr := getSzconfigmgr(ctx)
	err := g2configmgr.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
