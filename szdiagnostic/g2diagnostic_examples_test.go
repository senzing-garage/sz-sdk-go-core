//go:build linux

package szdiagnostic

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzdiagnostic_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	result := g2diagnostic.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzdiagnostic_CheckDatabasePerformance() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	secondsToRun := 1
	result, err := g2diagnostic.CheckDatabasePerformance(ctx, secondsToRun)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 25))
	// Output: {"numRecordsInserted":...
}

func ExampleSzdiagnostic_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := &Szdiagnostic{}
	err := g2diagnostic.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzdiagnostic_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := &Szdiagnostic{}
	instanceName := "Test module name"
	configId := int64(0) // "0" means use the current default configuration ID.
	settings, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	err = g2diagnostic.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzdiagnostic_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	configId := getDefaultConfigId()
	err := g2diagnostic.Reinitialize(ctx, configId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzdiagnostic_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	err := g2diagnostic.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
