//go:build linux

package szproduct_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
)

const AllLines = -1

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzproduct_GetLicense() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szProduct.Destroy(ctx)) }()

	result, err := szProduct.GetLicense(ctx)
	if err != nil {
		handleError(err)
		return
	}

	redactKeys := []string{"issueDate", "expireDate", "BUILD_VERSION"}
	fmt.Println(jsonutil.PrettyPrint(jsonutil.Truncate(result, AllLines, redactKeys...), jsonIndentation))
	// Output:
	// {
	//    "advSearch": 0,
	//    "billing": "YEARLY",
	//    "contract": "Senzing Public Test License",
	//    "customer": "Senzing Public Test License",
	//    "licenseLevel": "STANDARD",
	//    "licenseType": "EVAL (Solely for non-productive use)",
	//    "recordLimit": 50000
	// }
}

func ExampleSzproduct_GetVersion() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szProduct.Destroy(ctx)) }()

	result, err := szProduct.GetVersion(ctx)
	if err != nil {
		handleError(err)
		return
	}

	redactKeys := []string{"BUILD_DATE", "BUILD_NUMBER", "BUILD_VERSION", "ENGINE_SCHEMA_VERSION", "VERSION"}
	fmt.Println(jsonutil.PrettyPrint(jsonutil.Truncate(result, AllLines, redactKeys...), jsonIndentation))
	// Output:
	// {
	//     "COMPATIBILITY_VERSION": {
	//         "CONFIG_VERSION": "11"
	//     },
	//     "PRODUCT_NAME": "Senzing SDK",
	//     "SCHEMA_VERSION": {
	//         "MAXIMUM_REQUIRED_SCHEMA_VERSION": "4.99",
	//         "MINIMUM_REQUIRED_SCHEMA_VERSION": "4.0"
	//     }
	// }
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzproduct_SetLogLevel() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)

	err := szProduct.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
		return
	}
	// Output:
}

func ExampleSzproduct_SetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzproduct_GetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	result := szProduct.GetObserverOrigin(ctx)

	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}
