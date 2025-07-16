//go:build linux

package szconfig_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfig_RegisterDataSource() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
	}

	dataSourceCode := "GO_TEST"

	result, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "DSRC_ID": 1001
	// }
}

func ExampleSzconfig_UnregisterDataSource() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
	}

	dataSourceCode := "GO_TEST"

	result, err := szConfig.UnregisterDataSource(ctx, dataSourceCode)
	if err != nil {
		handleError(err)
	}

	fmt.Println(result)
	// Output:
}

func ExampleSzconfig_Export() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

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

func ExampleSzconfig_GetDataSourceRegistry() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		handleError(err)
	}

	result, err := szConfig.GetDataSourceRegistry(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCES": [
	//         {
	//             "DSRC_ID": 1,
	//             "DSRC_CODE": "TEST"
	//         },
	//         {
	//             "DSRC_ID": 2,
	//             "DSRC_CODE": "SEARCH"
	//         }
	//     ]
	// }
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzconfig_SetLogLevel() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)

	err := szConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfig_SetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfig_GetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	result := szConfig.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}
