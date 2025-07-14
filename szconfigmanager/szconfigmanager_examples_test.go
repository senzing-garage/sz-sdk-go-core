//go:build linux

package szconfigmanager_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_CreateConfigFromConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}

	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, configID)
	if err != nil {
		handleError(err)
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.Truncate(configDefinition, 7))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,...
}

func ExampleSzconfigmanager_CreateConfigFromString() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	szConfigFromTemplate, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
	}

	configDefinitionFromTemplate, err := szConfigFromTemplate.Export(ctx)
	if err != nil {
		handleError(err)
	}

	szConfig, err := szConfigManager.CreateConfigFromString(ctx, configDefinitionFromTemplate)
	if err != nil {
		handleError(err)
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.Truncate(configDefinition, 7))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,...
}

func ExampleSzconfigmanager_CreateConfigFromTemplate() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.Truncate(configDefinition, 7))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,...
}

func ExampleSzconfigmanager_GetConfigRegistry() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	configList, err := szConfigManager.GetConfigRegistry(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.Truncate(configList, 3))
	// Output: {"CONFIGS":[{...
}

func ExampleSzconfigmanager_GetDefaultConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_RegisterConfig() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
	}

	configComment := "Example configuration"

	configID, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
	}

	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_ReplaceDefaultConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}

	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, currentDefaultConfigID)
	if err != nil {
		handleError(err)
	}

	dataSourceCodes := []string{"TEST_DATASOURCE"}
	for _, dataSource := range dataSourceCodes {
		_, err = szConfig.RegisterDataSource(ctx, dataSource)
		if err != nil {
			handleError(err)
		}
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
	}

	configComment := "Configuration with TEST_DATASOURCE"

	newConfigID, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
	}

	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newConfigID)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfigmanager_SetDefaultConfig() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	defaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}

	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, defaultConfigID)
	if err != nil {
		handleError(err)
	}

	dataSourceCode := "GO_TEST"

	_, err = szConfig.RegisterDataSource(ctx, dataSourceCode)
	if err != nil {
		handleError(err)
	}

	configComment := "Added datasource: " + dataSourceCode

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
	}

	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
	}

	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_SetDefaultConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	configID, err := szConfigManager.GetDefaultConfigID(
		ctx,
	) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		handleError(err)
	}

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		handleError(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_SetLogLevel() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)

	err := szConfigManager.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfigmanager_SetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfigmanager_GetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	result := szConfigManager.GetObserverOrigin(ctx)

	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}
