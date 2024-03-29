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

func ExampleSzproduct_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzproduct_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2product_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
	result := g2product.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzproduct_Initialize() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	instanceName := "Test module name"
	settings, err := getIniParams()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	g2product.Initialize(ctx, instanceName, settings, verboseLogging)
	// Output:
}

func ExampleSzproduct_GetLicense() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	result, err := g2product.GetLicense(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jutil.Flatten(jutil.Redact(result, "customer", "contract", "issueDate", "licenseLevel", "billing", "licenseType", "expireDate", "recordLimit")))
	// Output: {"billing":null,"contract":null,"customer":null,"expireDate":null,"issueDate":null,"licenseLevel":null,"licenseType":null,"recordLimit":null}
}

func ExampleSzproduct_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	err := g2product.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzproduct_GetVersion() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	result, err := g2product.GetVersion(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 43))
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":...
}

func ExampleSzproduct_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	err := g2product.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
