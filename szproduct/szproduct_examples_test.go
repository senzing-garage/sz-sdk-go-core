//go:build linux

package szproduct

import (
	"context"
	"fmt"

	jutil "github.com/senzing-garage/go-common/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzProduct_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzProduct_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	result := szProduct.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzProduct_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	instanceName := "Test module name"
	settings, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	// Output:
}

func ExampleSzProduct_GetLicense() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	result, err := szProduct.GetLicense(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Redact(result, "customer", "contract", "issueDate", "licenseLevel", "billing", "licenseType", "expireDate", "recordLimit")))
	// Output: {"billing":null,"contract":null,"customer":null,"expireDate":null,"issueDate":null,"licenseLevel":null,"licenseType":null,"recordLimit":null}
}

func ExampleSzProduct_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	err := szProduct.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzProduct_GetVersion() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	result, err := szProduct.GetVersion(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 43))
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":...
}

func ExampleSzProduct_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	err := szProduct.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
