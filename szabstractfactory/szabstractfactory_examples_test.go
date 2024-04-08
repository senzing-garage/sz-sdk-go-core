//go:build linux

package szabstractfactory

import (
	"context"
	"fmt"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzAbstractFactory_CreateConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	config, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer config.Destroy(ctx)
	// Output:
}

func ExampleSzAbstractFactory_CreateConfigManager() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	configManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer configManager.Destroy(ctx)
	// Output:
}

func ExampleSzAbstractFactory_CreateDiagnostic() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	diagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer diagnostic.Destroy(ctx)
	// Output:
}

func ExampleSzAbstractFactory_CreateEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	engine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer engine.Destroy(ctx)
	// Output:
}

func ExampleSzAbstractFactory_CreateProduct() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	product, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer product.Destroy(ctx)
	// Output:
}
