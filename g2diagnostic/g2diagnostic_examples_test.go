//go:build linux

package g2diagnostic

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2diagnostic_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2diagnostic_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	result := g2diagnostic.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2diagnostic_CheckDBPerf() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	secondsToRun := 1
	result, err := g2diagnostic.CheckDBPerf(ctx, secondsToRun)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 25))
	// Output: {"numRecordsInserted":...
}

func ExampleG2diagnostic_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := &G2diagnostic{}
	err := g2diagnostic.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2diagnostic_Init() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := &G2diagnostic{}
	moduleName := "Test module name"
	iniParams, err := getIniParams()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	err = g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2diagnostic_InitWithConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := &G2diagnostic{}
	moduleName := "Test module name"
	iniParams, err := getIniParams()
	if err != nil {
		fmt.Println(err)
	}
	initConfigID := int64(1)
	verboseLogging := int64(0)
	err = g2diagnostic.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2diagnostic_Reinit() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	configID := getDefaultConfigID()
	err := g2diagnostic.Reinit(ctx, configID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2diagnostic_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	err := g2diagnostic.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
