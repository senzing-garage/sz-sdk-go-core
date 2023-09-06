package main

import (
	"context"
	"fmt"

	"github.com/senzing/g2-sdk-go-base/g2config"
	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-common/truthset"
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func getG2config(ctx context.Context) (g2api.G2config, error) {
	result := g2config.G2config{}
	moduleName := "Test module name"
	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	fmt.Printf(">>>>>> iniParams: %s\n", iniParams)
	if err != nil {
		return &result, err
	}
	err = result.Init(ctx, moduleName, iniParams, verboseLogging)
	return &result, err
}

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {
	var err error = nil
	ctx := context.TODO()

	// Get Senzing objects for installing a Senzing Engine configuration.

	fmt.Printf(">>>>>> Step 1.0\n")
	g2Config, err := getG2config(ctx)
	if err != nil {
		panic(err.Error())
	}

	// Using G2Config: Create a default configuration in memory.

	fmt.Printf(">>>>>> Step 2.0\n")
	configHandle, err := g2Config.Create(ctx)
	if err != nil {
		panic(err.Error())
	}

	// Using G2Config: Add data source to in-memory configuration.

	fmt.Printf(">>>>>> Step 3.0\n")
	ix := 0
	for _, testDataSource := range truthset.TruthsetDataSources {
		ix += 1
		fmt.Printf(">>>>>>    Step 3.%d: %s\n", ix, testDataSource.Json)
		_, err := g2Config.AddDataSource(ctx, configHandle, testDataSource.Json)
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
}
