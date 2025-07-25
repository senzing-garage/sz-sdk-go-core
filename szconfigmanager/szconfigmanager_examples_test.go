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

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
		return
	}

	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, configID)
	if err != nil {
		handleError(err)
		return
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.Truncate(configDefinition, 7))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,...
}

func ExampleSzconfigmanager_CreateConfigFromString() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	szConfigFromTemplate, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
		return
	}

	configDefinitionFromTemplate, err := szConfigFromTemplate.Export(ctx)
	if err != nil {
		handleError(err)
		return
	}

	szConfig, err := szConfigManager.CreateConfigFromString(ctx, configDefinitionFromTemplate)
	if err != nil {
		handleError(err)
		return
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.Truncate(configDefinition, 7))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,...
}

func ExampleSzconfigmanager_CreateConfigFromTemplate() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
		return
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.Truncate(configDefinition, 7))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,...
}

func ExampleSzconfigmanager_GetConfigRegistry() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	configList, err := szConfigManager.GetConfigRegistry(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.Truncate(configList, 3))
	// Output: {"CONFIGS":[{...
}

func ExampleSzconfigmanager_GetDefaultConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_RegisterConfig() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
		return
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
		return
	}

	configComment := "Example configuration"

	configID, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_ReplaceDefaultConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
		return
	}

	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, currentDefaultConfigID)
	if err != nil {
		handleError(err)
		return
	}

	dataSourceCodes := []string{"TEST_DATASOURCE"}
	for _, dataSource := range dataSourceCodes {
		_, err = szConfig.RegisterDataSource(ctx, dataSource)
		if err != nil {
			handleError(err)
			return
		}
	}

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
		return
	}

	configComment := "Configuration with TEST_DATASOURCE"

	newConfigID, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
		return
	}

	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newConfigID)
	if err != nil {
		handleError(err)
		return
	}
	// Output:
}

func ExampleSzconfigmanager_SetDefaultConfig() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	defaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
		return
	}

	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, defaultConfigID)
	if err != nil {
		handleError(err)
		return
	}

	dataSourceCode := "GO_TEST"

	_, err = szConfig.RegisterDataSource(ctx, dataSourceCode)
	if err != nil {
		handleError(err)
		return
	}

	configComment := "Added datasource: " + dataSourceCode

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		handleError(err)
		return
	}

	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_SetDefaultConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	configID, err := szConfigManager.GetDefaultConfigID(
		ctx,
	) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		handleError(err)
		return
	}

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		handleError(err)
		return
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
		return
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
