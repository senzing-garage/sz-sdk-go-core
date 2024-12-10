//go:build linux

package szabstractfactory

import (
	"context"
	"fmt"

	"github.com/senzing-garage/sz-sdk-go/senzing"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzabstractfactory_CreateConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactoryExample(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	_ = szConfig // szConfig can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateConfigManager() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactoryExample(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		fmt.Println(err)
	}
	_ = szConfigManager // szConfigManager can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateDiagnostic() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactoryExample(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		fmt.Println(err)
	}
	_ = szDiagnostic // szDiagnostic can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactoryExample(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	_ = szEngine // szEngine can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateProduct() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactoryExample(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		fmt.Println(err)
	}
	_ = szProduct // szProduct can now be used.
	// Output:
}

func ExampleSzabstractfactory_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactoryExample(ctx)
	err := szAbstractFactory.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzabstractfactory_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactoryExample(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = szAbstractFactory.Reinitialize(ctx, configID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getSzAbstractFactoryExample(ctx context.Context) senzing.SzAbstractFactory {
	result, err := getSzAbstractFactory(ctx)
	if err != nil {
		panic(err)
	}
	return result
}
